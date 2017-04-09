package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/topicai/candy"
)

const (
	defaultContainerPath = "/var/lib/docker/"
)

//Record record reading position
type Record struct {
	ContainerID string `json:"container_id"`
	Seek        int64  `json:"seek"`
}

func (r *Record) reset() {
	r.ContainerID = ""
	r.Seek = 0
}
func (r *Record) writeRecord(f string) error {
	b, err := json.Marshal(r)
	candy.Must(err)
	log.Printf("write record: %s to file: %s", string(b), f)
	err = ioutil.WriteFile(f, b, 0644)
	return err
}

func readBlock(rs io.ReadSeeker, seek int64, max int) ([]string, int64, error) {
	log.Printf("start reading block: seek: %d, max: %d", seek, max)
	res := []string{}
	if _, err := rs.Seek(seek, 0); err != nil {
		return res, seek, err
	}
	r := bufio.NewReader(rs)
	ln := 0
	pos := seek
	for {
		data, err := r.ReadBytes('\n')
		ln++
		pos += int64(len(data))
		if err == nil {
			if len(data) > 0 && data[len(data)-1] == '\n' {
				data = data[:len(data)-1]
			}
		} else if err == io.EOF {
			log.Println("read the last of file, close read block.")
			return res, pos, nil
		} else {
			return res, pos, err
		}
		if len(data) != 0 {
			res = append(res, string(data))
		}
		log.Printf("read log message: %s, pos: %d, ln: %d", data, pos, ln)
		if ln >= max {
			return res, pos, nil
		}
	}
}

func main() {
	recordFile := flag.String("record-file", "./log.pos", "record file position")
	containerName := flag.String("container-name", "registry", "container name to be monitor")
	containerPath := flag.String("container-path", defaultContainerPath, "docker path")
	dbConn := flag.String("dbconnect", "user:pass@localhost/registry?sslmode=disable", "the database connect string of Paddle Stat.")
	blockSize := flag.Int("block-size", 5, "reading how many lines every time")
	flag.Parse()
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s", *dbConn))
	candy.Must(err)
	defer db.Close()
	// if no record file, create it
	if _, err := os.Stat(*recordFile); os.IsNotExist(err) {
		log.Printf("record file %s does not exists, create it.", *recordFile)
		_, err = os.Create(*recordFile)
		candy.Must(err)
	}

	record := loadRecord(*recordFile)
	defer record.writeRecord(*recordFile)
	for {
		name2id, err := fetchName2ID(*containerPath)
		candy.Must(err)
		if id, ok := name2id[*containerName]; ok {
			record.ContainerID = id
			fn := fetchLogFileByContainerID(*containerPath, id)
			fr, err := os.Open(fn)
			if err != nil {
				log.Printf("open log file %s failed, reset record, slepp 10 seconds...\n", fn)
				record.reset()
				candy.Must(record.writeRecord(*recordFile))
				time.Sleep(10 * time.Second)
				continue
			}
			block, seek, err := readBlock(fr, record.Seek, *blockSize)
			if err != nil {
				log.Printf("read block faild, %s", err.Error())
				continue
			}
			record.Seek = seek
			//update record file
			candy.Must(record.writeRecord(*recordFile))
			if len(block) == 0 {
				log.Printf("block size == 0, sleep 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			}
			// process every log message
			for _, line := range block {
				log.Println("process log: ", line)
				message, err := ParseMessage(line)
				if err != nil {
					log.Println("pass invalied log message.")
					continue
				}
				log.Printf("insert message: %+v", message)
				_, err = db.Exec(
					"INSERT INTO request(requestID, timestamp, remoteAddr, imageName,"+
						"imageTag) VALUES($1, $2, $3, $4, $5)"+
						"ON CONFLICT (requestID) "+
						"DO UPDATE SET timestamp=$2, remoteAddr=$3, imageName=$4, imageTag=$5",
					message.RequestID, message.Timestamp,
					message.RemoteAddr, message.ImageName, message.ImageTag)
				candy.Must(err)
				log.Println("insert message into postgresql successed!")
			}
		} else {
			log.Printf("Can not find container id by name: %s, sleep 10 seconds...\n", *containerName)
			time.Sleep(10 * time.Second)
			continue
		}
	}

}
