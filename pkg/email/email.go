package email

import (
	"fmt"
	"github.com/go-mail/mail"
	"github.com/spf13/cast"
	"math/rand"
	"med/internal/configs"
)

type Email interface {
	SendConfirmation(receiver, pass string) error
}

type emailOpts struct {
	config *configs.Config
}

func NewEmail(config *configs.Config) Email {
	return &emailOpts{config: config}
}

func (e *emailOpts) SendConfirmation(receiver, pass string) error {
	from := e.config.From
	password := e.config.SenderPass

	smtpHost := e.config.SMTPHost
	smtpPort := e.config.SMTPPort

	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "MED  RECOMMEND")

	msg := fmt.Sprintf(`<!DOCTYPE html>
				<html>
				<head>
					<meta charset="UTF-8">
					<title>Confirmation Email</title>
					<style>
						body {
							font-family: Arial, sans-serif;
							background-color: #f4f4f4;
							padding: 20px;
						}
						.container {
							background-color: #ffffff;
							border-radius: 5px;
							box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
							padding: 20px;
							max-width: 600px;
							margin: auto;
						}
						h1 {
							color: #333333;
						}
						p {
							color: #555555;
						}
						.password {
							font-size: 20px; 
							color: #d9534f; 
							font-weight: bold; 
						}
						.footer {
							margin-top: 20px;
							font-size: 12px;
							color: #aaaaaa;
						}
					</style>
				</head>
				<body>
					<div class="container">
						<h1>ASSALMU ALLAYKUM HURMATLI MIJOZ</h1>
						<p>XURMATLI MIJOZ BIZNI TANLAGANGiZDAN JUDA XURSADNMIZ: <strong>MED SAYTI</strong></p>
						<p> TIBBIY KO"RIK ": <span class="ESLATMA: ">%s</span></p>
						<p>Thank you for choose our on <strong>med.uz</strong>!</p>
						<div class="footer">
							<p>This is an automated message, please do not reply.</p>
						</div>
					</div>
				</body>
				</html>
				`, pass)

	m.SetBody("text/html", msg)

	d := mail.NewDialer(smtpHost, cast.ToInt(smtpPort), from, password)

	err := d.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
