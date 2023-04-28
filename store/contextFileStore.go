package store

import (
	_struct "chat-from-file/struct"
	"sync"
)

var (
	Files     = make(map[string]_struct.ChatGptUploadRespStruct)
	FilesSync *sync.RWMutex
)
