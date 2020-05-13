#!/bin/bash

source ../conf.env

CSP=${1}
POSTFIX=${2:-developer}
if [ "${CSP}" == "aws" ]; then
	echo "[Test for AWS]"
	INDEX=1
elif [ "${CSP}" == "azure" ]; then
	echo "[Test for Azure]"
	INDEX=2
elif [ "${CSP}" == "gcp" ]; then
	echo "[Test for GCP]"
	INDEX=3
else
	echo "[No acceptable argument was provided (aws, azure, gcp, ..). Default: Test for AWS]"
	CSP="aws"
	INDEX=1
fi


curl -sX GET http://localhost:1024/spider/securitygroup/SG-$CSP-$POSTFIX -H 'Content-Type: application/json' -d '{ "ConnectionName": "'${CONN_CONFIG[INDEX]}'"}' | json_pp