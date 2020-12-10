#!/bin/bash

ARN=$(cat .layer.arn.txt)

. 00-vars.sh

echo "Updating function"
aws lambda update-function-code \
	--function-name "${FUNCTION_NAME}" \
	--zip-file "fileb://${OUTPUT_DIR}/function.zip" \
	--region ${AWS_REGION}
aws lambda update-function-configuration \
	--function-name "${FUNCTION_NAME}" \
	--region "${AWS_REGION}" \
	--environment "Variables={SECRET_NAME=${SECRET_NAME}}" \
	--layers "${ARN}"
aws lambda publish-version \
	--function-name "${FUNCTION_NAME}" \
	--description "extension test" \
	--region ${AWS_REGION}

