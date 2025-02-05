package main

import (
	dsc "github.com/bwmarrin/discordgo"
	info "info"
	"strconv"
	"time"
)

var tasks = []func(){
	func() {
		if cnf.TopCh == -1 {
			return
		}
		for {
			// ------ //

			if cnf.DelPrevTopListOnNewSend {
				channel, err := bot.Channel(strconv.FormatInt(cnf.TopCh, 10))
				Except("Rrror while getting channel: %s", err)

				if channel.Type != dsc.ChannelTypeGuildText {
					info.Fatal("channel 'TopMessagesChannelID' (from config) is not a text channel!")
				}

				// Check bot permissions
				perms, err := bot.UserChannelPermissions(bot.State.User.ID, strconv.FormatInt(cnf.TopCh, 10))
				Except("error while checking permissions: %s", err)

				if perms&dsc.PermissionManageMessages == 0 {
					info.Fatal("bot does not have permission to manage messages!")
				}

				var allMessages []*dsc.Message
				lastMessageID := ""
				for {
					messages, err := bot.ChannelMessages(strconv.FormatInt(cnf.TopCh, 10), 100, lastMessageID, "", "")
					Except(err)
					if len(messages) == 0 {
						break
					}
					allMessages = append(allMessages, messages...)
					lastMessageID = messages[len(messages)-1].ID
				}

				for _, msg := range allMessages {
					err := bot.ChannelMessageDelete(strconv.FormatInt(cnf.TopCh, 10), msg.ID)
					Except(err)
				}
			}
			// ----- //

			ShowTop(strconv.FormatInt(cnf.TopCh, 10))
			time.Sleep(time.Duration(cnf.DelayInSendingTopList) * time.Second)
		}
	},
}

func onReady(_ *dsc.Session, event *dsc.Ready) {
	info.Infof("The bot successfully connected to the discord API as a %s (%s). SessionID: %s", botUser.Username, botUser.ID, event.SessionID)
	for _, event := range tasks {
		go event()
	}
}

func onConnectionResumed(_ *dsc.Session, event *dsc.Resumed) {
	info.Info("The bot has regained its connection to Discord!")
}
