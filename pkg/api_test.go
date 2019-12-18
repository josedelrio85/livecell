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
	}

	queue := Queue{
		QueID:     244,
		QueDesc:   "",
		QueType:   2,
		QueSource: 73,
		QueActive: 1,
	}

	tests := []struct {
		Description    string
		Live           LeadLive
		Queue          Queue
		ExpectedResult LeadLive
		Result         bool
	}{
		{
			Description: "When queu_id param returns expected values",
			Live: LeadLive{
				QueueID: 244,
			},
			Queue: queue,
			ExpectedResult: LeadLive{
				SouID:  73,
				TypeID: 2,
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

			mock.ExpectQuery("SELECT (.+)").
				WithArgs(test.Live.QueueID).
				WillReturnRows(rs)

			err := client.getValues()
			assert.NoError(err)

			if test.Result {
				assert.Equal(test.ExpectedResult.SouID, client.Live.SouID)
				assert.Equal(test.ExpectedResult.TypeID, client.Live.TypeID)
			} else {
				assert.NotEqual(test.ExpectedResult.SouID, client.Live.SouID)
				assert.NotEqual(test.ExpectedResult.TypeID, client.Live.TypeID)
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

	tests := []struct {
		Description    string
		Live           LeadLive
		ExpectedRow    string
		Row            *sqlmock.Rows
		ExpectedResult LeadLive
		Result         bool
	}{
		{
			Description: "When queueID does not exists",
			Live: LeadLive{
				QueueID: 244765756,
			},
			ExpectedRow: "",
			Row:         mock.NewRows([]string{"que_id", "que_description", "que_type", "que_source", "que_active"}),
			ExpectedResult: LeadLive{
				SouID:  73,
				TypeID: 2,
			},
			Result: false,
		},
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

	tests := []struct {
		Description    string
		Live           LeadLive
		ExpectedResult ExpectedResult
	}{
		{
			Description: "When a GET request reach the endpoint",
			Live: LeadLive{
				Phone:         &phone,
				Wsid:          12345,
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
			q.Add("wsid", strconv.FormatInt(test.Live.Wsid, 10))
			q.Add("queue", strconv.FormatInt(test.Live.QueueID, 10))

			req.URL.RawQuery = q.Encode()

			http := &http.Client{}
			resp, err := http.Do(req)
			assert.NoError(err)

			assert.Equal(resp.StatusCode, test.ExpectedResult.Status)
		})
	}
}
