#  this script fetches all credentials to support the building of a Kraken cluster
#  for now this includes:
#   - an ssh key pair
#   - aws credentials file
#
#  Needed:
#   - gcloud service account file

#  we will use the IAM role of a kubelet to fetch this information from s3
set -x
SECRETS_ROOT=$1

#  ssh keys
mkdir ${SECRETS_ROOT}/.ssh/
aws s3 cp --recursive s3://sundry-automata/keys/common-tools-jenkins/ ${SECRETS_ROOT}/.ssh/
chmod 600 ${SECRETS_ROOT}/.ssh/*

#  aws configs
mkdir ${SECRETS_ROOT}/.aws/
aws s3 cp --recursive s3://sundry-automata/credentials/common-tools-jenkins/aws/ ${SECRETS_ROOT}/.aws/

#  gcloud configs
mkdir -p ${SECRETS_ROOT}/.config/gcloud/
aws s3 cp s3://sundry-automata/credentials/common-tools-jenkins/gke/patrickRobot.json ${SECRETS_ROOT}/.config/gcloud/
