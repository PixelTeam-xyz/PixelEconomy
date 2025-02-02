package main

import (
	"msg"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	MoneyIcon                        string  `toml:"MoneyIcon"`
	WorkDelay                        int     `toml:"WorkDelay"`
	DatabasePath                     string  `toml:"DatabasePath"`
	WorkMin                          int     `toml:"WorkEarningsMin"`
	WorkMax                          int     `toml:"WorkEarningsMax"`
	TopCh                            int64   `toml:"TopMessagesChannelID"`
	DelPrevTopListOnNewSend          bool    `toml:"DelPrevTopListOnNewSend"`
	DelayInSendingTopList            int     `toml:"DelayInSendingTopList"`
	CommandPrefix                    string  `toml:"CommandPrefix"`
	MainEmbedColor                   int     `toml:"MainEmbedColor"`
	ServerID                         string  `toml:"ServerID"`
	NumberOfUsersInTopList           int     `toml:"NumberOfUsersInTopList"`
	AllowedChannels                  []int64 `toml:"AllowedChannels"`
	DisappearanceTimeOfErrorMessages int     `toml:"DisappearanceTimeOfErrorMessages"`
	AdminUsersIDs                    []int64 `toml:"AdminUsersIDs"`
	AdminRolesIDs                    []int64 `toml:"AdminRolesIDs"`
}

type GroupedConfig struct {
	General struct {
		MoneyIcon      string `toml:"MoneyIcon"`
		DatabasePath   string `toml:"DatabasePath"`
		CommandPrefix  string `toml:"CommandPrefix"`
		MainEmbedColor int    `toml:"MainEmbedColor"`
	} `toml:"General"`
	Economy struct {
		WorkMin   int `toml:"WorkEarningsMin"`
		WorkMax   int `toml:"WorkEarningsMax"`
		WorkDelay int `toml:"WorkDelay"`
	} `toml:"Economy"`
	TopList struct {
		TopCh                   int64 `toml:"TopMessagesChannelID"`
		DelPrevTopListOnNewSend bool  `toml:"DelPrevTopListOnNewSend"`
		DelayInSendingTopList   int   `toml:"DelayInSendingTopList"`
		NumberOfUsersInTopList  int   `toml:"NumberOfUsersInTopList"`
	} `toml:"TopList"`
	Server struct {
		ServerID        string  `toml:"ServerID"`
		AllowedChannels []int64 `toml:"AllowedChannels"`
		AdminUsersIDs   []int64 `toml:"AdminUsersIDs"`
		AdminRolesIDs   []int64 `toml:"AdminRolesIDs"`
	} `toml:"Server"`
	Messages struct {
		DisappearanceTimeOfErrorMessages int `toml:"DisappearanceTimeOfErrorMessages"`
	} `toml:"Messages"`
}

func (self *GroupedConfig) toConfig() Config {
	return Config{
		MoneyIcon:                        self.General.MoneyIcon,
		DatabasePath:                     self.General.DatabasePath,
		CommandPrefix:                    self.General.CommandPrefix,
		MainEmbedColor:                   self.General.MainEmbedColor,
		WorkDelay:                        self.Economy.WorkDelay,
		WorkMin:                          self.Economy.WorkMin,
		WorkMax:                          self.Economy.WorkMax,
		TopCh:                            self.TopList.TopCh,
		DelPrevTopListOnNewSend:          self.TopList.DelPrevTopListOnNewSend,
		DelayInSendingTopList:            self.TopList.DelayInSendingTopList,
		NumberOfUsersInTopList:           self.TopList.NumberOfUsersInTopList,
		ServerID:                         self.Server.ServerID,
		AllowedChannels:                  self.Server.AllowedChannels,
		AdminUsersIDs:                    self.Server.AdminUsersIDs,
		AdminRolesIDs:                    self.Server.AdminRolesIDs,
		DisappearanceTimeOfErrorMessages: self.Messages.DisappearanceTimeOfErrorMessages,
	}
}

func (self Config) ToGrouped() GroupedConfig {
	var gc GroupedConfig
	gc.General.MoneyIcon = self.MoneyIcon
	gc.General.DatabasePath = self.DatabasePath
	gc.General.CommandPrefix = self.CommandPrefix
	gc.General.MainEmbedColor = self.MainEmbedColor
	gc.Economy.WorkMin = self.WorkMin
	gc.Economy.WorkMax = self.WorkMax
	gc.Economy.WorkDelay = self.WorkDelay
	gc.TopList.TopCh = self.TopCh
	gc.TopList.DelPrevTopListOnNewSend = self.DelPrevTopListOnNewSend
	gc.TopList.DelayInSendingTopList = self.DelayInSendingTopList
	gc.TopList.NumberOfUsersInTopList = self.NumberOfUsersInTopList
	gc.Server.ServerID = self.ServerID
	gc.Server.AllowedChannels = self.AllowedChannels
	gc.Messages.DisappearanceTimeOfErrorMessages = self.DisappearanceTimeOfErrorMessages
	return gc
}

var defaultConfig = Config{
	MoneyIcon:                        "💴",
	WorkDelay:                        30,
	DatabasePath:                     "economy.db",
	WorkMin:                          50,
	WorkMax:                          200,
	TopCh:                            -1,
	DelPrevTopListOnNewSend:          true,
	DelayInSendingTopList:            3600,
	CommandPrefix:                    "!",
	MainEmbedColor:                   colors["skyblue"],
	NumberOfUsersInTopList:           10,
	AllowedChannels:                  []int64{},
	DisappearanceTimeOfErrorMessages: 5,
}

func loadCnf() Config {
	var gc GroupedConfig
	if _, err := toml.DecodeFile("config.toml", &gc); err != nil {
		msg.Fatalf("Error reading config.toml or file does not exist! Using the default configuration")
		return defaultConfig
	}
	return gc.toConfig()
}

func createDefault() error {
	file, err := os.Create("config.toml")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := toml.NewEncoder(file)
	gc := defaultConfig.ToGrouped()
	if err := encoder.Encode(gc); err != nil {
		return err
	}
	return nil
}
