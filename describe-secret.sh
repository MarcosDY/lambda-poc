#!/bin/bash

SECRET_NAME=$1

. 00-vars.sh

# Get last hour logs
aws secretsmanager describe-secret \
	--region $AWS_REGION \
	--secret-id $SECRET_NAME
