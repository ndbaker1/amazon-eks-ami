package kubelet

import (
	_ "embed"
	"os"
	"path"

	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/json"
	"golang.org/x/mod/semver"
)

const (
	kubeletConfigRoot = "/etc/kubernetes/kubelet"
	kubeletConfigFile = "config.json"
	kubeletConfigDir  = "config.json.d"
	kubeletConfigPerm = 0644
)

//go:embed config.json
var defaultKubeletConfig string

func writeKubeletConfig(c *api.NodeConfig) error {
	kubeletVersion, err := GetKubeletVersion()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(kubeletConfigRoot, kubeletConfigPerm); err != nil {
		return err
	}
	if semver.Compare(*kubeletVersion, "v1.28.0") < 0 {
		return writeConfig(c)
	} else {
		return writeKubeletConfigToDir(c)
	}
}

// WriteConfig writes the kubelet config to a file.
// This should only be used for kubelet versions < 1.28.
func writeConfig(c *api.NodeConfig) error {
	var config string
	if c.Spec.Kubelet.Config.MergeWithDefaults {
		if c.Spec.Kubelet.Config.Inline != "" {
			mergedConfig, err := json.Merge(defaultKubeletConfig, c.Spec.Kubelet.Config.Inline)
			if err != nil {
				return err
			}
			config = *mergedConfig
		} else {
			config = defaultKubeletConfig
		}
	} else {
		config = c.Spec.Kubelet.Config.Inline
	}
	path := path.Join(kubeletConfigRoot, kubeletConfigFile)
	return os.WriteFile(path, []byte(config), kubeletConfigPerm)
}

// WriteKubeletConfigToDir writes the kubelet config to a directory.
// This is only supported on kubelet versions >= 1.28.
func writeKubeletConfigToDir(cfg *api.NodeConfig) error {
	configs := make(map[string]string)
	if cfg.Spec.Kubelet.Config.MergeWithDefaults {
		configs["10-defaults.json"] = defaultKubeletConfig
	}
	if cfg.Spec.Kubelet.Config.Inline != "" {
		configs["20-inline.json"] = cfg.Spec.Kubelet.Config.Inline
	}
	dirPath := path.Join(kubeletConfigRoot, kubeletConfigDir)
	if err := os.MkdirAll(dirPath, kubeletConfigPerm); err != nil {
		return err
	}
	for fileName, data := range configs {
		filePath := path.Join(dirPath, fileName)
		err := os.WriteFile(filePath, []byte(data), kubeletConfigPerm)
		if err != nil {
			return err
		}
	}
	return nil
}
