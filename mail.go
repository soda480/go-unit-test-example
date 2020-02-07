package sandbox

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

var (
	// patchable variables
	netLookupMX   = net.LookupMX
	netDial       = net.Dial
	smtpNewClient = smtp.NewClient
	getSmtpClient = _getSmtpClient
)

type smtpDialer interface {
	Close() error
	Hello(localName string) error
	Mail(from string) error
	Rcpt(to string) error
}

func ValidateHost(email string) (err error) {
	mx, err := netLookupMX(host(email))
	if err != nil {
		return err
	}

	smtpClient, err := getSmtpClient(fmt.Sprintf("%s:%d", mx[0].Host, 25))
	if err != nil {
		return err
	}

	defer func() {
		if er := smtpClient.Close(); er != nil {
			err = er
		}
	}()

	if err = smtpClient.Hello("checkmail.me"); err != nil {
		return err
	}

	if err = smtpClient.Mail("testing-email-host@gmail.com"); err != nil {
		return err
	}

	return smtpClient.Rcpt(email)
}

func host(email string) (host string) {
	index := strings.LastIndexByte(email, '@')
	return email[index+1:]
}

func _getSmtpClient(address string) (smtpDialer, error) {
	connection, err := netDial("tcp", address)
	if err != nil {
		return nil, err
	}

	smtpClient, err := smtpNewClient(connection, address)
	if err != nil {
		return nil, err
	}

	return smtpClient, nil
}
