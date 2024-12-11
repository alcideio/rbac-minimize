module github.com/alcideio/rbac-tool

go 1.22

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/antonmedv/expr v1.15.5
	github.com/emicklei/dot v1.6.2
	github.com/fatih/color v1.16.0
	github.com/fatih/structs v1.1.0
	github.com/google/cel-go v0.20.1
	github.com/kylelemons/godebug v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/cobra v1.8.0
	google.golang.org/protobuf v1.34.0
	k8s.io/api v0.26.15
	k8s.io/apimachinery v0.26.15
	k8s.io/apiserver v0.26.15
	k8s.io/client-go v0.26.15
	k8s.io/component-helpers v0.26.15
	k8s.io/klog v1.0.0
	k8s.io/kubernetes v1.26.15
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/oauth2 v0.19.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/term v0.27.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240429193739-8cf5692501f6 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240429193739-8cf5692501f6 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280 // indirect
	k8s.io/utils v0.0.0-20240423183400-0849a56e8f22 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
)

replace (
	//
	// k8s.io/kubernetes this is evil....but nessecary for audit2rbac
	//
	k8s.io/api => k8s.io/api v0.26.15
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.26.15
	k8s.io/apimachinery => k8s.io/apimachinery v0.26.15
	k8s.io/apiserver => k8s.io/apiserver v0.26.15
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.26.15
	k8s.io/client-go => k8s.io/client-go v0.26.15
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.26.15
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.26.15
	k8s.io/code-generator => k8s.io/code-generator v0.26.15
	k8s.io/component-base => k8s.io/component-base v0.26.15
	k8s.io/component-helpers => k8s.io/component-helpers v0.26.15
	k8s.io/controller-manager => k8s.io/controller-manager v0.26.15
	k8s.io/cri-api => k8s.io/cri-api v0.26.15
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.26.15
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.26.15
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.26.15
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.26.15
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.26.15
	k8s.io/kubectl => k8s.io/kubectl v0.26.15
	k8s.io/kubelet => k8s.io/kubelet v0.26.15
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.26.15
	k8s.io/metrics => k8s.io/metrics v0.26.15
	k8s.io/mount-utils => k8s.io/mount-utils v0.26.15
	k8s.io/node-api => k8s.io/node-api v0.26.15
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.26.15
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.26.15
	k8s.io/sample-controller => k8s.io/sample-controller v0.26.15
)
