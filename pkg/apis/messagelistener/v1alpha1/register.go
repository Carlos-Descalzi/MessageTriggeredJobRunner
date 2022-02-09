package v1alpha1

import (
	apis "github.com/Carlos-Descalzi/MessageTriggeredJobRunner/pkg/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Kind(kind string) schema.GroupKind {
	return apis.SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return apis.SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(apis.SchemeGroupVersion,
		&MessageListener{},
		&MessageListenerList{},
	)
	scheme.AddKnownTypes(apis.InternalSchemeGroupVersion,
		&MessageListener{},
		&MessageListenerList{},
	)
	metav1.AddToGroupVersion(scheme, apis.SchemeGroupVersion)
	return nil
}
