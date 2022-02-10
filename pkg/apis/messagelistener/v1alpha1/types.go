package v1alpha1

import (
	apis "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const ResourceType string = "messagelisteners"

var kind = apis.NewObjectKind(apis.SchemeGroupVersion, "MessageListener")

var listKind = apis.NewObjectKind(apis.SchemeGroupVersion, "MessageListenerList")

// Represents a message listener. Message triggered jobs will be dispatched based
// on the mssages received from a given queue
type MessageListener struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// message listener specification.
	Spec MessageListenerSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type MessageListenerSpec struct {
	// queue topic
	Topics []string `json:"topics" protobuf:"bytes,1,opt,name=topics"`

	// kafka configuration
	Kafka KafkaListener `json:"kafka" protobuf:"bytes,2,opt,name=kafka"`

	// RabbitMQ configuration
	RabbitMQ RabbitMQListener `json:"rabbitmq" protobuf:"bytes,3,opt,name=rabbitmq"`

	// ActiveMQ configuration
	ActiveMQ ActiveMQListener `json:"activemq" protobuf:"bytes,4,opt,name=activemq"`
}

func (j MessageListener) GetObjectKind() schema.ObjectKind {
	return kind
}
func (j MessageListener) DeepCopyObject() runtime.Object {
	return nil
}

// collection of message triggered jobs
type MessageListenerList struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []MessageListener `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func (j MessageListenerList) GetObjectKind() schema.ObjectKind {
	return listKind
}
func (j MessageListenerList) DeepCopyObject() runtime.Object {
	return nil
}

// Kafka consumer configuration
type KafkaListener struct {
	Config map[string]string `json:"config,omitempty" protobuf:"bytes,1,rep,name=config"`
}

func (l KafkaListener) IsSet() bool {
	return l.Config != nil
}

type RabbitMQListener struct {
	Url string `json:"url" protobuf:"bytes,1,opt,name=url"`
}

func (r RabbitMQListener) IsSet() bool {
	return r.Url != ""
}

type ActiveMQListener struct {
	Network string
	Address string
	Options map[string]string
}

func (a ActiveMQListener) IsSet() bool {
	return a.Network != "" && a.Address != ""
}
