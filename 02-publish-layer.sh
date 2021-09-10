#!/bin/bash

. 00-vars.sh

echo "publishing layer"
aws lambda publish-layer-version \
	--layer-name "${EXTENSION}" \
	--compatible-runtimes go1.x ruby2.7 python3.7 \
	--region "${AWS_REGION}" \
	--zip-file  "fileb://bin/extensions.zip" > .response.json
cat .response.json
read -n 1 -r -s -p $'Press enter to continue...\n'

cat .response.json | jq  -r '.LayerVersionArn' > .layer.arn.txt

