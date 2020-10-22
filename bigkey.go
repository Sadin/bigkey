package main

import (
	"flag"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string

func main() {

	logger := zap.NewExample()
	defer logger.Sync()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Info("failed auth to discord",
			zap.String("err", err.Error()),
		)
	}

	dg.Close()
}
