.PHONY: sign 
sign: ## Generate signature 
	@openssl dgst -sha1 -sign privatekey.pem sign.txt | base64
