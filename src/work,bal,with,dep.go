package main

import (
	"fmt"
	dsc "github.com/bwmarrin/discordgo"
	"time"
	. "utils"
)

func balCommand(msg *dsc.MessageCreate, userID string, cmd []string) {
	var usr string
	switch len(cmd) {
	case 1:
		usr = userID
	case 2:
		if HasPrefix(cmd[1], "<@") && HasSuffix(cmd[1], ">") {
			usr = TrimPrefix(TrimSuffix(cmd[1], ">"), "<@")
		} else {
			sendErrf(msg.ChannelID, "Niepoprawny format komendy bal!")
			sendTip(msg.ChannelID, "Użycie:", &dsc.MessageEmbedField{
				Name:  "1. Sprawdzanie swojego stanu konta",
				Value: fmt.Sprintf("%sbal", cnf.CommandPrefix),
			}, &dsc.MessageEmbedField{
				Name:  "2. Sprawdzanie stanu konta innego użytkownika",
				Value: fmt.Sprintf("%sbal @user", cnf.CommandPrefix),
			})
		}
	}
	sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
		Title: fmt.Sprintf("Stan konta @%s", func() string {
			x, err := bot.GuildMember(msg.GuildID, usr)
			Except(err)
			if err == nil {
				return x.User.Username
			} else {
				return usr
			}
		}()),
		Color: cnf.MainEmbedColor,
		Fields: []*dsc.MessageEmbedField{
			{
				Name:   "👛 Portfel",
				Value:  fmt.Sprintf("%d %s", getBal(usr), cnf.MoneyIcon),
				Inline: true,
			},
			{
				Name:   "🏦 Bank",
				Value:  fmt.Sprintf("%d %s", getBank(usr), cnf.MoneyIcon),
				Inline: true,
			},
			{
				Name:   "💰 Łącznie",
				Value:  fmt.Sprintf("%d %s", getBal(usr)+getBank(usr), cnf.MoneyIcon),
				Inline: true,
			},
		},
	})
}

func workCommand(msg *dsc.MessageCreate, userID string, cmd []string) {
	if can, remaining := canWork(userID); !can {
		nextWorkTime := time.Now().Add(time.Duration(remaining) * time.Second)
		sendErr(msg.ChannelID,
			fmt.Sprintf("Będziesz mógł pracować dopiero <t:%d:R> 🕒", nextWorkTime.Unix()),
		)
		return
	}
	income := int64(randInt(cnf.WorkMax, cnf.WorkMin))
	changeBal(userID, getBal(userID)+income)
	db.Exec("UPDATE users SET lastWork = ? WHERE id = ?", time.Now(), userID)
	sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
		Title:       "💼 Pracowałeś!",
		Description: fmt.Sprintf("Zarobiłeś %d%s!", income, cnf.MoneyIcon),
		Color:       cnf.MainEmbedColor,
	})
	//balCommand(msg, userID, []string{"bal"})
}

