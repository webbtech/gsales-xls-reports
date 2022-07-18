# if the KEY environment variable is not set to either stage or prod, makefile will fail
# KEY is confirmed below in the check_env directive
# example:
# for stage run: ENV=stage make
# for production run: ENV=prod make
include .env

# found yolo at: https://azer.bike/journal/a-good-makefile-for-go/

default: check_env build \
	local-api

deploy: check_env build \
	upload-defaults \
	dev-cloud

check_env:
	@echo -n "Your environment is $(ENV)? [y/N] " && read ans && [ $${ans:-N} = y ]

upload-defaults:
	@ aws s3 cp ./config/xls-reports-defaults.yml s3://$(AWS_LAMBDA_BUCKET)/public/
	@ aws s3api put-object-tagging \
  --bucket $(AWS_LAMBDA_BUCKET) \
  --key public/xls-reports-defaults.yml \
  --tagging '{"TagSet": [{"Key": "public", "Value": "true"}]}' && \
	echo "defaults file uploaded and tagged"

build:
	sam build

# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-local-start-api.html
local-api:
	sam local start-api --env-vars env.json --profile $(PROFILE)

dev-cloud:
	sam sync --stack-name $(STACK_NAME) --profile $(PROFILE) \
	--s3-prefix $(AWS_DEPLOYMENT_PREFIX) \
	--parameter-overrides \
		ParamCertificateArn=$(CERTIFICATE_ARN) \
		ParamCustomDomainName=$(CUSTOM_DOMAIN_NAME) \
		ParamENV=$(ENV) \
		ParamHostedZoneId=$(HOSTED_ZONE_ID) \
		ParamReportBucket=$(S3_REPORT_BUCKET) \
		ParamSSMPath=$(SSM_PARAM_PATH) \
		ParamUserPoolArn=$(USER_POOL_ARN)

dev-cloud-watch:
	sam sync --stack-name $(STACK_NAME) --watch --profile $(PROFILE) \
	--s3-prefix $(AWS_DEPLOYMENT_PREFIX) \
	--parameter-overrides \
		ParamCertificateArn=$(CERTIFICATE_ARN) \
		ParamCustomDomainName=$(CUSTOM_DOMAIN_NAME) \
		ParamENV=$(ENV) \
		ParamHostedZoneId=$(HOSTED_ZONE_ID) \
		ParamReportBucket=$(S3_REPORT_BUCKET) \
		ParamUserPoolArn=$(USER_POOL_ARN)

tail-logs:
	sam logs -n ReportsFunction --profile $(PROFILE) \
	--stack-name $(STACK_NAME) --tail

tail-logs-trace:
	sam logs -n PdfUrlFunction --profile $(PROFILE) \
	--stack-name $(STACK_NAME) --tail --include-traces

validate:
	sam validate

test:
	@go test -v ./...

# ========================== non-used methods =================================
clean:
	@rm -rf dist
	@mkdir -p dist

# "create-bucket-configuration" line only if not in us-east-1, or apparently some other regions as well...
configure:
	@ aws s3api create-bucket \
		--bucket $(AWS_BUCKET_NAME) \
		--region $(AWS_REGION) \
		--create-bucket-configuration LocationConstraint=$(AWS_REGION)

watch:
	@yolo -i . -e vendor -e dist -c "make build"

describe:
	@aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(AWS_STACK_NAME)

outputs:
	@ make describe \
		| jq -r '.Stacks[0].Outputs'
