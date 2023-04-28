package main

import (
	"chat-from-file/chatGpt"
	"chat-from-file/store"
	"github.com/Serendipity-sw/gutil"
	"github.com/swgloomy/gutil/glog"
)

func fileWatch() {
	var err error
	fsWatch, err = gutil.WatchFile(fileDirPath, "", removeFile, nil, nil, createFile)
	if err != nil {
		glog.Error("main fileWatch run err! webTemplate: %s err: %+v \n", fileDirPath, err)
		return
	}
	glog.Info("main fileWatch file watch start! \n")
}

func removeFile(filePath string) {
	store.FilesSync.Lock()
	defer store.FilesSync.Unlock()
	delete(store.Files, filePath)
	glog.Info("main removeFile success! fileName: %s \n", filePath)
}

func createFile(filePath string) {
	go func(path string) {
		model, err := chatGpt.UploadFile(filePath)
		if err != nil {
			glog.Error("main createFile UploadFile run err! filePath: %s err: %+v \n", path, err)
			return
		}
		store.FilesSync.Lock()
		defer store.FilesSync.Unlock()
		store.Files[path] = *model
	}(filePath)
	glog.Info("main createFile success! fileName: %s \n", filePath)
}
