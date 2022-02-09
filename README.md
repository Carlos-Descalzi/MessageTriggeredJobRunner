# MessageTriggeredJobRunner
A kubernetes controller for dispatching jobs from events coming from services like Kafka

For the moment only supports listening Kafka messages.

## WIP.

# How it works
First thing to do is to create a listener like this
```yaml
apiVersion: mtjobrunner.io/v1alpha1
kind: MessageListener
metadata:
  name: test-listener
spec:
  topic: events
  kafka:
    config:
      bootstrap.servers: 10.74.1.187:8080
      group.id: 'g1'
```

Then, create a message triggered job
```yaml
apiVersion: mtjobrunner.io/v1alpha1
kind: MessageTriggeredJob
metadata:
  name: test-job
spec:
  trigger:
    listenerName: test-listener
    topic: events
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: test-job
              image: ubuntu:20.04
              command: ["echo","hello world"]
          restartPolicy: Never
```

Message contents are passed as environment variables to the pod.
- TRIGGERED_JOB_MESSAGE_TOPIC: The topic from where the message came from
- TRIGGERED_JOB_MESSAGE_BODY: The message body encoded in base64
- TRIGGERED_JOB_MESSAGE_*: Any other message property.
The format depends on if it is kafka or rabbitmq
