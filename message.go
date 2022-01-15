package sns

import (
	"bytes"
	"reflect"
)

var (
	notificationSignKeys = []string{
		"Message",
		"MessageId",
		"Subject",
		"Timestamp",
		"TopicArn",
		"Type",
	}
	subscriptionConfirmationSignKeys = []string{
		"Message",
		"MessageId",
		"SubscribeURL",
		"Timestamp",
		"Token",
		"TopicArn",
		"Type",
	}
)

type MessageType int

const (
	MessageTypeUnknown MessageType = iota + 1
	MessageTypeNotification
	MessageTypeSubscriptionConfirmation
	MessageTypeUnsubscribeConfirmation
)

var messageTypeStrings = map[MessageType]string{
	MessageTypeUnknown:                  "",
	MessageTypeNotification:             "Notification",
	MessageTypeSubscriptionConfirmation: "SubscriptionConfirmation",
	MessageTypeUnsubscribeConfirmation:  "UnsubscribeConfirmation",
}

func NewMessageType(str string) MessageType {
	for k, v := range messageTypeStrings {
		if v == str {
			return k
		}
	}
	return MessageTypeUnknown
}

func (t MessageType) signKeys() []string {
	switch t {
	case MessageTypeSubscriptionConfirmation, MessageTypeUnsubscribeConfirmation:
		return subscriptionConfirmationSignKeys
	case MessageTypeNotification:
		return notificationSignKeys
	}
	return nil
}

func (t MessageType) sign(m interface{}) []byte {
	buf := &bytes.Buffer{}
	v := reflect.ValueOf(m)
	for _, key := range t.signKeys() {
		field := reflect.Indirect(v).FieldByName(key)
		val := field.String()
		if !field.IsValid() || val == "" {
			continue
		}
		buf.WriteString(key + "\n")
		buf.WriteString(val + "\n")
	}
	return buf.Bytes()
}

type MessageSignature struct {
	Signed           []byte
	SignatureVersion string
	Signature        string
	SigningCertURL   string
}

type SubscriptionConfirmation struct {
	Type             string
	MessageId        string
	Token            string
	TopicArn         string
	Message          string
	SubscribeURL     string
	Timestamp        string
	SignatureVersion string
	Signature        string
	SigningCertURL   string
}

func (m SubscriptionConfirmation) MessageSignature() MessageSignature {
	return MessageSignature{
		Signed:           NewMessageType(m.Type).sign(m),
		SignatureVersion: m.SignatureVersion,
		Signature:        m.Signature,
		SigningCertURL:   m.SigningCertURL,
	}
}

type Notification struct {
	Type              string
	MessageId         string
	TopicArn          string
	Subject           string
	Message           string
	Timestamp         string
	SignatureVersion  string
	Signature         string
	SigningCertURL    string
	UnsubscribeURL    string
	ReceiptHandle     *string
	MessageAttributes map[string]struct {
		Type  string
		Value string
	}
}

func (m Notification) MessageSignature() MessageSignature {
	return MessageSignature{
		Signed:           NewMessageType(m.Type).sign(m),
		SignatureVersion: m.SignatureVersion,
		Signature:        m.Signature,
		SigningCertURL:   m.SigningCertURL,
	}
}
