#!/usr/bin/env bash

set -euo pipefail

TEST_CONFIG_PATH=bootstrapv2-test/test-config.yaml
LT_DATA_PATH=bootstrapv2-test/bootstrap-v2-launch-template.json

CLUSTER_NAME=bootstrap-v2-cluster
NODEGROUP_NAME=bootstrap-v2-nodegroup
LT_NAME=bootstrap-v2-launch-template
userdata=$(cat bootstrapv2-test/bootstrap-v2-userdata.txt | base64 -w 0)

PROJECT_DIR=$(pwd)
BOOSTRAP_PROJECT_DIR=$(pwd)/nodeadm

# build bootstrap

cd $BOOSTRAP_PROJECT_DIR
make build
make dist

# build the ami with the new bootstrap
cd $PROJECT_DIR

ami_id=${1:-ami-0df65f87dada26fce}

# launch a nodegroup using this ami and make sure that is properly configured

cat << EOF > $LT_DATA_PATH
{
  "ImageId": "$ami_id",
  "InstanceType": "m5.large",
  "UserData": "$userdata"
}
EOF

aws ec2 delete-launch-template --launch-template-name $LT_NAME || true

LT_ID=$(aws ec2 create-launch-template \
  --launch-template-name $LT_NAME \
  --launch-template-data file://$LT_DATA_PATH \
  --query LaunchTemplate.LaunchTemplateId \
  --output text)

echo "LT_ID=$LT_ID"

cat << EOF > $TEST_CONFIG_PATH
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: $CLUSTER_NAME
  region: us-west-2

managedNodeGroups:
  - name: $NODEGROUP_NAME
    launchTemplate:
      id: $LT_ID
EOF

sed -i -e "s/ami: .*/ami: ${ami_id}/" $TEST_CONFIG_PATH
cat $TEST_CONFIG_PATH
eksctl create nodegroup -f $TEST_CONFIG_PATH 

read -p "pausing.." 

eksctl delete nodegroup -f $TEST_CONFIG_PATH --approve --wait
