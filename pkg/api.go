package noname

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// Client is a struct that I don't have no idea about what will it do
type Client struct {
	Storer Storer
	Live   LeadLive
}

// HandleFunction receives a GET request and decode the querystring values
// into LiveLead struct.
// Returns always a 200 OK status.
// Also invokes a goroutine to process data in an isolated process.
func (c *Client) HandleFunction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c.Live = LeadLive{}
		err := decoder.Decode(&c.Live, r.URL.Query())
		if err != nil {
			msg := "Error decoding query string"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			return
		}

		log.Printf("status OK => %d - %s", c.Live.SmartcenterID, time.Now().Format("2006-01-02 15:04:05"))
		w.WriteHeader(http.StatusOK)

		c.process()
		// result := c.process()
		// result.printout()
	})
}

// getValues is a function to retrieve the souid and typeid values from queueid
func (c *Client) getValues() error {
	queue := Queue{}
	db := c.Storer.Instance()
	live := &c.Live

	if result := db.Where("que_id = ?", live.QueueID).First(&queue); result.Error != nil {
		return fmt.Errorf("Error querying Queue registry: %#v", result.Error)
	}

	live.SouID = queue.QueSource
	live.TypeID = queue.QueType

	return nil
}

// ResultError is a struct used to return a response from a goroutine using a channel
type ResultError struct {
	res string
	err error
}

// process function tries to recover missing data and store Live Lead into DB
// uses a goroutine returning ResultError struct throw a channel
func (c *Client) process() ResultError {
	outputChannel := make(chan ResultError)

	go func() {
		if c.Live.QueueID > 0 {
			if err := c.getValues(); err != nil {
				msg := "Error retrieving Queue registry"
				e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
				e.sendAlarm()
				outputChannel <- ResultError{res: msg, err: err}
			}
		}

		if c.Live.SouID == 0 || c.Live.TypeID == 0 {
			msg := "Error retrieving source and type values"
			err := fmt.Errorf("%s", msg)
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			outputChannel <- ResultError{res: msg, err: err}
		}

		if err := c.Storer.Insert(&c.Live); err != nil {
			msg := "Error inserting live lead"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			outputChannel <- ResultError{res: msg, err: err}
		}
		outputChannel <- ResultError{res: "Succesful", err: nil}
	}()
	return <-outputChannel
}

// printout function prints the data returned from the goroutine
// used as example for future cases, not used in production environment
func (r *ResultError) printout() {
	if r.err != nil {
		fmt.Printf("Failed: %s\n", r.err.Error())
	}
	fmt.Printf("Name: \"%s\" has occurred \n", r.res)
}
