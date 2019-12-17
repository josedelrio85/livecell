package noname

import (
	"github.com/jinzhu/gorm"
)

// LeadLive represents the structure of Live Lead entity
type LeadLive struct {
	gorm.Model
	SmartcenterID int64   `schema:"lea_id"`
	SouID         int64   `schema:"-"`
	TypeID        int64   `schema:"-"`
	CatID         int64   `schema:"cat"`
	SubcatID      int64   `schema:"subcat"`
	QueueID       int64   `schema:"queue"`
	CliID         int64   `schema:"cli"`
	User          int64   `schema:"user,omitempty"`
	Closed        int64   `schema:"closed,omitempty"`
	Phone         *string `schema:"phone,omitempty"`
	Name          *string `schema:"name,omitempty"`
	Surname       *string `schema:"surn1,omitempty"`
	SecondSurname *string `schema:"surn2,omitempty"`
	DNI           *string `schema:"dni,omitempty"`
	Obs           *string `sql:"type:text" schema:"obs,omitempty"`
	FullName      *string `schema:"fullname,omitempty"`
	URL           *string `sql:"type:text" schema:"url,omitempty"`
	Wsid          int64   `schema:"wsid,omitempty"`
	IP            *string `schema:"ip,omitempty"`
	Mail          *string `schema:"mail,omitempty"`

	OrdID  int64 `schema:"ord_id,omitempty"`
	OrdSub int64 `schema:"ord_sub,omitempty"`

	ClientType    *string `schema:"clienttyp,omitempty"`
	Town          *string `schema:"town,omitempty"`
	State         *string `schema:"state,omitempty"`
	Street        *string `schema:"street,omitempty"`
	StreetAlt     *string `schema:"street2,omitempty"`
	PostalCode    *string `schema:"cp,omitempty"`
	Number        *string `schema:"number,omitempty"`
	FiberCompany  *string `schema:"fibercomp,omitempty"`
	MobileCompany *string `schema:"mobilecomp,omitempty"`

	// SouIDLeontel       int64      `sql:"-" schema:"sou_id_leontel"`
	// SouDescLeontel     string     `sql:"-" schema:"sou_desc_leontel"`
	// LeatypeIDLeontel   int64      `sql:"-" schema:"lea_type_leontel"`
	// LeatypeDescLeontel string     `sql:"-" schema:"lea_type_desc_leontel"`
}

// TableName sets the default table name
func (LeadLive) TableName() string {
	return "leadlive"
}
