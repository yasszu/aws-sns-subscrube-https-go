package sns

import (
	"context"
	"errors"
	"net/http"
)

const (
	ContextKeyNotification string = "sns.notification"
)

var (
	ErrNotFoundNotification = errors.New("not found Notification")
)

func SetNotification(r *http.Request, msg Notification) context.Context {
	return context.WithValue(r.Context(), ContextKeyNotification, msg)
}

func GetNotification(r *http.Request) (Notification, error) {
	if msg, ok := r.Context().Value(ContextKeyNotification).(Notification); ok {
		return msg, nil
	}
	return Notification{}, ErrNotFoundNotification
}
