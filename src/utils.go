package main

import (
	"fmt"
	dsc "github.com/bwmarrin/discordgo"
	"info"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
				info.Errorf(f, err.Error())
				os.Exit(*exitCode)
			} else {
				info.Fatalf(f, err.Error())
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

func randBool(chance int8) bool { return rand.Intn(100) < int(chance) }

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

func ToMoneyStr(x any) string {
	var num float64
	switch v := x.(type) {
	case int:
		num = float64(v)
	case int64:
		num = float64(v)
	case float32:
		num = float64(v)
	case float64:
		num = v
	default:
		return "..."
	}

	if num < 1000 {
		formatted := strconv.FormatFloat(num, 'f', -1, 64)
		formatted = strings.Replace(formatted, ".", ",", 1)
		return formatted + cnf.MoneyIcon
	}

	var suffixes = []struct {
		divisor float64
		suffix  string
	}{
		{1e15, "P"},
		{1e12, "T"},
		{1e9, "B"},
		{1000000, "M"},
		{1000, "k"},
	}

	var value float64
	var suffix string
	for _, s := range suffixes {
		if num >= s.divisor {
			value = num / s.divisor
			suffix = s.suffix
			break
		}
	}

	if value < 10 {
		rounded := math.Floor(value*10+0.5) / 10
		formatted := strconv.FormatFloat(rounded, 'f', 1, 64)
		formatted = strings.TrimSuffix(formatted, ".0")
		formatted = strings.Replace(formatted, ".", ",", 1)
		return formatted + suffix + cnf.MoneyIcon
	}

	formatted := strconv.FormatFloat(math.Floor(value+0.5), 'f', 0, 64)
	return formatted + suffix + cnf.MoneyIcon
}

func defaultIncorrect(msg *dsc.MessageCreate, cmd string) func(why string) {
	return func(why string) {
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       fmt.Sprintf("Nie poprawny format komendy %s%s! (%s)", cnf.CommandPrefix, cmd, why),
			Description: fmt.Sprintf("Poprawne użycie: %s%s *(brak argumentów)*", cnf.CommandPrefix, cmd),
			Color:       colors["red"],
		})
	}
}

func HasRole(serverID, userID, roleID string) bool {
	member, err := bot.GuildMember(serverID, userID)
	if err != nil {
		Except(err)
		return false
	}

	for _, role := range member.Roles {
		if role == roleID {
			return true
		}
	}

	return false
}
