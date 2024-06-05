package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/libcodebr/didyoumean/config"
	"github.com/libcodebr/didyoumean/infra/mongodb"
	"github.com/spf13/viper"
	"log"
	"os"
)

var configFile string

func init() {
	flag.StringVar(
		&configFile,
		"config",
		"./examples/.config-example.yaml",
		"config file (default is $HOME/examples/.config-example.yaml)",
	)

	flag.Parse()
}

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig(viper.GetViper(), configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client, err := mongodb.NewMongoDB(ctx, cfg.Mongo)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close(ctx)
	log.Println("connected to mongo db")
}
