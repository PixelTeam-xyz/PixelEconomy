package main

import (
    "fmt"
    dsc "github.com/bwmarrin/discordgo"
    "info"
    "strconv"
    "time"
    . "utils"
)

func balCommand(msg *dsc.MessageCreate, userID string, cmd []string) {
    var usr string
    incorrect := func() {
        sendErr(msg.ChannelID, "Niepoprawny format komendy "+cnf.CommandPrefix+"bal!")
        sendTip(msg.ChannelID, "Użycie:", &dsc.MessageEmbedField{
            Name:  "1. Sprawdzanie swojego stanu konta",
            Value: fmt.Sprintf("%sbal", cnf.CommandPrefix),
        }, &dsc.MessageEmbedField{
            Name:  "2. Sprawdzanie stanu konta innego użytkownika",
            Value: fmt.Sprintf("%sbal @user", cnf.CommandPrefix),
        })
    }
    switch len(cmd) {
    case 1:
        usr = userID
    case 2:
        if HasPrefix(cmd[1], "<@") && HasSuffix(cmd[1], ">") {
            usr = TrimPrefix(TrimSuffix(cmd[1], ">"), "<@")
        } else {
            incorrect()
            return
        }
    default:
        incorrect()
        return
    }

    user, err := bot.User(usr)
    Except(err)

    sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
        Author: &dsc.MessageEmbedAuthor{
            Name:    fmt.Sprintf("%s    - Stan konta", user.Username),
            IconURL: user.AvatarURL("512"),
        },
        //Title: fmt.Sprintf("Stan konta"),
        Color: cnf.MainEmbedColor,
        Fields: []*dsc.MessageEmbedField{
            {
                Name:   "👛 Portfel:  ",
                Value:  fmt.Sprintf("%s", ToMoneyStr(getBal(usr))),
                Inline: true,
            },
            {
                Name:   "🏦 Bank:  ",
                Value:  fmt.Sprintf("%s", ToMoneyStr(getBank(usr))),
                Inline: true,
            },
            {
                Name:   "💰 Łącznie:  ",
                Value:  fmt.Sprintf("%s", ToMoneyStr(getBal(usr)+getBank(usr))),
                Inline: true,
            },
        },
    })
}

func workCommand(msg *dsc.MessageCreate, userID string, cmd []string) {
    if can, remaining := canWork(userID); !can {
        info.Debugf("can: %v, remaining: %v...", can, remaining)
        nextWorkTime := time.Now().Add(time.Duration(remaining) * time.Second)
        sendErr(msg.ChannelID,
            fmt.Sprintf("Będziesz mógł pracować dopiero <t:%d:R> 🕒", nextWorkTime.Unix()),
        )
        return
    }
    income := int64(randInt(cnf.WorkMax, cnf.WorkMin))
    changeBal(userID, getBal(userID)+income)
    refresh(userID)
    sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
        Title:       "💼 Pracowałeś!",
        Description: fmt.Sprintf("Zarobiłeś %s!", ToMoneyStr(income)),
        Color:       cnf.MainEmbedColor,
    })
}

func depCommand(msg *dsc.MessageCreate, userID string, cmd []string) (success bool) {
    var toDep int64
    switch len(cmd) {
    case 1:
        toDep = getBal(userID)
    case 2:
        if cmd[1] == "all" {
            toDep = getBal(userID)
        } else if x, err := strconv.Atoi(cmd[1]); err != nil {
            sendErrf(msg.ChannelID, "Niepoprawna kwota! Podaj poprawną liczbe po poleceniu %sdep", cnf.CommandPrefix)
            return
        } else {
            toDep = int64(x)
        }
    }

    if toDep < 0 {
        sendErr(msg.ChannelID, "Nie możesz wpłacić ujemnej kwoty!")
        return
    }

    if getBal(userID) < toDep {
        sendErrf(msg.ChannelID, "Nie masz tyle %s w portfelu! (%d > %d)", cnf.MoneyIcon, toDep, getBal(userID))
        return
    }

    changeBal(userID, getBal(userID)-toDep)
    changeBank(userID, getBank(userID)+toDep)

    sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
        Title:       "💼 Wpłata",
        Description: fmt.Sprintf("Pomyślnie wpłacono %s na konto bankowe!", ToMoneyStr(toDep)),
        Color:       colors["green"],
    })
    return true
}

func withCommand(msg *dsc.MessageCreate, userID string, cmd []string) (success bool) {
    var toWith int64
    switch len(cmd) {
    case 1:
        toWith = getBank(userID)
    case 2:
        if cmd[1] == "all" {
            toWith = getBank(userID)
        } else if x, err := strconv.Atoi(cmd[1]); err != nil {
            sendErrf(msg.ChannelID, "Niepoprawna kwota! Podaj poprawną liczbę po poleceniu %swith", cnf.CommandPrefix)
            return
        } else {
            toWith = int64(x)
        }
    }

    if toWith < 0 {
        sendErr(msg.ChannelID, "Nie możesz wypłacić ujemnej kwoty!")
        return
    }

    if getBank(userID) < toWith {
        sendErrf(msg.ChannelID, "Nie masz tyle %s na koncie bankowym! (%d > %d)", cnf.MoneyIcon, toWith, getBank(userID))
        return
    }

    changeBank(userID, getBank(userID)-toWith)
    changeBal(userID, getBal(userID)+toWith)

    sendEmbed(msg.ChannelID, &dsc.MessageEmbed{
        Title:       "🏦 Wypłata",
        Description: fmt.Sprintf("Pomyślnie wypłacono %s z konta bankowego!", ToMoneyStr(toWith)),
        Color:       colors["green"],
    })
    return true
}
