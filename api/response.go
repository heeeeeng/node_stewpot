package api

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	RespFail    = 500
	RespSuccess = 200
)

type Resp struct {
	Status  int         `json:"status"`
	Error   string      `json:"error"`
	Content interface{} `json:"content"`
}

func Response(w http.ResponseWriter, content interface{}) {
	response(w, RespSuccess, "", content)
}

func ErrorResponse(w http.ResponseWriter, error string) {
	response(w, RespFail, error, nil)
}

func response(w http.ResponseWriter, status int, error string, content interface{}) {
	resp := Resp{}
	resp.Status = status
	resp.Error = error
	resp.Content = content

	data, err := json.Marshal(&resp)
	if err != nil {
		log.Fatalf("marshal responss to json error: %v", err)
		return
	}
	w.Write(data)
}
