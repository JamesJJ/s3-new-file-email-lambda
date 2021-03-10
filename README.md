# New file email notification

This is an AWS lambda function written in Go, that takes S3 object creation events and sends a templated email using AWS SES to a specified recepient. The email includes a presigned link to the object in S3.

## Deployment using AWS SAM

The [aws-serverless](./aws-serverless) directory includes a [template](./aws-serverless/template.yml) and [config](./aws-serverless/samconfig.toml) to deploy this app using the [AWS Serverless Application Model](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html).

### Prerequisites

You'll need:

 - `go` installed
 - [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) installed
 - AWS CLI credentials set-up

### Usage

 1. Execute `./aws-serverless/deploy.sh` 
 2. Enter parameter values when prompted:
   - `SesIdentityArn` The ARN of the SES domain identity that corresponds to the `MailFrom` address you want to use ([performing the SES domain / from address validation](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/verify-addresses-and-domains.html) is outside the scope of this project)
   - `MailFrom` The email address to send email from
   - `MailTo` The email address of the notification recepient (Must be 'verified' in SES if you are in the [AWS Sandbox](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/request-production-access.html)) 
   - `S3Bucket` The name of the S3 bucket that will invoke this lambda function
 3. When the stack has finished deploying, it will show the `LambdaARN` in the outputs. Please [set the S3 bucket to send](https://docs.aws.amazon.com/AmazonS3/latest/userguide/enable-event-notifications.html) `s3:ObjectCreated:Put` events to this Lambda ARN.

