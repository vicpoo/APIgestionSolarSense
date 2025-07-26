// src/email/email_service.go
package email

import (
    "crypto/tls"
    "fmt"
  
    "net/smtp"
)

type EmailService struct {
    smtpHost     string
    smtpPort     int
    smtpUsername string
    smtpPassword string
    fromEmail    string
}

func NewEmailService(host string, port int, username, password, from string) *EmailService {
    return &EmailService{
        smtpHost:     host,
        smtpPort:     port,
        smtpUsername: username,
        smtpPassword: password,
        fromEmail:    from,
    }
}

func (es *EmailService) SendAlertEmail(to, subject, body string) error {
    auth := smtp.PlainAuth("", es.smtpUsername, es.smtpPassword, es.smtpHost)
    
    msg := []byte(fmt.Sprintf(
        "From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
        es.fromEmail, to, subject, body,
    ))

    // Configuración alternativa para conexión segura
    tlsconfig := &tls.Config{
        InsecureSkipVerify: false,
        ServerName:         es.smtpHost,
    }

    // Primero intentar con STARTTLS
    c, err := smtp.Dial(fmt.Sprintf("%s:%d", es.smtpHost, es.smtpPort))
    if err != nil {
        return fmt.Errorf("error dialing server: %v", err)
    }
    defer c.Close()

    if err = c.StartTLS(tlsconfig); err != nil {
        return fmt.Errorf("error starting TLS: %v", err)
    }

    if err = c.Auth(auth); err != nil {
        return fmt.Errorf("error authenticating: %v", err)
    }

    if err = c.Mail(es.fromEmail); err != nil {
        return fmt.Errorf("error setting sender: %v", err)
    }

    if err = c.Rcpt(to); err != nil {
        return fmt.Errorf("error setting recipient: %v", err)
    }

    w, err := c.Data()
    if err != nil {
        return fmt.Errorf("error preparing data: %v", err)
    }
    defer w.Close()

    if _, err = w.Write(msg); err != nil {
        return fmt.Errorf("error writing message: %v", err)
    }

    return nil
}