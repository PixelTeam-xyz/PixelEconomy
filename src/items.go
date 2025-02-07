package main

import (
	"encoding/json"
	"io/ioutil"
)

type Item struct {
	Name        string  `json:"Name"`
	Description string  `json:"Description"`
	Price       int64   `json:"Price"`
	RoleID      string  `json:"RoleID"`
	Multiplier  float64 `json:"Multiplier"`
}

var DefaultItem = Item{
	Name:        "-",
	Description: "Wygląda na to że te pole nie zostało uzupełnione w items.json, prosze powiadomić właściciela by to naprawił!",
	Price:       0,
	Multiplier:  1,
}

func getItems() (items []Item) {
	data, err := ioutil.ReadFile("items.json")
	Except(err)
	if err != nil {
		return
	}

	var rawItems []map[string]interface{}
	if err := json.Unmarshal(data, &rawItems); err != nil {
		return
	}

	for _, raw := range rawItems {
		itm := DefaultItem
		if x, ok := raw["Name"].(string); ok {
			itm.Name = x
		}
		if x, ok := raw["Description"].(string); ok {
			itm.Description = x
		}
		if x, ok := raw["Price"].(float64); ok {
			itm.Price = int64(x)
		}
		if x, ok := raw["RoleID"].(string); ok {
			itm.RoleID = x
		}
		if x, ok := raw["Multiplier"].(float64); ok {
			itm.Multiplier = x
		}
		items = append(items, itm)
	}

	return items
}
