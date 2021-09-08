#!/bin/bash

FUNCTION_NAME=$1

. 00-vars.sh

# Get last hour logs
aws logs tail /aws/lambda/$FUNCTION_NAME \
	--region $AWS_REGION \
	--since 1h
