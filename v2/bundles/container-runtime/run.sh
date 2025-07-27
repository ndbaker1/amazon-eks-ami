#!/usr/bin/env bash

set -o pipefail
set -o nounset
set -o errexit

###############################################################################
### Containerd setup ##########################################################
###############################################################################

sudo dnf install -y runc-${RUNC_VERSION}
sudo dnf install -y containerd-${CONTAINERD_VERSION}
sudo dnf versionlock containerd-*

sudo systemctl enable ebs-initialize-bin@containerd

###############################################################################
### Nerdctl setup #############################################################
###############################################################################

sudo dnf install -y nerdctl

# TODO: are these necessary? What do they do?
sudo dnf install -y device-mapper-persistent-data lvm2
