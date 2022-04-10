package sns

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

var (
	signingCertHostRegexp = regexp.MustCompile(`^sns\.[a-zA-Z0-9\-]{3,}\.amazonaws\.com(\.cn)?$`)
	signingCertURLSchema  = "https"
	signatureVersion      = "1"
	signatureAlgorithm    = x509.SHA1WithRSA
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) ConfirmSubscription(msg SubscriptionConfirmation) (string, error) {
	resp, err := http.Get(msg.SubscribeURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", ErrConfirmSubscription
	}

	return string(body), nil
}

func (c *Client) ValidateCertURL(certURL string) error {
	u, err := url.Parse(certURL)
	if err != nil {
		return ErrInvalidCertURL
	}
	if u.Scheme != signingCertURLSchema {
		return ErrInvalidCertURLSchema
	}
	if !signingCertHostRegexp.MatchString(u.Host) {
		return ErrInvalidCertURLHost
	}
	return nil
}

func (c *Client) CheckSignature(ms MessageSignature) error {
	if ms.SignatureVersion != signatureVersion {
		return ErrInvalidSignatureVersion
	}

	signature, err := base64.StdEncoding.DecodeString(ms.Signature)
	if err != nil {
		return err
	}

	res, err := http.Get(ms.SigningCertURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	p, _ := pem.Decode(body)
	if p == nil {
		return ErrInvalidCertBody
	}

	cert, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		return err
	}

	if err := cert.CheckSignature(signatureAlgorithm, ms.Signed, signature); err != nil {
		return ErrInvalidSignature
	}

	return nil
}
