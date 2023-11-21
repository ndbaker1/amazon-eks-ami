#!/usr/bin/env bash

set -euo pipefail

TEST_DIR=bootstrapv2-test

# Dependency
userdata=$(cat $TEST_DIR/bootstrap-v2-userdata.txt | base64 -w 0)
ami_id=${1}

# Generated
TEST_CONFIG_PATH=$TEST_DIR/bootstrap-v2-config.yaml
LT_DATA_PATH=$TEST_DIR/bootstrap-v2-launch-template.json

CLUSTER_NAME=bootstrap-v2-cluster
NODEGROUP_NAME=bootstrap-v2-nodegroup
LT_NAME=bootstrap-v2-launch-template

# launch a nodegroup using this ami and make sure that is properly configured

cat << EOF > $LT_DATA_PATH
{
  "ImageId": "$ami_id",
  "InstanceType": "m5.large",
  "UserData": "$userdata"
}
EOF

aws ec2 delete-launch-template --launch-template-name $LT_NAME 2>&1>/dev/null || true

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

eksctl create nodegroup -f $TEST_CONFIG_PATH || true
read -p "pausing.."
eksctl delete nodegroup -f $TEST_CONFIG_PATH --approve --wait
