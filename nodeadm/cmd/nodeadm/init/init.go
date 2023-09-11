package init

import (
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/cli"
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/configprovider"
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/containerd"
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/daemon"
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/kubelet"
	"github.com/integrii/flaggy"
	"go.uber.org/zap"
)

func NewInitCommand() cli.Command {
	cmd := flaggy.NewSubcommand("init")
	cmd.Description = "Initialize this instance as a node in an EKS cluster"
	return &initCmd{
		cmd: cmd,
	}
}

type initCmd struct {
	cmd *flaggy.Subcommand
}

func (c *initCmd) Flaggy() *flaggy.Subcommand {
	return c.cmd
}

func (c *initCmd) Run(log *zap.Logger, opts *cli.GlobalOptions) error {
	root, err := cli.IsRunningAsRoot()
	if err != nil {
		return err
	}
	if !root {
		return cli.ErrMustRunAsRoot
	}
	log.Info("Loading configuration", zap.String("configSource", opts.ConfigSource))
	provider, err := configprovider.BuildConfigProvider(opts.ConfigSource)
	if err != nil {
		return err
	}

	config, err := provider.Provide()
	if err != nil {
		return err
	}
	log.Info("Loaded configuration", zap.Reflect("config", config))

	daemonManager, err := daemon.NewDaemonManager()
	if err != nil {
		return err
	}
	defer daemonManager.Close()

	daemons := []daemon.Daemon{
		containerd.NewContainerdDaemon(daemonManager),
		kubelet.NewKubeletDaemon(daemonManager),
	}

	for _, daemon := range daemons {
		nameField := zap.String("name", daemon.Name())
		log.Info("Configuring daemon", nameField)
		if err := daemon.Configure(config); err != nil {
			return err
		}
		log.Info("Configured daemon", nameField)
		log.Info("Ensuring daemon is running", nameField)
		if err := daemon.EnsureRunning(); err != nil {
			return err
		}
		log.Info("Daemon is running", nameField)
	}

	return nil
}
