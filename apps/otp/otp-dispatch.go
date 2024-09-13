package otp

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thrillee/triq/internals/emails"
	"github.com/thrillee/triq/internals/sms"
)

type dispatchHandler func(*DispatchInfo, string) error

var dispatchFactory = map[OTPDispatchType]dispatchHandler{
	OTP_DISPATCH_SMS: handleSMSDispatch,
}

type dispatchPayload struct {
	dispatchType []OTPDispatchType
	dispatchInfo DispatchInfo
	code         string
	codeType     OTPCodeType
	expiration   string
}

func getDispatcher(dispatchType OTPDispatchType, codeType OTPCodeType) dispatchHandler {
	dispatcher, ok := dispatchFactory[dispatchType]
	if !ok {
		switch {
		case dispatchType == OTP_DISPATCH_EMAIL && codeType == OTP_LONG:
			return handleLongCodeEmailDispatch
		case dispatchType == OTP_DISPATCH_EMAIL && codeType == OTP_SHORT:
			return handleShortCodeEmailDispatch
		}
	}

	return dispatcher
}

func handleDispatch(d *dispatchPayload) error {
	for _, dt := range d.dispatchType {
		dispatch := getDispatcher(dt, d.codeType)
		err := dispatch(&d.dispatchInfo, d.code)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Dispatch-Handler": dt,
				"Dispatch-Info":    d.dispatchInfo,
				"Code":             d.code,
			}).Error(err)
		}
	}
	return nil
}

func handleSMSDispatch(di *DispatchInfo, code string) error {
	message := fmt.Sprintf("Dear customer, Your OTP is %s", code)
	return sms.SendSMS(di.Phone, message)
}

func handleLongCodeEmailDispatch(di *DispatchInfo, code string) error {
	greeting := fmt.Sprintf("Hi %s", di.Fullname)

	verifyLink := os.Getenv("VERIFY_LINK")
	strings.ReplaceAll(verifyLink, "<TOKEN>", code)

	appName := os.Getenv("APP_NAME")

	emailPayload := emails.EmailLinkPayload{
		EmailStandard: emails.EmailStandard{
			PageTitle: "Account Verification",
			Greeting:  greeting,
			APP_NAME:  appName,
		},
		MainAction: "Verify Your Account",
		Message: fmt.Sprintf(`
			We received a request to reset your password for your %s account. 
			If you initiated this request, please click the link below to reset your password:
		`, appName),
		ActionMessage: "Verify Account",
		URL:           verifyLink,
	}

	r := emails.NewEmailRequest([]string{di.Email, "bellotobiloba01@gmail.com"}, fmt.Sprintf("Verify Your %s Account", appName))
	err := r.ParseTemplate(emailPayload, "templates/layout.tmpl", "templates/verify-link.tmpl")
	if err != nil {
		return err
	}

	ok, err := r.SendEmail()
	log.Printf("Email Sent: %v \t Error: %v\n", ok, err)
	return err
}

func handleShortCodeEmailDispatch(di *DispatchInfo, code string) error {
	greeting := fmt.Sprintf("Hi %s", di.Fullname)
	appName := os.Getenv("APP_NAME")

	emailPayload := emails.EmailPayloadOTP{
		EmailStandard: emails.EmailStandard{
			PageTitle: "Account Verification",
			Greeting:  greeting,
			APP_NAME:  appName,
		},
		MainAction: "Verify Your Account",
		EXPIRATION: di.Expiration,
		OTP:        code,
	}

	r := emails.NewEmailRequest([]string{di.Email, "bellotobiloba01@gmail.com"}, fmt.Sprintf("Verify Your %s Account", appName))
	err := r.ParseTemplate(emailPayload, "templates/layout.tmpl", "templates/verify-otp.tmpl")
	if err != nil {
		return err
	}

	ok, err := r.SendEmail()
	log.Printf("Email Sent: %v \t Error: %v\n", ok, err)
	return err
}
