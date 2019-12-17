package noname

import (
	"fmt"
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

			if test.Result{
				assert.Equal(test.ExpectedResult.SouID, client.Live.SouID)
				assert.Equal(test.ExpectedResult.TypeID, client.Live.TypeID)
			} else {
				assert.NotEqual(test.ExpectedResult.SouID, client.Live.SouID)
				assert.NotEqual(test.ExpectedResult.TypeID, client.Live.TypeID)
			}
		})
	}
}
