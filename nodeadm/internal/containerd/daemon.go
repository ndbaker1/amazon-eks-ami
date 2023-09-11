package containerd

import (
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/daemon"
)

const containerdDaemonName = "containerd"

var _ daemon.Daemon = &containerd{}

type containerd struct {
	daemonManager daemon.DaemonManager
}

func NewContainerdDaemon(daemonManager daemon.DaemonManager) daemon.Daemon {
	return &containerd{
		daemonManager: daemonManager,
	}
}

func (cd *containerd) Configure(c *api.NodeConfig) error {
	return writeContainerdConfig(c)
}

func (cd *containerd) EnsureRunning() error {
	return cd.daemonManager.StartDaemon(containerdDaemonName)
}

func (cd *containerd) Name() string {
	return containerdDaemonName
}
