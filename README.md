[![Go Report Card](https://goreportcard.com/badge/github.com/b-b3rn4rd/acm-approver-lambda)](https://goreportcard.com/report/github.com/b-b3rn4rd/acm-approver-lambda) AWS ACM Approver Lambda
=====================

AWS Lambda function &mdash; approves ACM certificates issued with DNS validation option.
Optionally, following lambda can be used as custom resource to automate certificate approval process.

Installation
----------------------------
`git clone https://github.com/b-b3rn4rd/acm-approver-lambda.git`

```
$ S3_BUCKET_MAME=my-bucket-name make deploy
... ouput ....
Waiting for changeset to be created..
Waiting for stack create/update to complete
Successfully created/updated stack - acm-approver-lamda
```

Usage Examples
-----------------------------
Once the function has been deployed, it can be executed manually, or as a custom resource

*Adhoc run*
```bash
$ aws lambda invoke \
    --function-name acm-approver-lamda \
    --payload '{"certificate-arn":"arn:aws:acm:us-west-2:1234567890123:certificate/cb4c2da3-cc3d-4142-8177-04f519117b33", "ttl":30}' \
    response.txt
```

*Custom Resource*
```bash

```


Known issues
---------------------
I have not found a way to 100%  accurately identify hosted zone id based on certificate's domain name, currently I'm using longest match suffix approach.
