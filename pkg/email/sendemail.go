package email

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

const ROOT_EMAIL = "2718114795@qq.com"
const SMTP_KEY = "qhgxotidoomcdhbe"

func SendEmail(targetEmail string, content string) error {
	smtpServer := "smtp.qq.com"
	emailAddr := ROOT_EMAIL
	smtpKey := SMTP_KEY

	em := email.NewEmail()
	em.From = fmt.Sprintf("GoMeeting <%s>", emailAddr)
	em.To = []string{targetEmail}
	em.Subject = "GoMeeting 验证码"
	em.Text = []byte(fmt.Sprintf("您的验证码是：%s，有效期5分钟。", content))
	//em.Text = []byte(content)

	// 使用QQ邮箱587端口 - 最常用的配置
	auth := smtp.PlainAuth("", emailAddr, smtpKey, smtpServer)

	// 直接使用Send方法，让库自动处理STARTTLS
	err := em.Send(smtpServer+":587", auth)
	if err != nil {
		//// 如果587端口失败，尝试25端口作为备选
		//err2 := em.Send(smtpServer+":25", auth)
		//if err2 != nil {
		//	return fmt.Errorf("发送邮件失败 (587端口: %w, 25端口: %w)", err, err2)
		//}
		//return fmt.Errorf("发送邮件失败 (587端口: %w)", err)
		return nil
	}
	return nil
}
