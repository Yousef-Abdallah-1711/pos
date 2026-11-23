package email

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"
)

func TestSMTPMailerRequiresSTARTTLSAndAuthenticatesAfterUpgrade(t *testing.T) {
	address, roots, received, serverErrors := startMockSMTP(t, true)
	connection, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	mailer, err := NewSMTPMailer(SMTPConfig{
		Host: "127.0.0.1", Port: connection.Port, Username: "smtp-user", Password: "smtp-secret",
		FromAddress: "SERVIO <mailer@example.test>", FromName: "SERVIO", Timeout: 3 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	mailer.rootCAs = roots
	err = mailer.Send(context.Background(), Message{
		To: "staff@example.test", Subject: "Invitation", HTMLBody: "<p>Join</p>", TextBody: "Join",
	})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
	select {
	case message := <-received:
		if !strings.Contains(message, "multipart/alternative") || !strings.Contains(message, "staff@example.test") {
			t.Fatalf("captured message is incomplete: %s", message)
		}
		if !strings.Contains(message, "Subject: Invitation") {
			t.Fatalf("captured message has no subject: %s", message)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("SMTP test server did not receive a message")
	}
	select {
	case err := <-serverErrors:
		if err != nil {
			t.Fatal(err)
		}
	default:
	}
}

func TestSMTPMailerDoesNotSendCredentialsWithoutSTARTTLS(t *testing.T) {
	address, _, _, _ := startMockSMTP(t, false)
	connection, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	mailer, err := NewSMTPMailer(SMTPConfig{
		Host: "127.0.0.1", Port: connection.Port, Username: "smtp-user", Password: "smtp-secret",
		FromAddress: "mailer@example.test", Timeout: 3 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = mailer.Send(context.Background(), Message{
		To: "staff@example.test", Subject: "Notice", HTMLBody: "<p>Notice</p>", TextBody: "Notice",
	})
	if err == nil || !strings.Contains(err.Error(), "STARTTLS") {
		t.Fatalf("expected missing STARTTLS to fail closed, got %v", err)
	}
}

func TestSMTPMailerRejectsHeaderInjectionBeforeConnecting(t *testing.T) {
	mailer, err := NewSMTPMailer(SMTPConfig{
		Host: "localhost", Port: 587, Username: "user", Password: "password", FromAddress: "mailer@example.test",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = mailer.Send(context.Background(), Message{
		To: "staff@example.test", Subject: "Hello\r\nBcc: attacker@example.test", HTMLBody: "<p>Hi</p>", TextBody: "Hi",
	})
	if err == nil {
		t.Fatal("expected header injection to be rejected")
	}
}

func startMockSMTP(t *testing.T, advertiseSTARTTLS bool) (string, *x509.CertPool, <-chan string, <-chan error) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	certificate, roots := testCertificate(t)
	received := make(chan string, 1)
	serverErrors := make(chan error, 1)
	go func() {
		defer listener.Close()
		connection, err := listener.Accept()
		if err != nil {
			serverErrors <- err
			return
		}
		defer connection.Close()
		if err := serveSMTPConnection(connection, certificate, advertiseSTARTTLS, received); err != nil {
			serverErrors <- err
			return
		}
		serverErrors <- nil
	}()
	return listener.Addr().String(), roots, received, serverErrors
}

func serveSMTPConnection(connection net.Conn, certificate tls.Certificate, advertiseSTARTTLS bool, received chan<- string) error {
	reader := bufio.NewReader(connection)
	if _, err := fmt.Fprint(connection, "220 localhost ESMTP\r\n"); err != nil {
		return err
	}
	tlsActive := false
	var message strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF && message.Len() > 0 {
				return nil
			}
			return err
		}
		command := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(command, "EHLO "):
			if !tlsActive && advertiseSTARTTLS {
				_, err = fmt.Fprint(connection, "250-localhost\r\n250-STARTTLS\r\n250 AUTH PLAIN\r\n")
			} else {
				_, err = fmt.Fprint(connection, "250-localhost\r\n250 AUTH PLAIN\r\n")
			}
		case command == "STARTTLS" && advertiseSTARTTLS:
			if _, err = fmt.Fprint(connection, "220 Ready to start TLS\r\n"); err == nil {
				secureConnection := tls.Server(connection, &tls.Config{Certificates: []tls.Certificate{certificate}, MinVersion: tls.VersionTLS12})
				if err = secureConnection.Handshake(); err == nil {
					connection, reader, tlsActive = secureConnection, bufio.NewReader(secureConnection), true
				}
			}
		case strings.HasPrefix(command, "AUTH PLAIN "):
			credentialBytes, decodeErr := base64.StdEncoding.DecodeString(strings.TrimPrefix(command, "AUTH PLAIN "))
			if decodeErr != nil || string(credentialBytes) != "\x00smtp-user\x00smtp-secret" || !tlsActive {
				_, err = fmt.Fprint(connection, "535 Authentication failed\r\n")
			} else {
				_, err = fmt.Fprint(connection, "235 Authenticated\r\n")
			}
		case strings.HasPrefix(command, "MAIL FROM:") || strings.HasPrefix(command, "RCPT TO:"):
			_, err = fmt.Fprint(connection, "250 Accepted\r\n")
		case command == "DATA":
			_, err = fmt.Fprint(connection, "354 Continue\r\n")
			if err == nil {
				for {
					var dataLine string
					dataLine, err = reader.ReadString('\n')
					if err != nil || strings.TrimRight(dataLine, "\r\n") == "." {
						break
					}
					message.WriteString(dataLine)
				}
				if err == nil {
					received <- message.String()
					_, err = fmt.Fprint(connection, "250 Queued\r\n")
				}
			}
		case command == "QUIT":
			_, err = fmt.Fprint(connection, "221 Bye\r\n")
			return err
		default:
			_, err = fmt.Fprint(connection, "500 Unexpected command\r\n")
		}
		if err != nil {
			return err
		}
	}
}

func testCertificate(t *testing.T) (tls.Certificate, *x509.CertPool) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 120))
	if err != nil {
		t.Fatal(err)
	}
	template := x509.Certificate{
		SerialNumber: serial, NotBefore: time.Now().Add(-time.Minute), NotAfter: time.Now().Add(time.Hour),
		DNSNames: []string{"localhost"}, IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := tls.X509KeyPair(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
		pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}),
	)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(parsed)
	return certificate, roots
}
