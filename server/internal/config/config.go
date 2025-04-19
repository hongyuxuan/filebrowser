package config

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	RoutePath string
	Database  string
}

func NewConfig(configFile, logLevel, routePath, dbfile *string, port *int) Config {
	var c = Config{}
	if *configFile != "" {
		conf.MustLoad(*configFile, &c)
	} else {
		c.Name = "FileBrowser"
		c.Host = "0.0.0.0"
		c.Port = 9121
		c.Timeout = 60000
		c.Mode = "prd"
		c.Log = logx.LogConf{
			Mode:     "console",
			Encoding: "plain",
			Level:    "info",
			Stat:     true,
		}
		c.Middlewares = rest.MiddlewaresConf{
			Log: true,
		}
		c.RoutePath = "/"
		c.Auth.AccessSecret = "wLnOk8keh/WO5u7lX8H1dB1/mcuHvnI/jfWCMXMPg9o="
		c.Auth.AccessExpire = 86400
	}
	if *logLevel != "" {
		c.Log.Level = *logLevel
	}
	if port != nil && *port != 0 {
		c.Port = *port
	}
	if *routePath != "" {
		c.RoutePath = *routePath
	}
	if *dbfile != "" {
		c.Database = *dbfile
	}
	logx.DisableStat()
	logx.MustSetup(c.Log)
	logx.Infof("Using config: %+v", c)
	return c

}
