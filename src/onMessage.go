package main

import (
    dsc "github.com/bwmarrin/discordgo"
    "strings"
    //log "msg"
    . "utils"
)

func onMessage(bot *dsc.Session, msg *dsc.MessageCreate) {
    if msg.Author.Bot {
        return
    }

    if HasPrefix(msg.Content, cnf.CommandPrefix) {
        cmd := strings.Split(TrimPrefix(msg.Content, cnf.CommandPrefix), " ")
        if parts := len(cmd); parts < 1 {
            return
        }

        switch cmd[0] {
        case "work":
            // TODO
        case "bal", "balance":
            // TODO
        case "top":
            // TODO
        case "buy":
            // TODO
        case "shop":
            // TODO
        case "test":
            _, err := bot.ChannelMessageSend(msg.ChannelID, "Test!")
            Except(err)
        }
    }
}
