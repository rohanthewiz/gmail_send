package gmail_send

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"

	"github.com/rohanthewiz/serr"
)

// Send HTML email via your Gmail account
// From Address will always be the Gmail Account email
func GmailSend(cfg GSMTPConfig) (err error) {
	const servername = "smtp.gmail.com:465" // 465 is the TLS only port
	rpl := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	host, _, _ := net.SplitHostPort(servername)

	tlsconfig := &tls.Config {
		InsecureSkipVerify: true,
		ServerName: host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require a TLS connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return serr.Wrap(err, "Error dialing mail server", "mailServer", servername)
	}

	// TODO - Wrap and return all error below
	cl, err := smtp.NewClient(conn, host)
	if err != nil {
		return serr.Wrap(err, "Error creating smtp client")
	}

	// Auth
	auth := smtp.PlainAuth("", cfg.AccountEmail, cfg.Word, host)
	if err = cl.Auth(auth); err != nil {
		return serr.Wrap(err, "Error setting up smtp client auth", "account", cfg.AccountEmail)
	}

	// From is required in the client
	if err = cl.Mail(cfg.AccountEmail); err != nil {
		return serr.Wrap(err, "Error applying From address to smtp client")
	}

	// Bcc - these must not be included in the mail body
	for _, bcc := range cfg.BCCs {
		eml := rpl.Replace(bcc)
		if err := cl.Rcpt(eml); err != nil {
			log.Println("Error applying bcc address: " + eml, err)
			continue
		}
	}

	// To
	var strTos []string
	for _, to := range cfg.ToAddrs {
		if err = cl.Rcpt(to); err != nil {
			log.Println("Error applying to address: " + to, err)
			continue
		}
		strTos = append(strTos, rpl.Replace(to))
	}

	// Data
	w, err := cl.Data() // get a data writer from the server
	if err != nil {
		return serr.Wrap(err, "Error obtaining a writer from the smtp server")
	}

	// Setup Mail headers
	headers := make(map[string]string)
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""
	headers["Content-Transfer-Encoding"] = "base64"
	if cfg.FromName != "" {
		headers["From"] = cfg.FromName
	}
	headers["To"] = strings.Join(strTos, ",") // don't include bcc here
	headers["Subject"] = cfg.Subject

	// Setup message
	message := ""
	for k,v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(cfg.Body))

	_, err = w.Write([]byte(message))
	if err != nil {
		return serr.Wrap(err, "Error writing email message to server")
	}

	err = w.Close()
	if err != nil {
		log.Println("Error closing email message writer", err)
		// allow to fall through
	}

	err = cl.Quit()
	if err != nil {
		return serr.Wrap(err, "Error quitting client")
	}

	return
}

