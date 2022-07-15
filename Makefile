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

# "go.useLanguageServer": false
# gopls -rpc.trace -v check path/to/file.go



# build: clean
# 	@for dir in `ls handler`; do \
# 		GOOS=linux go build -o dist/$$dir github.com/pulpfree/$(PROJECT_NAME)/handler/$$dir; \
# 	done
# 	@GOOS=linux go build -o dist/authorizer github.com/pulpfree/$(PROJECT_NAME)/authorizer;
# 	@cp ./config/defaults.yml dist/
# 	@echo "build successful"

# "create-bucket-configuration" line only if not in us-east-1, or apparently some other regions as well...
configure:
	@ aws s3api create-bucket \
		--bucket $(AWS_BUCKET_NAME) \
		--region $(AWS_REGION) \
		--create-bucket-configuration LocationConstraint=$(AWS_REGION)

# watch: Run given command when code changes. e.g; make watch run="echo 'hey'"
# @yolo -i . -e vendor -e bin -e dist -c $(run)

watch:
	@yolo -i . -e vendor -e dist -c "make build"

# run: build
# 	sam local start-api -n env.json



# awspackage:
# 	@aws cloudformation package \
#    --template-file ${FILE_TEMPLATE} \
#    --output-template-file ${FILE_PACKAGE} \
#    --s3-bucket $(AWS_LAMBDA_BUCKET) \
#    --s3-prefix $(AWS_BUCKET_PREFIX) \
#    --profile $(AWS_PROFILE) \
#    --region $(AWS_REGION)

# awsdeploy:
# 	@aws cloudformation deploy \
# 	--template-file ${FILE_PACKAGE} \
# 	--region $(AWS_REGION) \
# 	--stack-name $(AWS_STACK_NAME) \
# 	--capabilities CAPABILITY_IAM \
# 	--profile $(AWS_PROFILE) \
# 	--force-upload \
# 	--parameter-overrides \
# 		ParamCertificateArn=$(CERTIFICATE_ARN) \
# 		ParamCustomDomainName=$(CUSTOM_DOMAIN_NAME) \
# 		ParamHostedZoneId=$(HOSTED_ZONE_ID) \
# 		ParamKMSKeyID=$(KMS_KEY_ID) \
# 		ParamProjectName=$(PROJECT_NAME) \
# 		ParamReportBucket=${AWS_REPORT_BUCKET} \
# 		ParamSecurityGroupIds=$(SECURITY_GROUP_IDS) \
# 		ParamSubnetIds=$(SUBNET_IDS)

describe:
	@aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(AWS_STACK_NAME)

outputs:
	@ make describe \
		| jq -r '.Stacks[0].Outputs'
