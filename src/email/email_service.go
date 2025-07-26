// src/email/email_service.go
package email

import (
    "bytes"
    "crypto/tls"
    "fmt"
    "net/smtp"
    "net/textproto"
    "mime/multipart"
    "path/filepath"
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

// Versión mejorada de SendAlertEmail
func (es *EmailService) SendAlertEmail(to, subject, body string) error {
    auth := smtp.PlainAuth("", es.smtpUsername, es.smtpPassword, es.smtpHost)
    
    msg := []byte(fmt.Sprintf(
        "From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
        es.fromEmail, to, subject, body,
    ))

    // Configuración TLS mejorada
    tlsconfig := &tls.Config{
        InsecureSkipVerify: true, // Cambiado a true para pruebas
        ServerName:         es.smtpHost,
    }

    // Conexión directa con TLS (SMTPS)
    conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", es.smtpHost, es.smtpPort), tlsconfig)
    if err != nil {
        return fmt.Errorf("error connecting to SMTP server with TLS: %v", err)
    }
    defer conn.Close()

    client, err := smtp.NewClient(conn, es.smtpHost)
    if err != nil {
        return fmt.Errorf("error creating SMTP client: %v", err)
    }
    defer client.Close()

    // Autenticación
    if err = client.Auth(auth); err != nil {
        return fmt.Errorf("error authenticating: %v", err)
    }

    // Configurar remitente y destinatario
    if err = client.Mail(es.fromEmail); err != nil {
        return fmt.Errorf("error setting sender: %v", err)
    }
    if err = client.Rcpt(to); err != nil {
        return fmt.Errorf("error setting recipient: %v", err)
    }

    // Enviar el mensaje
    w, err := client.Data()
    if err != nil {
        return fmt.Errorf("error preparing email data: %v", err)
    }
    defer w.Close()

    if _, err = w.Write(msg); err != nil {
        return fmt.Errorf("error writing email content: %v", err)
    }

    return nil
}

// Versión mejorada de SendAlertEmailWithAttachment
func (es *EmailService) SendAlertEmailWithAttachment(to, subject, body string, attachment *Attachment) error {
    // Crear el buffer para el mensaje MIME
    var buf bytes.Buffer
    writer := multipart.NewWriter(&buf)

    // Encabezados del mensaje
    headers := map[string]string{
        "From":         es.fromEmail,
        "To":          to,
        "Subject":     subject,
        "MIME-Version": "1.0",
        "Content-Type": "multipart/mixed; boundary=" + writer.Boundary(),
    }

    // Escribir encabezados
    for k, v := range headers {
        fmt.Fprintf(&buf, "%s: %s\r\n", k, v)
    }
    fmt.Fprintf(&buf, "\r\n")

    // Parte del texto
    part, err := writer.CreatePart(textproto.MIMEHeader{
        "Content-Type": []string{"text/plain; charset=utf-8"},
    })
    if err != nil {
        return fmt.Errorf("error creating text part: %v", err)
    }
    part.Write([]byte(body))

    // Parte del adjunto
    if attachment != nil {
        part, err := writer.CreatePart(textproto.MIMEHeader{
            "Content-Type":        []string{attachment.ContentType},
            "Content-Disposition": []string{fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(attachment.Filename))},
        })
        if err != nil {
            return fmt.Errorf("error creating attachment part: %v", err)
        }
        part.Write(attachment.Data)
    }

    writer.Close()

    // Configuración TLS mejorada
    tlsconfig := &tls.Config{
        InsecureSkipVerify: true, // Cambiado a true para pruebas
        ServerName:         es.smtpHost,
    }

    // Conexión directa con TLS (SMTPS)
    conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", es.smtpHost, es.smtpPort), tlsconfig)
    if err != nil {
        return fmt.Errorf("error connecting to SMTP server with TLS: %v", err)
    }
    defer conn.Close()

    client, err := smtp.NewClient(conn, es.smtpHost)
    if err != nil {
        return fmt.Errorf("error creating SMTP client: %v", err)
    }
    defer client.Close()

    // Autenticación
    auth := smtp.PlainAuth("", es.smtpUsername, es.smtpPassword, es.smtpHost)
    if err = client.Auth(auth); err != nil {
        return fmt.Errorf("error authenticating: %v", err)
    }

    // Configurar remitente y destinatario
    if err = client.Mail(es.fromEmail); err != nil {
        return fmt.Errorf("error setting sender: %v", err)
    }
    if err = client.Rcpt(to); err != nil {
        return fmt.Errorf("error setting recipient: %v", err)
    }

    // Enviar el mensaje
    w, err := client.Data()
    if err != nil {
        return fmt.Errorf("error preparing email data: %v", err)
    }
    defer w.Close()

    if _, err = buf.WriteTo(w); err != nil {
        return fmt.Errorf("error writing email content: %v", err)
    }

    return nil
}

type Attachment struct {
    Data        []byte
    Filename    string
    ContentType string
}