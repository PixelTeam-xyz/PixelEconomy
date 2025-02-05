package main

import (
	"fmt"
	dsc "github.com/bwmarrin/discordgo"
	"strings"
)

func robCommand(msg *dsc.MessageCreate, userID string, cmd []string) {
	incorrect := func(why string) {
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       fmt.Sprintf("Nie poprawny format komendy %srob! (%s)", cnf.CommandPrefix, why),
			Description: fmt.Sprintf("Poprawne użycie: %srob <user>", cnf.CommandPrefix),
			Fields: []*dsc.MessageEmbedField{
				{
					Name:  "**<user>**:",
					Value: "Dowolny użytkownik na serwerze w formie pingu",
				},
			},
		})
	}

	if len(cmd) != 2 {
		incorrect("Nie poprawna ilość argumentów")
		return
	}

	target, err := bot.User(strings.TrimSuffix(strings.TrimPrefix(cmd[1], "<@"), ">"))
	if err != nil {
		incorrect("Podany <user> nie jest poprawnym pingiem discord!")
		return
	}

	targetBal := getBal(target.ID)

	if randBool(cnf.RobberySuccesChance) {
		income := targetBal * int64(randInt(cnf.RobSuccessMax, cnf.RobSuccessMin)/100)
		changeBal(target.ID, targetBal-income)
		changeBal(userID, getBal(userID)+income)
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       "💰 Sukces",
			Description: fmt.Sprintf("Udało ci się okraść <@%s> i zdobyłeś %s!", target.ID, ToMoneyStr(income)),
		})
	} else {

	}

}

func crimeCommand(msg *dsc.MessageCreate, userID string, cmd []string) {
	// TODO
}
