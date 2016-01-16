package wechat

import (
	"git.oschina.net/xuebing1110/queryapi"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var APIQuery queryapi.API
var WechatConf *Configure

type Configure struct {
	Port           int               `yaml:"Port"`
	Token          string            `yaml:"Token"`
	EncodingAESKey string            `yaml:"EncodingAESKey"`
	AppID          string            `yaml:"AppID"`
	OriID          string            `yaml:"OriID"`
	AppSecret      string            `yaml:"AppSecret"`
	UrlBase        string            `yaml:"UrlBase"`
	DBTYPE         string            `yaml:"DBTYPE"`
	DBString       string            `yaml:"DBString"`
	CookiePath     string            `yaml:"CookiePath"`
	LoginMail      string            `yaml:"LoginMail"`
	LoginPwd       string            `yaml:"LoginPwd"`
	APIKey         map[string]string `yaml:"APIKey"`
}

func LoadConfig(filename string) (*Configure, error) {
	WechatConf = &Configure{}

	fi, err := os.Open(filename)
	if err != nil {
		return WechatConf, err
	}
	defer fi.Close()

	f_bytes, read_err := ioutil.ReadAll(fi)
	if read_err != nil {
		return WechatConf, read_err
	}

	yaml_err := yaml.Unmarshal(f_bytes, WechatConf)

	APIQuery.KEY = WechatConf.APIKey
	return WechatConf, yaml_err
}
