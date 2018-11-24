
HANDLER ?= handler
PACKAGE ?= $(HANDLER)
GOPATH  ?= $(HOME)/go

WORKDIR = $(CURDIR:$(GOPATH)%=/go%)
ifeq ($(WORKDIR),$(CURDIR))
	WORKDIR = /tmp
endif

ROLE_ARN=`aws iam get-role --role-name lambda_basic_execution --query 'Role.Arn' --output text`

all: build pack

build:
	@GOARCH=amd64 GOOS=linux go build -o $(HANDLER)

pack:
	@zip $(PACKAGE).zip $(HANDLER)

clean:
	@rm -rf $(HANDLER) $(PACKAGE).zip

create:
	@aws lambda create-function                                                  \
	  --function-name mibebe                                                \
	  --zip-file fileb://handler.zip                                             \
	  --role $(ROLE_ARN)                                                         \
	  --runtime go1.x                                                       \
	  --handler handler

update:
	@aws lambda update-function-code                                             \
	  --function-name mibebe                                                 \
	  --zip-file fileb://handler.zip

invoke:
	@aws lambda invoke                                                           \
	  --function-name mibebe invoke.out

.PHONY: all build pack clean create update invoke