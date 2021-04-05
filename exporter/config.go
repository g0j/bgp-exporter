package exporter

import (
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Asn                int
	Rid                net.IP
	ListenAddr         string
	MetricsPort        int
	DeleteOnDisconnect bool
	LogLevel           log.Level
}

func NewConfig(testConfig bool) Config {
	c := Config{}
	var configPath string

	pflag.StringVarP(&configPath, "config", "c", "", "Config file path")
	pflag.Parse()

	viper.SetDefault("asn", "64512")
	viper.SetDefault("router_id", "1.1.1.1")
	viper.SetDefault("listen_address", "0.0.0.0")
	viper.SetDefault("metrics_port", "9179")
	viper.SetDefault("delete_on_disconnect", false)
	viper.SetDefault("log_level", "debug")

	if configPath != "" {
		log.Infof("Parsing config: %s", configPath)
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Error("Unable to read config file: %s", err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("bgpexp")

	switch strings.ToLower(viper.GetString("log_level")) {
	case "panic":
		c.LogLevel = log.PanicLevel
	case "fatal":
		c.LogLevel = log.FatalLevel
	case "error":
		c.LogLevel = log.ErrorLevel
	case "wran":
		c.LogLevel = log.WarnLevel
	case "info":
		c.LogLevel = log.InfoLevel
	case "debug":
		c.LogLevel = log.DebugLevel
	default:
		c.LogLevel = log.InfoLevel
	}

	c.Asn = viper.GetInt("asn")

	c.Rid = net.ParseIP(viper.GetString("router_id"))
	if c.Rid.To4() == nil {
		log.Fatal("Router ID is invalid")
	}

	c.MetricsPort = viper.GetInt("metrics_port")

	if _, err := net.ResolveTCPAddr("tcp", viper.GetString("listen_address")+":"+BGP_TCP_PORT); err != nil {
		log.Fatal("Listen addres is invalid: ", err)
	} else {
		c.ListenAddr = viper.GetString("listen_address")
	}

	c.DeleteOnDisconnect = viper.GetBool("del_on_disconnect")

	return c
}
