package sns

import (
	"net/http"
	"reflect"
	"testing"
)

func TestGetNotification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want Notification
	}{
		{
			name: "success",
			want: Notification{
				Type:             "Notification",
				MessageId:        "2e41209f-2772-4a8d-8014-ed1fc296499d",
				TopicArn:         "arn:aws:sns:ap-northeast-1:000000000000:en-topic",
				Message:          "test",
				Timestamp:        "2021-12-17T02:28:11.491Z",
				SignatureVersion: "1",
				Signature:        "EXAMPLEpH+DcEwjAPg8O9mY8dReBSwksfg2S7WKQcikcNKWLQjwu6A4VbeS0QHVCkhRS7fUQvi2egU3N858fiTDN6bkkOxYDVrY0Ad8L10Hs3zH81mtnPk5uvvolIC1CXGu43obcgFxeL3khZl8IKvO61GWB6jI9b5+gLPoBc1Q=",
				SigningCertURL:   "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-0000000000000000000000.pem",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &http.Request{}
			ctx := SetNotification(r, tt.want)
			r = r.WithContext(ctx)

			got, err := GetNotification(r)
			if err != nil {
				t.Errorf("err should be nil, but got %q", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNotification() got = %v, want %v", got, tt.want)
			}
		})
	}
}
