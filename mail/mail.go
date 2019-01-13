package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

func SendToMail(to, subject, body, mailtype string) error {

	user := "linuxMonitor@loex.com"
	password := "loexloex"
	host := "47.244.111.89:25"


	auth := LoginAuth(user, password)
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " +subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}


type loginAuth struct {
	username, password string
}
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	command := string(fromServer)
	command = strings.TrimSpace(command)
	command = strings.TrimSuffix(command, ":")
	command = strings.ToLower(command)
	if more {
		if (command == "username") {
			return []byte(fmt.Sprintf("%s", a.username)), nil
		} else if (command == "password") {
			return []byte(fmt.Sprintf("%s", a.password)), nil
		} else {
			// We've already sent everything.
			return nil, fmt.Errorf("unexpected server challenge: %s", command)
		}
	}
	return nil, nil
}
