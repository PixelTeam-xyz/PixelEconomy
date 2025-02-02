module PixelEconomy

go 1.23.4

require (
	github.com/BurntSushi/toml v1.4.0
	github.com/bwmarrin/discordgo v0.28.1
	github.com/mattn/go-sqlite3 v1.14.24
	msg v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)

replace (
	msg => ./Include/msg
	utils => ./Include/utils
)
