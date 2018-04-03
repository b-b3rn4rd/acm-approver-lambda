FUNCTION_ALIAS ?= prd
S3_BUCKET_MAME ?= ""
STACK_NAME ?= acm-approver-lamda

install:
	go get gopkg.in/alecthomas/gometalinter.v2
	gometalinter.v2 --install
.PHONY: install

mocks:
	go get -u -v github.com/aws/aws-sdk-go/...
	mockery -name ACMAPI -dir ../../aws/aws-sdk-go/service/acm/acmiface -recursive
	mockery -name Route53API -dir ../../aws/aws-sdk-go/service/route53/route53iface -recursive
	mockery -name Certificate -dir pkg/certificate -recursive
.PHONY: mocks

pre_build: install mocks
	go get -t ./cmd/lambda
	gometalinter.v2 ./...
	go test -v ./...
.PHONY: pre_build

build: pre_build
	GOOS=linux GOARCH=amd64 go build -o main ./cmd/lambda
	@zip -9 -r ./handler.zip main
.PHONY: build

deploy: build
	aws cloudformation package \
		--template-file cfn.yaml \
		--output-template-file cfn.out.yaml \
		--s3-bucket $(S3_BUCKET_MAME) \
		--s3-prefix cfn

	aws cloudformation deploy \
		--template-file cfn.out.yaml \
		--capabilities CAPABILITY_IAM \
		--stack-name $(STACK_NAME) \
        --parameter-overrides \
        	FunctionAlias=$(FUNCTION_ALIAS) \
        	FunctionName=$(STACK_NAME)
.PHONY: deploy