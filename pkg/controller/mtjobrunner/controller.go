package mtjobrunner

import (
	"context"
	"fmt"
	"time"

	ifces "github.com/Carlos-Descalzi/mtjobrunner/client/v1alpha1"
	ltypes "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/messagelistener/v1alpha1"
	jtypes "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/mtjob/v1alpha1"
	"go.uber.org/zap"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
)

type inputMessage struct {
	namespace string
	trigger   jtypes.Trigger
	message   *Message
}

type jobReference struct {
	namespace string
	name      string
}

type MessageTriggeredJobController struct {
	kubeClient   *clientset.Clientset
	active       bool
	workerCount  int
	lIfce        ifces.MessageListenerInterface
	jIfce        ifces.MessageTriggeredJobInterface
	ctx          context.Context
	handlers     map[string]*SubscriberHandler
	jobSpecs     map[string]jobReference
	logger       *zap.SugaredLogger
	inputChannel chan inputMessage
	ltnrWatcher  watch.Interface
	jobWatcher   watch.Interface
}

func MessageTriggeredJobControllerNew(
	client *clientset.Clientset,
	lIfce ifces.MessageListenerInterface,
	jIfce ifces.MessageTriggeredJobInterface,
	woerkerCount int,
	logger *zap.SugaredLogger,
) (*MessageTriggeredJobController, error) {

	var controller = MessageTriggeredJobController{
		kubeClient:   client,
		lIfce:        lIfce,
		jIfce:        jIfce,
		handlers:     make(map[string]*SubscriberHandler),
		jobSpecs:     make(map[string]jobReference),
		logger:       logger,
		workerCount:  woerkerCount,
		ctx:          context.Background(),
		inputChannel: make(chan inputMessage),
	}

	return &controller, nil
}

func (c *MessageTriggeredJobController) Close() {
	c.logger.Info("Closing controller")
	c.active = false
	close(c.inputChannel)
	if c.jobWatcher != nil {
		c.jobWatcher.Stop()
	}
	if c.ltnrWatcher != nil {
		c.ltnrWatcher.Stop()
	}
	for key, handler := range c.handlers {
		c.logger.Infof("Stopping handler for %s", key)
		handler.Stop()
	}
}

func (c *MessageTriggeredJobController) Start() {
	c.logger.Info("Starting controller")
	c.active = true
	go c.watchJobs()
	go c.watchListeners()
	for i := 0; i < c.workerCount; i++ {
		c.logger.Infof("Creating worker #%d", i)
		go c.loop()
	}
}

func (c *MessageTriggeredJobController) MessageReceived(namespace string, trigger jtypes.Trigger, message *Message) {
	c.logger.Infof("Received message from ns:%s, listener:%s, topic:%s", namespace, trigger.ListenerName, trigger.Topic)
	c.inputChannel <- inputMessage{namespace: namespace, trigger: trigger, message: message}
}

func (c *MessageTriggeredJobController) watchJobs() {
	c.logger.Info("Starting Job Watcher")

	result, err := c.jIfce.List(c.ctx, metav1.NamespaceAll, metav1.ListOptions{})

	if err != nil {
		c.logger.Error("Error fetching jobs:", err)

	} else {
		for i := 0; i < len(result.Items); i++ {
			c.addJob(&result.Items[i])
		}
	}

	watcher, err := c.jIfce.Watch(c.ctx, metav1.NamespaceAll, metav1.ListOptions{})

	if err != nil {
		c.logger.Error("Unable to watch for message triggered jobs:", err)
	} else {
		c.jobWatcher = watcher
		channel := watcher.ResultChan()

		for c.active {
			evt, ok := <-channel

			if ok {
				switch evt.Type {
				case watch.Added:
					job := evt.Object.(*jtypes.MessageTriggeredJob)
					c.addJob(job)
				case watch.Deleted:
					job := evt.Object.(*jtypes.MessageTriggeredJob)
					c.deleteJob(job)
				case watch.Modified:
					job := evt.Object.(*jtypes.MessageTriggeredJob)
					c.modifyJob(job)
				case watch.Error:
					status := evt.Object.(*metav1.Status)
					c.logger.Error("Error watching for message triggered jobs:", status)
				default:
					c.logger.Infof("Watch event %s: %s", evt.Type, evt.Object)
				}
			} else {
				time.Sleep(time.Duration(1 * time.Second))
			}

		}
	}
}

func (c *MessageTriggeredJobController) addJob(job *jtypes.MessageTriggeredJob) {
	key := c.makeKey(job.Namespace, job.Spec.Trigger)

	if _, ok := c.jobSpecs[key]; !ok {
		c.logger.Infof("New message triggered job %s/%s", job.Namespace, job.Name)
		c.jobSpecs[key] = jobReference{
			namespace: job.Namespace,
			name:      job.Name,
		}
	}
}
func (c *MessageTriggeredJobController) deleteJob(job *jtypes.MessageTriggeredJob) {
	key := c.makeKey(job.Namespace, job.Spec.Trigger)

	delete(c.jobSpecs, key)
}
func (c *MessageTriggeredJobController) modifyJob(job *jtypes.MessageTriggeredJob) {
	c.deleteJob(job)
	c.addJob(job)
}

func (c *MessageTriggeredJobController) watchListeners() {
	c.logger.Info("Starting Listener Watcher")
	result, err := c.lIfce.List(c.ctx, metav1.NamespaceAll, metav1.ListOptions{})

	if err != nil {
		c.logger.Error("Error fetching jobs:", err)
	} else {
		for i := 0; i < len(result.Items); i++ {
			c.addSubscriberHandler(&result.Items[i])
		}
	}

	watcher, err := c.lIfce.Watch(c.ctx, metav1.NamespaceAll, metav1.ListOptions{})

	if err != nil {
		c.logger.Error("Unable to watch for listeners: ", err)
	} else {
		c.ltnrWatcher = watcher
		channel := watcher.ResultChan()
		for c.active {

			evt, ok := <-channel

			if ok {
				switch evt.Type {
				case watch.Added:
					listener := evt.Object.(*ltypes.MessageListener)
					c.addSubscriberHandler(listener)
				case watch.Deleted:
					listener := evt.Object.(*ltypes.MessageListener)
					c.deleteSubscriberHandler(listener)
				case watch.Modified:
					listener := evt.Object.(*ltypes.MessageListener)
					c.modifySubscriberHandler(listener)
				case watch.Error:
					status := evt.Object.(*metav1.Status)
					c.logger.Error("Error watching for listeners:", status)
				}
			} else {
				time.Sleep(time.Duration(1 * time.Second))
			}

		}
	}
}

func (c *MessageTriggeredJobController) addSubscriberHandler(listener *ltypes.MessageListener) {

	key := string(listener.UID)

	if _, ok := c.handlers[key]; !ok {
		c.logger.Infof("New message listener %s/%s", listener.Namespace, listener.Name)
		handler := SubscriberHandlerNew(listener.Namespace, *listener, c.logger)
		handler.AddListener(c)
		c.handlers[string(listener.UID)] = handler
		err := handler.Start()

		if err != nil {
			c.logger.Error(err)
		}
	}

}

func (c *MessageTriggeredJobController) deleteSubscriberHandler(listener *ltypes.MessageListener) {

	key := string(listener.UID)
	handler, ok := c.handlers[key]

	if ok {
		handler.Stop()
		delete(c.handlers, key)
	}
}

func (c *MessageTriggeredJobController) modifySubscriberHandler(listener *ltypes.MessageListener) {
	c.deleteSubscriberHandler(listener)
	c.addSubscriberHandler(listener)
}

func (c MessageTriggeredJobController) makeKey(namespace string, trigger jtypes.Trigger) string {
	return fmt.Sprintf("%s-%s-%s", namespace, trigger.ListenerName, trigger.Topic)
}

func (c *MessageTriggeredJobController) loop() {
	for c.active {
		inputMsg, ok := <-c.inputChannel

		if ok {
			jobRef, ok := c.jobSpecs[c.makeKey(inputMsg.namespace, inputMsg.trigger)]

			if ok {
				jobSpec, err := c.fetchJob(jobRef)
				if err == nil {
					time := metav1.NewTime(time.Now())
					job, err := c.makeJob(jobRef.name, inputMsg.namespace, time, jobSpec.Spec.JobTemplate, inputMsg)

					if err == nil {
						_, err := c.kubeClient.
							BatchV1().
							Jobs(inputMsg.namespace).
							Create(c.ctx, job, metav1.CreateOptions{})

						if err != nil {
							c.logger.Error("Unable to create job", err)
						} else {
							c.logger.Infof("Created MessageTriggeredJob %s/%s", job.Namespace, job.Name)

							jobSpec.Status.LastScheduleTime = &time
							c.updateJob(jobSpec)
						}
					}
				} else {
					c.logger.Errorf("MessageTriggeredJob not found %s/%s: %s", inputMsg.namespace, jobSpec.Name, err)
				}
			} else {
				c.logger.Warnf("Job for listener %s and topic %s not not present",
					inputMsg.trigger.ListenerName, inputMsg.trigger.Topic,
				)
			}

		}

	}
}

func (c *MessageTriggeredJobController) fetchJob(jobRef jobReference) (*jtypes.MessageTriggeredJob, error) {
	return c.jIfce.Get(c.ctx, jobRef.namespace, jobRef.name, metav1.GetOptions{})
}

func (c *MessageTriggeredJobController) updateJob(job *jtypes.MessageTriggeredJob) {
	_, err := c.jIfce.Update(c.ctx, job)

	if err != nil {
		c.logger.Errorf("Error updating mtjob %s/%s", job.Namespace, job.Name)
	}
}

func (c *MessageTriggeredJobController) makeJob(
	namePrefix string,
	namespace string,
	time metav1.Time,
	template batchv1.JobTemplateSpec,
	message inputMessage) (*batchv1.Job, error) {

	job := batchv1.Job{
		Spec: *template.Spec.DeepCopy(),
	}
	job.Name = fmt.Sprintf("%s-%v", namePrefix, time.Unix())
	job.Namespace = namespace
	job.CreationTimestamp = time

	for i := range job.Spec.Template.Spec.Containers {
		job.Spec.Template.Spec.Containers[i].Env = append(
			job.Spec.Template.Spec.Containers[i].Env,
			corev1.EnvVar{Name: "TRIGGERED_JOB_MESSAGE_TOPIC", Value: message.message.Topic},
			corev1.EnvVar{Name: "TRIGGERED_JOB_MESSAGE_BODY", Value: message.message.String()},
		)
		for key, value := range message.message.Properties {
			job.Spec.Template.Spec.Containers[i].Env = append(
				job.Spec.Template.Spec.Containers[i].Env,
				corev1.EnvVar{Name: fmt.Sprintf("TRIGGERED_JOB_MESSAGE_%s", key), Value: value},
			)
		}
	}

	return &job, nil

}
