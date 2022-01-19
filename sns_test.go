package sns

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConfirmSubscription(t *testing.T) {
	tests := map[string]struct {
		msg     SubscriptionConfirmation
		handler func(w http.ResponseWriter, r *http.Request)
		want    string
		err     error
	}{
		"success": {
			msg: SubscriptionConfirmation{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "ConfirmSubscription")
			},
			want: "ConfirmSubscription",
			err:  nil,
		},
		"Not_Found": {
			msg: SubscriptionConfirmation{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "404 not found", http.StatusNotFound)
			},
			want: "",
			err:  ErrConfirmSubscription,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer srv.Close()

			tt.msg.SubscribeURL = srv.URL
			c := NewClient()
			got, err := c.ConfirmSubscription(tt.msg)
			if err != tt.err {
				t.Errorf("err = %v, want %v", err, tt.err)
			}
			if got != tt.want {
				t.Errorf("NewMessageType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_ValidateCertUrl(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		crtURL string
		want   error
	}{
		"success": {
			crtURL: "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
			want:   nil,
		},
		"invalid scheme": {
			crtURL: "http://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
			want:   ErrInvalidCertURLSchema,
		},
		"invalid host": {
			crtURL: "https://sns.us-west-2.example.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
			want:   ErrInvalidCertURLHost,
		},
		"empty": {
			crtURL: "",
			want:   ErrInvalidCertURLSchema,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c := NewClient()
			if got := c.ValidateCertURL(tt.crtURL); got != tt.want {
				t.Errorf("ValidateCertURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckSignature(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		sig         MessageSignature
		certificate string
		want        error
	}{
		"Notification": {
			sig: MessageSignature{
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
				Signature:        "cwMmnINV7NWn5wb4o1faQx9QZBOEpSaJaA86Asdkrpr9C0rdkI/RnyUNl5DrqmueaCiCImuy4Jh0CNeOzqXEdv6WuBjUPbQT/YyAb1h00VVqvjyOvsl2kq+7B3bTfNEahHFZJS2Xh0AtwtWENt159iNnlIRD5NSeVlRyicVv2mgCgK9qxLGGyOFESk43sqUnx5abr0mDR2oFRgbWgwHOly3bQjoaXCfrFYXbmEpz9mMScxoOcRgAUqGVkNLzNBDPU4d9OiBwHxifZBfA6AB3ZxoLm/IZXQJCoK7g44O3NjBCC5nnaMDnHJm1TeSqwVXx8MQQ+8LHhcLbghKkPvo33g==",
				SigningCertURL:   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
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
			want: nil,
		},
		"SubscriptionConfirmation": {
			sig: MessageSignature{
				Signed: []byte(strings.Join([]string{
					"Message",
					"You have chosen to subscribe to the topic arn:aws:sns:us-west-2:123456789012:MyTopic.\\nTo confirm the subscription, visit the SubscribeURL included in this message.",
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
				Signature:        "qVgHxCgP6yNeqk38SUK6bwpw8qkpDoltPdJD2uyMP8nP9VAtz0+Uw+QnOgQ6phAnV21iIPADUa3kDs+BZk5GJ6V1j0p2M1x+alAvQVWvHtbgvYP+dvU/4BFtyW+DEwgeObn3UsRJaWJE+j8e3ssQQu37+5XvBOPzn8h73Js7DxuU1gKdMBuNQNdoenJU6SgN6yVeyyGkqSGrVWJDR36ViwHHq9Sgy0bqV/axqlT7m/UURb2luRSBbIyeD0p5slOYKLLpdt7wyiWG/SjOvxhxo2IpJNTLDNBAVoOG2dynUnaFs1YMU1zz4BUxdoyx1QhUwiXA8HPL3kVSE7bzKNOFtg==",
				SigningCertURL:   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
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
			want: nil,
		},
		"invalid SignatureVersion": {
			sig: MessageSignature{
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
				SignatureVersion: "2",
				Signature:        "cwMmnINV7NWn5wb4o1faQx9QZBOEpSaJaA86Asdkrpr9C0rdkI/RnyUNl5DrqmueaCiCImuy4Jh0CNeOzqXEdv6WuBjUPbQT/YyAb1h00VVqvjyOvsl2kq+7B3bTfNEahHFZJS2Xh0AtwtWENt159iNnlIRD5NSeVlRyicVv2mgCgK9qxLGGyOFESk43sqUnx5abr0mDR2oFRgbWgwHOly3bQjoaXCfrFYXbmEpz9mMScxoOcRgAUqGVkNLzNBDPU4d9OiBwHxifZBfA6AB3ZxoLm/IZXQJCoK7g44O3NjBCC5nnaMDnHJm1TeSqwVXx8MQQ+8LHhcLbghKkPvo33g==",
				SigningCertURL:   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
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
			want: ErrInvalidSignatureVersion,
		},
		"invalid signature": {
			sig: MessageSignature{
				Signed: []byte(strings.Join([]string{
					"Message",
					"Invalid message", // invalid
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
				Signature:        "cwMmnINV7NWn5wb4o1faQx9QZBOEpSaJaA86Asdkrpr9C0rdkI/RnyUNl5DrqmueaCiCImuy4Jh0CNeOzqXEdv6WuBjUPbQT/YyAb1h00VVqvjyOvsl2kq+7B3bTfNEahHFZJS2Xh0AtwtWENt159iNnlIRD5NSeVlRyicVv2mgCgK9qxLGGyOFESk43sqUnx5abr0mDR2oFRgbWgwHOly3bQjoaXCfrFYXbmEpz9mMScxoOcRgAUqGVkNLzNBDPU4d9OiBwHxifZBfA6AB3ZxoLm/IZXQJCoK7g44O3NjBCC5nnaMDnHJm1TeSqwVXx8MQQ+8LHhcLbghKkPvo33g==",
				SigningCertURL:   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
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
			want: ErrInvalidSignature,
		},
		"invalid certificate body": {
			sig: MessageSignature{
				Signed: []byte(strings.Join([]string{
					"Message",
					"Invalid message", // invalid
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
				Signature:        "cwMmnINV7NWn5wb4o1faQx9QZBOEpSaJaA86Asdkrpr9C0rdkI/RnyUNl5DrqmueaCiCImuy4Jh0CNeOzqXEdv6WuBjUPbQT/YyAb1h00VVqvjyOvsl2kq+7B3bTfNEahHFZJS2Xh0AtwtWENt159iNnlIRD5NSeVlRyicVv2mgCgK9qxLGGyOFESk43sqUnx5abr0mDR2oFRgbWgwHOly3bQjoaXCfrFYXbmEpz9mMScxoOcRgAUqGVkNLzNBDPU4d9OiBwHxifZBfA6AB3ZxoLm/IZXQJCoK7g44O3NjBCC5nnaMDnHJm1TeSqwVXx8MQQ+8LHhcLbghKkPvo33g==",
				SigningCertURL:   "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
			},
			certificate: "",
			want:        ErrInvalidCertBody,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler := func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, tt.certificate)
			}

			srv := httptest.NewServer(http.HandlerFunc(handler))
			defer srv.Close()

			tt.sig.SigningCertURL = srv.URL
			c := NewClient()
			if got := c.CheckSignature(tt.sig); got != tt.want {
				t.Errorf("CheckSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
