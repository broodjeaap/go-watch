package notifiers

import (
	"fmt"
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

func (discord *DiscordNotifier) Open(configPath string) bool {
	userIDPath := fmt.Sprintf("%s.userID", configPath)
	serverPath := fmt.Sprintf("%s.server", configPath)
	if !viper.IsSet(userIDPath) && !viper.IsSet(serverPath) {
		log.Println("Net either 'serverID' or 'userID' for Discord")
		return false
	}
	bot, err := discordgo.New("Bot " + viper.GetString("notifiers.discord.token"))
	if err != nil {
		log.Println("Could not start Discord notifier:\n", err)
		return false
	}
	if viper.IsSet(userIDPath) {
		discord.UserID = viper.GetString(userIDPath)
		channel, err := bot.UserChannelCreate(discord.UserID)
		if err != nil {
			log.Println("Could not connect to user channel:", discord.UserID, err)
			return false
		}
		discord.UserChannel = channel
		log.Println("Authorized discord bot for:", channel.Recipients)
	}
	if viper.IsSet(serverPath) {
		serverIDPath := fmt.Sprintf("%s.server.ID", configPath)
		serverChannelPath := fmt.Sprintf("%s.server.channel", configPath)
		discord.ServerID = viper.GetString(serverIDPath)
		discord.ChannelID = viper.GetString(serverChannelPath)
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
	debugPath := fmt.Sprintf("%s.debug", configPath)
	discord.Debug = viper.GetBool(debugPath)
	if discord.Debug {
		bot.LogLevel = discordgo.LogDebug
	} else {
		bot.LogLevel = discordgo.LogInformational
	}
	discord.Bot = bot
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
