#!/bin/bash

ARN=$(cat .layer.arn.txt)

. 00-vars.sh

echo "Updating function"
aws lambda update-function-code \
	--function-name "${FUNCTION_DB}" \
	--zip-file "fileb://${OUTPUT_DIR}/function.zip" \
	--region ${AWS_REGION}
aws lambda update-function-configuration \
	--function-name "${FUNCTION_DB}" \
	--region "${AWS_REGION}" \
	--environment "Variables={SECRET_NAME=${SECRET_DB}}" \
	--layers "${ARN}"
aws lambda publish-version \
	--function-name "${FUNCTION_DB}" \
	--description "extension test" \
	--region ${AWS_REGION}

aws lambda update-function-code \
	--function-name "${FUNCTION_WEB}" \
	--zip-file "fileb://${OUTPUT_DIR}/function.zip" \
	--region ${AWS_REGION}
aws lambda update-function-configuration \
	--function-name "${FUNCTION_WEB}" \
	--region "${AWS_REGION}" \
	--environment "Variables={SECRET_NAME=${SECRET_WEB}}" \
	--layers "${ARN}"
aws lambda publish-version \
	--function-name "${FUNCTION_WEB}" \
	--description "extension test" \
	--region ${AWS_REGION}


