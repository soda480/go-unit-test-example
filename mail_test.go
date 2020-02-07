package sandbox

import (
	"errors"
	"net"
	"net/smtp"
	"testing"
	"time"
)

var ()

type smtpClientMock struct {
	HelloError error
	MailError  error
	RcptError  error
	CloseError error
}

func (s *smtpClientMock) Close() error {
	return s.CloseError
}
func (s *smtpClientMock) Hello(localName string) error {
	return s.HelloError
}
func (s *smtpClientMock) Mail(from string) error {
	return s.MailError
}
func (s *smtpClientMock) Rcpt(to string) error {
	return s.RcptError
}

type addrMock struct {
}

func (*addrMock) Network() string {
	return ""
}
func (*addrMock) String() string {
	return ""
}

type netConnMock struct {
}

func (*netConnMock) Read(b []byte) (int, error) {
	return 0, nil
}
func (*netConnMock) Write(b []byte) (n int, err error) {
	return 0, nil
}
func (*netConnMock) Close() error {
	return nil
}
func (*netConnMock) LocalAddr() net.Addr {
	return &addrMock{}
}
func (*netConnMock) RemoteAddr() net.Addr {
	return &addrMock{}
}
func (*netConnMock) SetDeadline(t time.Time) error {
	return nil
}
func (*netConnMock) SetReadDeadline(t time.Time) error {
	return nil
}
func (*netConnMock) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestValidateHost_NetLookupMxError(t *testing.T) {
	var netLookupMXError = errors.New("Net Lookup MX Error")
	netLookupMX = func(name string) ([]*net.MX, error) {
		return nil, netLookupMXError
	}
	gotError := ValidateHost("mail@company.com")
	if gotError != netLookupMXError {
		t.Errorf("Failed to return expected Net Lookup MX Error")
	}
}

func TestValidateHost_GetSmtpClientError(t *testing.T) {
	var getSmtpClientError = errors.New("Get SMTP Client Error")
	netLookupMX = func(name string) ([]*net.MX, error) {
		mxs := []*net.MX{
			{
				Host: "host.tld",
				Pref: 1,
			},
		}
		return mxs, nil
	}
	getSmtpClient = func(address string) (smtpDialer, error) {
		return nil, getSmtpClientError
	}
	gotError := ValidateHost("mail@company.com")
	if gotError != getSmtpClientError {
		t.Errorf("Failed to return expected Net Lookup MX Error")
	}
}

func TestValidateHost_SmtpClientHelloError(t *testing.T) {
	var smtpClientHelloError = errors.New("SMTP Client Hello Error")
	netLookupMX = func(name string) ([]*net.MX, error) {
		mxs := []*net.MX{
			{
				Host: "host.tld",
				Pref: 1,
			},
		}
		return mxs, nil
	}
	getSmtpClient = func(address string) (smtpDialer, error) {
		client := &smtpClientMock{
			HelloError: smtpClientHelloError,
		}
		return client, nil
	}
	gotError := ValidateHost("mail@company.com")
	if gotError != smtpClientHelloError {
		t.Errorf("Failed to return expected SMTP Client Hello Error")
	}
}

func TestValidateHost_SmtpClientMailError(t *testing.T) {
	var smtpClientMailError = errors.New("SMTP Client Mail Error")
	netLookupMX = func(name string) ([]*net.MX, error) {
		mxs := []*net.MX{
			{
				Host: "host.tld",
				Pref: 1,
			},
		}
		return mxs, nil
	}
	getSmtpClient = func(address string) (smtpDialer, error) {
		client := &smtpClientMock{
			MailError: smtpClientMailError,
		}
		return client, nil
	}
	gotError := ValidateHost("mail@company.com")
	if gotError != smtpClientMailError {
		t.Errorf("Failed to return expected SMTP Client Mail Error")
	}
}

func TestValidateHost_SmtpClientRcptError(t *testing.T) {
	var smtpClientRcptError = errors.New("SMTP Client Rcpt Error")
	netLookupMX = func(name string) ([]*net.MX, error) {
		mxs := []*net.MX{
			{
				Host: "host.tld",
				Pref: 1,
			},
		}
		return mxs, nil
	}
	getSmtpClient = func(address string) (smtpDialer, error) {
		client := &smtpClientMock{
			RcptError: smtpClientRcptError,
		}
		return client, nil
	}
	gotError := ValidateHost("mail@company.com")
	if gotError != smtpClientRcptError {
		t.Errorf("Failed to return expected SMTP Client Rcpt Error")
	}
}

func TestValidateHost_SmtpClientCloseError(t *testing.T) {
	var smtpClientCloseError = errors.New("SMTP Client Close Error")
	netLookupMX = func(name string) ([]*net.MX, error) {
		mxs := []*net.MX{
			{
				Host: "host.tld",
				Pref: 1,
			},
		}
		return mxs, nil
	}
	getSmtpClient = func(address string) (smtpDialer, error) {
		client := &smtpClientMock{
			CloseError: smtpClientCloseError,
		}
		return client, nil
	}
	gotError := ValidateHost("mail@company.com")
	if gotError != smtpClientCloseError {
		t.Errorf("Failed to return expected SMTP Client Close Error")
	}
}

func TestValidateHost_SmtpClientErrors(t *testing.T) {
	var clientError = errors.New("SMTP Client Error")

	tests := []struct {
		name       string
		smtpClient *smtpClientMock
		wantError  error
	}{
		{
			name: "SmtpClientHelloError",
			smtpClient: &smtpClientMock{
				HelloError: clientError,
			},
			wantError: clientError,
		}, {
			name: "SmtpClientMailError",
			smtpClient: &smtpClientMock{
				MailError: clientError,
			},
			wantError: clientError,
		}, {
			name: "SmtpClientRcptError",
			smtpClient: &smtpClientMock{
				RcptError: clientError,
			},
			wantError: clientError,
		}, {
			name: "SmtpClientCloseError",
			smtpClient: &smtpClientMock{
				CloseError: clientError,
			},
			wantError: clientError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			netLookupMX = func(name string) ([]*net.MX, error) {
				mxs := []*net.MX{
					{
						Host: "host.tld",
						Pref: 1,
					},
				}
				return mxs, nil
			}
			getSmtpClient = func(address string) (smtpDialer, error) {
				client := test.smtpClient
				return client, nil
			}
			gotError := ValidateHost("mail@company.com")
			if gotError != test.wantError {
				t.Errorf("Failed to return expected SMTP Client Error")
			}
		})
	}
}

func TestValidateHost(t *testing.T) {
	netLookupMX = func(name string) ([]*net.MX, error) {
		mxs := []*net.MX{
			{
				Host: "host.tld",
				Pref: 1,
			},
		}
		return mxs, nil
	}
	getSmtpClient = func(address string) (smtpDialer, error) {
		client := &smtpClientMock{}
		return client, nil
	}
	gotError := ValidateHost("mail@company.com")
	if gotError != nil {
		t.Errorf("Failed error returned was not expected")
	}
}

func TestGetSmtpClient_NetDialError(t *testing.T) {
	var netDialError = errors.New("Net Dial Error")
	netDial = func(network, address string) (net.Conn, error) {
		return nil, netDialError
	}
	_, gotError := _getSmtpClient("some address")
	if gotError != netDialError {
		t.Errorf("Failed to return expected Net Dial Error")
	}
}

func TestGetSmtpClient_SmtpNewClientError(t *testing.T) {
	var smtpNewClientError = errors.New("SMTP New Client Error")
	netDial = func(network, address string) (net.Conn, error) {
		return &netConnMock{}, nil
	}
	smtpNewClient = func(conn net.Conn, host string) (*smtp.Client, error) {
		return nil, smtpNewClientError
	}
	_, gotError := _getSmtpClient("some address")
	if gotError != smtpNewClientError {
		t.Errorf("Failed to return expected SMTP New Client Error")
	}
}

func TestGetSmtpClient(t *testing.T) {
	netDial = func(network, address string) (net.Conn, error) {
		return &netConnMock{}, nil
	}
	wantClient := &smtp.Client{}
	smtpNewClient = func(conn net.Conn, host string) (*smtp.Client, error) {
		return &smtp.Client{}, nil
	}
	gotClient, gotError := _getSmtpClient("some address")
	if gotError != nil {
		t.Errorf("Failed error returned was not expected")
	}
	if gotClient == wantClient {
		t.Errorf("Failed to return expected SMTP client")
	}
}
