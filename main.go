package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/broodjeaap/go-watch/web"
)

func main() {
	writeConfFlag := flag.String("writeConfig", "-", "Path to write template config to")
	var printConfigFlag bool
	flag.BoolVar(&printConfigFlag, "printConfig", false, "Print the template config to stdout")
	flag.Parse()

	if *writeConfFlag != "-" {
		conf, err := web.EMBED_FS.ReadFile("config.tmpl")
		if err != nil {
			log.Fatalln("Could not read config.tmpl")
		}
		os.WriteFile(*writeConfFlag, conf, 0666)
		log.Println("Wrote template config to:", *writeConfFlag)
		return
	}

	if printConfigFlag {
		conf, err := web.EMBED_FS.ReadFile("config.tmpl")
		if err != nil {
			log.Fatalln("Could not read config.tmpl")
		}
		log.SetFlags(0)
		log.Println(string(conf))
		return
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("GOWATCH")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Could not load config file, using env/defaults")
	}

	web := web.NewWeb()
	web.Run()
}
