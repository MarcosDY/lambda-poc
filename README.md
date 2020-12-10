# lambda-poc

This POC is aimed to demonstrate that is possible to consume an X509 SVID persisted in a Secret inside AWS.

This secrets has a marshaled `workload.X509SVIDResponse` proto added and updated for SPIRE Agent.

## Components:  

### Extension:

Simple extension that reads a Secret from Secret manager

*NOTE: secret MUST have a workload.X509SVIDResponse as binary.*

Secret name or ARN must be provided using envvar "SECRET_NAME" on function. 
Once secret is parsed svids is persisted on '/tmp' folder.

### Function:

Python function that read SVID from disk and print it into log.

## Scripts:

* [00-var.sh](./00-vars.sh): Contains all variables used to run POC, it is important to update 'FUNCTION_ROLE' with a valid execution role, it MUST have access to Secret
* [01-build.sh](./01-build.sh): Build extension and function
* [02-publish-layer.sh](./02-publish-layer.sh): publish a new version of SPIRE Extension
* [03-create-function](./03-create-function.sh): create POC function into configure AWS Region
* [04-run.sh](./04-run.sh): invoke POC function
* [05-cleanup.sh](./05-cleanup.sh): remove function, extension and logs, from AWS

Util scripts
* [describe-secret.sh](./describe-secret.sh): describes secret that function is using
* [get-logs.sh](./get-logs.sh): tail the last hour logs
* [update-function](./update-function.sh): update function in case it exists

## SPIRE changes

SPIRE Agent was updated to be able to update secrets on AWS, it is using a new kind of plugins 'SVIDStore' that is notified every time that an 'storable' SVID is updated, and update Secret Binary, with an updated plugin.

![SPIRE Diagram](./images/agent-pusher-pipe.png)

### Entry example:

```
Entry ID      : e1d52419-a2aa-4ca1-abd6-e714a46ae77d
SPIFFE ID     : spiffe://example.org/aws/workload1
Parent ID     : spiffe://example.org/agent
Revision      : 6
TTL           : default 
Selector      : aws:name:SVID_Example
Selector      : store:aws_secretsmanager
FederatesWith : federated.td2
```

* `store:aws_secretsmanager`: `store` is the key for all storable SVIDs, `aws_secretsmanager` is the plugin name used to `store` SVID
* `aws:name:SVID_Example`: is the revelant information to update secret, in this case we say that the secrete name is `SVID_Example`, it is possibe to use `aws:name` or `aws:arn` as secret ID.
