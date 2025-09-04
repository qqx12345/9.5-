package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Payload struct {
	DATA []string `json:"data"`
	ID string `json:"id"`
}
type Qet struct {
	ID string `json:"id"`
}

func Put(writer http.ResponseWriter, request *http.Request){
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if request.Method == http.MethodOptions {
		writer.WriteHeader(http.StatusOK)
		return
	}
	httpBody, _ := io.ReadAll(request.Body)
	defer request.Body.Close()
    payload := &Payload{}
	json.Unmarshal(httpBody, payload)
	err:=Insert(payload.ID,payload.DATA)
	if err!=nil {
		log.Println("插入失败")
	}
	json.NewEncoder(writer).Encode(payload.ID)
}

func Get(writer http.ResponseWriter, request *http.Request) {
	httpBody, _ := io.ReadAll(request.Body)
	defer request.Body.Close()
	get:=&Qet{}
	json.Unmarshal(httpBody, get)
	data,err:=Query(get.ID)
	if err!=nil {
		json.NewEncoder(writer).Encode("查询失败")
	} else {
		json.NewEncoder(writer).Encode(data)
	}
}

func main() {
	http.HandleFunc("/put", Put)
	http.HandleFunc("/get", Get)
	http.ListenAndServe(":1234", nil)
}