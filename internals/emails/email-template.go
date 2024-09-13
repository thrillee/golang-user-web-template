package emails

type EmailStandard struct {
	AppURL        string
	AppLogo       string
	PageTitle     string
	MainAction    string
	Greeting      string
	APP_NAME      string
	SUPPORT_PHONE string
	SUPPORT_EMAIL string
}

type EmailPayloadOTP struct {
	EmailStandard
	MainAction string
	OTP        string
	EXPIRATION string
}

type EmailLinkPayload struct {
	EmailStandard
	MainAction    string
	Message       string
	ActionMessage string
	URL           string
}
