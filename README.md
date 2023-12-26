# Allen_Challenge

## SRE/SED Challenge Execution

For this challenge I'm going to use Pulumi. Why Pulumi? Pulumi is an IaC service, similar to AWS CDK, that lets user code deployment stacks. This has the benefit of being about to create and validate multiple stacks with a shared set of logic amongst other benefits. Also, I've never used Pulumi before and it is fun to learn new things.

### Setup `pulumi` Project

Install `pulumi`:
```shell
brew install pulumi/tap/pulumi
```

Set AWS credentials:
```shell
export AWS_ACCESS_KEY_ID=<YOUR_ACCESS_KEY_ID>
export AWS_SECRET_ACCESS_KEY=<YOUR_SECRET_ACCESS_KEY>
```

Setup a new project:
```shell
pulumi new aws-go
```

### Design

First let's borrow a nice starting [template](https://github.com/pulumi/templates/tree/master/static-website-aws-go). This template uses S3's static website functionality and CloudFront to host a static website. S3 has the advantage of being globally reachable and CloudFront brings cached resources closer to the user for better latency.

This template seems to have some downsides. While CloudFront is able to pull and cache the S3 bucket data, S3 is set for public consumption. Let's change this so access is only allowed through CloudFront.

#### AWS Changes
- Set S3 bucket to `private`.
- Creating CloudFront Origin Access Control (OAC). This allows CloudFront to access the S3 bucket without using an S3 website endpoint.
- Update CloudFront configuration to use OAC.
- Update CloudFront to have a default root object. This will display the `index.html` data if/when no path is specified.
- Add a Bucket Policy to the S3 bucket to that allows CloudFront to access its resources.

All this results in a secure cloudfront link prefixed with some random identifier, but that isn't very easy to remember. Instead, let's use a CNAME using a subdomain of a domain I already own. The new site will be <https://sed.flick.is>.

#### Certificate Changes
- Request a certificate for `sed.flick.is` via AWS Certificate Manager (ACM).
- Validate the certificate request through DNS. This will be done manually given the DNS provider I currently use, but can easily be glued together. Some solutions would be to use AWS hosted zones in Route 53 and a CrossDomainDelegation policy for validation or leveraging APIs for Comcast's internal DNS service.
- Wait for validation. This should take a couple minutes and then <https://sed.flick.is> is accessible.

### Helpful Commands
 
Uncache CDN:
```shell
aws cloudfront create-invalidation --distribution-id=DISTRIBUTION_ID --paths /
```
> I had to do this a few times to verify permission changes were working as intended.

### Closing Thoughts

The challenge asked for a secure static web application in AWS. The defining reason for choosing S3 and CloudFront stems from the word `static`. More dynamic websites or backend services would have been approached differently. The challenge mentions to secure the application to only the appropriate ports which indicates that managing and attaching security groups would be necessary, but was not needed here. In short, the challenge seems to want to steer someone to use services like EC2 with a ALB, but I took a different approach.

## Coding Challenge Execution

The coding challenge can be found in the `cc` folder. Inside is a `cc` package with included test file. In the example folder is a `main.go` that will loop through a slice and print `Valid` or `Invalid` based on the input.

```shell
$ cd cc
$ go test -v 
=== RUN   TestIsValid
=== RUN   TestIsValid/6_validity_is_false
=== RUN   TestIsValid/4123456789123456_validity_is_true
=== RUN   TestIsValid/5123456789123456_validity_is_true
=== RUN   TestIsValid/6123456789123456_validity_is_true
=== RUN   TestIsValid/5123-4567-8912-3456_validity_is_true
=== RUN   TestIsValid/61234-567-8912-3456_validity_is_false
=== RUN   TestIsValid/5100-0067-8912-3456_validity_is_false
=== RUN   TestIsValid/5111-1167-8912-3456_validity_is_false
=== RUN   TestIsValid/5122-2267-8912-3456_validity_is_false
=== RUN   TestIsValid/5133-3367-8912-3456_validity_is_false
=== RUN   TestIsValid/5144-4467-8912-3456_validity_is_false
=== RUN   TestIsValid/5155-5567-8912-3456_validity_is_false
=== RUN   TestIsValid/5166-6667-8912-3456_validity_is_false
=== RUN   TestIsValid/5177-7767-8912-3456_validity_is_false
=== RUN   TestIsValid/5188-8867-8912-3456_validity_is_false
=== RUN   TestIsValid/5199-9967-8912-3456_validity_is_false
=== RUN   TestIsValid/5123_-_3567_-_8912_-_3456_validity_is_false
=== RUN   TestIsValid/5133-336789123456_validity_is_false
--- PASS: TestIsValid (0.00s)
    --- PASS: TestIsValid/6_validity_is_false (0.00s)
    --- PASS: TestIsValid/4123456789123456_validity_is_true (0.00s)
    --- PASS: TestIsValid/5123456789123456_validity_is_true (0.00s)
    --- PASS: TestIsValid/6123456789123456_validity_is_true (0.00s)
    --- PASS: TestIsValid/5123-4567-8912-3456_validity_is_true (0.00s)
    --- PASS: TestIsValid/61234-567-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5100-0067-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5111-1167-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5122-2267-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5133-3367-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5144-4467-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5155-5567-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5166-6667-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5177-7767-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5188-8867-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5199-9967-8912-3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5123_-_3567_-_8912_-_3456_validity_is_false (0.00s)
    --- PASS: TestIsValid/5133-336789123456_validity_is_false (0.00s)
PASS
ok      Allen_Challenge/cc      0.378s

$ cd example
$ go run *.go             
Invalid
Valid
Valid
Invalid
Invalid
Invalid
Valid
Invalid
Invalid
Invalid
Invalid
Invalid
Valid
Invalid
Invalid
Invalid
Invalid
Invalid
Invalid
Invalid
```