package email

import (
	"github.com/wscherphof/expeertise/db"
	"log"
	"time"
)

const QUEUE_TABLE string = "email_queue"
const QUEUE_TIMEOUT time.Duration = 60 * time.Second

type emailJob struct {
	ID         string `gorethink:"id,omitempty"`
	Created    time.Time
	Subject    string
	Message    string
	Recipients []string
}

// TODO: init()?
func initQueue() {
	if cursor, _ := db.TableCreate(QUEUE_TABLE); cursor != nil {
		log.Println("INFO: table created:", QUEUE_TABLE)
	}
	go func() {
		for {
			processQueue()
			time.Sleep(QUEUE_TIMEOUT)
		}
	}()
}

func enQueue(subject, message string, recipients ...string) (err error) {
	_, err = db.Insert(QUEUE_TABLE, &emailJob{
		Created:    time.Now(),
		Subject:    subject,
		Message:    message,
		Recipients: recipients,
	})
	return
}

func deQueue(job *emailJob) {
	db.Delete(QUEUE_TABLE, job.ID)
}

func processQueue() {
	if cursor, err := db.All(QUEUE_TABLE); err != nil {
		log.Println("ERROR: reading"+QUEUE_TABLE+":", err)
	} else {
		job := new(emailJob)
		for cursor.Next(job) {
			processJob(job)
		}
		if cursor.Err() != nil {
			log.Println("ERROR: looping through"+QUEUE_TABLE+":", err)
		}
	}
}

func processJob(job *emailJob) {
	if err := send(job.Subject, job.Message, job.Recipients...); err == nil {
		deQueue(job)
	}
}
