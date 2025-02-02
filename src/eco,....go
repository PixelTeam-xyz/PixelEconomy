package main

import (
	"fmt"
	dsc "github.com/bwmarrin/discordgo"
	log "msg"
	"strconv"
	"strings"
)

func isAdmin(user dsc.User) (isAdm bool) {
	for _, adminID := range cnf.AdminUsersIDs {
		if user.ID == strconv.FormatInt(adminID, 10) {
			isAdm = true
			break
		}
	}

	if !isAdm {
		usr, err := bot.GuildMember(cnf.ServerID, user.ID)
		if err != nil {
			Except(err)
			return false
		}

		for _, userRoleID := range usr.Roles {
			for _, adminRoleID := range cnf.AdminRolesIDs {
				if userRoleID == strconv.FormatInt(adminRoleID, 10) {
					isAdm = true
					break
				}
			}
			if isAdm {
				break
			}
		}
	}

	return
}

func ecoCommand(msg *dsc.MessageCreate, userID string, cmd []string) {
	usr, err := bot.GuildMember(cnf.ServerID, msg.Author.ID)
	if err != nil {
		log.Debugf("AdminUsers: %v, AdminRoles: %v, UserRoles: %v, isAdmin?: %v", cnf.AdminUsersIDs, cnf.AdminRolesIDs, usr.Roles, isAdmin(*msg.Author))
	}

	if !isAdmin(*msg.Author) {
		_, err := bot.ChannelMessageSendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       "🛑 Brak uprawnień!",
			Description: "Aby użyć tej komendy musisz mieć uprawnienia administratora!",
			Color:       colors["red"],
		})
		Except(err)
		return
	}

	incorrect := func(why string) {
		sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
			Title:       fmt.Sprintf("Nie poprawny format komendy %seco! (%s)", cnf.CommandPrefix, why),
			Description: fmt.Sprintf("Poprawne użycie: `%seco <user> <operation> <amount> (opcjonalnie) <target>`", cnf.CommandPrefix),
			Fields: []*dsc.MessageEmbedField{
				{
					Name:  "<user>:",
					Value: "Dowolny użytkownik tego servera jako ping",
				},
				{
					Name: "<operation>: ",
					Value: "***dozwolone operacje:***\n" + // .eco <@1274610053843783768> set 100
						"    **+=/add**: Dodaje <amount> do <target> użytkownika\n" +
						"    **-=/deduct**: Odejmuje <amount> od <target> użytkownika\n" +
						"    **=/set**: Ustawia <target> użytkownika na <amout>\n",
				},
				{
					Name:  "<amount>: ",
					Value: "Dowolna liczba",
				},
				{
					Name:  "<target>: ",
					Value: "Ta wartość jest opcjonalna, dozwolone wartości: **portfel**, **bank**. Domyślnie portfel",
				},
			},
		})
	}

	// true == target is bank
	var target bool

	switch len(cmd) {
	case 4:
		target = false
	case 5:
		if cmd[4] == "bank" {
			target = true
		} else if cmd[4] == "portfel" {
			target = false
		} else {
			incorrect("Nie poprawny argument <target> (Dopuszczane wartości: **portfel** i **bank**)")
			return
		}
	default:
		incorrect("Nie poprawna ilośc argumentów")
		return
	}

	user, err := bot.User(strings.TrimSuffix(strings.TrimPrefix(cmd[1], "<@"), ">"))
	if err != nil {
		incorrect("podany <user> nie jest poprawnym pingiem discord/user nie istnieje")
		return
	}

	amount, err := strconv.ParseInt(cmd[3], 10, 64)
	if err != nil {
		incorrect("podany <amount> nie jest poprawną liczbą!")
		return
	}

	var opStr string

	switch cmd[2] {
	case "+=", "add":
		opStr = "dodano %d do %s %s"
		if target {
			changeBank(user.ID, getBank(user.ID)+amount)
		} else {
			changeBal(user.ID, getBal(user.ID)+amount)
		}
	case "-=", "deduct":
		opStr = "odjęto %d od %s %s"
		if target {
			changeBank(user.ID, getBank(user.ID)-amount)
		} else {
			changeBal(user.ID, getBal(user.ID)-amount)
		}
	case "=", "set":
		opStr = "ustawiono wartość %d do %s %s"
		if target {
			changeBank(user.ID, amount)
		} else {
			changeBal(user.ID, amount)
		}
	default:
		incorrect("Nie poprawna operacja!")
		return
	}

	sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
		Title: "💾 Sukces",
		Description: fmt.Sprintf("Pomyślnie "+opStr, amount, func() string {
			if target {
				return "banku"
			} else {
				return "portfela"
			}
		}(), user.Username),
		Fields: []*dsc.MessageEmbedField{
			{
				Name: fmt.Sprintf("Aktualny stan konta %s:", user.Username),
			},
		},
	})
	balCommand(msg, userID, []string{"bal", fmt.Sprintf("<@%s>", user.ID)})
}
