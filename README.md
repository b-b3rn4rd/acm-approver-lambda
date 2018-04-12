[![Build Status](https://travis-ci.org/b-b3rn4rd/acm-approver-lambda.svg?branch=master)](https://travis-ci.org/b-b3rn4rd/acm-approver-lambda) [![Go Report Card](https://goreportcard.com/badge/github.com/b-b3rn4rd/acm-approver-lambda)](https://goreportcard.com/report/github.com/b-b3rn4rd/acm-approver-lambda) AWS CloudFormation ACM Approver Golang Custom Resource
=====================

AWS Lambda function &mdash; approves ACM certificates issued with DNS validation option.
Following lambda is written as a custom resource to automate certificate approval process in a stack.


Installation & Usage
----------------------------
Download code:

`git clone https://github.com/b-b3rn4rd/acm-approver-lambda.git`


Create CloudFormation stack
```
$ S3_BUCKET_NAME=bucket-name DOMAIN_NAME=www.example.net make deploy

... ouput ....
Waiting for changeset to be created..
Waiting for stack create/update to complete
Successfully created/updated stack - acm-approver-lamda
```

Following command will create CloudFormation stack, which provisions lambda function and invokes it as a custom resource
to request and confirm required certificate.

Known issues
---------------------
I have not found a way to 100%  accurately identify hosted zone id based on certificate's domain name, currently I'm using longest match suffix approach.
