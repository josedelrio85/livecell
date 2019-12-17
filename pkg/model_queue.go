package noname

// Queue is a struct that represents the data of que_queues Leontel table.
type Queue struct {
	QueID     int64  `sql:"name:que_id"`
	QueDesc   string `sql:"name:que_description"`
	QueType   int64  `sql:"name:que_type"`
	QueSource int64  `sql:"name:que_source"`
	QueActive int64  `sql:"name:que_active"`
}

// TableName sets the default table name
func (Queue) TableName() string {
	return "que_queues"
}
