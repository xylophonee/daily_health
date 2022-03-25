package mail

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
)

const (
	SMTPHost     = "smtp.qq.com"
	SMTPPort     = ":25"
	SMTPUsername = "******@qq.com"
	SMTPPassword = "******"
)

func SendMail(title, to, subject, body, format string) error {
	bs64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	header := make(map[string]string)
	header["From"] = title + "<" + SMTPUsername + ">"
	header["To"] = to
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", bs64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/" + format + "; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"
	data := ""
	for k, v := range header {
		data += k + ": " + v + "\r\n"
	}
	data += "\r\n" + bs64.EncodeToString([]byte(body))
	err := smtp.SendMail(SMTPHost+SMTPPort, smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost), SMTPUsername, []string{to}, []byte(data))
	return err
}
