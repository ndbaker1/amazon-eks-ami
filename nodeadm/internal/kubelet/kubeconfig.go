package kubelet

import (
	"bytes"
	_ "embed"
	"os"
	"path"
	"text/template"

	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"
)

const (
	kubeconfigDir  = "/var/lib/kubelet"
	kubeconfigPerm = 0644
)

var (
	//go:embed kubeconfig.tpl
	kubeconfigTemplateData string
	kubeconfigTemplate     = template.Must(template.New("kubeconfig").Parse(kubeconfigTemplateData))
)

func writeKubeconfig(c *api.NodeConfig) error {
	var buf bytes.Buffer
	if err := kubeconfigTemplate.Execute(&buf, c); err != nil {
		return err
	}
	if err := os.MkdirAll(kubeconfigDir, kubeconfigPerm); err != nil {
		return err
	}
	path := path.Join(kubeconfigDir, "kubeconfig")
	return os.WriteFile(path, buf.Bytes(), kubeconfigPerm)
}
