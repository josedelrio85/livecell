package noname

// Subcat is a struct that represents the data of sub_subcategories Leontel table.
type Subcat struct {
	SubID       int64  `sql:"name:sub_id"`
	SubDesc     string `sql:"name:sub_description"`
	SubCat      int    `sql:"name:sub_cat"`
	SubAction   string `sql:"name:sub_action"`
	SubSystem   int    `sql:"name:sub_system"`
	SubActive   int    `sql:"name:sub_active"`
	SubClosing  int    `sql:"name:sub_closing"`
	SubUtil     int    `sql:"name:sub_util"`
	SubChAuto   int    `sql:"name:sub_sch_auto"`
	SubAux      int    `sql:"name:sub_aux"`
	SubCallback string `sql:"name:sub_callback"`
}

// Result blablabla
type Result struct {
	Result        string `json:"result,omitempty"`
	AgendaTipo    string `json:"agendaTipo,omitempty"`
	AgendaDestino string `json:"agendaDestino,omitempty"`
	AgendaMinutos string `json:"agendaMinutos,omitempty"`
	CierreTipo    string `json:"cierreTipo,omitempty"`
}

// TableName sets the default table name
func (Subcat) TableName() string {
	return "sub_subcategories"
}
