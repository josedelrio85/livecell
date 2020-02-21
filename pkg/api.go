package noname

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// Client is a struct that I don't have no idea about what will it do
type Client struct {
	Storer  Storer
	Live    LeadLive
	Payload LeadPayload
	Result  Result
}

// HandleFunction receives a GET request and decode the querystring values
// into LiveLead struct.
// Returns always a 200 OK status.
// Also invokes a goroutine to process data in an isolated process.
func (c *Client) HandleFunction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c.Live = LeadLive{}
		c.Payload = LeadPayload{}
		err := decoder.Decode(&c.Payload, r.URL.Query())
		if err != nil {
			msg := "Error decoding query string"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			return
		}

		log.Printf("status OK => %d - %s", c.Live.SmartcenterID, time.Now().Format("2006-01-02 15:04:05"))
		w.WriteHeader(http.StatusOK)

		// c.process()
		result := c.process()
		result.printout()
	})
}

// normalizeValues assigns payload values to LeadLive struct
func (c *Client) normalizeValues() error {
	payload := &c.Payload

	c.Live.QueueID = c.Payload.QueueID
	c.Live.SmartcenterID = c.Payload.SmartcenterID
	c.Live.CatID = c.Payload.CatID
	c.Live.SubcatID = c.Payload.SubcatID
	c.Live.Phone = c.Payload.Phone
	c.Live.IsClient = c.Payload.IsClient
	c.Live.URL = c.Payload.URL

	if !strings.HasPrefix(*payload.Wsid, "{{") {
		n, err := strconv.ParseInt(*payload.Wsid, 10, 64)
		if err != nil {
			return err
		}
		c.Live.Wsid = n
	}

	if !strings.HasPrefix(*payload.OrdID, "{{") {
		n, err := strconv.ParseInt(*payload.OrdID, 10, 64)
		if err != nil {
			return err
		}
		c.Live.OrdID = n
	}

	if strings.HasPrefix(*payload.URL, "{{") {
		c.Live.URL = nil
	}
	return nil
}

// getValues is a function to retrieve the souid and typeid values from queueid
// and closed value from sub_categories table
func (c *Client) getValues() error {
	queue := Queue{}
	subcat := Subcat{}
	db := c.Storer.Instance()
	live := &c.Live

	// TODO this a hard dependency from another database. Maybe create an endpoint over crmti to handle this.
	if result := db.Raw("SELECT que_source, que_type FROM crmti.que_queues WHERE que_id = ?", live.QueueID).Scan(&queue); result.Error != nil {
		return fmt.Errorf("Error querying Queue registry: %#v", result.Error)
	}
	live.SouID = queue.QueSource
	live.TypeID = queue.QueType

	if result2 := db.Raw("SELECT sub_action FROM crmti.sub_subcategories WHERE sub_id = ?", live.SubcatID).Scan(&subcat); result2.Error != nil {
		return fmt.Errorf("Error querying Subcategories registry: %#v", result2.Error)
	}

	r := Result{}
	if subcat.SubAction != "" {
		body := []byte(subcat.SubAction)
		if err := json.Unmarshal(body, &r); err != nil {
			return fmt.Errorf("Error unmarshaling sub_action field: %#v", err)
		}
	} else {
		// sqlmock does not handle json strings, so this is a hack
		r.Result = c.Result.Result
	}

	if r.Result == "2-cierre" {
		live.Closed = 1
	}

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
		if c.Payload.QueueID > 0 {
			if err := c.normalizeValues(); err != nil {
				msg := "Error normalizing values"
				e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
				e.sendAlarm()
				outputChannel <- ResultError{res: msg, err: err}
			}

			if err := c.getValues(); err != nil {
				msg := "Error getting values"
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
