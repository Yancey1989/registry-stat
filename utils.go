package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

	"github.com/topicai/candy"
)

//ReadRecord read record message from file
func loadRecord(f string) Record {
	var record Record
	err := json.NewDecoder(candy.ReadAll(f)).Decode(&record)
	if err != nil {
		record = Record{
			ContainerID: "",
			Seek:        0,
		}
		return record
	}
	return record
}
func fetchLogFileByContainerID(containerPath, containerID string) string {
	return path.Join(containerPath, "containers", containerID, containerID+"-json.log")
}

func fetchName2ID(containerPath string) (map[string]string, error) {
	fs, err := ioutil.ReadDir(path.Join(containerPath, "containers"))
	candy.Must(err)
	id2name := make(map[string]string)
	for _, f := range fs {
		id := f.Name()
		conf := path.Join(containerPath, "containers", id, "config.v2.json")
		j := make(map[string]interface{})
		candy.Must(json.NewDecoder(candy.ReadAll(conf)).Decode(&j))
		name := strings.TrimLeft(j["Name"].(string), "/")
		id2name[name] = id
	}
	return id2name, nil
}
