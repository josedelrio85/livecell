package noname

import (
	"fmt"
	"log"
	"net/http"

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

		keys := r.URL.Query()

		// TODO remove it
		for i, k := range keys {
			log.Printf("key %v => %s", keys[i], k)
		}

		c.Live = LeadLive{}
		err := decoder.Decode(&c.Live, r.URL.Query())
		if err != nil {
			msg := "Error decoding query string"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			return
		}

		w.WriteHeader(http.StatusOK)

		go c.process()
	})
}

// process function tries to recover missing data and store Live Lead into DB
func (c *Client) process() {
	if c.Live.QueueID > 0 {
		if err := c.getValues(); err != nil {
			msg := "Error retrieving Queue registry"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
		}
	}

	// TODO remove it
	log.Printf("Live SouID %d", c.Live.SouID)
	log.Printf("Live TypeID %d", c.Live.TypeID)

	if c.Live.SouID == 0 || c.Live.TypeID == 0 {
		msg := "Error retrieving source and type values"
		err := fmt.Errorf("%s", msg)
		e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
		e.sendAlarm()
		return
	}

	if err := c.Storer.Insert(&c.Live); err != nil {
		msg := "Error inserting live lead"
		e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
		e.sendAlarm()
	}

	// TODO remove it
	log.Println("done...")
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
