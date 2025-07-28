package mail

import (
    "fmt"
    "net/smtp"
)

type SMTPMailer struct {
    Host     string
    Port     string
    Username string
    Password string
    From     string
}

func NewSMTPMailer(host, port, username, password, from string) *SMTPMailer {
    return &SMTPMailer{
        Host:     host,
        Port:     port,
        Username: username,
        Password: password,
        From:     from,
    }
}

func (m *SMTPMailer) SendVerificationEmail(to, link string) error {
    auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
    msg := fmt.Sprintf("To: %s\r\nSubject: Подтверждение почты\r\n\r\n"+
        "Привет!\nПодтверди почту, перейдя по ссылке:\n%s", to, link)

    addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
    err := smtp.SendMail(addr, auth, m.From, []string{to}, []byte(msg))
    if err != nil {
        fmt.Println("SMTP send error:", err)
    }
    return err
}

