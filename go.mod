module PixelEconomy

go 1.23.4

require (
	github.com/BurntSushi/toml v1.4.0
	github.com/bwmarrin/discordgo v0.28.1
	github.com/mattn/go-sqlite3 v1.14.24
	golang.org/x/exp v0.0.0-20250128182459-e0ece0dbea4c
	msg v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)

replace (
	msg => ./Include/msg
	utils => ./Include/utils
)
