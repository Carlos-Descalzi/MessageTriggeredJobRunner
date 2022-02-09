package v1alpha1

import (
	apis "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis"
	listenertypes "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/messagelistener/v1alpha1"
	mtjobtypes "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/mtjob/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	Scheme             = runtime.NewScheme()
	Codecs             = serializer.NewCodecFactory(Scheme)
	ParameterCodec     = runtime.NewParameterCodec(Scheme)
	localSchemeBuilder = runtime.SchemeBuilder{
		mtjobtypes.AddToScheme,
		listenertypes.AddToScheme,
	}
	AddToScheme = localSchemeBuilder.AddToScheme
)

func init() {
	metav1.AddToGroupVersion(Scheme, apis.SchemeGroupVersion)
	utilruntime.Must(AddToScheme(Scheme))
}
