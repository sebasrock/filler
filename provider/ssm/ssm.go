package ssm

import (
	"errors"

	"fmt"

	"github.com/sebasrock/filler/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"strings"
)

type ssmTag struct {
	skip       bool
	optional   bool
	defaultVal string
	customPath string
}

type SsmProvider struct {
	ssmClient *ssm.SSM
}

func NewSsmProvider() (*SsmProvider, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	svc := ssm.New(sess)
	return &SsmProvider{
		ssmClient: svc,
	}, nil
}

func (p *SsmProvider) Execute(ctx *config.Context) (string, error) {
	envTag := parseSSM(ctx.TagValue)
	if !envTag.skip {
		if strings.TrimSpace(envTag.customPath) == "" {
			return "", errors.New("ssm-provider: The customNamePath is required")
		}

		parameter, err := p.ssmClient.GetParameter(
			&ssm.GetParameterInput{
				Name:           aws.String(envTag.customPath),
				WithDecryption: aws.Bool(true),
			},
		)
		if err != nil {
			return "", err
		}

		str := aws.StringValue(parameter.Parameter.Value)

		if strings.TrimSpace(envTag.defaultVal) != "" && strings.TrimSpace(str) == "" {
			str = envTag.defaultVal
		}

		if !envTag.optional && strings.TrimSpace(str) == "" {
			return "", fmt.Errorf("ssm-provider: The %s is required", envTag.customPath)
		}
		return str, nil
	}
	return "", nil
}

func parseSSM(s string) *ssmTag {
	var t ssmTag
	tokens := strings.Split(s, ",")
	for _, v := range tokens {
		switch {
		case v == "-":
			t.skip = true
		case v == "optional":
			t.optional = true
		case strings.HasPrefix(v, "default="):
			t.defaultVal = strings.TrimPrefix(v, "default=")
		default:
			t.customPath = v
		}
	}
	return &t
}
