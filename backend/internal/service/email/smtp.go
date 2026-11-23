package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/config"
)

const maxMessageBytes = 1 << 20

type Message struct {
	To       string
	Subject  string
	HTMLBody string
	TextBody string
}

type SMTPConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
	Timeout     time.Duration
}

type SMTPMailer struct {
	config  SMTPConfig
	rootCAs *x509.CertPool
}

func NewSMTPMailer(config SMTPConfig) (*SMTPMailer, error) {
	if err := validateSMTPConfig(config); err != nil {
		return nil, err
	}
	if config.Timeout <= 0 {
		config.Timeout = 15 * time.Second
	}
	return &SMTPMailer{config: config}, nil
}

func NewSMTPMailerFromConfig(config config.EmailConfig) (*SMTPMailer, error) {
	return NewSMTPMailer(SMTPConfig{
		Host: config.SMTPHost, Port: config.SMTPPort, Username: config.SMTPUsername,
		Password: config.SMTPPassword, FromAddress: config.FromEmail, FromName: config.FromName,
	})
}

func validateSMTPConfig(config SMTPConfig) error {
	if strings.TrimSpace(config.Host) == "" || strings.ContainsAny(config.Host, "\r\n :/") {
		return fmt.Errorf("SMTP host is invalid")
	}
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("SMTP port is invalid")
	}
	if strings.TrimSpace(config.Username) == "" || config.Password == "" {
		return fmt.Errorf("SMTP credentials are required")
	}
	if _, err := parseMailbox(config.FromAddress); err != nil {
		return fmt.Errorf("SMTP sender address is invalid")
	}
	if hasHeaderBreak(config.FromName) {
		return fmt.Errorf("SMTP sender name contains an invalid line break")
	}
	return nil
}

func (mailer *SMTPMailer) Send(ctx context.Context, message Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	recipient, err := parseMailbox(message.To)
	if err != nil || message.Subject == "" || hasHeaderBreak(message.Subject) {
		return fmt.Errorf("email recipient or subject is invalid")
	}
	if len(message.HTMLBody)+len(message.TextBody) > maxMessageBytes || message.HTMLBody == "" || message.TextBody == "" {
		return fmt.Errorf("email body is missing or exceeds the size limit")
	}
	payload, err := mailer.mimeMessage(recipient, message)
	if err != nil {
		return err
	}
	address := net.JoinHostPort(mailer.config.Host, strconv.Itoa(mailer.config.Port))
	connection, err := (&net.Dialer{Timeout: mailer.config.Timeout}).DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connect to SMTP server: %w", err)
	}
	stopClose := context.AfterFunc(ctx, func() { _ = connection.Close() })
	defer stopClose()
	defer connection.Close()
	deadline := time.Now().Add(mailer.config.Timeout)
	if contextDeadline, ok := ctx.Deadline(); ok && contextDeadline.Before(deadline) {
		deadline = contextDeadline
	}
	_ = connection.SetDeadline(deadline)
	return mailer.sendOverTLS(connection, recipient.Address, payload)
}

func (mailer *SMTPMailer) sendOverTLS(connection net.Conn, recipient string, payload []byte) error {
	client, err := smtp.NewClient(connection, mailer.config.Host)
	if err != nil {
		return fmt.Errorf("start SMTP session: %w", err)
	}
	defer client.Close()
	if err := client.Hello("servio"); err != nil {
		return fmt.Errorf("identify SMTP client: %w", err)
	}
	if supported, _ := client.Extension("STARTTLS"); !supported {
		return fmt.Errorf("SMTP server does not offer required STARTTLS")
	}
	tlsConfig := &tls.Config{ServerName: mailer.config.Host, MinVersion: tls.VersionTLS12, RootCAs: mailer.rootCAs}
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("upgrade SMTP connection to TLS: %w", err)
	}
	if err := client.Auth(smtp.PlainAuth("", mailer.config.Username, mailer.config.Password, mailer.config.Host)); err != nil {
		return fmt.Errorf("authenticate with SMTP server: %w", err)
	}
	from, _ := parseMailbox(mailer.config.FromAddress)
	if err := client.Mail(from.Address); err != nil {
		return fmt.Errorf("set SMTP sender: %w", err)
	}
	if err := client.Rcpt(recipient); err != nil {
		return fmt.Errorf("set SMTP recipient: %w", err)
	}
	data, err := client.Data()
	if err != nil {
		return fmt.Errorf("start SMTP message: %w", err)
	}
	if _, err := data.Write(payload); err != nil {
		_ = data.Close()
		return fmt.Errorf("write SMTP message: %w", err)
	}
	if err := data.Close(); err != nil {
		return fmt.Errorf("finish SMTP message: %w", err)
	}
	return nil
}

func (mailer *SMTPMailer) mimeMessage(recipient *mail.Address, message Message) ([]byte, error) {
	from, _ := parseMailbox(mailer.config.FromAddress)
	if mailer.config.FromName != "" {
		from.Name = mailer.config.FromName
	}
	var body bytes.Buffer
	multipartWriter := multipart.NewWriter(&body)
	if err := writeTextPart(multipartWriter, "text/plain", message.TextBody); err != nil {
		return nil, err
	}
	if err := writeTextPart(multipartWriter, "text/html", message.HTMLBody); err != nil {
		return nil, err
	}
	if err := multipartWriter.Close(); err != nil {
		return nil, fmt.Errorf("build email body: %w", err)
	}
	var output bytes.Buffer
	fmt.Fprintf(&output, "From: %s\r\n", from.String())
	fmt.Fprintf(&output, "To: %s\r\n", recipient.String())
	fmt.Fprintf(&output, "Subject: %s\r\n", mime.QEncoding.Encode("UTF-8", message.Subject))
	fmt.Fprintf(&output, "MIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=%q\r\n\r\n", multipartWriter.Boundary())
	if _, err := io.Copy(&output, &body); err != nil {
		return nil, fmt.Errorf("assemble email message: %w", err)
	}
	return output.Bytes(), nil
}

func writeTextPart(writer *multipart.Writer, contentType, value string) error {
	part, err := writer.CreatePart(map[string][]string{
		"Content-Type":              {contentType + "; charset=UTF-8"},
		"Content-Transfer-Encoding": {"base64"},
	})
	if err != nil {
		return fmt.Errorf("create email body part: %w", err)
	}
	encoder := base64.NewEncoder(base64.StdEncoding, part)
	if _, err := encoder.Write([]byte(value)); err != nil {
		return fmt.Errorf("encode email body: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("finish email body encoding: %w", err)
	}
	return nil
}

func parseMailbox(value string) (*mail.Address, error) {
	address, err := mail.ParseAddress(strings.TrimSpace(value))
	if err != nil || address.Address == "" || hasHeaderBreak(address.Address) {
		return nil, fmt.Errorf("invalid email address")
	}
	return address, nil
}

func hasHeaderBreak(value string) bool {
	return strings.ContainsAny(value, "\r\n")
}
