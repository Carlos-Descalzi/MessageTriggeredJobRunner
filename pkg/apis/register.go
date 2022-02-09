package v1alpha1

import "k8s.io/apimachinery/pkg/runtime/schema"

const (
	GroupName = "mtjobrunner.io"
)

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}
var InternalSchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "__internal"}
