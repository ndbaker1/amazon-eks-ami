#!/usr/bin/env bash

set -o pipefail
set -o nounset
set -o errexit

if [ "$ENABLE_ACCELERATOR" != "rocm" ]; then
  exit 0
fi


echo "Installing AMD ROCm drivers..."

################################################################################
### Install drivers ############################################################
################################################################################

sudo dnf -y install kernel-modules-extra-common-$(uname -r)

if [[ "$(uname -r)" == 6.12.* ]]; then
  sudo dnf -y install kernel6.12-modules-extra-$(uname -r)
else
  sudo dnf -y install kernel-modules-extra-$(uname -r)
fi

sudo dnf -y install kernel-devel-$(uname -r) kernel-headers-$(uname -r)

version=6.4.1

# add the repo gpg key
sudo rpm --import https://repo.radeon.com/rocm/rocm.gpg.key
# add the ROCm and AMDGPU repositories
dnf config-manager --add-repo https://repo.radeon.com/rocm/el9/$version/main
dnf config-manager --add-repo https://repo.radeon.com/amdgpu/$version/el/9.6/main/$(uname -m)/
# update all current .repo sources to enable gpgcheck
sudo dnf config-manager --save --setopt=*.gpgcheck=1

sudo dnf install -y amdgpu rocm

########################################################################
### Post-install #######################################################
########################################################################

# inform the linker for ROCm applications
cat << EOF | sudo tee -a /etc/ld.so.conf.d/rocm.conf
/opt/rocm/lib
/opt/rocm/lib64
EOF
