.PHONY: sign-notification
sign-notification: ## Generate signature
	@openssl dgst -sha1 -sign testdata/privatekey.pem testdata/sign_notification | base64

.PHONY: sign-confirm-subscription
sign-confirm-subscription: ## Generate signature
	@openssl dgst -sha1 -sign testdata/privatekey.pem testdata/sign_confirm_subscription | base64