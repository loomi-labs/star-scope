package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/commchannel"
	"github.com/loomi-labs/star-scope/ent/state"
	"github.com/loomi-labs/star-scope/ent/user"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/microcosm-cc/bluemonday"
	"github.com/shifty11/go-logger/log"
	"time"
)

const maxMsgLength = 4096

type DeleteTelegramUser struct {
	ChatId        int64
	EventListener *ent.EventListener
}

func (client *TelegramBot) deleteCommChannels(toBeDeleted []DeleteTelegramUser) {
	ctx := context.Background()
	for _, data := range toBeDeleted {
		users, err := data.EventListener.QueryUsers().Where(user.HasCommChannelsWith(commchannel.TelegramChatIDEQ(data.ChatId))).All(ctx)
		if err != nil {
			log.Sugar.Errorf("while querying users with telegram chat id %v: %v", data.ChatId, err)
			continue
		}
		if len(users) == 0 {
			log.Sugar.Errorf("no users with telegram chat id %v", data.ChatId)
			continue
		}
		for _, u := range users {
			err := client.userManager.DeleteTelegramCommChannel(ctx, u.TelegramUserID, data.ChatId, false)
			if err != nil {
				log.Sugar.Errorf("while deleting telegram comm channel for user %v: %v", u.TelegramUsername, err)
				break
			}
		}
	}
}

func (client *TelegramBot) sendNewEvents() {
	notifierState, err := client.eventListenerManager.QueryNotifierState(context.Background(), state.EntityTelegram)
	if err != nil {
		log.Sugar.Panicf("while querying notifier state: %v", err)
	}
	ctx := context.Background()
	toBeDeleted := make([]DeleteTelegramUser, 0)
	endTime := time.Now()
	events, err := client.eventListenerManager.QueryEventsSince(ctx, notifierState.LastEventTime, endTime, notifierState.Entity)
	if err != nil {
		log.Sugar.Panicf("while querying events since %v: %v", notifierState.LastEventTime, err)
	}
	if len(events) > 0 {
		log.Sugar.Infof("sending %v events", len(events))

		p := bluemonday.StripTagsPolicy()

		for _, entEvent := range events {
			pbEvent, err := kafka_internal.EntEventToProto(entEvent, entEvent.Edges.EventListener.Edges.Chain)
			if err != nil {
				log.Sugar.Errorf("while converting event to proto: %v", err)
				continue
			}
			for _, tg := range entEvent.Edges.EventListener.Edges.CommChannels {
				log.Sugar.Debugf("sending event to telegram chat %v (%v)", tg.Name, tg.TelegramChatID)
				var textMsgs []string
				text := fmt.Sprintf("<b>%v</b>\n\n%v", p.Sanitize(pbEvent.Title), p.Sanitize(pbEvent.Description))

				if len(text) <= maxMsgLength {
					textMsgs = append(textMsgs, text)
				} else {
					textMsgs = append(textMsgs, text[:maxMsgLength-4]+"</i>")
					text = text[:len(text)-4] // remove the last 4 characters which are "</i>"
					for _, chunk := range common.Chunks(text[maxMsgLength-4:], maxMsgLength-7) {
						textMsgs = append(textMsgs, fmt.Sprintf("<i>%v</i>", chunk))
					}
				}

				for _, textMsg := range textMsgs {
					msg := tgbotapi.NewMessage(tg.TelegramChatID, textMsg)
					msg.ParseMode = "html"
					msg.DisableWebPagePreview = true

					_, err := client.api.Send(msg)
					if err != nil {
						if client.shouldDeleteUser(err) {
							toBeDeleted = append(toBeDeleted, DeleteTelegramUser{
								EventListener: entEvent.Edges.EventListener,
								ChatId:        tg.TelegramChatID,
							})
						} else {
							log.Sugar.Errorf("Error sending event to telegram chat %v (%v): %v", tg.Name, tg.TelegramChatID, err)
						}
						continue
					}
				}
			}
		}
	}
	if len(toBeDeleted) > 0 {
		log.Sugar.Infof("deleting %v users", len(toBeDeleted))
		client.deleteCommChannels(toBeDeleted)
	}
	_, err = client.eventListenerManager.UpdateNotifierState(context.Background(), notifierState, endTime)
	if err != nil {
		log.Sugar.Errorf("while updating notifier state: %v", err)
	}
}

func (client *TelegramBot) startTelegramEventNotifier() {
	log.Sugar.Info("Start Telegram event notifier")
	for {
		client.sendNewEvents()
		time.Sleep(1 * time.Minute)
	}
}
