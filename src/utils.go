package main

import (
	"fmt"
	dsc "github.com/bwmarrin/discordgo"
	"math/rand"
	"msg"
	"os"
)

func Except(args ...any) {
	var err error
	var exitCode *int
	var f string

	go func() {
		for _, arg := range args {
			switch v := arg.(type) {
			case error:
				err = v
			case int:
				exitCode = &v
			case string:
				f = v
			}
		}

		if f == "" {
			f = "%s"
		}

		if err != nil {
			if exitCode != nil {
				msg.Errorf(f, err.Error())
				os.Exit(*exitCode)
			} else {
				msg.Fatalf(f, err.Error())
			}
		}
	}()
}

func SendMsg(channelID string, f string, x ...any) {
	_, err := bot.ChannelMessageSend(channelID, fmt.Sprintf(f, x...))
	Except(err)
}

func sendEmbed(channelID string, embed *dsc.MessageEmbed) {
	_, err := bot.ChannelMessageSendEmbed(channelID, embed)
	Except(err)
}

func randInt(x, y int) int {
	if x > y {
		x, y = y, x
	}
	return x + rand.Intn(y-x+1)
}

func randFloat(x, y float64) float64 {
	return x + rand.Float64()*(y-x)
}

func sendErrf(channelID string, f string, args ...any) {
	text := fmt.Sprintf(f, args...)
	sendEmbed(channelID, &dsc.MessageEmbed{
		Title:       "🛑 Błąd",
		Description: text,
		Color:       colors["red"],
	})
}

func sendErr(channelID, text string, fields ...*dsc.MessageEmbedField) {
	sendEmbed(channelID, &dsc.MessageEmbed{
		Title:       "🛑 Błąd",
		Description: text,
		Color:       colors["red"],
		Fields:      fields,
	})
}

func sendWarnf(channelID string, f string, args ...any) {
	text := fmt.Sprintf(f, args...)
	sendEmbed(channelID, &dsc.MessageEmbed{
		Title:       "⚠️ Ostrzeżenie",
		Description: text,
		Color:       colors["yellow"],
	})
}

func sendWarn(channelID, text string, fields ...*dsc.MessageEmbedField) {
	sendEmbed(channelID, &dsc.MessageEmbed{
		Title:       "⚠️ Ostrzeżenie",
		Description: text,
		Color:       colors["yellow"],
		Fields:      fields,
	})
}

func sendTipf(channelID string, f string, args ...any) {
	text := fmt.Sprintf(f, args...)
	sendEmbed(channelID, &dsc.MessageEmbed{
		Title:       "💡 Wskazówka",
		Description: text,
		Color:       colors["blue"],
	})
}

func sendTip(channelID, text string, fields ...*dsc.MessageEmbedField) {
	sendEmbed(channelID, &dsc.MessageEmbed{
		Title:       "💡 Wskazówka",
		Description: text,
		Color:       colors["blue"],
		Fields:      fields,
	})
}
