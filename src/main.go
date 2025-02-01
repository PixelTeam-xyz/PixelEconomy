package main

import (
    "database/sql"
    dsc "github.com/bwmarrin/discordgo"
    "io/ioutil"
    "msg"
    "os"

    _ "github.com/mattn/go-sqlite3"
    . "utils"
)

var (
    cnf Config = loadCnf()
    db  *sql.DB
    bot *dsc.Session
)

func ShowError(msgs ...any) {
    msg.Error(msgs...)
    db.Close()
    os.Exit(1)
}

func main() {
    defer db.Close()
    defer stmt.Close()

    if Contains("--initConfig", os.Args) {
        err := createDefault()
        Except("Creating default config file failed (%s)", err)
    }

    data, err := ioutil.ReadFile("token.txt")
    Except("Reading token.txt failed (%s)", err, ConfigErrorExit)

    tk := string(data)
    if tk == "" {
        ShowError("Put your bot token in the token.txt file and then restart the bot")
    }

    bot, err = dsc.New("Bot " + tk)
    Except("Opening connection to discord API failed (%s)", err, ConnectToDiscordAPIErrorExit)

    botUser, err := bot.User("@me")

    bot.AddHandler(onMessage)

    err = bot.Open()
    Except("Opening connection to discord API failed (%s)", err, ConnectToDiscordAPIErrorExit)

    msg.Infof("The bot successfully connected to the discord API as a %s (%s)", botUser.Username, botUser.ID)
    select {}
}
