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

	dg.AddHandler(channelUpdate)

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates)

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

	logger.Info("updating status")
	s.UpdateStatus(0, "beeg keys")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	logger.Info("message sent",
		zap.String("content", m.Content),
		zap.String("channelid", m.ChannelID),
		zap.String("userid", m.Message.Author.ID),
		zap.String("guildId", m.GuildID),
	)

}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	logger.Info("guild joined",
		zap.String("guildName", event.Guild.Name),
		zap.String("guildId", event.Guild.ID),
		zap.Int("memberCount", event.Guild.MemberCount),
		zap.String("region", event.Guild.Region),
	)

	/*
		for _, channel := range event.Guild.Channels {
			logger.Info(channel.Name)
			logger.Info(channel.)
		}
	*/
}

func channelUpdate(s *discordgo.Session, event *discordgo.ChannelUpdate) {
	logger.Info("channel updated",
		zap.String("channelId", event.Channel.ID),
		zap.String("guildId", event.Channel.GuildID),
	)
}
