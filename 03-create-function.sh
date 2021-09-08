#!/bin/bash

. 00-vars.sh

ARN=$(cat .layer.arn.txt)
echo "creating db-client function"
aws lambda create-function \
	--function-name "${FUNCTION_DB}" \
	--runtime "python3.7" \
	--role "${FUNCTION_ROLE}" \
	--layers ${ARN} \
	--region "${AWS_REGION}" \
	--handler "main.pop_handler" \
	--zip-file "fileb://${OUTPUT_DIR}/function.zip" \
	--environment "Variables={SECRET_NAME=${SECRET_DB}}"
	
echo "creating web-client function"
aws lambda create-function \
	--function-name "${FUNCTION_WEB}" \
	--runtime "python3.7" \
	--role "${FUNCTION_ROLE}" \
	--layers ${ARN} \
	--region "${AWS_REGION}" \
	--handler "main.pop_handler" \
	--zip-file "fileb://${OUTPUT_DIR}/function.zip" \
	--environment "Variables={SECRET_NAME=${SECRET_WEB}}"

