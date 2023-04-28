package chatGpt

import (
	"bytes"
	"chat-from-file/store"
	_struct "chat-from-file/struct"
	"encoding/json"
	"fmt"
	"github.com/swgloomy/gutil/glog"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func UploadFile(filePathStr string) (*_struct.ChatGptUploadRespStruct, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePathStr))
	if err != nil {
		glog.Error("chatGpt UploadFile CreateFormFile err! filePath: %s err: %+v \n", filePathStr, err)
		return nil, err
	}
	file, err := os.Open(filePathStr)
	if err != nil {
		glog.Error("chatGpt UploadFile Open file err! filePath: %s err: %+v \n", filePathStr, err)
		return nil, err
	}
	defer file.Close()
	_, err = io.Copy(part, file)
	if err != nil {
		glog.Error("chatGpt UploadFile file copy err! filePath: %s err: %+v \n", filePathStr, err)
		return nil, err
	}
	err = writer.WriteField("purpose", "fine-tune")
	if err != nil {
		glog.Error("chatGpt UploadFile WriteField err! filePath: %s err: %+v \n", filePathStr, err)
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/files", body)
	if err != nil {
		glog.Error("chatGpt UploadFile NewRequest err! filePath: %s err: %+v \n", filePathStr, err)
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", store.ChatGPTToken))
	uri, err := url.Parse("socket5://127.0.0.1:10808")
	if err != nil {
		glog.Error("chatGpt UploadFile Parse err! filePath: %s err: %+v \n", filePathStr, err)
		return nil, err
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		glog.Error("chatGpt UploadFile http send err! filePath: %s err: %+v \n", filePathStr, err)
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			glog.Error("chatGpt UploadFile body close err! filePath: %s err: %+v \n", filePathStr, err)
		}
	}()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		glog.Error("chatGpt UploadFile body read err! filePath: %s err: %+v \n", filePathStr, err)
		return nil, err
	}

	var model _struct.ChatGptUploadRespStruct

	err = json.Unmarshal(respBody, &model)
	if err != nil {
		glog.Error("chatGpt UploadFile body Unmarshal err! filePath: %s body: %s err: %+v \n", filePathStr, string(respBody), err)
		return nil, err
	}
	glog.Info("chatGpt UploadFile run success! filePath: %s respModel: %+v body: %s \n", filePathStr, model, string(respBody))
	return &model, nil
}
