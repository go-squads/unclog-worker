package alertic

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/smtp"
	"strconv"
	"strings"
)

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func buildMail(appName, logLevel string, currentCount int, receivers []string) (mail Mail) {
	mail = Mail{}
	mail.senderId = viper.GetString("EMAIL_SENDER")
	mail.toIds = receivers
	mail.subject = "[Barito Alert] " + appName
	mail.body = appName + " had " + strconv.Itoa(currentCount) + " " + logLevel

	mail.body = mail.BuildMessage()

	return
}

func setupConnection(mail Mail) (client *smtp.Client) {
	smtpServer := SmtpServer{host: viper.GetString("SMTP_SERVER"), port: viper.GetString("SMTP_PORT")}

	log.Println(smtpServer.host)

	auth := smtp.PlainAuth("", mail.senderId, viper.GetString("EMAIL_SENDER_PASS"), smtpServer.host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err = smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Panic(err)
	}

	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	return
}

func send(client smtp.Client, mail Mail) {
	if err := client.Mail(mail.senderId); err != nil {
		log.Panic(err)
	}

	for _, k := range mail.toIds {
		if err := client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(mail.body))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}
}

func SendEmail(appName, logLevel string, currentCount int, receivers []string) {
	mail := buildMail(appName, logLevel, currentCount, receivers)
	client := setupConnection(mail)
	send(*client, mail)

	log.Println("Mail sent successfully to ", mail.toIds)
}
