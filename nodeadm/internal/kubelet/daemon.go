package kubelet

import (
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/daemon"
)

const kubeletDaemonName = "kubelet"

var _ daemon.Daemon = &kubelet{}

type kubelet struct {
	daemonManager daemon.DaemonManager
}

func NewKubeletDaemon(daemonManager daemon.DaemonManager) daemon.Daemon {
	return &kubelet{
		daemonManager: daemonManager,
	}
}

func (k *kubelet) Configure(c *api.NodeConfig) error {
	err := writeKubeconfig(c)
	if err != nil {
		return err
	}
	return writeKubeletConfig(c)
}

func (k *kubelet) EnsureRunning() error {
	return k.daemonManager.StartDaemon(kubeletDaemonName)
}

func (k *kubelet) Name() string {
	return kubeletDaemonName
}
