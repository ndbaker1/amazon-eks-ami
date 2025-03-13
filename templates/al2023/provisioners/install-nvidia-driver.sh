#!/usr/bin/env bash

set -o pipefail
set -o nounset
set -o errexit

if [ "$ENABLE_ACCELERATOR" != "nvidia" ]; then
  exit 0
fi

echo "Installing NVIDIA drivers..."

################################################################################
### Install drivers ############################################################
################################################################################
sudo mv ${WORKING_DIR}/gpu/gpu-ami-util /usr/bin/
sudo mv ${WORKING_DIR}/gpu/kmod-util /usr/bin/

sudo dnf -y install dkms
sudo dnf -y install kernel-modules-extra

function install-nvidia-driver() {
  local runfile_url=$1
  local workdir=NVIDIA-Linux

  sudo mkdir -p ${workdir}
  pushd ${workdir}

  sudo curl -O ${runfile_url}

  # should only be one file for the runfile.
  local package=$(ls)
  sudo sh ${package} --extract-only
  sudo rm -rf ${package}
  # should only be one directory for the extracted contents.
  cd $(ls)

  # we want differet names for the dkms module between proprietary and open
  # builds. this enables us to load them by name at boot time.
  sudo sed -i 's/PACKAGE_NAME="nvidia"/PACKAGE_NAME="nvidia-open"/g' kernel-open/dkms.conf

  # proprietary kmod installation
  sudo ./nvidia-installer --kernel-module-type proprietary --dkms --silent || sudo cat /var/log/nvidia-installer.log
  sudo kmod-util archive nvidia
  sudo kmod-util remove nvidia

  # open kmod installation
  sudo ./nvidia-installer --kernel-module-type open --dkms --silent || sudo cat /var/log/nvidia-installer.log
  sudo kmod-util archive nvidia-open
  sudo kmod-util remove nvidia-open
  
  # uninstall everything before doing a clean install of just the nvidia
  # userspace drivers.
  sudo ./nvidia-installer --uninstall --silent
  sudo ./nvidia-installer --no-kernel-modules --silent

  # assemble the list of supported nvidia devices for the open kernel modules
  echo "# This file was generated from supported-gpus/supported-gpus.json\n$(sed -e 's/^/# /g' supported-gpus/LICENSE)" \
    | sudo tee -a /etc/eks/nvidia-supported-devices.txt

  cat supported-gpus/supported-gpus.json \
    | jq -r '.chips[] | select(.features[] | contains("kernelopen")) | "\(.devid) \(.name)"' \
    | sort -u \
    | sudo tee -a /etc/eks/nvidia-supported-devices.txt

  popd

  sudo rm -rf $workdir
}

function install-fabricmanager() {
  local asset_url=$1

  sudo mkdir -p staging
  sudo curl ${asset_url} | sudo tar -J -xvf - --strip-components 1 -C staging

  pushd staging
  # see: https://github.com/NVIDIA/yum-packaging-fabric-manager/blob/main/fabricmanager.spec#L83-L109
  sudo rsync -alvK bin/ /usr/bin/
  sudo rsync -alvK sbin/ /usr/sbin/
  sudo rsync -alvK include/ /usr/include/
  sudo rsync -alvK lib/ /usr/lib/
  sudo rsync -alvK systemd/ /usr/lib/systemd/system/
  sudo rsync -alvK share/ /usr/share/
  sudo rsync -alvK etc/ /usr/share/nvidia/nvswitch/
  popd

  sudo rm -rf staging
}

function install-imex() {
  local asset_url=$1

  sudo mkdir -p staging
  sudo curl ${asset_url} | sudo tar -J -xvf - --strip-components 1 -C staging

  pushd staging
  sudo rsync -alvK usr/bin/ /usr/bin/
  sudo rsync -alvK lib/ /usr/lib/
  sudo rsync -alvK etc/ /usr/etc/
  popd

  sudo rm -rf staging
}

function install-nvidia-container-toolkit() {
  local asset_url=$1

  sudo mkdir -p staging
  sudo curl -L ${asset_url} | sudo tar -zxvf - --strip-components 4 -C staging

  pushd staging
  sudo dnf install -y ./libnvidia-container1-*
  sudo dnf install -y ./libnvidia-container-tools-*
  sudo dnf install -y ./nvidia-container-toolkit-base-*
  sudo dnf install -y ./nvidia-container-toolkit-*
  popd

  sudo rm -rf staging
}

sudo dnf install -y rsync

install-nvidia-driver             ${NVIDIA_DRIVER_RUNFILE_URL}
install-fabricmanager             ${NVIDIA_FABRICMANAGER_URL}
install-nvidia-container-toolkit  ${NVIDIA_CONTAINER_TOOLKIT_URL}
install-imex                      ${NVIDIA_IMEX_URL}

sudo dnf remove -y rsync

sudo systemctl enable nvidia-fabricmanager
sudo systemctl enable nvidia-persistenced

################################################################################
### Prepare for nvidia init ####################################################
################################################################################
sudo mv ${WORKING_DIR}/gpu/nvidia-kmod-load.sh /etc/eks/
sudo mv ${WORKING_DIR}/gpu/nvidia-kmod-load.service /etc/systemd/system/nvidia-kmod-load.service
sudo mv ${WORKING_DIR}/gpu/set-nvidia-clocks.service /etc/systemd/system/set-nvidia-clocks.service
sudo systemctl daemon-reload
sudo systemctl enable nvidia-kmod-load.service
sudo systemctl enable set-nvidia-clocks.service
