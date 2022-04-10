package sns

import "errors"

var (
	ErrConfirmSubscription     = errors.New("error confirm subscription")
	ErrInvalidCertURL          = errors.New("error invalid cert url")
	ErrInvalidCertURLSchema    = errors.New("error invalid cert url scheme")
	ErrInvalidCertURLHost      = errors.New("error invalid cert url host")
	ErrInvalidCertBody         = errors.New("error invalid cert body")
	ErrInvalidSignatureVersion = errors.New("error invalid signature version")
	ErrInvalidSignature        = errors.New("error invalid signature")
)
