package main

import (
	"flag"
	"strings"

	"github.com/17TheWord/zerobot-plugin-mcqq/mcqq"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
	})
}

func main() {
	configPath := flag.String("config", "config.yml", "YAML 配置文件路径")
	flag.Parse()

	appConfig, err := loadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	level, err := log.ParseLevel(strings.ToLower(appConfig.Log.Level))
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)

	drivers, err := buildDrivers(appConfig.Zero.Drivers)
	if err != nil {
		log.Fatal(err)
	}

	mcqq.PluginConfig = appConfig.Mcqq
	mcqq.InitPlugin()
	zero.RunAndBlock(&zero.Config{
		NickName:      appConfig.Zero.NickName,
		CommandPrefix: appConfig.Zero.CommandPrefix,
		SuperUsers:    appConfig.Zero.SuperUsers,
		Driver:        drivers,
	}, nil)
}
