package sns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func prepare(t *testing.T) {
	signingCertHostRegexp = regexp.MustCompile("")
	signingCertURLSchema = "http"

	t.Cleanup(func() {
		signingCertHostRegexp = regexp.MustCompile(`^sns\.[a-zA-Z0-9\-]{3,}\.amazonaws\.com(\.cn)?$`)
		signingCertURLSchema = "https"
	})
}

func TestMiddleware_Notification(t *testing.T) {
	prepare(t)

	tests := []struct {
		name           string
		topicARN       string
		messageType    string
		body           map[string]interface{}
		certificate    string
		wantStatusCode int
	}{
		{
			name:        "Success",
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
			certificate: strings.Join([]string{
				"-----BEGIN CERTIFICATE-----",
				"MIIDyDCCArACCQDWjKayfhZXGDANBgkqhkiG9w0BAQUFADCBpDELMAkGA1UEBhMC",
				"VVMxEzARBgNVBAgMCldhc2hpbmd0b24xEDAOBgNVBAcMB1NlYXR0bGUxHDAaBgNV",
				"BAoME0V4YW1wbGUgQ29ycG9yYXRpb24xEjAQBgNVBAsMCU1hcmtldGluZzEYMBYG",
				"A1UEAwwPd3d3LmV4YW1wbGUuY29tMSIwIAYJKoZIhvcNAQkBFhNzb21lb25lQGV4",
				"YW1wbGUuY29tMCAXDTIyMDExNTA2MjcxNVoYDzIxMjExMjIyMDYyNzE1WjCBpDEL",
				"MAkGA1UEBhMCVVMxEzARBgNVBAgMCldhc2hpbmd0b24xEDAOBgNVBAcMB1NlYXR0",
				"bGUxHDAaBgNVBAoME0V4YW1wbGUgQ29ycG9yYXRpb24xEjAQBgNVBAsMCU1hcmtl",
				"dGluZzEYMBYGA1UEAwwPd3d3LmV4YW1wbGUuY29tMSIwIAYJKoZIhvcNAQkBFhNz",
				"b21lb25lQGV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC",
				"AQEAuPysDbyweqaP99HQJM1jP3jXrbvndetPXnHxoxsg2vlLsbZ9lcH3KqqEUTd7",
				"8JgulOWF6mtcBpIPEdJtXkw2wAFDz2AokCJ49QaNUEn79p2yrdNzvZNWS+S2X53Q",
				"g8Bjq0amFnqx9x4R2po4NqZcgBu3f1Pc3vQ0z4eKagW7OmGudxatx0A6jXV4U2bF",
				"8zZrwWtYjCkhsy5hNgnxiANR14AxP2N14GlWl1fl3o7EZye2Z8KV7QeuUy4HSnMB",
				"+Nv5lvbYWaUxUSf130Ls/8LIzQWA58WozyTERYGkeG+NWq2vdquDEF6iPBSYTYZi",
				"l8bzq8ovgI5SCCxDSCuvsJvnuwIDAQABMA0GCSqGSIb3DQEBBQUAA4IBAQAof9y/",
				"A2F6qpxVQDJAtAKHRJRXdeZKdhUyAIYMzCVDJJD4vdr8mpg1AnXgUu4ilLJgyJ3e",
				"9ZOpuvfIVZ4R/GzL58Stb+4EiKIoZnFse1zlQRgHj96J9RD8Bov1RwBmNpxZYoVv",
				"o8qjEJfnB9OVfb5ISX/KmArL3Z+uxZ29Iosm04lLVxukeiIccbD6/24d75ptjrSo",
				"253nyYGaLiATF35xTgu9DDHwNwG1vgGxsZ3g0Uio7/34uVUWa9LsZ08Vjtjm0GYr",
				"/pq3fArHBzkGiwy+l7akZ+C4tK68Vyk4Un+uCzG0nVqaODADeKFSC/E7OL3Gee8x",
				"aG+fmXds0GMne+zb",
				"-----END CERTIFICATE-----\n",
			}, "\n"),
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, tt.certificate)
			}))
			defer srv.Close()

			tt.body["SigningCertURL"] = srv.URL

			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "OK")
			}

			b, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("GET", "/", bytes.NewReader(b))
			req.Header.Set(XAmzSnsTopicArn, tt.topicARN)
			req.Header.Set(XAmzSnsMessageType, tt.messageType)
			w := httptest.NewRecorder()

			h := Middleware(tt.topicARN)(http.HandlerFunc(handler))
			h.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatusCode {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
				}
				resp.Body.Close()

				t.Errorf("Middleware() = %v: %s, want %v", resp.StatusCode, string(body), tt.wantStatusCode)
			}
		})
	}
}
