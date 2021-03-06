#!/bin/bash

ARN=$(cat .layer.arn.txt)

. 00-vars.sh

# Get last hour logs
aws logs tail /aws/lambda/$FUNCTION_NAME \
	--region $AWS_REGION \
	--since 1h
