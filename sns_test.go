package sns

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConfirmSubscription(t *testing.T) {
	t.Parallel()

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
			got, err := ConfirmSubscription(tt.msg)
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

			if got := ValidateCertURL(tt.crtURL); got != tt.want {
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
		"success": {
			sig: MessageSignature{
				Signed: []byte(strings.Join([]string{
					"Message",
					"hello world",
					"MessageId",
					"2e41209f-2772-4a8d-8014-ed1fc296499d",
					"Timestamp",
					"2021-12-17T02:28:11.491Z",
					"TopicArn",
					"arn:aws:sns:ap-northeast-1:000000000000:ls-topic",
					"Type",
					"Notification\n",
				}, "\n")),
				SignatureVersion: "1",
				Signature:        "PkFWsirSdfl64d0m/uEuwacw2z4Otq0w8aqc1mp9sPVO471kZZolMAdtkCQw3BqEvE/ZTdmnBVitsSti5TZTJmpB+Xch4T+mqwY17f3fIzP6kYs8GcGD3GKmgktn2MK/nW19QASQXK2XsWHh3ldz0YKh9gooPwh35Kk1avs7TbGsE1uPEqqNBTAPBi+1yFP6tAmsDL2b/Aji63/vrO7J35YixfaTm1LFT0JoYCTWnlsnRA4KA6sVLVHtRfxGSEXMpQhwjaFP3dmyMR/0LFq+lYWUv/mgX1bSX72CTNyzpSG9MeP62SSGLdTTnWNu+sddbSvxtF7UUEAaeOt/ZodUww==",
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
				"-----END CERTIFICATE-----",
			}, "\n"),
			want: nil,
		},
		"invalid SignatureVersion": {
			sig: MessageSignature{
				Signed: []byte(strings.Join([]string{
					"Message",
					"hello world",
					"MessageId",
					"2e41209f-2772-4a8d-8014-ed1fc296499d",
					"Timestamp",
					"2021-12-17T02:28:11.491Z",
					"TopicArn",
					"arn:aws:sns:ap-northeast-1:000000000000:ls-topic",
					"Type",
					"Notification\n",
				}, "\n")),
				SignatureVersion: "2",
				Signature:        "PkFWsirSdfl64d0m/uEuwacw2z4Otq0w8aqc1mp9sPVO471kZZolMAdtkCQw3BqEvE/ZTdmnBVitsSti5TZTJmpB+Xch4T+mqwY17f3fIzP6kYs8GcGD3GKmgktn2MK/nW19QASQXK2XsWHh3ldz0YKh9gooPwh35Kk1avs7TbGsE1uPEqqNBTAPBi+1yFP6tAmsDL2b/Aji63/vrO7J35YixfaTm1LFT0JoYCTWnlsnRA4KA6sVLVHtRfxGSEXMpQhwjaFP3dmyMR/0LFq+lYWUv/mgX1bSX72CTNyzpSG9MeP62SSGLdTTnWNu+sddbSvxtF7UUEAaeOt/ZodUww==",
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
				"-----END CERTIFICATE-----",
			}, "\n"),
			want: ErrInvalidSignatureVersion,
		},
		"invalid signature": {
			sig: MessageSignature{
				Signed: []byte(strings.Join([]string{
					"Message",
					"hello world",
					"MessageId",
					"2e41209f-2772-4a8d-8014-ed1fc296499d",
					"Timestamp",
					"2021-12-17T02:28:11.491Z",
					"TopicArn",
					"arn:aws:sns:ap-northeast-1:000000000000:xx-topic", // wrong arn
					"Type",
					"Notification\n",
				}, "\n")),
				SignatureVersion: "1",
				Signature:        "PkFWsirSdfl64d0m/uEuwacw2z4Otq0w8aqc1mp9sPVO471kZZolMAdtkCQw3BqEvE/ZTdmnBVitsSti5TZTJmpB+Xch4T+mqwY17f3fIzP6kYs8GcGD3GKmgktn2MK/nW19QASQXK2XsWHh3ldz0YKh9gooPwh35Kk1avs7TbGsE1uPEqqNBTAPBi+1yFP6tAmsDL2b/Aji63/vrO7J35YixfaTm1LFT0JoYCTWnlsnRA4KA6sVLVHtRfxGSEXMpQhwjaFP3dmyMR/0LFq+lYWUv/mgX1bSX72CTNyzpSG9MeP62SSGLdTTnWNu+sddbSvxtF7UUEAaeOt/ZodUww==",
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
				"-----END CERTIFICATE-----",
			}, "\n"),
			want: ErrInvalidSignature,
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
			if got := CheckSignature(tt.sig); got != tt.want {
				t.Errorf("CheckSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
