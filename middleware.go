package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	XAmzSnsMessageType string = "x-amz-sns-message-type"
	XAmzSnsTopicArn    string = "x-amz-sns-topic-arn"
)

type client interface {
	ConfirmSubscription(msg SubscriptionConfirmation) (string, error)
	ValidateCertURL(certURL string) error
	CheckSignature(ms MessageSignature) error
}

type Middleware struct {
	client client
}

func NewMiddleware() *Middleware {
	return &Middleware{
		client: NewClient(),
	}
}

func (m *Middleware) Subscribe(snsTopicARN string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			topicArn := r.Header.Get(XAmzSnsTopicArn)
			if topicArn != snsTopicARN {
				http.Error(w, "invalid SNS TopicArn", http.StatusForbidden)
				return
			}

			var ctx context.Context
			switch NewMessageType(r.Header.Get(XAmzSnsMessageType)) {
			case MessageTypeSubscriptionConfirmation:
				var msg SubscriptionConfirmation
				if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if err := m.client.ValidateCertURL(msg.SigningCertURL); err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
				if err := m.client.CheckSignature(msg.MessageSignature()); err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
				body, err := m.client.ConfirmSubscription(msg)
				if err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, body)
				return
			case MessageTypeNotification:
				var msg Notification
				if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if err := m.client.ValidateCertURL(msg.SigningCertURL); err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
				if err := m.client.CheckSignature(msg.MessageSignature()); err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
				ctx = SetNotification(r, msg)
			default:
				http.Error(w, "unexpected message type", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
