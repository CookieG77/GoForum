package emailsHandlers

import (
	f "GoForum/functions"
	"fmt"
	"html/template"
)

func SendResetPasswordMail(email string) {
	tmpl, err := template.ParseFiles("templates/emails/resetPasswordEmail.html")
	if err != nil {
		f.ErrorPrintf("An error occurred while trying to parse the template -> %v\n", err)
		return
	}
	emailLinkID, err := f.CreateEmailIdentificationLink(email, f.ResetPasswordEmail)
	if err != nil {
		f.ErrorPrintf("Error while creating the email identification link: %s\n", err)
		return
	}
	interfaceContent := make(map[string]interface{})
	if f.IsCertified() {
		interfaceContent["Url"] = fmt.Sprintf("https://localhost/reset-password?id=%s", emailLinkID)
	} else {
		interfaceContent["Url"] = fmt.Sprintf("http://localhost/reset-password?id=%s", emailLinkID)
	}
	mailContent := f.TemplateToText(tmpl, interfaceContent)
	if mailContent == "" {
		// No need to resent an error email, the error is already logged
		// We just need to inform the user that an error occurred
		mailContent = "An error occurred while trying to create your email. Please try again later. If the problem persists, please contact the administrator."
	}
	f.SendMail(email, "Reset your password", mailContent)

}
