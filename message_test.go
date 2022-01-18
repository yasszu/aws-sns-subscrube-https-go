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
				MessageId:        "165545c9-2a5c-472c-8df2-7ff2be2b3b1b",
				Token:            "Ethevee8dae4mie3",
				TopicArn:         "arn:aws:sns:us-west-2:123456789012:MyTopic",
				Message:          "You have chosen to subscribe to the topic arn:aws:sns:us-west-2:123456789012:MyTopic.\nTo confirm the subscription, visit the SubscribeURL included in this message.",
				SubscribeURL:     "https://sns.us-west-2.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:us-west-2:123456789012:MyTopic&Token=Ethevee8dae4mie3",
				Timestamp:        "2012-04-26T20:45:04.751Z",
				SignatureVersion: "1",
				Signature:        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				SigningCertURL:   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
			},
			want: MessageSignature{
				Signed: []byte(strings.Join([]string{
					"Message",
					"You have chosen to subscribe to the topic arn:aws:sns:us-west-2:123456789012:MyTopic.\nTo confirm the subscription, visit the SubscribeURL included in this message.",
					"MessageId",
					"165545c9-2a5c-472c-8df2-7ff2be2b3b1b",
					"SubscribeURL",
					"https://sns.us-west-2.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:us-west-2:123456789012:MyTopic&Token=Ethevee8dae4mie3",
					"Timestamp",
					"2012-04-26T20:45:04.751Z",
					"Token",
					"Ethevee8dae4mie3",
					"TopicArn",
					"arn:aws:sns:us-west-2:123456789012:MyTopic",
					"Type",
					"SubscriptionConfirmation\n",
				}, "\n")),
				SignatureVersion: "1",
				Signature:        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				SigningCertURL:   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
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
				MessageId:        "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
				TopicArn:         "arn:aws:sns:us-west-2:123456789012:MyTopic",
				Subject:          "My First Message",
				Message:          "Hello world!",
				Timestamp:        "2012-05-02T00:54:06.655Z",
				SignatureVersion: "1",
				Signature:        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				SigningCertURL:   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
				UnsubscribeURL:   "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96",
			},
			want: MessageSignature{
				Signed: []byte(strings.Join([]string{
					"Message",
					"Hello world!",
					"MessageId",
					"22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
					"Subject",
					"My First Message",
					"Timestamp",
					"2012-05-02T00:54:06.655Z",
					"TopicArn",
					"arn:aws:sns:us-west-2:123456789012:MyTopic",
					"Type",
					"Notification\n",
				}, "\n")),
				SignatureVersion: "1",
				Signature:        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				SigningCertURL:   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
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
