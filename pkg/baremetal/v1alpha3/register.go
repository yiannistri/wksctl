package v1alpha3

import (
	"github.com/weaveworks/wksctl/pkg/baremetal"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +k8s:deepcopy-gen=false
type BareMetalProviderSpecCodec struct {
	encoder runtime.Encoder
	decoder runtime.Decoder
}

const GroupName = "cluster.weave.works"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha3"}

var (
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

func init() {
	localSchemeBuilder.Register(addKnownTypes)
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&BareMetalMachine{},
	)
	scheme.AddKnownTypes(SchemeGroupVersion,
		&BareMetalCluster{},
	)
	return nil
}

func NewScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := baremetalproviderspec.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}