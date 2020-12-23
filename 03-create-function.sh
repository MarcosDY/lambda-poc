#!/bin/bash

. 00-vars.sh

ARN=$(cat .layer.arn.txt)
echo "creating function"
aws lambda create-function \
	--function-name "${FUNCTION_NAME}" \
	--runtime "python3.7" \
	--role "${FUNCTION_ROLE}" \
	--layers ${ARN} \
	--region "${AWS_REGION}" \
	--handler "main.pop_handler" \
	--zip-file "fileb://${OUTPUT_DIR}/function.zip" \
	--environment "Variables={SECRET_NAME=${SECRET_NAME}}"
	