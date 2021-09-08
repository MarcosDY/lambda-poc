#!/bin/bash

. 00-vars.sh

aws lambda invoke \
	--function-name "${FUNCTION_DB}" \
	--region "${AWS_REGION}" \
	.db-lambda-response > /dev/null

cat .db-lambda-response | jq .cert -r | openssl x509 -text -noout | grep "URI:"

aws lambda invoke \
	--function-name "${FUNCTION_WEB}" \
	--region "${AWS_REGION}" \
	.web-lambda-response > /dev/null

cat .web-lambda-response | jq .cert -r | openssl x509 -text -noout | grep "URI:"

