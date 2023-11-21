package kubelet

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"io"
	"os"
	"path"
	"text/template"

	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"
)

const (
	caCertificatePath = "/etc/kubernetes/pki"
	kubeconfigDir     = "/var/lib/kubelet"
	kubeconfigPerm    = 0644
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

func writeClusterCaCert(c *api.NodeConfig) error {
	caDecoded := base64.NewDecoder(&base64.Encoding{}, bytes.NewReader(c.Spec.Cluster.CertificateAuthority))
	caDecodedStr, err := io.ReadAll(caDecoded)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(path.Dir(caCertificatePath), kubeletConfigPerm); err != nil {
		return err
	}
	return os.WriteFile(caCertificatePath, caDecodedStr, kubeletConfigPerm)
}
