#!/bin/bash

. 00-vars.sh

mkdir -p "${EXTENSION_BUILD_DIR}"
mkdir -p "${FUNCTION_BUILD_DIR}"
mkdir -p "${OUTPUT_DIR}"

echo "Building extension"
GOOS=linux go build -o "${EXTENSION_BUILD_DIR}/extensions/${EXTENSION}" .
cd "${EXTENSION_BUILD_DIR}" && zip -r ../../${OUTPUT_DIR}/extensions.zip . -x *.DS_Store *.gitkeep

cd ../../

echo "Building function"
cp "${FUNCTION_DIR}/main.py" "${FUNCTION_BUILD_DIR}"
cd "${FUNCTION_BUILD_DIR}" && zip -j -r ../../${OUTPUT_DIR}/function.zip main.py

