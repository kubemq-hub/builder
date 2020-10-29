package cluster

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kubemq-hub/builder/survey"
)

var jWTSignMethods = map[string]jwt.SigningMethod{
	"HS256": jwt.SigningMethodHS256,
	"HS384": jwt.SigningMethodHS384,
	"HS512": jwt.SigningMethodHS512,
	"RS256": jwt.SigningMethodRS256,
	"RS384": jwt.SigningMethodRS384,
	"RS512": jwt.SigningMethodRS512,
	"ES256": jwt.SigningMethodES256,
	"ES384": jwt.SigningMethodES384,
	"ES512": jwt.SigningMethodES512,
}

const authenticationTml = `
<red>authentication:</>
  <yellow>key:</> |-<white>{{ .Key | indent 4}}</>
  <yellow>type:</> <white>{{ .Type}}</>
`

type Authentication struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

func NewAuthentication() *Authentication {
	return &Authentication{}
}
func (a *Authentication) Clone() *Authentication {
	return &Authentication{
		Key:  a.Key,
		Type: a.Type,
	}
}
func (a *Authentication) askKey() error {
	err := survey.NewMultiline().
		SetKind("multiline").
		SetName("key").
		SetMessage("Load JWT Authentication verification key").
		SetDefault(a.Key).
		SetHelp("Load JWT Authentication verification key").
		SetRequired(true).
		Render(&a.Key)
	if err != nil {
		return err
	}
	return nil
}
func (a *Authentication) askType() error {
	err := survey.NewString().
		SetKind("string").
		SetName("type").
		SetMessage("Set JWT signing method").
		SetDefault(a.Type).
		SetOptions([]string{"HS256", "HS384", "HS512", "RS256", "RS384", "RS512", "ES256", "ES384", "ES512"}).
		SetHelp("Set JWT signing Method").
		SetRequired(true).
		Render(&a.Type)
	if err != nil {
		return err
	}
	return nil
}

func (a *Authentication) Validate() error {
	var err error
	signType, ok := jWTSignMethods[a.Type]
	if !ok {
		return fmt.Errorf("invalid jwt singing method: %s", a.Type)
	}

	switch signType {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:

	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
		_, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(a.Key))
		if err != nil {
			return err
		}
	case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		_, err = jwt.ParseECPrivateKeyFromPEM([]byte(a.Key))
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Authentication) Render() (*Authentication, error) {
	if err := a.askKey(); err != nil {
		return nil, err
	}
	if err := a.askType(); err != nil {
		return nil, err
	}
	return a, nil
}
func (a *Authentication) ColoredYaml() (string, error) {
	t := NewTemplate(authenticationTml, a)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring authentication spec,%s", err.Error()), nil
	}
	return string(b), nil
}

var _ Validator = NewAuthentication()
