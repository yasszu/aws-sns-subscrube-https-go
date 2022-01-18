.PHONY: test-sign
test-sign: ## Generate signature
	@openssl dgst -sha1 -sign testdata/privatekey.pem testdata/sign.txt | base64
