package helper

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/jessevdk/go-flags"
	"os"
)

type FlagOptions struct {
	ClientId string `short:"c" long:"client-id" description:"mqtt client ID" default:"go_mqtt_client"`
	Host     string `short:"h" long:"host" description:"mqtt broker host" default:"localhost"`
	Port     int    `short:"p" long:"port" description:"mqtt broker port" default:"1883"`
	Username string `short:"U" long:"username" description:"mqtt broker username"`
	Password string `short:"P" long:"password" description:"mqtt broker password"`
}

func GetBaseClientOpt() *mqtt.ClientOptions {
	var fOpts FlagOptions
	_, err := flags.ParseArgs(&fOpts, os.Args)
	if err != nil {
		panic(err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", fOpts.Host, fOpts.Port))

	opts.SetClientID(fOpts.ClientId + "_" + uuid.New().String())

	if fOpts.Username != "" {
		opts.SetUsername(fOpts.Username)
	}

	if fOpts.Password != "" {
		opts.SetPassword(fOpts.Password)
	}
	return opts
}
