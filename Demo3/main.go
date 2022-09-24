package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Router struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	User     string `json:"user"`
	Pwd      string `json:"pwd"`
}

type Inventory struct {
	Routers []Router `json:"routers"`
}

func main() {
	file, err := os.Open("data.json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	data := json.NewDecoder(file)

	var inv Inventory

	err = data.Decode(&inv)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(inv.Routers); i++ {
		fmt.Printf("%+v\n", inv.Routers[i].Hostname)
	}

}
