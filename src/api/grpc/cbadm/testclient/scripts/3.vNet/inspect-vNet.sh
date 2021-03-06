#!/bin/bash

#function create_vNet() {


	TestSetFile=${4:-../testSet.env}
    if [ ! -f "$TestSetFile" ]; then
        echo "$TestSetFile does not exist."
        exit
    fi
	source $TestSetFile
    source ../conf.env
	
	echo "####################################################################"
	echo "## 3. vNet: Status"
	echo "####################################################################"

	CSP=${1}
	REGION=${2:-1}
	POSTFIX=${3:-developer}

	source ../common-functions.sh
	getCloudIndex $CSP

    $CBTUMBLEBUG_ROOT/src/api/grpc/cbadm/cbadm util inspect-mcir --config $CBTUMBLEBUG_ROOT/src/api/grpc/cbadm/grpc_conf.yaml -o json --cc ${CONN_CONFIG[$INDEX,$REGION]} --type vNet | jq ''
    echo ""
#}

#create_vNet
