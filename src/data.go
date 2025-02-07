package main

import (
	"database/sql"
	dsc "github.com/bwmarrin/discordgo"
)

const (
	_ = iota
	ConnectToDiscordAPIErrorExit
	BotErrorExit
	CommandErrorExit
	DatabaseErrorExit
	ConfigErrorExit
	ErrorExit
)

var (
	cnf     Config = loadCnf()
	db      *sql.DB
	items   []Item = getItems()
	bot     *dsc.Session
	botUser *dsc.User
)

var colors = map[string]int{
	"red":       0xFF0000,
	"green":     0x00FF00,
	"blue":      0x0000FF,
	"yellow":    0xFFFF00,
	"purple":    0x800080,
	"cyan":      0x00FFFF,
	"orange":    0xFFA500,
	"pink":      0xFFC0CB,
	"violet":    0x8A2BE2,
	"turquoise": 0x40E0D0,
	"gray":      0x808080,
	"skyblue":   0x87CEEB,
}
