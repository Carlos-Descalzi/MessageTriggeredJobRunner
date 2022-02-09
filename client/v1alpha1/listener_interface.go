package v1alpha1

import (
	"context"
	"net/http"
	"time"

	apis "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis"
	ltypes "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/messagelistener/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type MessageListenerInterface interface {
	Create(ctx context.Context, messagelistener *ltypes.MessageListener) (*ltypes.MessageListener, error)
	Update(ctx context.Context, messagelistener *ltypes.MessageListener) (*ltypes.MessageListener, error)
	Delete(ctx context.Context, namespace string, name string, options *metav1.DeleteOptions) error
	Get(ctx context.Context, namespace string, name string, options metav1.GetOptions) (*ltypes.MessageListener, error)
	List(ctx context.Context, namespace string, options metav1.ListOptions) (*ltypes.MessageListenerList, error)
	Watch(ctx context.Context, namespace string, options metav1.ListOptions) (watch.Interface, error)
}

type messageListenerImpl struct {
	client rest.Interface
}

func MessageListenerInterfaceNew(config *rest.Config, httpClient *http.Client) (MessageListenerInterface, error) {
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

	return &messageListenerImpl{client: restClient}, nil
}

func (c *messageListenerImpl) Create(ctx context.Context, listener *ltypes.MessageListener) (*ltypes.MessageListener, error) {
	result := &ltypes.MessageListener{}
	err := c.client.Post().
		Namespace(listener.Namespace).
		Resource(ltypes.ResourceType).
		Body(listener).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *messageListenerImpl) Update(ctx context.Context, listener *ltypes.MessageListener) (*ltypes.MessageListener, error) {
	result := &ltypes.MessageListener{}
	err := c.client.Put().
		Namespace(listener.Namespace).
		Resource(ltypes.ResourceType).
		Name(listener.Name).
		Body(listener).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *messageListenerImpl) Delete(ctx context.Context, namespace string, name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(namespace).
		Resource(ltypes.ResourceType).
		Name(name).
		Body(options).
		Do(ctx).
		Error()
}

func (c *messageListenerImpl) Get(ctx context.Context, namespace string, name string, options metav1.GetOptions) (*ltypes.MessageListener, error) {
	result := &ltypes.MessageListener{}
	err := c.client.Get().
		Namespace(namespace).
		Resource(ltypes.ResourceType).
		Name(name).
		VersionedParams(&options, ParameterCodec).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *messageListenerImpl) List(ctx context.Context, namespace string, options metav1.ListOptions) (*ltypes.MessageListenerList, error) {
	var timeout time.Duration
	if options.TimeoutSeconds != nil {
		timeout = time.Duration(*options.TimeoutSeconds) * time.Second
	}
	result := &ltypes.MessageListenerList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource(ltypes.ResourceType).
		VersionedParams(&options, ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return result, err

}

func (c *messageListenerImpl) Watch(ctx context.Context, namespace string, options metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if options.TimeoutSeconds != nil {
		timeout = time.Duration(*options.TimeoutSeconds) * time.Second
	}
	options.Watch = true
	return c.client.Get().
		Namespace(namespace).
		Resource(ltypes.ResourceType).
		VersionedParams(&options, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}
