package main

import (
	"fmt"
	dsc "github.com/bwmarrin/discordgo"
	"strings"
)

func ShowTop(ch string) {
	top, err := getTop(cnf.NumberOfUsersInTopList)
	Except("Failed to retrieve top list from database (%s)", err)
	topEmbed := dsc.MessageEmbed{
		Title: fmt.Sprintf("Top %d użytkowników", func() int {
			if cnf.NumberOfUsersInTopList < len(top) {
				return cnf.NumberOfUsersInTopList
			}
			return len(top)
		}()),
		Description: "Ci użytkownicy mają najwięcej pieniędzy na serwerze",
		Color:       cnf.MainEmbedColor,
		Fields:      make([]*dsc.MessageEmbedField, 0),
	}
	for i, userID := range top {
		user, err := bot.User(userID)
		Except("Failed to retrieve user from discord API (%s)", err)
		topEmbed.Fields = append(topEmbed.Fields, &dsc.MessageEmbedField{
			Name:   fmt.Sprintf("%d. @%s", i+1, user.Username),
			Value:  fmt.Sprintf("👛: %d%s,\n 🏦: %d%s", getBal(userID), cnf.MoneyIcon, getBank(userID), cnf.MoneyIcon),
			Inline: false,
		})
	}

	sendEmbed(ch, &topEmbed)
}

func shopCommand(msg *dsc.MessageCreate, _ string, cmd []string) {
	incorrect := defaultIncorrect(msg, "shop")

	if len(cmd) != 1 {
		incorrect("Nie poprawna liczba argumentów")
	}

	embed := &dsc.MessageEmbed{
		Title:       "🛒 Sklep",
		Description: fmt.Sprintf("Te przedmioty są dostępne w sklepie tego servera, jeśli chcesz je kupić użyj %sbuy <itemName>", cnf.CommandPrefix),
		Color:       cnf.MainEmbedColor,
	}

	for _, item := range items {
		embed.Fields = append(embed.Fields, &dsc.MessageEmbedField{
			Name:  item.Name,
			Value: fmt.Sprintf("%s\n**Cena:** %s\n**Mnożnik:** %.1f", item.Description, ToMoneyStr(item.Price), item.Multiplier),
		})
	}

	sendEmbed(msg.ChannelID, embed)
}

func buyCommand(msg *dsc.MessageCreate, userID string, cmd []string) (success bool) {
	incorrect := func(why string) {
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       fmt.Sprintf("Nie poprawne użycie komendy %sbuy! (%s)", cnf.CommandPrefix, why),
			Description: fmt.Sprintf("Poprawne użycie: %sbuy <item>", cnf.CommandPrefix),
			Fields: []*dsc.MessageEmbedField{
				{
					Name:  "<item>",
					Value: fmt.Sprintf("Dowolny przedmiot ze sklepu. Jeśli chcesz zobaczyć wszytskie dostępne w sklepie przedmioty użyj %sshop", cnf.CommandPrefix),
				},
			},
		})
	}
	if len(cmd) != 2 {
		incorrect("Nie poprawna ilość argumentów")
		return
	}

	itemName := cmd[1]
	userBal := getBal(userID)
	exists, item := func() (bool, Item) {
		for _, item_ := range items {
			if strings.ToLower(item_.Name) == strings.ToLower(itemName) {
				return true, item_
			}
		}
		return false, DefaultItem
	}()

	if !exists {
		incorrect("Podany <item> nie istnieje w sklepie")
		return
	}

	if !(userBal > float64(item.Price)) {
		sendErrf(msg.ChannelID, "Nie masz wystarczająco pieniędzy by kupić %s! (%s > %s)", item.Name, ToMoneyStr(item.Price), ToMoneyStr(userBal))
		return
	}

	if HasRole(msg.GuildID, msg.Author.ID, item.RoleID) {
		sendTipf(msg.ChannelID, "Już posiadasz role %s!", item.Name)
		return
	}

	err := bot.GuildMemberRoleAdd(msg.GuildID, userID, item.RoleID)
	Except(err)
	changeBal(userID, userBal-float64(item.Price))
	sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
		Title:       "🛒 Zakup udany",
		Description: fmt.Sprintf("Pomyślnie zakupiono %s za %s! **Aktualny stan konta:**", item.Name, ToMoneyStr(item.Price)),
		Color:       colors["green"],
	})

	return true
}
