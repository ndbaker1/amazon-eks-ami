package containerd

import (
	_ "embed"
	"os"
	"path"

	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/toml"
)

const (
	containerdConfigFile = "/etc/containerd/config.toml"
	containerdConfigPerm = 0644
)

//go:embed config.toml
var defaultContainerdConfig string

func writeContainerdConfig(c *api.NodeConfig) error {
	var config string
	if c.Spec.Containerd.Config.MergeWithDefaults {
		if c.Spec.Containerd.Config.Inline != "" {
			mergedConfig, err := toml.Merge(defaultContainerdConfig, c.Spec.Containerd.Config.Inline)
			if err != nil {
				return err
			}
			config = *mergedConfig
		} else {
			config = defaultContainerdConfig
		}
	} else {
		config = c.Spec.Containerd.Config.Inline
	}
	err := os.MkdirAll(path.Dir(containerdConfigFile), containerdConfigPerm)
	if err != nil {
		return err
	}
	return os.WriteFile(containerdConfigFile, []byte(config), containerdConfigPerm)
}
