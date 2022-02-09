package v1alpha1

import (
	"context"
	"net/http"

	"time"

	apis "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis"
	jtypes "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/mtjob/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type MessageTriggeredJobInterface interface {
	Create(ctx context.Context, job *jtypes.MessageTriggeredJob) (*jtypes.MessageTriggeredJob, error)
	Update(ctx context.Context, job *jtypes.MessageTriggeredJob) (*jtypes.MessageTriggeredJob, error)
	Delete(ctx context.Context, namespace string, name string, options *metav1.DeleteOptions) error
	Get(ctx context.Context, namespace string, name string, options metav1.GetOptions) (*jtypes.MessageTriggeredJob, error)
	List(ctx context.Context, namespace string, options metav1.ListOptions) (*jtypes.MessageTriggeredJobList, error)
	Watch(ctx context.Context, namespace string, options metav1.ListOptions) (watch.Interface, error)
}

type messageTriggeredJobImpl struct {
	client rest.Interface
}

func MessageTriggeredJobInterfaceNew(config *rest.Config, httpClient *http.Client) (MessageTriggeredJobInterface, error) {
	cfgCopy := *config
	gv := apis.SchemeGroupVersion
	cfgCopy.GroupVersion = &gv
	cfgCopy.APIPath = "/apis"
	cfgCopy.NegotiatedSerializer = Codecs

	if cfgCopy.UserAgent == "" {
		cfgCopy.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	restClient, err := rest.RESTClientForConfigAndClient(&cfgCopy, httpClient)

	if err != nil {
		return nil, err
	}

	return &messageTriggeredJobImpl{client: restClient}, nil
}

func (c *messageTriggeredJobImpl) Create(ctx context.Context, mtjob *jtypes.MessageTriggeredJob) (*jtypes.MessageTriggeredJob, error) {
	result := &jtypes.MessageTriggeredJob{}
	err := c.client.Post().
		Namespace(mtjob.Namespace).
		Resource(jtypes.ResourceType).
		Body(mtjob).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *messageTriggeredJobImpl) Update(ctx context.Context, job *jtypes.MessageTriggeredJob) (*jtypes.MessageTriggeredJob, error) {
	result := &jtypes.MessageTriggeredJob{}
	err := c.client.Put().
		Namespace(job.Namespace).
		Resource(jtypes.ResourceType).
		Name(job.Name).
		Body(job).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *messageTriggeredJobImpl) Delete(ctx context.Context, namespace string, name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(namespace).
		Resource(jtypes.ResourceType).
		Name(name).
		Body(options).
		Do(ctx).
		Error()
}

func (c *messageTriggeredJobImpl) Get(
	ctx context.Context,
	namespace string,
	name string,
	options metav1.GetOptions,
) (*jtypes.MessageTriggeredJob, error) {

	result := &jtypes.MessageTriggeredJob{}
	err := c.client.Get().
		Resource(jtypes.ResourceType).
		Namespace(namespace).
		Name(name).
		VersionedParams(&options, ParameterCodec).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *messageTriggeredJobImpl) List(ctx context.Context, namespace string, options metav1.ListOptions) (*jtypes.MessageTriggeredJobList, error) {
	var timeout time.Duration
	if options.TimeoutSeconds != nil {
		timeout = time.Duration(*options.TimeoutSeconds) * time.Second
	}
	result := &jtypes.MessageTriggeredJobList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource(jtypes.ResourceType).
		VersionedParams(&options, ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return result, err

}

func (c *messageTriggeredJobImpl) Watch(ctx context.Context, namespace string, options metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if options.TimeoutSeconds != nil {
		timeout = time.Duration(*options.TimeoutSeconds) * time.Second
	}
	options.Watch = true

	return c.client.Get().
		Namespace(namespace).
		Resource(jtypes.ResourceType).
		VersionedParams(&options, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}
