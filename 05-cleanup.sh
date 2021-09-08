#!/bin/bash

. 00-vars.sh
ARN=$(cat .layer.arn.txt)

get_versions () {
  echo $(aws lambda list-layer-versions --layer-name "$EXTENSION" --region "$AWS_REGION" --output text --query LayerVersions[].Version | tr '[:blank:]' '\n')
}

echo "Deleting db-client Function"
aws lambda delete-function \
	--function-name "${FUNCTION_DB}" \
	--region "${AWS_REGION}" > /dev/null

echo "Deleting web-client Function"
aws lambda delete-function \
	--function-name "${FUNCTION_WEB}" \
	--region "${AWS_REGION}" > /dev/null

echo "Deleting Layers"
versions=$(get_versions)
for version in $versions;
do
    echo "deleting arn:aws:lambda:$AWS_REGION:*:layer:$EXTENSION:$version"
    aws lambda delete-layer-version --region "$AWS_REGION" --layer-name "$EXTENSION" --version-number "$version" > /dev/null
done

echo "Deleting logs"
while true; do
    read -p "Delete function log group (/aws/lambda/*? (y/n)" response
    case $response in
	[Yy]* ) aws logs delete-log-group --log-group-name /aws/lambda/$FUNCTION_DB > /dev/null; aws logs delete-log-group --log-group-name /aws/lambda/$FUNCTION_WEB > /dev/null; break;;
	[Nn]* ) break;;
	* ) echo "Response must start with y or n.";;
    esac
done

