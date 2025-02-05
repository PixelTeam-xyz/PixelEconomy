package main

import (
	dsc "github.com/bwmarrin/discordgo"
	"time"

	//log "msg"
	"strconv"
	"strings"
	. "utils"
)

func onMessage(_ *dsc.Session, msg *dsc.MessageCreate) {
	if msg.Author.Bot {
		return
	}

	var allowed = make([]string, len(cnf.AllowedChannels))
	for _, ch := range cnf.AllowedChannels {
		allowed = append(allowed, strconv.FormatInt(ch, 10))
	}

	//log.Debugf("msg.ChannelID: %s, allowed: %v, contains: %v", msg.ChannelID, allowed, In(allowed, msg.ChannelID))

	userID := msg.Author.ID

	if HasPrefix(msg.Content, cnf.CommandPrefix) {
		if !In(allowed, msg.ChannelID) && (cnf.AllowedChannels != nil && len(cnf.AllowedChannels) > 0) {
			embed := &dsc.MessageEmbed{
				Title:       "Nie poprawny kanał",
				Description: "Nie możesz używać komend ekonomi na tym kanale, użyj odpowiedniego kanału!",
				Color:       colors["red"],
			}

			components := []dsc.MessageComponent{
				dsc.ActionsRow{
					Components: []dsc.MessageComponent{dsc.Button{
						Label:    "Nie ważne",
						Style:    dsc.DangerButton,
						CustomID: "delete_message",
					}},
				},
			}

			msg, err := bot.ChannelMessageSendComplex(msg.ChannelID, &dsc.MessageSend{
				Embed:      embed,
				Components: components,
				Reference: &dsc.MessageReference{
					MessageID: msg.ID,
					ChannelID: msg.ChannelID,
					GuildID:   msg.GuildID,
				},
			})
			Except(err)

			time.AfterFunc(time.Duration(cnf.DisappearanceTimeOfErrorMessages)*time.Second, func() { bot.ChannelMessageDelete(msg.ChannelID, msg.ID) })
			return
		}
		commands := strings.Split(msg.Content, "\n")
		for _, cc := range commands {
			cmd := strings.Split(TrimPrefix(strings.TrimSpace(cc), cnf.CommandPrefix), " ")
			if parts := len(cmd); parts < 1 {
				return
			}

			switch strings.ToLower(strings.TrimSpace(cmd[0])) {
			case "work":
				go workCommand(msg, userID, cmd)
			case "bal", "balance":
				balCommand(msg, userID, cmd)
			case "dep", "deposit":
				s := depCommand(msg, userID, cmd)
				if s {
					go balCommand(msg, userID, []string{"bal"})
				}
			case "with", "withdraw":
				s := withCommand(msg, userID, cmd)
				if s {
					go balCommand(msg, userID, []string{"bal"})
				}
			case "rob", "robbery":
				s := robCommand(msg, msg.Author.ID, cmd)
				if s {
					go balCommand(msg, userID, []string{"bal"})
				}
			case "crime":
				s := crimeCommand(msg, msg.Author.ID, cmd)
				if s {
					go balCommand(msg, userID, []string{"bal"})
				}
			case "top":
				ShowTop(msg.ChannelID)
			case "eco", "economy":
				ecoCommand(msg, userID, cmd)
			case "restart":
				restartCommand(msg, userID, cmd)
			case "buy":
				// TODO
			case "shop":
				// TODO
			case "help", "commands", "cmds":
				go help(msg)
			case "refresh":
				refresh(userID)
			}
		}
	}
}

func help(msg *dsc.MessageCreate) {
	embed := &dsc.MessageEmbed{
		Title:       "📟 Lista komend",
		Description: "Lista komend ekonomii dostępnych na tym serwerze",
		Color:       cnf.MainEmbedColor,
		Fields: []*dsc.MessageEmbedField{
			{Name: "work", Value: "Pozwala ci pracować, aby zarobić pieniądze"},
			{Name: "bal/balance", Value: "Pozwala ci sprawdzić swój stan konta"},
			{Name: "dep/deposit", Value: "Pozwala ci wpłacić pieniądze na konto bankowe"},
			{Name: "with/withdraw", Value: "Pozwala ci wypłacić pieniądze z konta bankowego"},
			{Name: "top", Value: "Pokazuje top listę użytkowników, którzy mają najwięcej pieniędzy na serwerze"},
			{Name: "shop", Value: "Pokazuje sklep, w którym możesz kupić różne przedmioty **(Nie działa: TODO)**"},
			{Name: "buy", Value: "Pozwala ci kupić przedmiot z sklepu **(Nie działa: TODO)**"},
			{Name: "bj/blackjack", Value: "Pozwala ci zagrać w blackjacka **(Nie działa: TODO)**"},
			{Name: "rullette/rl", Value: "Pozwala ci zagrać w ruletkę **(Nie działa: TODO)**"},
			{Name: "help/commands/cmds", Value: "Pokazuje tą listę komend"},
		},
	}
	for _, field := range embed.Fields {
		field.Name = cnf.CommandPrefix + field.Name
	}
	sendEmbed(msg.ChannelID, embed)

	go func() {
		if isAdmin(msg.GuildID, *msg.Author) {
			admEmbed := &dsc.MessageEmbed{
				Title:       "🛡️ Lista komend administracyjnych",
				Description: "Te komendy są dostępne tylko dla administratorów *(czyli też dla ciebie)*",
				Color:       cnf.MainEmbedColor,
				Fields: []*dsc.MessageEmbedField{
					{Name: "eco/economy", Value: "Pozwala kontrolować balansem użytkowników"},
					{Name: "restart", Value: "Usuwa CAŁĄ BAZE DANYCH, w tym stan konta WSZYTSKICH użytkowników. **UŻYWAJ OSTROŻNIE, nie da się tego cofnąć!**"},
				},
			}
			for _, field := range admEmbed.Fields {
				field.Name = cnf.CommandPrefix + field.Name
			}
			sendEmbed(msg.ChannelID, admEmbed)
		}
	}()
}
