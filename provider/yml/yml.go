package yml

import (
	"errors"
	"fmt"
	"github.com/sebasrock/filler/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"reflect"

	"strings"
)

type envTag struct {
	skip       bool
	optional   bool
	defaultVal string
	customName string
}

type YmlProvider struct {
	filesValues map[string]string
}

func NewYmlProvider(file string) *YmlProvider {
	filesValues := make(map[string]interface{})
	data, err := ioutil.ReadFile(fmt.Sprintf("%s", file))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &filesValues)
	if err != nil {
		panic(err)
	}
	return &YmlProvider{filesValues: parseMapString(filesValues)}
}

func (p *YmlProvider) Execute(ctx *config.Context) (string, error) {
	envTag := parseTag(ctx.TagValue)
	if !envTag.skip {
		if strings.TrimSpace(envTag.customName) == "" {
			return "", errors.New("yml-provider: The customName is required")
		}
		str := p.filesValues[envTag.customName]
		if strings.TrimSpace(envTag.defaultVal) != "" && strings.TrimSpace(str) == "" {
			str = envTag.defaultVal
		}
		if !envTag.optional && strings.TrimSpace(str) == "" {
			return "", fmt.Errorf("yml-provider: The %s is required", envTag.customName)
		}
		return str, nil
	}
	return "", nil
}

func parseMapString(kvs map[string]interface{}) map[string]string {
	var giveback = make(map[string]string)
	mapToMap(kvs, ".", giveback)
	return giveback
}

func mapToMap(kvs map[string]interface{}, key string, giveback map[string]string) {
	for k, v := range kvs {
		kt := reflect.ValueOf(k)
		vt := reflect.ValueOf(v)
		nk := key + "." + k
		if kt.Kind() == reflect.String && vt.Kind() == reflect.Map {
			mapToMap(v.(map[string]interface{}), nk, giveback)
		} else if kt.Kind() == reflect.String && vt.Kind() != reflect.Map {
			nk = strings.ReplaceAll(nk, "..", "")
			giveback[nk] = parserToString(v)
		}
	}
}

func parserToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
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
