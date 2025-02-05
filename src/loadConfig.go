package main

import (
	"info"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	MoneyIcon                        string
	WorkDelay                        int
	RobDelay                         int
	CrimeDelay                       int
	DatabasePath                     string
	WorkMin                          int
	WorkMax                          int
	TopCh                            int64
	DelPrevTopListOnNewSend          bool
	DelayInSendingTopList            int
	CommandPrefix                    string
	MainEmbedColor                   int
	ServerID                         string
	NumberOfUsersInTopList           int
	AllowedChannels                  []int64
	DisappearanceTimeOfErrorMessages int
	AdminUsersIDs                    []int64
	AdminRolesIDs                    []int64
	RobberySuccessChance             int8
	RobSuccessEarningsMin            int
	RobSuccessEarningsMax            int
	RobFailureLossMin                int
	RobFailureLossMax                int
	CrimeSuccessEarningsMax          int
	CrimeSuccessEarningsMin          int
	CrimeFailureLossMin              int
	CrimeFailureLossMax              int
	CrimeSuccessChance               int8
}

type GroupedConfig struct {
	General struct {
		MoneyIcon      string `toml:"MoneyIcon"`
		DatabasePath   string `toml:"DatabasePath"`
		CommandPrefix  string `toml:"CommandPrefix"`
		MainEmbedColor int    `toml:"MainEmbedColor"`
	} `toml:"General"`
	Economy struct {
		WorkMin                 int  `toml:"WorkEarningsMin"`
		WorkMax                 int  `toml:"WorkEarningsMax"`
		WorkDelay               int  `toml:"WorkDelay"`
		RobDelay                int  `toml:"RobDelay"`
		CrimeDelay              int  `toml:"CrimeDelay"`
		RobberySuccessChance    int8 `toml:"RobberySuccessChance"`
		RobSuccessEarningsMin   int  `toml:"RobSuccessEarningsMin"`
		RobSuccessEarningsMax   int  `toml:"RobSuccessEarningsMax"`
		RobFailureLossMin       int  `toml:"RobFailureLossMin"`
		RobFailureLossMax       int  `toml:"RobFailureLossMax"`
		CrimeSuccessEarningsMax int  `toml:"CrimeSuccessEarningsMax"`
		CrimeSuccessEarningsMin int  `toml:"CrimeSuccessEarningsMin"`
		CrimeFailureLossMin     int  `toml:"CrimeFailureLossMin"`
		CrimeFailureLossMax     int  `toml:"CrimeFailureLossMax"`
		CrimeSuccessChance      int8 `toml:"CrimeSuccessChance"`
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
		RobDelay:                         self.Economy.RobDelay,
		CrimeDelay:                       self.Economy.CrimeDelay,
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
		RobberySuccessChance:             self.Economy.RobberySuccessChance,
		RobSuccessEarningsMax:            self.Economy.RobSuccessEarningsMax,
		RobSuccessEarningsMin:            self.Economy.RobSuccessEarningsMin,
		RobFailureLossMax:                self.Economy.RobFailureLossMax,
		RobFailureLossMin:                self.Economy.RobFailureLossMin,
		CrimeSuccessEarningsMax:          self.Economy.CrimeSuccessEarningsMax,
		CrimeSuccessEarningsMin:          self.Economy.CrimeSuccessEarningsMin,
		CrimeFailureLossMin:              self.Economy.CrimeFailureLossMin,
		CrimeFailureLossMax:              self.Economy.CrimeFailureLossMax,
		CrimeSuccessChance:               self.Economy.CrimeSuccessChance,
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

	// ... //

	gc.Economy.RobDelay = self.RobDelay
	gc.Economy.RobberySuccessChance = self.RobberySuccessChance
	gc.Economy.RobSuccessEarningsMin = self.RobSuccessEarningsMin
	gc.Economy.RobSuccessEarningsMax = self.RobSuccessEarningsMax
	gc.Economy.RobFailureLossMin = self.RobFailureLossMin
	gc.Economy.RobFailureLossMax = self.RobFailureLossMax

	// ... //

	gc.Economy.CrimeDelay = self.CrimeDelay
	gc.Economy.CrimeSuccessChance = self.CrimeSuccessChance
	gc.Economy.CrimeSuccessEarningsMin = self.CrimeSuccessEarningsMin
	gc.Economy.CrimeSuccessEarningsMax = self.CrimeSuccessEarningsMax
	gc.Economy.CrimeFailureLossMin = self.CrimeFailureLossMin
	gc.Economy.CrimeFailureLossMax = self.CrimeFailureLossMax
	return gc
}

var defaultConfig = Config{
	MoneyIcon:                        "💴",
	WorkDelay:                        20,
	DatabasePath:                     "economy.db",
	WorkMin:                          100,
	WorkMax:                          300,
	TopCh:                            -1,
	DelPrevTopListOnNewSend:          true,
	DelayInSendingTopList:            3600,
	CommandPrefix:                    "!",
	MainEmbedColor:                   colors["skyblue"],
	NumberOfUsersInTopList:           10,
	AllowedChannels:                  []int64{},
	DisappearanceTimeOfErrorMessages: 5,

	// ... //

	RobDelay:              300,
	RobberySuccessChance:  50,
	RobSuccessEarningsMax: 400,
	RobSuccessEarningsMin: 50,
	RobFailureLossMax:     150,
	RobFailureLossMin:     50,

	// ... //

	CrimeDelay:              200,
	CrimeSuccessChance:      50,
	CrimeSuccessEarningsMax: 400,
	CrimeSuccessEarningsMin: 50,
	CrimeFailureLossMax:     150,
	CrimeFailureLossMin:     50,
}

func loadCnf() Config {
	var gc GroupedConfig
	if _, err := toml.DecodeFile("config.toml", &gc); err != nil {
		info.Fatalf("Error reading config.toml or file does not exist! Using the default configuration")
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
