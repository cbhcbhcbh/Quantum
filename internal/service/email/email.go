package email

import (
	"context"
	"fmt"
	"net/smtp"
	"time"

	"github.com/cbhcbhcbh/Quantum/internal/config"
	"github.com/cbhcbhcbh/Quantum/pkg/redis"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
)

const (
	REGISTERED_CODE = 1
	RESET_PS_CODE   = 2
)

type IEmailService interface {
	SendEmail(code string, emailType int, email string, subject string, body string) error
	getCacheFix(email string, emailType int) string
	CheckCode(email string, code string, emailType int) bool
}

type EmailService struct {
	emailConfig *config.EmailConfig
}

func NewEmailService() EmailService {
	emailService := config.InitEmailConfig()

	return EmailService{
		emailConfig: emailService,
	}
}

func (e EmailService) SendEmail(code string, emailType int, email string, subject string, body string) error {
	header := make(map[string]string)

	header["From"] = "im-service:" + "<" + e.emailConfig.Name + ">"
	header["To"] = email
	header["Subject"] = subject
	header["Content-Type"] = "text/html;chartset=UTF-8"

	message := ""

	for k, v := range header {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}

	message += "\r\n" + body

	auth := smtp.PlainAuth(
		"",
		e.emailConfig.Name,
		e.emailConfig.Password,
		e.emailConfig.Host,
	)

	err := sendMailUsingTLS(
		fmt.Sprintf("%s:%d", e.emailConfig.Host, e.emailConfig.Port),
		auth,
		e.emailConfig.Name,
		[]string{email},
		[]byte(message),
	)

	if err != nil {
		return err
	}
	redis.R.SetKey(e.getCacheFix(email, emailType), code, time.Minute*1)
	return nil
}

func (e EmailService) getCacheFix(email string, emailType int) string {
	switch emailType {
	case REGISTERED_CODE:
		return fmt.Sprintf("%s.%d", email, REGISTERED_CODE)
	case RESET_PS_CODE:
		return fmt.Sprintf("%s.%d", email, RESET_PS_CODE)
	default:
		return fmt.Sprintf("%s.%d", email, REGISTERED_CODE)
	}
}

func (e EmailService) CheckCode(email string, code string, emailType int) bool {
	val := redis.R.GetKey(e.getCacheFix(email, emailType))

	if val == code {
		redis.R.DelKey(e.getCacheFix(email, emailType))
		return true
	}

	return false
}

func sendMailUsingTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		log.C(context.TODO()).Errorw("Failed to dial SMTP server", "error", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.C(context.TODO()).Errorw("Error during AUTH")
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addrs := range to {
		if err = c.Rcpt(addrs); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
