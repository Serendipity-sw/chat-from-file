package main

import (
	"chat-from-file/store"
	"github.com/guotie/config"
)

func readConfig() {
	store.ChatGPTToken = config.GetString("token")
	store.ChatGPTModelId = config.GetString("modelId")
}
