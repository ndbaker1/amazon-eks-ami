---
apiVersion: v1
kind: Config
clusters:
  - name: kubernetes
    cluster:
      certificate-authority: /etc/kubernetes/pki/ca.crt
      server: {{.Spec.Cluster.APIServerEndpoint}}
current-context: kubelet
contexts:
  - name: kubelet
    context:
      cluster: kubernetes
      user: kubelet
users:
  - name: kubelet
    user:
      exec:
        apiVersion: client.authentication.k8s.io/v1beta1
        command: aws
        args:
          - "eks"
          - "get-token"
          - "--cluster-name"
          - "{{.Spec.Cluster.Name}}"
          - "--region"
          - "{{.Spec.Cluster.Region}}"
