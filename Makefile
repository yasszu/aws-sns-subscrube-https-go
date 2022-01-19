.PHONY: sign-notification
sign-notification: ## Generate signature
	@openssl dgst -sha1 -sign testdata/privatekey.pem testdata/sign_notification | base64

.PHONY: sign-subscription-confirmation
sign-subscription-confirmation: ## Generate signature
	@openssl dgst -sha1 -sign testdata/privatekey.pem testdata/sign_subscription_confirmation | base64