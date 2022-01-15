package sns

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewMessageType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
		want MessageType
	}{
		{
			name: "Notification",
			str:  "Notification",
			want: MessageTypeNotification,
		}, {
			name: "SubscriptionConfirmation",
			str:  "SubscriptionConfirmation",
			want: MessageTypeSubscriptionConfirmation,
		},
		{
			name: "UnsubscribeConfirmation",
			str:  "UnsubscribeConfirmation",
			want: MessageTypeUnsubscribeConfirmation,
		},
		{
			name: "Empty",
			str:  "",
			want: MessageTypeUnknown,
		},
		{
			name: "Unknown",
			str:  "unknown-key",
			want: MessageTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewMessageType(tt.str); got != tt.want {
				t.Errorf("NewMessageType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscriptionConfirmation_MessageSignature(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		message SubscriptionConfirmation
		want    MessageSignature
	}{
		"success": {
			message: SubscriptionConfirmation{
				Type:             "SubscriptionConfirmation",
				MessageId:        "5d81b8c4-f374-4d13-9c8b-0455470c5ae2",
				Token:            "2ba91ca0",
				TopicArn:         "arn:aws:sns:ap-northeast-1:000000000000:ls-topic",
				Message:          "You have chosen to subscribe to the topic arn:aws:sns:ap-northeast-1:000000000000:ls-topic.\nTo confirm the subscription, visit the SubscribeURL included in this message.",
				SubscribeURL:     "http://localstack:4566/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:ap-northeast-1:000000000000:ls-topic&Token=2ba91ca0",
				Timestamp:        "2021-12-17T02:14:06.999Z",
				SignatureVersion: "1",
				Signature:        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				SigningCertURL:   "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-0000000000000000000000.pem",
			},
			want: MessageSignature{
				Signed: []byte(strings.Join([]string{
					"Message",
					"You have chosen to subscribe to the topic arn:aws:sns:ap-northeast-1:000000000000:ls-topic.\nTo confirm the subscription, visit the SubscribeURL included in this message.",
					"MessageId",
					"5d81b8c4-f374-4d13-9c8b-0455470c5ae2",
					"SubscribeURL",
					"http://localstack:4566/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:ap-northeast-1:000000000000:ls-topic&Token=2ba91ca0",
					"Timestamp",
					"2021-12-17T02:14:06.999Z",
					"Token",
					"2ba91ca0",
					"TopicArn",
					"arn:aws:sns:ap-northeast-1:000000000000:ls-topic",
					"Type",
					"SubscriptionConfirmation\n",
				}, "\n")),
				SignatureVersion: "1",
				Signature:        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				SigningCertURL:   "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-0000000000000000000000.pem",
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := tt.message.MessageSignature(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MessageSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotification_MessageSignature(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		message Notification
		want    MessageSignature
	}{
		"success": {
			message: Notification{
				Type:             "Notification",
				MessageId:        "2e41209f-2772-4a8d-8014-ed1fc296499d",
				TopicArn:         "arn:aws:sns:ap-northeast-1:000000000000:en-topic",
				Message:          "test",
				Timestamp:        "2021-12-17T02:28:11.491Z",
				SignatureVersion: "1",
				Signature:        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				SigningCertURL:   "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-0000000000000000000000.pem",
			},
			want: MessageSignature{
				Signed: []byte(strings.Join([]string{
					"Message",
					"test",
					"MessageId",
					"2e41209f-2772-4a8d-8014-ed1fc296499d",
					"Timestamp",
					"2021-12-17T02:28:11.491Z",
					"TopicArn",
					"arn:aws:sns:ap-northeast-1:000000000000:en-topic",
					"Type",
					"Notification\n",
				}, "\n")),
				SignatureVersion: "1",
				Signature:        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				SigningCertURL:   "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-0000000000000000000000.pem",
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := tt.message.MessageSignature(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MessageSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
