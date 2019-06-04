package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/widuu/goini"
	"strconv"
)

// 配置项
type Options struct {
	Port           int
	TlsCertFile    string
	TlsKeyFile     string
	SidecarCfgFile string
}

// 解析配置文件
func ParseConf(c *cli.Context) (*Options, error) {
	options := &Options{}
	var err error
	var conf *goini.Config
	if c.IsSet("configure") || c.IsSet("C") {
		if c.IsSet("configure") {
			conf = goini.SetConfig(c.String("configure"))
		} else {
			conf = goini.SetConfig(c.String("C"))
		}
		port := conf.GetValue("common", "port")
		if options.Port, err = strconv.Atoi(port); nil != err {
			log.Errorf("%v", err)
			options.Port = 443
		}
		options.TlsCertFile = conf.GetValue("common", "tlsCertFile")
		if "" == options.TlsCertFile {
			options.TlsCertFile = "/etc/webhook/certs/cert.pem"
		}
		options.TlsKeyFile = conf.GetValue("common", "tlsKeyFile")
		if "" == options.TlsKeyFile {
			options.TlsKeyFile = "/etc/webhook/certs/key.pem"
		}
		options.SidecarCfgFile = conf.GetValue("common", "sidecarCfgFile")
		if "" == options.SidecarCfgFile {
			options.SidecarCfgFile = "/etc/webhook/config/sidecarconfig.yaml"
		}
		return options, nil
	} else {
		options.Port = 443
		options.TlsCertFile = "/etc/webhook/certs/cert.pem"
		options.TlsKeyFile = "/etc/webhook/certs/key.pem"
		options.SidecarCfgFile = "/etc/webhook/config/sidecarconfig.yaml"
		return options, nil
	}
}
