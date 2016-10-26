package email

import (
	"github.com/wscherphof/entity"
	"log"
	"time"
)

const (
	timeout = 60 * time.Second
)

type job struct {
	*entity.Base
	Subject    string
	Message    string
	Recipients []string
}

func init() {
	entity.Register(&job{})
	go func() {
		for {
			processQueue()
			time.Sleep(timeout)
		}
	}()
}

func initJob() *job {
	return &job{Base: &entity.Base{}}
}

func enQueue(subject, message string, recipients ...string) (err error) {
	j := initJob()
	j.Subject = subject
	j.Message = message
	j.Recipients = recipients
	return j.Update(j)
}

func processQueue() {
	if cursor, err := entity.ReadAll(&job{}); err != nil {
		log.Println("ERROR: reading email queue:", err)
	} else {
		defer cursor.Close()
		j := initJob()
		for cursor.Next(j) {
			processJob(j)
		}
		if cursor.Err() != nil {
			log.Println("ERROR: looping through email queue:", cursor.Err())
		}
	}
}

func processJob(j *job) {
	if err := send(j.Subject, j.Message, j.Recipients...); err == nil {
		j.Delete(j)
	}
}
