package notifiers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

type DiscordNotifier struct {
	Bot           *discordgo.Session
	Token         string
	UserID        string
	UserChannel   *discordgo.Channel
	ServerID      string
	ChannelID     string
	ServerChannel *discordgo.Channel
	Debug         bool
}

func (discord *DiscordNotifier) Open() bool {
	if !viper.IsSet("notifiers.discord.userID") && !viper.IsSet("notifiers.discord.server") {
		log.Println("Net either 'serverID' or 'userID' for Discord")
		return false
	}
	bot, err := discordgo.New("Bot " + viper.GetString("notifiers.discord.token"))
	if err != nil {
		log.Println("Could not start Discord notifier:\n", err)
		return false
	}
	if viper.IsSet("notifiers.discord.userID") {
		discord.UserID = viper.GetString("notifiers.discord.userID")
		channel, err := bot.UserChannelCreate(discord.UserID)
		if err != nil {
			log.Println("Could not connect to user channel:", discord.UserID, err)
			return false
		}
		discord.UserChannel = channel
		log.Println("Authorized discord bot for:", channel.Recipients)
	}
	if viper.IsSet("notifiers.discord.server") {
		discord.ServerID = viper.GetString("notifiers.discord.server.ID")
		discord.ChannelID = viper.GetString("notifiers.discord.server.channel")
		channels, err := bot.GuildChannels(discord.ServerID)
		if err != nil {
			log.Println("Could not connect to server channel:", discord.ServerID, err)
			return false
		}
		foundChannel := false
		for i := range channels {
			channel := channels[i]
			if channel.ID == discord.ChannelID {
				foundChannel = true
				discord.ServerChannel = channel
				break
			}
		}
		if !foundChannel {
			log.Println("Did not find channel with '"+discord.ChannelID+"' in server:", discord.ServerID)
			return false
		}
		log.Println("Authorized discord bot for:", discord.ServerChannel.Name)
	}
	discord.Bot = bot
	bot.Debug = viper.GetBool("notifiers.discord.debug")
	return true
}

func (discord *DiscordNotifier) Message(message string) bool {
	if discord.UserChannel != nil {
		discord.Bot.ChannelMessageSend(discord.UserChannel.ID, message)
	}
	if discord.ServerChannel != nil {
		discord.Bot.ChannelMessageSend(discord.ServerChannel.ID, message)
	}
	return true
}

func (discord *DiscordNotifier) Close() bool {
	discord.Bot.Close()
	return true
}
