#!/bin/sh
#  this script will update the generated config to have all necessary values set
#  expects first argument to be config file
#  expects second argument to be cluster name

set -x

FILE_PATH=$1
CLUSTER_NAME=$2
SECRETS_ROOT=$3

cluster_name=`echo ${CLUSTER_NAME} | tr -cd '[[:alnum:]]-' | tr '[:upper:]' '[:lower:]'`

#  new style config
sed -i -e "s/- name:[[:space:]]*$/- name: ${cluster_name}/" ${FILE_PATH}


# move regions and AZs to us-east-2. note that this is the CNCT CI region for
# API rate limit purposes.
sed -i -e "s/us-east-1/us-east-2/g" ${FILE_PATH}

# all secrets paths use $HOME, so replace that
sed -i -e "s#\$HOME#${SECRETS_ROOT}#" ${FILE_PATH}