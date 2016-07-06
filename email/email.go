// Copyright 2016  Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type SmtpConfig struct {
	Username string
	Password string
	Host     string
	Addr     string
}

type SmtpHandler struct {
	smtp.Auth
	Addr    string
	From    string
	To      []string
	Subject string
}

func NewSmtpHandler(config SmtpConfig, from string, to []string, subject string) (*SmtpHandler) {
	auth := smtp.PlainAuth(
		"",
		config.Username,
		config.Password,
		config.Host,
	)
	return &SmtpHandler{auth, config.Addr, from, to, subject}
}

func (s *SmtpHandler)SendMail(subject string, message string, isHtml bool) (error) {

	contentType := "text/plain"
	if isHtml {
		contentType = "text/html"
	}

	if ("" == subject) {
		subject = s.Subject
	}

	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: %s; charset=UTF-8\r\n\r\n%s", strings.Join(s.To, ";"), s.From, subject, contentType, message)
	return smtp.SendMail(s.Addr, s.Auth, s.From, s.To, []byte(msg))
}

func (s *SmtpHandler)AddToAccount(to... string) {
	s.To = append(s.To, to...)
}

func (s *SmtpHandler)SetFrom(from string) {
	s.From = from
}

func (s *SmtpHandler)SetSubject(subject string) {
	s.Subject = subject
}

//// send mail
//func SendMail(subject string, message string, from string, to []string, smtpConfig SmtpConfig, isHtml bool) error {
//	auth := smtp.PlainAuth(
//		"",
//		smtpConfig.Username,
//		smtpConfig.Password,
//		smtpConfig.Host,
//	)
//	contentType := "text/plain"
//	if isHtml {
//		contentType = "text/html"
//	}
//	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: %s; charset=UTF-8\r\n\r\n%s", strings.Join(to, ";"), from, subject, contentType, message)
//	return smtp.SendMail(smtpConfig.Addr, auth, from, to, []byte(msg))
//}

//// exec /usr/sbin/sendmail -t -i
//func SendMailExec(subject string, message string, from string, to []string, sendmailPath string, isHtml bool) error {
//	cmdArgs := strings.Fields(sendmailPath)
//	cmdArgs = append(cmdArgs, "-f", from)
//	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
//	cmdStdin, err := cmd.StdinPipe()
//	if err != nil {
//		return err
//	}
//	if err := cmd.Start(); err != nil {
//		return err
//	}
//	contentType := "text/plain"
//	if isHtml {
//		contentType = "text/html"
//	}
//	_, err = fmt.Fprintf(cmdStdin, "To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: %s; charset=UTF-8\r\n\r\n%s", strings.Join(to, ";"), from, subject, contentType, message)
//	err = cmdStdin.Close()
//	if err != nil {
//		return err
//	}
//	return cmd.Wait()
//}
