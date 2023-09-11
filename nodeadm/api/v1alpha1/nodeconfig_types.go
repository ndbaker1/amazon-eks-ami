package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&NodeConfig{}, &NodeConfigList{})
}

//+kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
//+kubebuilder:storageversion

// NodeConfig is the Schema for the nodeconfigs API
type NodeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NodeConfigSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// NodeConfigList contains a list of NodeConfig
type NodeConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeConfig `json:"items"`
}

type NodeConfigSpec struct {
	AWS        AWSOptions        `json:"aws,omitempty"`
	Cluster    ClusterDetails    `json:"cluster,omitempty"`
	Containerd ContainerdOptions `json:"containerd,omitempty"`
	Kubelet    KubeletOptions    `json:"kubelet,omitempty"`
	Storage    StorageOptions    `json:"storage,omitempty"`
}

type ClusterDetails struct {
	APIServerEndpoint    string   `json:"apiServerEndpoint,omitempty"`
	CertificateAuthority []byte   `json:"certificateAuthority,omitempty"`
	ID                   string   `json:"id,omitempty"`
	Name                 string   `json:"name,omitempty"`
	Region               string   `json:"region,omitempty"`
	DNSAddress           string   `json:"dnsAddress,omitempty"`
	IsOutpost            bool     `json:"isOutpost,omitempty"`
	IPFamily             IPFamily `json:"ipFamily,omitempty"`
	CIDR                 string   `json:"cidr,omitempty"`
}

type IPFamily string

const (
	IPFamilyIPv4 IPFamily = "ipv4"
	IPFamilyIPv6 IPFamily = "ipv6"
)

type DaemonConfigOptions struct {
	Arguments         map[string]string `json:"arguments,omitempty"`
	Source            string            `json:"source,omitempty"`
	Inline            string            `json:"inline,omitempty"`
	MergeWithDefaults bool              `json:"mergeWithDefaults,omitempty"`
}

type ContainerdOptions struct {
	Config DaemonConfigOptions `json:"config,omitempty"`
}

type KubeletOptions struct {
	Config DaemonConfigOptions `json:"config,omitempty"`
}

type StorageOptions struct {
	LocalDisk LocalDiskMode `json:"localDisk,omitempty"`
}

type LocalDiskMode string

const (
	LocalDiskModeMount LocalDiskMode = "mount"
	LocalDiskModeRaid0 LocalDiskMode = "raid0"
)

type GlobalOptions struct {
	PauseContainer ContainerCoordinates `json:"pauseContainer,omitempty"`
	LogLevel       LogLevel             `json:"logLevel,omitempty"`
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type ContainerCoordinates struct {
	Ref string `json:"ref,omitempty"`
}

type AWSOptions struct {
	Retry RetryOptions `json:"retry,omitempty"`
}

type RetryOptions struct {
	MaxAttempts         int    `json:"maxRetries,omitempty"`
	BackoffRate         string `json:"backoffRate,omitempty"`
	InitialDelaySeconds int    `json:"initialDelaySeconds,omitempty"`
}
