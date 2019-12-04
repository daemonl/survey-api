package awssecret

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"gopkg.daemonl.com/envconf"
)

type Translator struct {
	service *secretsmanager.SecretsManager
}

func NewTranslator(service *secretsmanager.SecretsManager) (*Translator, error) {
	return &Translator{
		service: service,
	}, nil
}

func (t Translator) Translate(in string) (string, error) {
	output, err := t.service.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(in),
	})
	if err != nil {
		return "", err
	}
	return aws.StringValue(output.SecretString), nil
}

func Default() error {
	translator, err := NewFromEnv()
	if err != nil {
		return err
	}
	envconf.DefaultParser.Translators["ssm"] = translator
	return nil
}

func NewFromEnv() (*Translator, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	service := secretsmanager.New(sess)
	return NewTranslator(service)
}
