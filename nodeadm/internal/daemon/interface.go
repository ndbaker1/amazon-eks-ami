package daemon

import "github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"

type Daemon interface {
	// Configure configures the daemon.
	Configure(*api.NodeConfig) error

	// EnsureRunning ensures that the daemon is running.
	// If the daemon is not running, it will be started.
	// If the daemon is already running, and has been re-configured, it will be restarted.
	EnsureRunning() error

	// Name returns the name of the daemon.
	Name() string
}
