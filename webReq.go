package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
)

type sampleResult struct {
	Message      string `json:"message"`
	Succeeded    bool   `json:"succeeded"`
	ErrorMessage string `json:"errorMessage"`
	name         string
}

func webRequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name, err := getParam(r, "name")
		if err != nil {
			log.Println(err.Error())
			fmt.Fprint(w, "The parameters have no names")
			return
		}
		cur, _ := os.Getwd()
		w.Header().Set("Content-Type", "application/octet-stream")
		file, fileErr := os.Open(fmt.Sprintf("%s/files/%s", cur, name))
		if fileErr != nil {
			log.Println(fileErr.Error())
			w.WriteHeader(404)
			return
		}
		fileinfo, staterr := file.Stat()
		if staterr != nil {
			log.Println(staterr.Error())
			w.WriteHeader(404)
			return
		}
		cd := mime.FormatMediaType("attachment", map[string]string{"filename": fileinfo.Name()})
		w.Header().Set("Content-Disposition", cd)
		http.ServeContent(w, r, fileinfo.Name(), fileinfo.ModTime(), file)
		/*
			fileData := make([]byte, fileinfo.Size())
			fileSize, loadErr := file.Read(fileData)
			if loadErr != nil {
				log.Println(loadErr.Error())
				w.WriteHeader(404)
				return
			}
			w.Header().Set("Content-Type", "image/png")
			w.Write(fileData)*/
		break
	case http.MethodPost:
		returnValue := &sampleResult{}
		w.Header().Set("Content-Type", "application/json")
		body, readBodyError := ioutil.ReadAll(r.Body)

		if readBodyError != nil {
			log.Println(readBodyError.Error())
			returnValue.Succeeded = false
			returnValue.ErrorMessage = "Failed reading values from body"
			failedReadingData, _ := json.Marshal(returnValue)
			w.Write(failedReadingData)
			return
		}
		wsMessage := &websocketMessage{}
		convertError := json.Unmarshal(body, &wsMessage)
		if convertError != nil {
			log.Println(convertError.Error())
			returnValue.Succeeded = false
			returnValue.ErrorMessage = "Failed converting to WebSocketMessage"
			failedConvertingData, _ := json.Marshal(returnValue)
			w.Write(failedConvertingData)
			return
		}
		w.WriteHeader(200)
		returnValue.Message = fmt.Sprintf("%s_1", wsMessage.Data)
		returnValue.Succeeded = true
		returnValue.name = "hello world"
		log.Println(returnValue)
		data, _ := json.Marshal(returnValue)

		sssmple := &sampleResult{}
		samError := json.Unmarshal(data, &sssmple)
		if samError == nil {
			log.Println(samError)
		} else {
			log.Println(samError.Error())
		}
		w.Write(data)
		break
	}
}
