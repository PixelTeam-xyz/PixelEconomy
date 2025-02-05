package main

import (
	"fmt"
	dsc "github.com/bwmarrin/discordgo"
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
