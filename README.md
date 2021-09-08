# AWS Lambda Example

This is a proof of concept to demonstrate how serverless computing can be supported in SPIRE through the introduction of an `SVIDStore` agent plugin.

The model leverages the use of secret management services offered by cloud providers to store and retrieve the SVIDs and keys in a secure way, inside the cloud infrastructure.

The serverless functions are registered in SPIRE in the same way that regular workloads are registered through registration entries. The `svidstore` key is used to distinguish the "storable" entries, and `SVIDStore` plugins receive updates of those entries only, which indicates that the issued SVID and key must be securely stored in a location accessible by the serverless function, like AWS Secrets Manager. This way, selectors provide a flexible way to describe the attributes needed to store the corresponding issued SVID and key, like the type of store, name to provide to the secret, and any specific attribute needed by the specific service used.

## Components

### AWS Lambda Extension

Simple extension that reads a secret from AWS Secrets Manager.

*NOTE: it is expected that the secret is a binary JSON message.*

```
{
	"spiffeId": "spiffe://example.org/dbuser",
	"x509Svid": "PEM_CERTIFICATES",
	"x509SvidKey": "PEM_KEY",
	"bundle": "PEM_BUNDLE",
	"federatedBundles": {
		"spiffe://federated.org": "PEM_CERTIFICATE"
	}
}
```

The secret name or ARN must be provided using the environment variable "SECRET_NAME" in the function. 
The secret is unmarshaled and the X509-SVID, bundle and key are persisted in the `/tmp` folder.

### Function

The function itself is a Python function that reads the stored SVID from disk prints it so it can be read from the log and returns it in a JSON response.

## Scripts

* [00-var.sh](./00-vars.sh): Contains all the variables used to run this POC. 'FUNCTION_ROLE' must be updated with a valid execution role, ensuring that is has access to AWS Secrets Manager.
* [01-build.sh](./01-build.sh): Build extension and function.
* [02-publish-layer.sh](./02-publish-layer.sh): Publish a new version of the extension.
* [03-create-function](./03-create-function.sh): Create the AWS Lambda functions in the configured AWS region.
* [04-run.sh](./04-run.sh): Invoke the AWS Lambda functions.
* [05-cleanup.sh](./05-cleanup.sh): Cleanup the AWS resources associated with this project (functions, layers, logs).
Util scripts
* [describe-secret.sh](./describe-secret.sh): Describe the stored secret, it requires secrets name.
* [get-logs.sh](./get-logs.sh): Tail the last hour logs, it requires function name.
* [update-function](./update-function.sh): Update the already created AWS Lambda functions, if exists.

## SPIRE changes

The SPIRE Agent cache manager was updated to be able to identify "storable" entries and notify the corresponding plugin when the entries are updated. A new `SVIDStore` agent plugin is introduced for this.

![SPIRE Diagram](./images/svid-store-aws.png)

### Entry example

```
Entry ID         : 7c141f95-db0e-4968-b10c-628d1f7fe5d7
SPIFFE ID        : spiffe://example.org/db
Parent ID        : spiffe://example.org/agent
Revision         : 0
TTL              : default
Selector         : aws_secretsmanager:secretname:db-svid
StoreSvid        : true

Entry ID         : 14a25fd1-d7c0-4cbd-b0f6-936df944ba59
SPIFFE ID        : spiffe://example.org/web
Parent ID        : spiffe://example.org/agent
Revision         : 0
TTL              : default
Selector         : aws_secretsmanager:secretname:web-svid
StoreSvid        : true
```

* `StoreSvid` indicate that the issued SVID and key must be stored in a secure store.
* `aws_secretsmanager:secretname:web-svid` indicates that the entry must be stored using `aws_secretsmanager` plugin that creates Secrets on AWS Secret Manager with name `web-svid`.

A log must be displayed on SPIRE Agent when it creates a secret on AWS

Creation:
```
DEBU[0001] Secret created                                arn="SOME_ARN" external=false name=web-svid plugin_name=aws_secretsmanager plugin_type=SVIDStore subsystem_name=catalog version_id=e884ff64-a50a-470b-9b4c-ad15e41a78d4
```

Updated:
```
DEBU[0002] Secret value updated                          arn="SOME_ARN" external=false name=web-svid plugin_name=aws_secretsmanager plugin_type=SVIDStore subsystem_name=catalog version_id=640e7123-318b-48f6-acf1-9e9a5d306016
```

