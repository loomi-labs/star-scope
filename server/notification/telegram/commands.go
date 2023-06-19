package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/icza/gog"
	"github.com/shifty11/go-logger/log"
)

type MessageCommand string

const (
	MessageCmdStart         MessageCommand = "start"
	MessageCmdSubscriptions MessageCommand = "subscriptions"
	MessageCmdStop          MessageCommand = "stop"
)

func (client TelegramBot) handleCommand(update *tgbotapi.Update) {
	switch MessageCommand(update.Message.Command()) {
	case MessageCmdStart, MessageCmdSubscriptions:
		client.handleStart(update)
	case MessageCmdStop:
		client.handleStop(update)
	}
}

const subscriptionsMsg = `🚀 Star Scope bot started.
%v
🔔 Active subscriptions: %v

<b>How does it work?</b>
- You subscribe to on-chain events
- An on-chain event happens
- A notification is sent to this chat

To stop the bot send the command /stop
`

func (client TelegramBot) handleStart(update *tgbotapi.Update) {
	userId := getUserIdX(update)
	userName := getUserName(update)
	chatId := getChatIdX(update)
	chatName := getChatName(update)
	isGroup := isGroupX(update)

	log.Sugar.Debugf("Send start to %v %v (%v)", gog.If(isGroup, "group", "user"), chatName, chatId)

	text := ""
	ctx := context.Background()
	err := client.UserManager.CreateOrUpdateForTelegramUser(ctx, userId, userName, chatId, chatName, isGroup)
	if err != nil {
		text = "There was an error registering your user. Please try again later."
	} else {
		adminText := ""
		if isGroup {
			adminText += "\n👮‍♂ Bot admins in this chat\n"
			for _, user := range client.UserManager.QueryUsersForTelegramChat(ctx, chatId) {
				adminText += fmt.Sprintf("- @%v\n", user.Name)
			}
		}
		cnt := client.EventListenerManager.QuerySubscriptionsCountForTelegramChat(ctx, chatId)
		text = fmt.Sprintf(subscriptionsMsg, adminText, cnt)
	}

	var buttons [][]Button
	buttons = append(buttons, client.getSubscriptionButtonRow(update))
	replyMarkup := createKeyboard(buttons)

	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = replyMarkup
	msg.ParseMode = "html"
	msg.DisableWebPagePreview = true
	_, err = client.api.Send(msg)
	if err != nil {
		log.Sugar.Errorf("Error while sending /start response for user %v (%v): %v", chatName, chatId, err)
	}
}

func (client TelegramBot) getSubscriptionButtonRow(_ *tgbotapi.Update) []Button {
	var buttonRow []Button
	button := NewButton("🔔 Subscriptions")
	button.LoginURL = &tgbotapi.LoginURL{URL: client.webAppUrl, RequestWriteAccess: true}
	buttonRow = append(buttonRow, button)
	return buttonRow
}

func (client TelegramBot) handleStop(update *tgbotapi.Update) {
	userId := getUserIdX(update)
	chatId := getChatIdX(update)
	chatName := getChatName(update)
	isGroup := isGroupX(update)

	log.Sugar.Debugf("Send stop to %v %v (%v)", gog.If(isGroup, "group", "user"), chatName, chatId)

	text := ""
	err := client.UserManager.DeleteTelegramCommChannel(context.Background(), userId, chatId)
	if err != nil {
		text = "There was an error unregistering your user. Please try again later."
	} else {
		text = "😴 Bot stopped. Send /start to start it again."
	}

	var buttons [][]Button
	buttons = append(buttons, client.getSubscriptionButtonRow(update))
	replyMarkup := createKeyboard(buttons)

	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = replyMarkup
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	_, err = client.api.Send(msg)
	if err != nil {
		log.Sugar.Errorf("Error while sending /stop response for user %v (%v): %v", chatName, chatId, err)
	}
}
