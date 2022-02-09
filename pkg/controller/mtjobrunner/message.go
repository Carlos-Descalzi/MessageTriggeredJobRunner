package mtjobrunner

import (
	"encoding/base64"

	types "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/mtjob/v1alpha1"
)

type Message struct {
	Topic      string
	Properties map[string]string
	Payload    []byte
}

func (m Message) String() string {
	return base64.StdEncoding.EncodeToString(m.Payload)
}

type Listener interface {
	MessageReceived(namespace string, trigger types.Trigger, message *Message)
}

type Subscriber interface {
	Next(timeout uint32) (*Message, error)
	Stop()
}
