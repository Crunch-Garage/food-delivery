package helper

import (
	"bytes"
	"crunchgarage/restaurant-food-delivery/models"
	"crypto/tls"
	"html/template"
	"os"

	gomail "gopkg.in/gomail.v2"
)

var smtp_email = os.Getenv("SMTP_EMAIL")
var smtp_password = os.Getenv("SMTP_EMAIL_PASSWORD")
var err error

func RegisterEmailAccount(user models.User) (string, error) {

	t := template.New("register.html")

	t, err = t.ParseFiles("html/register/register.html")

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, struct {
		Name    string
		Message string
	}{
		Name:    user.First_name,
		Message: "Welcome to EatsFood. Please take a sec to confirm your email.",
	}); err != nil {
		return "", err
	}

	result := tpl.String()

	msg := ""

	msg, err = HandleSendEmail(user.Email, "Please Confirm Your E-mail Address", result)

	if err != nil {
		return "", err
	}

	return msg, nil

}

func HandleSendEmail(to, subject, body string) (string, error) {

	m := gomail.NewMessage()

	/*set email sender*/
	m.SetHeader("From", smtp_email)

	/*set email receiver*/
	m.SetHeader("To", to)

	/*set email subject*/
	m.SetHeader("Subject", subject)

	/*set email body*/
	m.SetBody("text/html", body)

	/*settings for SMTP server*/
	d := gomail.NewDialer("smtp.gmail.com", 587, smtp_email, smtp_password)

	/*this is only needed when SSL/TLS certificateis not valid on server*/
	/*In production this should be set to false*/
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	/*send email*/
	if err := d.DialAndSend(m); err != nil {
		return "", err
	}

	msg := "Email sent successfully"
	return msg, nil

}
