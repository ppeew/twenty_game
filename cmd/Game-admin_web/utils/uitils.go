package utils

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
)

func CreateTable() {
}

// SendMessage 发送消息
func SendMessage(toEmail []string, msg string) error {
	if msg == "" {
		return nil
	}
	e := email.NewEmail()
	e.From = "<2069234934@qq.com>"
	e.To = toEmail
	if len(toEmail) == 0 {
		e.To = []string{"3194044365@qq.com"}
	}
	e.Subject = "twenty_game服务挂机通知"
	//e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<b>" + msg + "</b>")
	return e.SendWithTLS("smtp.qq.com:465",
		smtp.PlainAuth("", "2069234934@qq.com", "hzhqhrqvubaodcbi", "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
}
