package noname

import (
	"github.com/jinzhu/gorm"
)

// LeadPayload represents the received data structure from Leontel environment
// Equivalences to get values from Leontel
// lea_id => leadId		|| 	cat => cat_id
// subcat => sub_id  	|| 	queue => queId
// phone => TELEFONO 	|| 	url => url
// is_client => cli 	|| 	ws_id => wsid
// ord_id => isOpenOrder
type LeadPayload struct {
	QueueID       int64   `schema:"queue"`
	SmartcenterID int64   `schema:"lea_id"`
	CatID         int64   `schema:"cat"`
	SubcatID      int64   `schema:"subcat"`
	OrdID         *string `schema:"ord_id,omitempty"` // {{
	Wsid          *string `schema:"ws_id,omitempty"`  // {{
	Phone         *string `schema:"phone"`
	URL           *string `schema:"url"` // {{
	IsClient      bool    `schema:"is_client"`
}

// LeadLive represents the structure of Live Lead entity
type LeadLive struct {
	gorm.Model
	SouID         int64
	TypeID        int64
	QueueID       int64
	SmartcenterID int64
	CatID         int64
	SubcatID      int64
	OrdID         int64
	Wsid          int64
	Closed        int64
	Phone         *string
	IsClient      bool
	URL           *string
	// User   int64 `schema:"user,omitempty"`
}

// TableName sets the default table name
func (LeadLive) TableName() string {
	return "leadlive"
}
