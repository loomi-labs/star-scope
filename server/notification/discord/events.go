package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/commchannel"
	"github.com/loomi-labs/star-scope/ent/state"
	"github.com/loomi-labs/star-scope/ent/user"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/microcosm-cc/bluemonday"
	"github.com/shifty11/go-logger/log"
	"strconv"
	"time"
)

const maxMsgLength = 2000

type DeleteDiscordUser struct {
	ChannelId     int64
	EventListener *ent.EventListener
}

func (dc *DiscordBot) deleteCommChannels(toBeDeleted []DeleteDiscordUser) {
	ctx := context.Background()
	for _, data := range toBeDeleted {
		users, err := data.EventListener.QueryUsers().Where(user.HasCommChannelsWith(commchannel.DiscordChannelID(data.ChannelId))).All(ctx)
		if err != nil {
			log.Sugar.Errorf("while querying users with discord channel id %v: %v", data.ChannelId, err)
			continue
		}
		if len(users) == 0 {
			log.Sugar.Errorf("no users with discord channel id %v", data.ChannelId)
			continue
		}
		for _, u := range users {
			err := dc.userManager.DeleteDiscordCommChannel(ctx, u.DiscordUserID, data.ChannelId, false)
			if err != nil {
				log.Sugar.Errorf("while deleting discord comm channel for user %v: %v", u.DiscordUsername, err)
				break
			}
		}
	}
}

func (dc *DiscordBot) sendNewEvents() {
	notifierState, err := dc.eventListenerManager.QueryNotifierState(context.Background(), state.EntityDiscord)
	if err != nil {
		log.Sugar.Panicf("while querying notifier state: %v", err)
	}
	ctx := context.Background()
	toBeDeleted := make([]DeleteDiscordUser, 0)
	endTime := time.Now()
	events, err := dc.eventListenerManager.QueryEventsSince(ctx, notifierState.LastEventTime, endTime, notifierState.Entity)
	if err != nil {
		log.Sugar.Panicf("while querying events since %v: %v", notifierState.LastEventTime, err)
	}
	if len(events) > 0 {
		log.Sugar.Infof("sending %v events", len(events))
		session := dc.startDiscordSession()
		defer dc.closeDiscordSession(session)

		p := bluemonday.StripTagsPolicy()

		for _, entEvent := range events {
			pbEvent, err := kafka_internal.EntEventToProto(entEvent, entEvent.Edges.EventListener.Edges.Chain)
			if err != nil {
				log.Sugar.Errorf("while converting event to proto: %v", err)
				continue
			}
			for _, commChannel := range entEvent.Edges.EventListener.Edges.CommChannels {
				var textMsgs []string
				text := fmt.Sprintf("**%v**\n\n%v", p.Sanitize(pbEvent.Title), sanitizeUrls(p.Sanitize(pbEvent.Description)))
				if len(text) <= maxMsgLength {
					textMsgs = append(textMsgs, text)
				} else {
					textMsgs = append(textMsgs, text[:maxMsgLength-1]+"*")
					text = text[:len(text)-1] // remove the last character which is a *
					for _, chunk := range common.Chunks(text[maxMsgLength-1:], maxMsgLength-2) {
						textMsgs = append(textMsgs, fmt.Sprintf("*%v*", chunk))
					}
				}
				for _, textMsg := range textMsgs {
					var _, err = session.ChannelMessageSendComplex(strconv.FormatInt(commChannel.DiscordChannelID, 10),
						&discordgo.MessageSend{
							Content: textMsg,
						})
					if err != nil {
						if shouldDeleteUser(err) {
							toBeDeleted = append(toBeDeleted, DeleteDiscordUser{
								EventListener: entEvent.Edges.EventListener,
								ChannelId:     commChannel.DiscordChannelID,
							})
						} else {
							log.Sugar.Errorf("Error sending event to discord channel %v (%v): %v", commChannel.Name, commChannel.DiscordChannelID, err)
						}
						continue
					}
				}
			}
		}
	}
	if len(toBeDeleted) > 0 {
		log.Sugar.Infof("deleting %v users", len(toBeDeleted))
		dc.deleteCommChannels(toBeDeleted)
	}
	_, err = dc.eventListenerManager.UpdateNotifierState(context.Background(), notifierState, endTime)
	if err != nil {
		log.Sugar.Errorf("while updating notifier state: %v", err)
	}
}

func (dc *DiscordBot) startDiscordEventNotifier() {
	log.Sugar.Info("Start Discord event notifier")
	//cr := cron.New()
	//_, err := cr.AddFunc("* * * * *", func() { client.sendNewEvents() }) // every minute
	//if err != nil {
	//	log.Sugar.Errorf("while executing 'addOrUpdateChains' via cron: %v", err)
	//}
	//cr.Start()
	for {
		dc.sendNewEvents()
		time.Sleep(1 * time.Minute)
	}
}
