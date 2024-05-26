package mailer

type MailConfigEnv struct {
	// MAIL_HOST is the host of the mail server.
	MAIL_HOST string `json:"mail_host,omitempty"`

	// MAIL_PORT is the port of the mail server.
	MAIL_PORT string `validate:"number" json:"mail_port,omitempty"`

	// MAIL_HOST_USER is the username for the mail server.
	MAIL_HOST_USER string `json:"mail_host_user,omitempty"`

	// MAIL_HOST_PASSWORD is the password for the mail server.
	MAIL_HOST_PASSWORD string `json:"mail_host_password,omitempty"`

	// MAIL_USE_TLS is a boolean that determines whether to use TLS.
	MAIL_USE_TLS bool `validate:"boolean" json:"mail_use_tls,omitempty"`

	// MAIL_USE_SSL is a boolean that determines whether to use SSL.
	MAIL_USE_SS bool `validate:"boolean" json:"mail_use_ssl,omitempty"`

	// MAIL_TIMEOUT is the timeout for the mail server.
	MAIL_TIMEOUT int `validate:"number" json:"mail_timeout,omitempty"`

	// MAIL_API_SERVICE is the service to use for sending emails.
	MAIL_API_SERVICE string `json:"mail_api_service,omitempty"`

	// MAIL_API_KEY is the key to use for sending emails.
	MAIL_API_KEY string `json:"mail_api_key,omitempty"`

	// MAIL_FROM_NAME is the name that will be used as the sender.
	MAIL_FROM_NAME string `json:"mail_from_name,omitempty"`

	// MAIL_REPLY_TO_EMAIL is the name that will be used as the sender.
	MAIL_REPLY_TO_EMAIL string `json:"mail_reply_to_email,omitempty"`
}
