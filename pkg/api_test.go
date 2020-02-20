package noname

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeValues(t *testing.T) {
	assert := assert.New(t)

	client := Client{}

	wsidok := "234"
	wsidko := "{{wsid}}"
	ordidok := "1"
	ordidko := "{{ordid}}"
	urlok := "adfasdf"
	urlko := "{{url}}"

	tests := []struct {
		Description    string
		Payload        LeadPayload
		Live           LeadLive
		ExpectedResult LeadLive
	}{
		{
			Description: "When wsid, ordid and url are valid values",
			Payload: LeadPayload{
				Wsid:  &wsidok,
				OrdID: &ordidok,
				URL:   &urlok,
			},
			Live: LeadLive{},
			ExpectedResult: LeadLive{
				Wsid:  234,
				OrdID: 1,
				URL:   &urlok,
			},
		},
		{
			Description: "When wsid is not a valid value",
			Payload: LeadPayload{
				Wsid:  &wsidko,
				OrdID: &ordidok,
				URL:   &urlok,
			},
			Live: LeadLive{},
			ExpectedResult: LeadLive{
				Wsid:  0,
				OrdID: 1,
				URL:   &urlok,
			},
		},
		{
			Description: "When ordid is not a valid value",
			Payload: LeadPayload{
				Wsid:  &wsidok,
				OrdID: &ordidko,
				URL:   &urlok,
			},
			Live: LeadLive{},
			ExpectedResult: LeadLive{
				Wsid:  234,
				OrdID: 0,
				URL:   &urlok,
			},
		},
		{
			Description: "When url is not a valid value",
			Payload: LeadPayload{
				Wsid:  &wsidok,
				OrdID: &ordidok,
				URL:   &urlko,
			},
			Live: LeadLive{},
			ExpectedResult: LeadLive{
				Wsid:  234,
				OrdID: 1,
				URL:   nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			client.Payload = test.Payload
			client.Live = test.Live

			err := client.normalizeValues()
			assert.NoError(err)

			assert.Equal(test.ExpectedResult.Wsid, client.Live.Wsid)
			assert.Equal(test.ExpectedResult.OrdID, client.Live.OrdID)
			assert.Equal(test.ExpectedResult.URL, client.Live.URL)
		})
	}
}

func TestGetValues(t *testing.T) {
	assert := assert.New(t)

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_1")
	assert.NoError(err)

	db, err = gorm.Open("sqlmock", "sqlmock_db_1")
	defer db.Close()

	fakedb := FakeDb{
		OpenFunc:     func() error { return nil },
		CloseFunc:    func() error { return nil },
		UpdateFunc:   func(a interface{}, wCond string, wFields []string) error { return nil },
		InsertFunc:   func(lead interface{}) error { return nil },
		InstanceFunc: func() *gorm.DB { return db },
	}

	client := Client{
		Storer: &fakedb,
		Result: Result{
			Result:        "2-cierre",
			AgendaTipo:    "",
			AgendaDestino: "",
			AgendaMinutos: "",
			CierreTipo:    "1-positivo",
		},
	}

	queue := Queue{
		QueID:     244,
		QueDesc:   "",
		QueType:   2,
		QueSource: 73,
		QueActive: 1,
	}

	subcat := Subcat{
		SubID: 341,
	}

	tests := []struct {
		Description    string
		Live           LeadLive
		Queue          Queue
		Subcat         Subcat
		ExpectedResult LeadLive
		ExpectedSubcat Subcat
		Result         bool
	}{
		{
			Description: "When queu_id param returns expected values",
			Live: LeadLive{
				QueueID:  244,
				SubcatID: 341,
			},
			Queue:  queue,
			Subcat: subcat,
			ExpectedResult: LeadLive{
				SouID:  73,
				TypeID: 2,
				Closed: 1,
			},
			Result: true,
		},
		{
			Description: "When queu_id param returns unexpected values",
			Live: LeadLive{
				QueueID: 244,
			},
			Queue: queue,
			ExpectedResult: LeadLive{
				SouID:  777,
				TypeID: 222,
				Closed: 1,
			},
			Result: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			client.Live = test.Live

			row := fmt.Sprintf("%d,%s,%d,%d,%d",
				test.Queue.QueID, test.Queue.QueDesc, test.Queue.QueType, test.Queue.QueSource, test.Queue.QueActive)

			rs := mock.NewRows([]string{"que_id", "que_description", "que_type", "que_source", "que_active"}).
				FromCSVString(row)

			mock.ExpectQuery("SELECT (.+)").WithArgs(test.Live.QueueID).WillReturnRows(rs)

			rowsub := fmt.Sprintf("%d", test.Subcat.SubID)
			rssub := mock.NewRows([]string{"sub_id"}).FromCSVString(rowsub)
			mock.ExpectQuery("SELECT sub_action FROM crmti.sub_subcategories").WithArgs(test.Live.SubcatID).WillReturnRows(rssub)

			err := client.getValues()
			assert.NoError(err)

			if test.Result {
				assert.Equal(test.ExpectedResult.SouID, client.Live.SouID)
				assert.Equal(test.ExpectedResult.TypeID, client.Live.TypeID)
				assert.Equal(test.ExpectedResult.Closed, client.Live.Closed)
			} else {
				assert.NotEqual(test.ExpectedResult.SouID, client.Live.SouID)
				assert.NotEqual(test.ExpectedResult.TypeID, client.Live.TypeID)
				assert.Equal(test.ExpectedResult.Closed, client.Live.Closed)
			}
		})
	}
}

func TestProcess(t *testing.T) {
	assert := assert.New(t)

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_2")
	assert.NoError(err)

	db, err = gorm.Open("sqlmock", "sqlmock_db_2")
	defer db.Close()

	fakedb := FakeDb{
		OpenFunc:     func() error { return nil },
		CloseFunc:    func() error { return nil },
		UpdateFunc:   func(a interface{}, wCond string, wFields []string) error { return nil },
		InsertFunc:   func(lead interface{}) error { return nil },
		InstanceFunc: func() *gorm.DB { return db },
	}

	client := Client{
		Storer: &fakedb,
	}

	queue := Queue{
		QueID:     244,
		QueDesc:   "",
		QueType:   2,
		QueSource: 73,
		QueActive: 1,
	}

	// subcat := Subcat{
	// 	SubID: 341,
	// }

	tests := []struct {
		Description    string
		Live           LeadLive
		ExpectedRow    string
		Row            *sqlmock.Rows
		ExpectedResult LeadLive
		Result         bool
	}{
		// {
		// 	Description: "When queueID does not exists",
		// 	Live: LeadLive{
		// 		QueueID: 244765756,
		// 	},
		// 	ExpectedRow: "",
		// 	Row:         mock.NewRows([]string{"que_id", "que_description", "que_type", "que_source", "que_active"}),
		// 	ExpectedResult: LeadLive{
		// 		SouID:  73,
		// 		TypeID: 2,
		// 	},
		// 	Result: false,
		// },
		{
			Description: "When queueID exists",
			Live: LeadLive{
				QueueID: 244,
			},
			ExpectedRow: fmt.Sprintf("%d,%s,%d,%d,%d", queue.QueID, queue.QueDesc, queue.QueType, queue.QueSource, queue.QueActive),
			Row:         mock.NewRows([]string{"que_id", "que_description", "que_type", "que_source", "que_active"}),
			ExpectedResult: LeadLive{
				SouID:  73,
				TypeID: 2,
			},
			Result: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			client.Live = test.Live

			if test.ExpectedRow != "" {
				test.Row.FromCSVString(test.ExpectedRow)
			}

			mock.ExpectQuery("SELECT (.+)").
				WithArgs(test.Live.QueueID).
				WillReturnRows(test.Row)

			err := client.process()

			if test.Result {
				assert.Nil(err.err)

				assert.Equal(test.ExpectedResult.SouID, client.Live.SouID)
				assert.Equal(test.ExpectedResult.TypeID, client.Live.TypeID)
			} else {
				assert.NotNil(err.err)

				assert.NotEqual(test.ExpectedResult.SouID, client.Live.SouID)
				assert.NotEqual(test.ExpectedResult.TypeID, client.Live.TypeID)
			}
		})
	}
}

type ExpectedResult struct {
	Status int
}

func TestHandler(t *testing.T) {
	assert := assert.New(t)

	phone := "123456789"
	// wsid := "12345"

	tests := []struct {
		Description    string
		Live           LeadPayload
		ExpectedResult ExpectedResult
	}{
		{
			Description: "When a GET request reach the endpoint",
			Live: LeadPayload{
				Phone: &phone,
				// Wsid:          &wsid,
				SmartcenterID: 1234,
				QueueID:       244,
			},
			ExpectedResult: ExpectedResult{
				Status: http.StatusOK,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				testresp, err := json.Marshal(test.ExpectedResult.Status)
				assert.Nil(err)

				res.WriteHeader(test.ExpectedResult.Status)
				res.Write(testresp)
			}))
			defer func() { ts.Close() }()

			req, err := http.NewRequest("GET", ts.URL, nil)
			assert.NoError(err)

			q := req.URL.Query()
			q.Add("lea_id", strconv.FormatInt(test.Live.SmartcenterID, 10))
			q.Add("phone", *test.Live.Phone)
			// q.Add("wsid", strconv.FormatInt(test.Live.Wsid, 10))
			// q.Add("ws_id", *test.Live.WsValue)
			q.Add("queue", strconv.FormatInt(test.Live.QueueID, 10))

			req.URL.RawQuery = q.Encode()

			http := &http.Client{}
			resp, err := http.Do(req)
			assert.NoError(err)

			assert.Equal(resp.StatusCode, test.ExpectedResult.Status)
		})
	}
}
