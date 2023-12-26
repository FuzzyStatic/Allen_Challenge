package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/acm"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	synced "github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const domain = "sed.flick.is" // This would be replaced by an environment variable to prefix to a proper domain. EX: dev.example.com / stage.example.com, etc.

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		path := "./www"
		if param := cfg.Get("path"); param != "" {
			path = param
		}
		indexDocument := "index.html"
		if param := cfg.Get("indexDocument"); param != "" {
			indexDocument = param
		}
		errorDocument := "error.html"
		if param := cfg.Get("errorDocument"); param != "" {
			errorDocument = param
		}

		// Create an S3 bucket.
		bucket, err := s3.NewBucket(ctx, "bucket", &s3.BucketArgs{
			Website: &s3.BucketWebsiteArgs{
				IndexDocument: pulumi.String(indexDocument),
				ErrorDocument: pulumi.String(errorDocument),
			},
		})
		if err != nil {
			return err
		}

		// Set ownership controls for the new S3 bucket
		ownershipControls, err := s3.NewBucketOwnershipControls(ctx, "ownership-controls", &s3.BucketOwnershipControlsArgs{
			Bucket: bucket.Bucket,
			Rule: &s3.BucketOwnershipControlsRuleArgs{
				ObjectOwnership: pulumi.String("ObjectWriter"),
			},
		})
		if err != nil {
			return err
		}

		// Configure public access block for the new S3 bucket
		publicAccessBlock, err := s3.NewBucketPublicAccessBlock(ctx, "public-access-block", &s3.BucketPublicAccessBlockArgs{
			Bucket:          bucket.Bucket,
			BlockPublicAcls: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// Use a synced folder to manage the files of the website
		_, err = synced.NewS3BucketFolder(ctx, "bucket-folder", &synced.S3BucketFolderArgs{
			Path:       pulumi.String(path),
			BucketName: bucket.Bucket,
			Acl:        pulumi.String("private"),
		}, pulumi.DependsOn([]pulumi.Resource{ownershipControls, publicAccessBlock}))
		if err != nil {
			return err
		}

		oac, err := cloudfront.NewOriginAccessControl(ctx, "oac", &cloudfront.OriginAccessControlArgs{
			Name:                          pulumi.String("S3OriginAccessControl"),
			OriginAccessControlOriginType: pulumi.String("s3"),
			SigningBehavior:               pulumi.String("always"),
			SigningProtocol:               pulumi.String("sigv4"),
		})
		if err != nil {
			return err
		}

		// Request a new certificate for the desired domain
		_, err = acm.NewCertificate(ctx, "certificate", &acm.CertificateArgs{
			DomainName:       pulumi.String(domain),
			ValidationMethod: pulumi.String("DNS"),
		})
		if err != nil {
			return err
		}

		// Insert DNS validation glue for certificate here.
		// This may take some time and will be a blocking operation.
		// Maybe this happens first before stack deployments.

		// Create a CloudFront CDN to distribute and cache the website
		cdn, err := cloudfront.NewDistribution(ctx, "cdn", &cloudfront.DistributionArgs{
			Aliases:           pulumi.ToStringArray([]string{domain}),
			Enabled:           pulumi.Bool(true),
			DefaultRootObject: pulumi.String("index.html"),
			Origins: cloudfront.DistributionOriginArray{
				&cloudfront.DistributionOriginArgs{
					OriginId:              bucket.Arn,
					DomainName:            bucket.BucketDomainName,
					OriginAccessControlId: oac.ID(),
				},
			},
			DefaultCacheBehavior: &cloudfront.DistributionDefaultCacheBehaviorArgs{
				TargetOriginId:       bucket.Arn,
				ViewerProtocolPolicy: pulumi.String("redirect-to-https"),
				AllowedMethods: pulumi.StringArray{
					pulumi.String("GET"),
					pulumi.String("HEAD"),
					pulumi.String("OPTIONS"),
				},
				CachedMethods: pulumi.StringArray{
					pulumi.String("GET"),
					pulumi.String("HEAD"),
					pulumi.String("OPTIONS"),
				},
				DefaultTtl: pulumi.Int(600),
				MaxTtl:     pulumi.Int(600),
				MinTtl:     pulumi.Int(600),
				ForwardedValues: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesArgs{
					QueryString: pulumi.Bool(true),
					Cookies: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesCookiesArgs{
						Forward: pulumi.String("all"),
					},
				},
			},
			PriceClass: pulumi.String("PriceClass_100"),
			Restrictions: &cloudfront.DistributionRestrictionsArgs{
				GeoRestriction: &cloudfront.DistributionRestrictionsGeoRestrictionArgs{
					RestrictionType: pulumi.String("none"),
				},
			},
			ViewerCertificate: &cloudfront.DistributionViewerCertificateArgs{
				AcmCertificateArn:      pulumi.String("arn:aws:acm:us-east-1:<redacted>:certificate/<redacted>"), // This would be done programmatically using the certificate created above
				MinimumProtocolVersion: pulumi.String("TLSv1.2_2021"),
				SslSupportMethod:       pulumi.String("sni-only"),
			},
		})
		if err != nil {
			return err
		}

		// Custom S3 policy granting the CDN access to the S3 bucket
		bucketPolicy := pulumi.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "AllowCloudFrontServicePrincipalReadOnly",
					"Effect": "Allow",
					"Principal": {
						"Service": "cloudfront.amazonaws.com"
					},
					"Action": "s3:GetObject",
					"Resource": "arn:aws:s3:::%s/*",
					"Condition": {
						"StringEquals": {
							"AWS:SourceArn": "%s"
						}
					}
				}
			]
		}`, bucket.Bucket, cdn.Arn)
		s3.NewBucketPolicy(ctx, "bucket-policy", &s3.BucketPolicyArgs{
			Bucket: bucket.Bucket,
			Policy: bucketPolicy,
		})

		// Export the URLs and hostnames of the bucket and distribution.
		ctx.Export("cdnURL", pulumi.Sprintf("https://%s", cdn.DomainName))
		ctx.Export("cdnHostname", cdn.DomainName)
		return nil
	})
}
