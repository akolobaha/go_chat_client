package main

import (
	"context"
	"go_chat_client/config"
	"go_chat_client/internal/app"
)

const configFilePath = "./config.toml"

func main() {
	ctx := context.Background()
	config.Parse(configFilePath)

	cmd := app.NewMenu()

	err := cmd.ExecuteContext(ctx)

	if err != nil {
		return
	}

}
