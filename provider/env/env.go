package env

import (
	"errors"
	"fmt"

	"github.com/sebasrock/filler/config"

	"os"
	"strings"
)

type envTag struct {
	skip       bool
	optional   bool
	defaultVal string
	customName string
}

type EnvProvider struct {
}

func NewEnvProvider() *EnvProvider {
	return &EnvProvider{}
}

func (p *EnvProvider) Execute(ctx *config.Context) (string, error) {
	envTag := parseTag(ctx.TagValue)
	if !envTag.skip {
		if strings.TrimSpace(envTag.customName) == "" {
			return "", errors.New("env-provider: The customName is required")
		}
		str := os.Getenv(envTag.customName)
		if strings.TrimSpace(envTag.defaultVal) != "" && strings.TrimSpace(str) == "" {
			str = envTag.defaultVal
		}
		if !envTag.optional && strings.TrimSpace(str) == "" {
			return "", fmt.Errorf("env-provider: The %s is required", envTag.customName)
		}
		return str, nil
	}
	return "", nil
}

func parseTag(s string) *envTag {
	var t envTag
	tokens := strings.Split(s, ",")
	for _, v := range tokens {
		switch {
		case v == "-":
			t.skip = true
		case v == "optional": // nolint
			t.optional = true
		case strings.HasPrefix(v, "default="):
			t.defaultVal = strings.TrimPrefix(v, "default=")
		default:
			t.customName = v
		}
	}
	return &t
}
