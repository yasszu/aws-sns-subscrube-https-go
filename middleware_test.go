package sns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var _ subscriber = (*mockSubscriber)(nil)

type mockSubscriber struct {
	ExpectConfirmSubscription func(msg SubscriptionConfirmation) (string, error)
	ExpectValidateCertURL     func(certURL string) error
	ExpectCheckSignature      func(ms MessageSignature) error
}

func (m *mockSubscriber) ConfirmSubscription(msg SubscriptionConfirmation) (string, error) {
	return m.ExpectConfirmSubscription(msg)
}

func (m *mockSubscriber) ValidateCertURL(certURL string) error {
	return m.ExpectValidateCertURL(certURL)
}

func (m *mockSubscriber) CheckSignature(ms MessageSignature) error {
	return m.ExpectCheckSignature(ms)
}

func TestMiddleware_Subscribe_Notification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		prepare        func() subscriber
		topicARN       string
		messageType    string
		body           map[string]interface{}
		wantStatusCode int
	}{
		{
			name: "it returns ok",
			prepare: func() subscriber {
				c := &mockSubscriber{
					ExpectConfirmSubscription: func(msg SubscriptionConfirmation) (string, error) {
						return "", nil
					},
					ExpectValidateCertURL: func(certURL string) error {
						return nil
					},
					ExpectCheckSignature: func(ms MessageSignature) error {
						return nil
					},
				}
				return c
			},
			topicARN:    "arn:aws:sns:us-west-2:123456789012:MyTopic",
			messageType: "Notification",
			body: map[string]interface{}{
				"Type":             "Notification",
				"MessageId":        "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
				"TopicArn":         "arn:aws:sns:us-west-2:123456789012:MyTopic",
				"Subject":          "My First Message",
				"Message":          "Hello world!",
				"Timestamp":        "2012-05-02T00:54:06.655Z",
				"SignatureVersion": "1",
				"Signature":        "cwMmnINV7NWn5wb4o1faQx9QZBOEpSaJaA86Asdkrpr9C0rdkI/RnyUNl5DrqmueaCiCImuy4Jh0CNeOzqXEdv6WuBjUPbQT/YyAb1h00VVqvjyOvsl2kq+7B3bTfNEahHFZJS2Xh0AtwtWENt159iNnlIRD5NSeVlRyicVv2mgCgK9qxLGGyOFESk43sqUnx5abr0mDR2oFRgbWgwHOly3bQjoaXCfrFYXbmEpz9mMScxoOcRgAUqGVkNLzNBDPU4d9OiBwHxifZBfA6AB3ZxoLm/IZXQJCoK7g44O3NjBCC5nnaMDnHJm1TeSqwVXx8MQQ+8LHhcLbghKkPvo33g==",
				"SigningCertURL":   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
				"UnsubscribeURL":   "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "it returns forbidden when ValidateCertURL failed",
			prepare: func() subscriber {
				c := &mockSubscriber{
					ExpectConfirmSubscription: func(msg SubscriptionConfirmation) (string, error) {
						return "", nil
					},
					ExpectValidateCertURL: func(certURL string) error {
						return ErrInvalidCertURLHost
					},
					ExpectCheckSignature: func(ms MessageSignature) error {
						return nil
					},
				}
				return c
			},
			topicARN:    "arn:aws:sns:us-west-2:123456789012:MyTopic",
			messageType: "Notification",
			body: map[string]interface{}{
				"Type":             "Notification",
				"MessageId":        "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
				"TopicArn":         "arn:aws:sns:us-west-2:123456789012:MyTopic",
				"Subject":          "My First Message",
				"Message":          "Hello world!",
				"Timestamp":        "2012-05-02T00:54:06.655Z",
				"SignatureVersion": "1",
				"Signature":        "cwMmnINV7NWn5wb4o1faQx9QZBOEpSaJaA86Asdkrpr9C0rdkI/RnyUNl5DrqmueaCiCImuy4Jh0CNeOzqXEdv6WuBjUPbQT/YyAb1h00VVqvjyOvsl2kq+7B3bTfNEahHFZJS2Xh0AtwtWENt159iNnlIRD5NSeVlRyicVv2mgCgK9qxLGGyOFESk43sqUnx5abr0mDR2oFRgbWgwHOly3bQjoaXCfrFYXbmEpz9mMScxoOcRgAUqGVkNLzNBDPU4d9OiBwHxifZBfA6AB3ZxoLm/IZXQJCoK7g44O3NjBCC5nnaMDnHJm1TeSqwVXx8MQQ+8LHhcLbghKkPvo33g==",
				"SigningCertURL":   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
				"UnsubscribeURL":   "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96",
			},
			wantStatusCode: http.StatusForbidden,
		},
		{
			name: "it returns forbidden when CheckSignature failed",
			prepare: func() subscriber {
				c := &mockSubscriber{
					ExpectConfirmSubscription: func(msg SubscriptionConfirmation) (string, error) {
						return "", nil
					},
					ExpectValidateCertURL: func(certURL string) error {
						return nil
					},
					ExpectCheckSignature: func(ms MessageSignature) error {
						return ErrInvalidSignature
					},
				}
				return c
			},
			topicARN:    "arn:aws:sns:us-west-2:123456789012:MyTopic",
			messageType: "Notification",
			body: map[string]interface{}{
				"Type":             "Notification",
				"MessageId":        "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
				"TopicArn":         "arn:aws:sns:us-west-2:123456789012:MyTopic",
				"Subject":          "My First Message",
				"Message":          "Hello world!",
				"Timestamp":        "2012-05-02T00:54:06.655Z",
				"SignatureVersion": "1",
				"Signature":        "cwMmnINV7NWn5wb4o1faQx9QZBOEpSaJaA86Asdkrpr9C0rdkI/RnyUNl5DrqmueaCiCImuy4Jh0CNeOzqXEdv6WuBjUPbQT/YyAb1h00VVqvjyOvsl2kq+7B3bTfNEahHFZJS2Xh0AtwtWENt159iNnlIRD5NSeVlRyicVv2mgCgK9qxLGGyOFESk43sqUnx5abr0mDR2oFRgbWgwHOly3bQjoaXCfrFYXbmEpz9mMScxoOcRgAUqGVkNLzNBDPU4d9OiBwHxifZBfA6AB3ZxoLm/IZXQJCoK7g44O3NjBCC5nnaMDnHJm1TeSqwVXx8MQQ+8LHhcLbghKkPvo33g==",
				"SigningCertURL":   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
				"UnsubscribeURL":   "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96",
			},
			wantStatusCode: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "OK")
			}

			b, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("GET", "/", bytes.NewReader(b))
			req.Header.Set(XAmzSnsTopicArn, tt.topicARN)
			req.Header.Set(XAmzSnsMessageType, tt.messageType)

			w := httptest.NewRecorder()
			m := NewMiddleware()
			m.subscriber = tt.prepare()
			h := m.Subscribe(tt.topicARN)(http.HandlerFunc(handler))
			h.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Subscribe() = %v, want %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestMiddleware_Subscribe_SubscriptionConfirmation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		prepare        func() subscriber
		topicARN       string
		messageType    string
		body           map[string]interface{}
		wantStatusCode int
	}{
		{
			name: "it returns ok",
			prepare: func() subscriber {
				c := &mockSubscriber{
					ExpectConfirmSubscription: func(msg SubscriptionConfirmation) (string, error) {
						return "success", nil
					},
					ExpectValidateCertURL: func(certURL string) error {
						return nil
					},
					ExpectCheckSignature: func(ms MessageSignature) error {
						return nil
					},
				}
				return c
			},
			topicARN:    "arn:aws:sns:us-west-2:123456789012:MyTopic",
			messageType: "SubscriptionConfirmation",
			body: map[string]interface{}{
				"Type":             "SubscriptionConfirmation",
				"MessageId":        "165545c9-2a5c-472c-8df2-7ff2be2b3b1b",
				"Token":            "Ethevee8dae4mie3",
				"TopicArn":         "arn:aws:sns:us-west-2:123456789012:MyTopic",
				"Message":          "You have chosen to subscribe to the topic arn:aws:sns:us-west-2:123456789012:MyTopic.\nTo confirm the subscription, visit the SubscribeURL included in this message.",
				"SubscribeURL":     "https://sns.us-west-2.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:us-west-2:123456789012:MyTopic&Token=Ethevee8dae4mie3",
				"Timestamp":        "2012-04-26T20:45:04.751Z",
				"SignatureVersion": "1",
				"Signature":        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				"SigningCertURL":   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "it returns forbidden when ConfirmSubscription failed",
			prepare: func() subscriber {
				c := &mockSubscriber{
					ExpectConfirmSubscription: func(msg SubscriptionConfirmation) (string, error) {
						return "", ErrConfirmSubscription
					},
					ExpectValidateCertURL: func(certURL string) error {
						return nil
					},
					ExpectCheckSignature: func(ms MessageSignature) error {
						return nil
					},
				}
				return c
			},
			topicARN:    "arn:aws:sns:us-west-2:123456789012:MyTopic",
			messageType: "SubscriptionConfirmation",
			body: map[string]interface{}{
				"Type":             "SubscriptionConfirmation",
				"MessageId":        "165545c9-2a5c-472c-8df2-7ff2be2b3b1b",
				"Token":            "Ethevee8dae4mie3",
				"TopicArn":         "arn:aws:sns:us-west-2:123456789012:MyTopic",
				"Message":          "You have chosen to subscribe to the topic arn:aws:sns:us-west-2:123456789012:MyTopic.\nTo confirm the subscription, visit the SubscribeURL included in this message.",
				"SubscribeURL":     "https://sns.us-west-2.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:us-west-2:123456789012:MyTopic&Token=Ethevee8dae4mie3",
				"Timestamp":        "2012-04-26T20:45:04.751Z",
				"SignatureVersion": "1",
				"Signature":        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				"SigningCertURL":   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
			},
			wantStatusCode: http.StatusForbidden,
		},
		{
			name: "it returns forbidden when ValidateCertURL failed",
			prepare: func() subscriber {
				c := &mockSubscriber{
					ExpectConfirmSubscription: func(msg SubscriptionConfirmation) (string, error) {
						return "ok", nil
					},
					ExpectValidateCertURL: func(certURL string) error {
						return ErrInvalidCertURLSchema
					},
					ExpectCheckSignature: func(ms MessageSignature) error {
						return nil
					},
				}
				return c
			},
			topicARN:    "arn:aws:sns:us-west-2:123456789012:MyTopic",
			messageType: "SubscriptionConfirmation",
			body: map[string]interface{}{
				"Type":             "SubscriptionConfirmation",
				"MessageId":        "165545c9-2a5c-472c-8df2-7ff2be2b3b1b",
				"Token":            "Ethevee8dae4mie3",
				"TopicArn":         "arn:aws:sns:us-west-2:123456789012:MyTopic",
				"Message":          "You have chosen to subscribe to the topic arn:aws:sns:us-west-2:123456789012:MyTopic.\nTo confirm the subscription, visit the SubscribeURL included in this message.",
				"SubscribeURL":     "https://sns.us-west-2.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:us-west-2:123456789012:MyTopic&Token=Ethevee8dae4mie3",
				"Timestamp":        "2012-04-26T20:45:04.751Z",
				"SignatureVersion": "1",
				"Signature":        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				"SigningCertURL":   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
			},
			wantStatusCode: http.StatusForbidden,
		},
		{
			name: "it returns forbidden when CheckSignature failed",
			prepare: func() subscriber {
				c := &mockSubscriber{
					ExpectConfirmSubscription: func(msg SubscriptionConfirmation) (string, error) {
						return "ok", nil
					},
					ExpectValidateCertURL: func(certURL string) error {
						return nil
					},
					ExpectCheckSignature: func(ms MessageSignature) error {
						return ErrInvalidSignature
					},
				}
				return c
			},
			topicARN:    "arn:aws:sns:us-west-2:123456789012:MyTopic",
			messageType: "SubscriptionConfirmation",
			body: map[string]interface{}{
				"Type":             "SubscriptionConfirmation",
				"MessageId":        "165545c9-2a5c-472c-8df2-7ff2be2b3b1b",
				"Token":            "Ethevee8dae4mie3",
				"TopicArn":         "arn:aws:sns:us-west-2:123456789012:MyTopic",
				"Message":          "You have chosen to subscribe to the topic arn:aws:sns:us-west-2:123456789012:MyTopic.\nTo confirm the subscription, visit the SubscribeURL included in this message.",
				"SubscribeURL":     "https://sns.us-west-2.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:us-west-2:123456789012:MyTopic&Token=Ethevee8dae4mie3",
				"Timestamp":        "2012-04-26T20:45:04.751Z",
				"SignatureVersion": "1",
				"Signature":        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				"SigningCertURL":   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
			},
			wantStatusCode: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "OK")
			}

			b, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("GET", "/", bytes.NewReader(b))
			req.Header.Set(XAmzSnsTopicArn, tt.topicARN)
			req.Header.Set(XAmzSnsMessageType, tt.messageType)

			w := httptest.NewRecorder()
			m := NewMiddleware()
			m.subscriber = tt.prepare()
			h := m.Subscribe(tt.topicARN)(http.HandlerFunc(handler))
			h.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Subscribe() = %v, want %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}
