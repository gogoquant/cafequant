package model

import (
	"strings"
	"time"
)

var (
	DATAGRAM_INFO  = "data_info"
	DATAGRAM_VAL   = "data_val"
	DATAGRAM_COLOR   = "data_color"
)

// Datagram
type Datagram struct {
	Uid       string    			   `json:"uuid"`
	Fields    map[string]interface{}   `json:"fields"`
	Tags      map[string]string        `json:"tags"`
	Time time.Time 					   `json:"time"`
}

// SetInfo
func (d *Datagram)SetInfo(m map[string]string){
	var info string
	for key, val := range(m){
		key  = strings.Replace(key, ",", " ", -1)
		val  = strings.Replace(val, ",", " ", -1)
		info = info + key + "=" + val + ","
	}
	d.Tags[DATAGRAM_INFO] = info
}

// SetColor
func (d *Datagram)SetColor(c string){
	d.Tags[DATAGRAM_COLOR] = c
}

// SetVal
func (d *Datagram)SetVal(v float64){
	d.Fields[DATAGRAM_VAL] = v
}