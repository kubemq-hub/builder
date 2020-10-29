package cluster

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

const tlsTml = `
<red>tls:</>
  <yellow>cert:</> |-<white>{{ .Cert | indent 4}}</>
  <yellow>key:</> |-<white>{{ .Key | indent 4}}</>
  <yellow>ca:</> |-<white>{{ .Ca | indent 4}}</>
`

type Tls struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
	Ca   string `json:"ca"`
}

func NewTls() *Tls {
	return &Tls{}
}
func (t *Tls) Clone() *Tls {
	return &Tls{
		Cert: t.Cert,
		Key:  t.Key,
		Ca:   t.Ca,
	}
}
func (t *Tls) askCert() error {
	err := survey.NewMultiline().
		SetKind("multiline").
		SetName("cert").
		SetMessage("Load TLS certificate data").
		SetDefault(t.Cert).
		SetHelp("Load TLS certificate data").
		SetRequired(true).
		Render(&t.Cert)
	if err != nil {
		return err
	}
	return nil
}
func (t *Tls) askKey() error {
	err := survey.NewMultiline().
		SetKind("multiline").
		SetName("key").
		SetMessage("Load TLS key data").
		SetDefault(t.Key).
		SetHelp("Load TLS key data").
		SetRequired(true).
		Render(&t.Key)
	if err != nil {
		return err
	}
	return nil
}
func (t *Tls) askCA() error {
	boolVal := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("ca").
		SetMessage("Would you like to load TLS certificate authority (CA) data").
		SetDefault("false").
		SetHelp("Load TLS certificate authority (CA) data").
		SetRequired(true).
		Render(&boolVal)
	if err != nil {
		return err
	}
	if !boolVal {
		t.Ca = ""
		return nil
	}
	err = survey.NewMultiline().
		SetKind("multiline").
		SetName("ca").
		SetMessage("Load TLS certificate authority (CA) data").
		SetDefault(t.Ca).
		SetHelp("Load TLS certificate authority (CA) data").
		SetRequired(false).
		Render(&t.Ca)
	if err != nil {
		return err
	}
	return nil
}
func (t *Tls) Validate() error {
	var err error
	cfg := new(tls.Config)
	cfg.Certificates = make([]tls.Certificate, 1)
	if cfg.Certificates[0], err = tls.X509KeyPair([]byte(t.Cert), []byte(t.Key)); err != nil {
		return err
	}
	if t.Ca != "" {
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM([]byte(t.Ca))
		if !ok {
			return fmt.Errorf("inavlid certificate authority (CA)")
		}
	}
	return nil
}

func (t *Tls) Render() (*Tls, error) {
	if err := t.askCert(); err != nil {
		return nil, err
	}
	if err := t.askKey(); err != nil {
		return nil, err
	}
	if err := t.askCA(); err != nil {
		return nil, err
	}
	return t, nil
}
func (t *Tls) ColoredYaml() (string, error) {
	tmpl := NewTemplate(tlsTml, t)
	b, err := tmpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring tls spec,%s", err.Error()), nil
	}
	return string(b), nil
}
