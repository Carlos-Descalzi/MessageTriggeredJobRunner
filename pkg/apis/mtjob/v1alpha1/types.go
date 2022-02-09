package v1alpha1

import (
	apis "github.com/Carlos-Descalzi/MessageTriggeredJobRunner/pkg/apis"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const ResourceType string = "messagetriggeredjobs"

var kind = apis.NewObjectKind(apis.SchemeGroupVersion, "MessageTriggeredJob")

var listKind = apis.NewObjectKind(apis.SchemeGroupVersion, "MessageTriggeredJobList")

// Represents a message triggered job
type MessageTriggeredJob struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// message triggered job specification.
	Spec MessageTriggeredJobSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// current status of the job
	Status MessageTriggeredJobStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

func (j MessageTriggeredJob) GetObjectKind() schema.ObjectKind {
	return kind
}
func (j MessageTriggeredJob) DeepCopyObject() runtime.Object {
	return nil
}

// collection of message triggered jobs
type MessageTriggeredJobList struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []MessageTriggeredJob `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func (j MessageTriggeredJobList) GetObjectKind() schema.ObjectKind {
	return listKind
}
func (j MessageTriggeredJobList) DeepCopyObject() runtime.Object {
	return nil
}

// job trigger settings.
type Trigger struct {
	// Name of the message listener
	ListenerName string `json:"listenerName" protobuf:"bytes,1,opt,name=listenerName"`
	// Topic to listen
	Topic string `json:"topic" protobuf:"bytes,2,opt,name=topic"`
}

type MessageTriggeredJobSpec struct {
	// Specifies whe source that triggers jobs
	Trigger Trigger `json:"trigger" protobuf:"bytes,1,opt,name=trigger"`

	// Specifies the job that will be created when executing a TZCronJob.
	JobTemplate batchv1.JobTemplateSpec `json:"jobTemplate" protobuf:"bytes,2,opt,name=jobTemplate"`

	// The number of successful finished jobs to retain.
	// This is a pointer to distinguish between explicit zero and not specified.
	// Defaults to 3.
	// +optional
	SuccessfulJobsHistoryLimit *int32 `json:"successfulJobsHistoryLimit,omitempty" protobuf:"varint,3,opt,name=successfulJobsHistoryLimit"`

	// The number of failed finished jobs to retain.
	// This is a pointer to distinguish between explicit zero and not specified.
	// Defaults to 1.
	// +optional
	FailedJobsHistoryLimit *int32 `json:"failedJobsHistoryLimit,omitempty" protobuf:"varint,4,opt,name=failedJobsHistoryLimit"`
}

type MessageTriggeredJobStatus struct {
	// A list of pointers to currently running jobs.
	// +optional
	Active []v1.ObjectReference `json:"active,omitempty" protobuf:"bytes,1,rep,name=active"`

	// Information when was the last time the job was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty" protobuf:"bytes,4,opt,name=lastScheduleTime"`
}
