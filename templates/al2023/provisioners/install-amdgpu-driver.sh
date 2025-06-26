#!/usr/bin/env bash

set -o pipefail
set -o nounset
set -o errexit

if [ "$ENABLE_ACCELERATOR" != "amdgpu" ]; then
  exit 0
fi

echo "Installing AMD GPU drivers..."

################################################################################
### Install Kernel Deps ########################################################
################################################################################

sudo dnf -y install kernel-modules-extra-common-$(uname -r)

if [[ "$(uname -r)" == 6.12.* ]]; then
  sudo dnf -y install kernel6.12-modules-extra-$(uname -r)
else
  sudo dnf -y install kernel-modules-extra-$(uname -r)
fi

sudo dnf -y install kernel-devel-$(uname -r) kernel-headers-$(uname -r)

################################################################################
### Install Drivers ############################################################
################################################################################

sudo rpm --import https://repo.radeon.com/rocm/rocm.gpg.key
sudo dnf config-manager --add-repo https://repo.radeon.com/amdgpu/$AMD_GPU_DRIVER_VERSION/el/9.6/main/$(uname -m)/
sudo dnf config-manager --save --setopt=*.gpgcheck=1
sudo dnf install -y amdgpu-dkms
