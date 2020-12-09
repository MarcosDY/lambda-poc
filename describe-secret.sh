#!/bin/bash

ARN=$(cat .layer.arn.txt)

. 00-vars.sh

# Get last hour logs
aws secretsmanager describe-secret \
	--region $AWS_REGION \
	--secret-id $SECRET_NAME
