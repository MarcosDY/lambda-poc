#!/bin/bash

# Temporal building directory
BUILD_DIR=.build
# Extension building directory
EXTENSION_BUILD_DIR="${BUILD_DIR}/extension"
# Function building directory
FUNCTION_BUILD_DIR="${BUILD_DIR}/function"
# Builded binaries
OUTPUT_DIR=bin
# Extension name used to create layer
EXTENSION=spire-extension
# Path to function code
FUNCTION_DIR=function
# AWS region where deploy function and extension
AWS_REGION=us-east-2
# Function name
FUNCTION_NAME=svid-client
# Execution role used on function, it must have access to secret
FUNCTION_ROLE=arn:aws:iam::529024819027:role/lambda-role
# Secret name, it must be updated for SPIRE to keep an udpated X509 SVID
SECRET_NAME=svid-dbuser

