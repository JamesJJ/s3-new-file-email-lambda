#!/bin/bash

set -eo pipefail

mkdir -p ./tmp

cd function
GOOS=linux go build -o ../tmp/main main.go
cd ../

CONFIG="./config.env"
if [ -f "$CONFIG" ] ; then
  source "$CONFIG"
fi

if [ -z "$ARTIFACT_BUCKET" ] ; then
	read -p 'Please type the S3 bucket name to store built artifacts and press enter: ' ARTIFACT_BUCKET
fi

if [ -z "$APP_MAILFROM" ] ; then
	read -p 'Please type the mail "from" address and press enter: ' APP_MAILFROM
fi

if [ -z "$APP_MAILTO" ] ; then
	read -p 'Please type the mail "to" and press enter: ' APP_MAILTO
fi

if [ -z "$APP_SESIDENTITYARN" ] ; then
	read -p 'Please type the SES sending identity ARN and press enter: ' APP_SESIDENTITYARN
fi

aws cloudformation package \
  --template-file template.yml \
  --s3-bucket "$ARTIFACT_BUCKET" \
  --output-template-file ./tmp/out.yml

aws cloudformation deploy \
  --template-file ./tmp/out.yml \
  --stack-name s3-new-file-email-lambda \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides "MailFrom=$APP_MAILFROM" "MailTo=$APP_MAILTO" "SesIdentityArn=$APP_SESIDENTITYARN"
