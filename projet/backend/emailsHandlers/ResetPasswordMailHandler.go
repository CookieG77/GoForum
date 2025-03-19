package emailsHandlers

import (
	f "GoForum/functions"
	"html/template"
)

func SendResetPasswordMail(email string) {
	tmpl, err := template.ParseFiles("templates/emails/resetPasswordEmail.html")
	if err != nil {
		f.ErrorPrintf("An error occurred while trying to parse the template -> %v\n", err)
		return
	}
	interfaceContent := make(map[string]interface{})
	interfaceContent["Url"] = "URL"
	interfaceContent["UrlSimplified"] = "SimplifiedURL"
	mailContent := f.TemplateToText(tmpl, interfaceContent)
	if mailContent == "" {
		// No need to resent an error email, the error is already logged
		// We just need to inform the user that an error occurred
		mailContent = "An error occurred while trying to create your email. Please try again later. If the problem persists, please contact the administrator."
	}
	f.SendMail(email, "Reset your password", mailContent)

}
