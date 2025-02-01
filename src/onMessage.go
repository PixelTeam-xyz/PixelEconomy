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
    userID := msg.Author.ID

    if HasPrefix(msg.Content, cnf.CommandPrefix) {
        cmd := strings.Split(TrimPrefix(msg.Content, cnf.CommandPrefix), " ")
        if parts := len(cmd); parts < 1 {
            return
        }

        switch cmd[0] {
        case "work":
            workCommand(msg, userID, cmd)
        case "bal", "balance":
            balCommand(msg, userID, cmd)
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
