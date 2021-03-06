#!/bin/bash

. 00-vars.sh

aws lambda invoke \
	--function-name "${FUNCTION_NAME}" \
	--region "${AWS_REGION}" \
	--invocation-type "RequestResponse" \
	.lambda-response

cat .lambda-response | jq .svid -r > svid.pem
