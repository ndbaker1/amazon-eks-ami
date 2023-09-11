module github.com/awslabs/amazon-eks-ami/nodeadm

go 1.21

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.5.0

require (
	github.com/RaveNoX/go-jsonmerge v1.0.0
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000
	github.com/integrii/flaggy v1.5.2
	github.com/pelletier/go-toml/v2 v2.1.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/zap v1.25.0
	golang.org/x/mod v0.12.0
	k8s.io/apimachinery v0.27.2
	sigs.k8s.io/controller-runtime v0.15.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
	k8s.io/utils v0.0.0-20230209194617-a36077c30491 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
