package email

import (
	"github.com/wscherphof/essix/db"
	"log"
	"time"
)

const (
	table   = "email_queue"
	timeout = 60 * time.Second
)

type emailJob struct {
	ID         string `gorethink:"id,omitempty"`
	Created    time.Time
	Subject    string
	Message    string
	Recipients []string
}

func initQueue() {
	if _, err := db.TableCreate(table); err == nil {
		log.Println("INFO: table created:", table)
	}
	go func() {
		for {
			processQueue()
			time.Sleep(timeout)
		}
	}()
}

func enQueue(subject, message string, recipients ...string) (err error) {
	_, err = db.Insert(table, &emailJob{
		Created:    time.Now(),
		Subject:    subject,
		Message:    message,
		Recipients: recipients,
	})
	return
}

func deQueue(job *emailJob) {
	db.Delete(table, job.ID)
}

func processQueue() {
	if cursor, err := db.All(table); err != nil {
		log.Println("ERROR: reading "+table+":", err)
	} else {
		defer cursor.Close()
		job := new(emailJob)
		for cursor.Next(job) {
			processJob(job)
		}
		if cursor.Err() != nil {
			log.Println("ERROR: looping through "+table+":", err)
		}
	}
}

func processJob(job *emailJob) {
	if err := send(job.Subject, job.Message, job.Recipients...); err == nil {
		deQueue(job)
	}
}
