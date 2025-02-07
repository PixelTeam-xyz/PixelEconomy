package main

import (
	"fmt"
	dsc "github.com/bwmarrin/discordgo"
	"info"
	"strings"
	"time"
)

func robCommand(msg *dsc.MessageCreate, userID string, cmd []string) (success bool) {
	incorrect := func(why string) {
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       fmt.Sprintf("Nie poprawny format komendy %srob! (%s)", cnf.CommandPrefix, why),
			Description: fmt.Sprintf("Poprawne użycie: %srob <user>", cnf.CommandPrefix),
			Fields: []*dsc.MessageEmbedField{
				{
					Name:  "**<user>**:",
					Value: "Dowolny użytkownik na serwerze *(oprócz ciebie samego!)* w formie pingu",
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
		incorrect("Podany <user> nie jest poprawnym pingiem discord lub nie istnieje na serverze!")
		return
	}

	if target.ID == userID {
		incorrect("Nie możesz okraść samego siebie!")
		return
	}

	targetBal := getBal(target.ID)

	if targetBal < 10 {
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       "🛑 Błąd",
			Description: fmt.Sprintf("Podany użytkownik nie ma przy sobie pieniędzy, spróbuj ponownie puźniej"),
			Color:       colors["red"],
		})
		return
	}

	if can, remaining := canRob(userID); !can {
		nextRobTime := time.Now().Add(time.Duration(remaining) * time.Second)
		sendErr(msg.ChannelID,
			fmt.Sprintf("Będziesz mógł kraść dopiero <t:%d:R> 🕒", nextRobTime.Unix()),
		)
		return
	}

	refresh(userID, "rob")
	if income := targetBal * float64(randInt(cnf.RobSuccessEarningsMax, cnf.RobSuccessEarningsMin)/100); randBool(cnf.RobberySuccessChance) && income > 0 {
		info.Debug(randBool(cnf.RobberySuccessChance) && income > 0)
		if targetBal-income < 0 {
			income = targetBal
		}
		changeBal(target.ID, targetBal-income)
		changeBal(userID, getBal(userID)+income)
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       "💰 Sukces",
			Description: fmt.Sprintf("Udało ci się okraść <@%s> i zdobyłeś %s!\n**Aktualny stan konta:**", target.ID, ToMoneyStr(income)),
			Color:       colors["green"],
		})
	} else {
		loss := targetBal * float64(randInt(cnf.RobFailureLossMin, cnf.RobFailureLossMax)) / 100
		changeBal(userID, getBal(userID)-loss)
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       "🚨 Porażka",
			Description: fmt.Sprintf("Przyłapali cię podczas kradzieży od <@%s>, i straciłeś %s\n**Aktualny stan konta:**", target.ID, ToMoneyStr(loss)),
			Color:       colors["red"],
		})
	}
	return true
}

func crimeCommand(msg *dsc.MessageCreate, userID string, cmd []string) (success bool) {
	incorrect := func(why string) {
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       fmt.Sprintf("Nie poprawny format komendy %scrime! (%s)", cnf.CommandPrefix, why),
			Description: fmt.Sprintf("Poprawne użycie: %scrime *(brak argumentów)*", cnf.CommandPrefix),
		})
	}

	if len(cmd) != 1 {
		incorrect("Nie poprawna ilość argumentów")
	}

	if can, remaining := canCrime(userID); !can {
		nextRobTime := time.Now().Add(time.Duration(remaining) * time.Second)
		sendErr(msg.ChannelID,
			fmt.Sprintf("Będziesz mógł pracować nielegalnie dopiero <t:%d:R> 🕒", nextRobTime.Unix()),
		)
		return
	}

	refresh(userID, "crime")
	if randBool(cnf.CrimeSuccessChance) {
		income := float64(randInt(cnf.CrimeSuccessEarningsMin, cnf.CrimeSuccessEarningsMax))
		changeBal(userID, getBal(userID)+income)

		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       "💰 Sukces",
			Description: fmt.Sprintf("Udało ci sie nielegalnie zarobić %s!\n**Aktualny stan konta:**", ToMoneyStr(income)),
			Color:       colors["green"],
		})
	} else {
		loss := float64(randInt(cnf.CrimeFailureLossMin, cnf.CrimeFailureLossMax))
		changeBal(userID, getBal(userID)-loss)
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       "🚨 Porażka",
			Description: fmt.Sprintf("Przyłapali cię na przestępstwie, i straciłeś %s\n**Aktualny stan konta:**", ToMoneyStr(loss)),
			Color:       colors["red"],
		})
	}
	return true
}
