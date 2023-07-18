package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/icza/gog"
	"github.com/shifty11/go-logger/log"
	"net/url"
)

var (
	startCmd = "start"
	stopCmd  = "stop"
	cmds     = []*discordgo.ApplicationCommand{
		{
			Name:        startCmd,
			Description: "Start the bot and receive notifications",
		},
		{
			Name:        stopCmd,
			Description: "Stop the bot",
		},
	}
	cmdHandlers = map[string]func(dc *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate){
		startCmd: func(dc *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !canInteractWithBot(s, i) {
				sendEmptyResponse(s, i)
				return
			}

			userId := getUserIdX(i)
			userName := getUserName(i)
			channelId := getChannelId(i)
			channelName := getChannelName(s, i)
			isGroup := isGroup(i)

			params := url.Values{}
			params.Add("client_id", dc.clientId)
			params.Add("redirect_uri", dc.webAppUrl)
			params.Add("response_type", "code")
			params.Add("scope", "identify")
			redirectUrl := fmt.Sprintf("https://discord.com/oauth2/authorize?%v", params.Encode())
			text := ""

			ctx := context.Background()
			_, err := dc.userManager.CreateOrUpdateByDiscordUser(ctx, userId, userName, &channelId, &channelName, &isGroup)
			if err != nil {
				log.Sugar.Errorf("Error while creating or updating user: %v", err)
				text = "There was an error registering your user. Please try again later."
			} else {
				//cntSubs := dc.eventListenerManager.QuerySubscriptionsCountForDiscordChannel(ctx, channelId)
				if isGroup {
					adminText := ""
					for _, user := range dc.userManager.QueryUsersForDiscordChannel(ctx, channelId) {
						adminText += fmt.Sprintf("- `%v`\n", user.DiscordUsername)
					}
					text = ":rocket: Star Scope bot started\n\n" +
						fmt.Sprintf(":police_officer: Bot admins in this channel:\n%v\n", adminText) +
						//fmt.Sprintf(":bell: Active subscriptions: %v\n\n", cntSubs) +
						fmt.Sprintf("Go to **[Star Scope](%v)** to change subscriptions for this channel.\n\n", redirectUrl) +
						"**How does it work?**\n" +
						"- You subscribe this channel to a on-chain events\n" +
						"- An on-chain event happens\n" +
						"- A notification is sent to this channel\n\n" +
						"To register another user as admin he has to send the command `/start` to the bot.\n" +
						"To stop the bot send the command `/stop`."
				} else {
					text = ":rocket: Star Scope bot started\n\n" +
						//fmt.Sprintf(":bell: Active subscriptions: %v\n\n", cntSubs) +
						fmt.Sprintf("Go to **[Star Scope](%v)** to change your subscriptions.\n\n", redirectUrl) +
						"**How does it work?**\n" +
						"- You subscribe to on-chain events\n" +
						"- An on-chain event happens\n" +
						"- A notification sent to you\n\n" +
						"To stop the bot send the command `/stop`."
				}
			}

			log.Sugar.Debugf("Send start to %v %v (%v)", gog.If(isGroup, "group", "user"), channelName, channelId)

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: text,
				},
			})
			if err != nil {
				log.Sugar.Errorf("Error while sending subscriptions: %v", err)
			}
		},
		stopCmd: func(dc *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !canInteractWithBot(s, i) {
				sendEmptyResponse(s, i)
				return
			}

			userId := getUserIdX(i)
			channelId := getChannelId(i)
			channelName := getChannelName(s, i)
			isGroup := isGroup(i)

			log.Sugar.Debugf("Send stop to %v %v (%v)", gog.If(isGroup, "group", "user"), channelName, channelId)

			text := ":sleeping: Bot stopped. Send `/start` to start it again."
			err := dc.userManager.DeleteDiscordCommChannel(context.Background(), userId, channelId, true)
			if err != nil {
				log.Sugar.Errorf("Error while deleting user: %v", err)
				text = "There was an error unregistering your user. Please try again later."
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: text,
				},
			})
			if err != nil {
				log.Sugar.Errorf("Error while sending subscriptions: %v", err)
			}
		},
	}
)
