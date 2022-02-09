package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type objectKind struct {
	kind schema.GroupVersionKind
}

func (o *objectKind) SetGroupVersionKind(kind schema.GroupVersionKind) {
	o.kind = kind
}

func (o objectKind) GroupVersionKind() schema.GroupVersionKind {
	return o.kind
}

func NewObjectKind(groupVersion schema.GroupVersion, kind string) schema.ObjectKind {

	return &objectKind{
		kind: schema.GroupVersionKind{
			Group:   groupVersion.Group,
			Version: groupVersion.Version,
			Kind:    kind,
		},
	}
}
