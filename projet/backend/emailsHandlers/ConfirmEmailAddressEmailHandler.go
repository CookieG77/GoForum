package emailsHandlers

import (
	f "GoForum/functions"
	"fmt"
	"html/template"
)

func SendConfirmEmail(email string) {
	tmpl, err := template.ParseFiles("templates/emails/confirmEmailAddressEmail.html")
	if err != nil {
		f.ErrorPrintf("An error occurred while trying to parse the template -> %v\n", err)
		return
	}
	// Before we create the email identification link, we need to remove the previous ones if they exist
	err = f.RemoveEmailIdentificationForUser(email, f.VerifyEmailEmail)
	if err != nil {
		f.ErrorPrintf("Error while removing the previous email identification links: %s\n", err)
		return
	}
	// Create the email identification link
	emailLinkID, err := f.CreateEmailIdentificationLink(email, f.VerifyEmailEmail)
	if err != nil {
		f.ErrorPrintf("Error while creating the email identification link: %s\n", err)
		return
	}
	interfaceContent := make(map[string]interface{})
	if f.IsCertified() {
		interfaceContent["Url"] = fmt.Sprintf("https://localhost/confirm-email-address?token=%s", emailLinkID)
	} else {
		interfaceContent["Url"] = fmt.Sprintf("http://localhost/confirm-email-address?token=%s", emailLinkID)
	}
	mailContent := f.TemplateToText(tmpl, interfaceContent)
	if mailContent == "" {
		// No need to resent an error email, the error is already logged
		// We just need to inform the user that an error occurred
		mailContent = "An error occurred while trying to create your email. Please try again later. If the problem persists, please contact the administrator."
	}
	f.SendMail(email, "Confirm your Email", mailContent)

}
