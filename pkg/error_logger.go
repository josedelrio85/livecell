package livelead

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/bysidecar/voalarm"
)

// errorLogger is a struct to handle error properties
type errorLogger struct {
	msg    string
	status int
	err    error
	log    string
}

// sendAlarm to VictorOps plattform and format the error for more info
func (e *errorLogger) sendAlarm() {
	e.msg = fmt.Sprintf("Livelead -> %s %v", e.msg, e.err)
	log.Println(e.log)

	mstype := voalarm.Acknowledgement
	switch e.status {
	case http.StatusInternalServerError:
		mstype = voalarm.Warning
	case http.StatusUnprocessableEntity:
		mstype = voalarm.Info
	}

	alarm := voalarm.NewClient("")
	_, err := alarm.SendAlarm(e.msg, mstype, e.err)
	if err != nil {
		log.Fatalf(e.msg)
	}
}

// logError obtains a trace of the line and file where the error happens
func logError(err error) string {
	pc, fn, line, _ := runtime.Caller(1)
	return fmt.Sprintf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
}
