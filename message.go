package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"
)

//Message record structed log message
type Message struct {
	RequestID  string
	ImageName  string
	ImageTag   string
	Timestamp  int64
	RemoteAddr string
}

func str2timestamp(datestr string) (int64, error) {
	//	d := string(datestr[:len(datestr)-11]) + "Z"
	t, err := time.Parse(time.RFC3339Nano, datestr)
	fmt.Println(t.Local())
	if err != nil {
		fmt.Println(err.Error())
		return -1, err
	}
	return t.Local().Unix(), nil
}

//ParseMessage parse log message to struct
func ParseMessage(message string) (Message, error) {
	var data map[string]interface{}
	var m Message
	if err := json.Unmarshal([]byte(message), &data); err != nil {
		panic(err)
	}
	l := data["log"].(string)
	t := data["time"].(string)
	timestamp, err := str2timestamp(t)
	if err != nil {
		log.Printf("parse datestr failed: %s", t)
	}
	r, err := regexp.Compile("http.request.id=([a-zA-Z0-9-]+) http.request.method=GET http.request.remoteaddr=\"([0-9.:]+)\" http.request.uri=\"/v2/([a-zA-Z0-9]+)/manifests/([a-zA-Z0-9-.]+)\" ")
	res := r.FindStringSubmatch(l)
	fmt.Println(res)
	if err != nil {
		return m, err
	}
	if len(res) != 6 {
		return m, errors.New("parse log message failed:" + message)
	}

	m = Message{
		RequestID:  res[1],
		RemoteAddr: res[2],
		ImageName:  res[3],
		ImageTag:   res[4],
		Timestamp:  timestamp,
	}
	return m, nil
}
