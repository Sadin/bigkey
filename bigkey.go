package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string
var logger = zap.NewExample()

func main() {

	defer logger.Sync()

	logger.Info(token)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Fatal("failed auth to discord",
			zap.String("err", err.Error()),
		)
	}

	dg.AddHandler(ready)

	dg.AddHandler(guildCreate)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		logger.Fatal("failed to open discord session",
			zap.String("err", err.Error()),
		)
	}

	logger.Info("bigkey is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {

	s.UpdateStatus(0, "beeg keys")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	logger.Info(m.Content,
		zap.String("channelid", m.ChannelID),
		zap.String("userid", m.Message.Author.ID),
	)

}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Airhorn is ready! Type !airhorn while in a voice channel to play a sound.")
			return
		}
	}
}
