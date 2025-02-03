package main

import (
	dsc "github.com/bwmarrin/discordgo"
	"os"
)

func onInteraction(_ *dsc.Session, i *dsc.InteractionCreate) {
	if i.Type != dsc.InteractionMessageComponent {
		return
	}

	switch i.MessageComponentData().CustomID {
	case "delete_message":
		err := bot.ChannelMessageDelete(i.ChannelID, i.Message.ID)
		Except(err)

		err = bot.InteractionRespond(i.Interaction, &dsc.InteractionResponse{
			Type: dsc.InteractionResponseUpdateMessage,
			Data: &dsc.InteractionResponseData{
				Components: []dsc.MessageComponent{},
			},
		})
		Except(err)

	case "BUTTON_TO_DELETE_THE_ENTIRE_DATABASE":
		if !isAdmin(i.GuildID, *i.Member.User) {
			bot.InteractionRespond(i.Interaction, &dsc.InteractionResponse{
				Type: dsc.InteractionResponseChannelMessageWithSource,
				Data: &dsc.InteractionResponseData{
					Flags:   dsc.MessageFlagsEphemeral,
					Content: "Nie masz uprawnień do wykonania tej operacji ||(kurwa debil, serio myślałeś że usuniesz całą baze danych bez uprawnień?)||",
				},
			})
			return
		}

		err := os.Remove(cnf.DatabasePath)
		Except(err, DatabaseErrorExit)

		bot.InteractionRespond(i.Interaction, &dsc.InteractionResponse{
			Type: dsc.InteractionResponseChannelMessageWithSource,
			Data: &dsc.InteractionResponseData{
				Content: "Cała baza danych została usunięta! Dla bezpieczeństwa bot musi teraz zostać wyłączony, włącz go ponownie aby używać dalej",
			},
		})

		db.Close()
		os.Exit(0)
	}
}
