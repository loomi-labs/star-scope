package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loomi-labs/star-scope/database"
	"github.com/shifty11/go-logger/log"
	"os"
	"os/signal"
	"syscall"
)

//goland:noinspection GoNameStartsWithPackageName
type DiscordBot struct {
	s                    *discordgo.Session
	userManager          *database.UserManager
	eventListenerManager *database.EventListenerManager
	botToken             string
	clientId             string
	webAppUrl            string
}

func NewDiscordBot(managers *database.DbManagers, botToken string, clientId string, webAppUrl string) *DiscordBot {
	return &DiscordBot{
		userManager:          managers.UserManager,
		eventListenerManager: managers.EventListenerManager,
		botToken:             botToken,
		clientId:             clientId,
		webAppUrl:            webAppUrl,
	}
}

func (dc *DiscordBot) initDiscord() *discordgo.Session {
	log.Sugar.Info("Init discord bot")

	var err error
	s, err := discordgo.New("Bot " + dc.botToken)
	if err != nil {
		log.Sugar.Fatalf("Invalid bot parameters: %v", err)
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		go func() {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				if h, ok := cmdHandlers[i.ApplicationCommandData().Name]; ok {
					h(dc, s, i)
				}
			}
		}()
	})
	return s
}

func (dc *DiscordBot) addCommands() {
	for _, v := range cmds {
		_, err := dc.s.ApplicationCommandCreate(dc.s.State.User.ID, "", v)
		if err != nil {
			log.Sugar.Panic("Cannot create '%v' command: %v", v.Name, err)
		}
	}
}

func (dc *DiscordBot) removeCommands() {
	registeredCommands, err := dc.s.ApplicationCommands(dc.s.State.User.ID, "")
	if err != nil {
		log.Sugar.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := dc.s.ApplicationCommandDelete(dc.s.State.User.ID, "", v.ID)
		if err != nil {
			log.Sugar.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func (dc *DiscordBot) startDiscordSession() *discordgo.Session {
	var err error
	session, err := discordgo.New("Bot " + dc.botToken)
	if err != nil {
		log.Sugar.Fatalf("Invalid bot parameters: %v", err)
	}

	err = session.Open()
	if err != nil {
		log.Sugar.Fatalf("Cannot open the s: %v", err)
	}
	return session
}

func (dc *DiscordBot) closeDiscordSession(session *discordgo.Session) {
	err := session.Close()
	if err != nil {
		log.Sugar.Errorf("Error while closing discord s: %v", err)
	}
}

func (dc *DiscordBot) Start() {
	dc.s = dc.initDiscord()
	log.Sugar.Info("Start discord bot")

	err := dc.s.Open()
	if err != nil {
		log.Sugar.Fatalf("Cannot open the s: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer dc.s.Close()

	dc.removeCommands()
	dc.addCommands()

	go dc.startDiscordEventNotifier()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	<-stop
}
