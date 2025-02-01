package main

import (
    "msg"
    "os"

    "github.com/BurntSushi/toml"
)

type Config struct {
    b             string `toml:"MoneyIcon"`
    WorkDelay     int    `toml:"WorkDelay"`
    database      string `toml:"DatabasePath"`
    WorkMin       int    `toml:"WorkEarningsMin"`
    WorkMax       int    `toml:"WorkEarningsMax"`
    TopCh         int64  `toml:"TopMessagesChannelID"`
    CommandPrefix rune   `toml:"CommandPrefix"`
    // If the value is -1, it will not send the top list ^^^^

    // TODO: Add more configuration options
}

var defaultConfig = Config{
    b:             "💴",
    WorkDelay:     30,
    database:      "economy.db",
    WorkMin:       50,
    WorkMax:       200,
    TopCh:         -1,
    CommandPrefix: '!',
}

func loadCnf() Config {
    var config Config

    if _, err := toml.DecodeFile("config.toml", &config); err != nil {
        msg.Fatalf("Error reading config.toml or file does not exist! Using the default configuration")
        return defaultConfig
    }

    return config
}

func createDefault() error {
    file, err := os.Create("config.toml")
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := toml.NewEncoder(file)
    if err := encoder.Encode(defaultConfig); err != nil {
        return err
    }
    return nil
}
