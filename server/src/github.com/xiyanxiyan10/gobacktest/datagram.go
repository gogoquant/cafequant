package gobacktest

import (
	"github.com/influxdata/influxdb/client/v2"
	"time"
)

var (
	DATAGRAM_VAL         = "data_val"
	DATAGRAM_COLOR       = "data_color"
	DATAGRAM_SYMBOL       = "data_symbol"
	DATAGRAM_USER_PREFIX = "data_user:"
)

// Datagram
type Datagram struct {
	Event

	uid    string                 `json:"uuid"`
	fields map[string]interface{} `json:"fields"`
	tags   map[string]string      `json:"tags"`
	time   time.Time              `json:"time"`
}

// SetId
func (d *Datagram) SetId(uid string) {
	d.uid = uid
}

// Id
func (d *Datagram) Id() string {
	return d.uid
}

// SetTag
func (d *Datagram) SetTag(key, val string) {
	d.tags[DATAGRAM_USER_PREFIX+key] = val
}

// Tags
func (d *Datagram) Tags() map[string]string {
	res := make(map[string]string)
	for key, val := range d.tags {
		res[key] = val
	}
	return res
}

// SetColor
func (d *Datagram) SetColor(c string) {
	d.tags[DATAGRAM_COLOR] = c
}

// Color
func (d *Datagram) Color(c string) string {
	return d.tags[DATAGRAM_COLOR]
}

// SetVal
func (d *Datagram) SetVal(key string, v interface{}) {
	d.fields[key] = v
}

// Val
func (d *Datagram) Val() map[string]interface{} {
	res := make(map[string]interface{})
	for key, val := range d.fields {
		res[key] = val
	}
	return res
}

// DataGramMaster ...
type DataGramMaster struct {
	influx InfluxdbHandler
}

// AddPoint ...
func (m *DataGramMaster)AddPoint(datagram DataGramEvent) error{
	point, err := client.NewPoint(
		datagram.Id(),
		datagram.Tags(),
		datagram.Val(),
		datagram.Time(),
	)
	if err != nil{
		return err
	}
	return m.influx.WritesPoints(point)
}

// Connect ...
func (m *DataGramMaster)Connect() error{
	return m.influx.Connect()
}

// Close ...
func (m *DataGramMaster)Close() error{
	return m.influx.Close()
}
