package gobacktest

import (
	"errors"
	"github.com/influxdata/influxdb/client/v2"
	"gopkg.in/logger.v1"
	"time"
)

var (
	DATAGRAM_VAL    = "data_val"
	DATAGRAM_COLOR  = "data_color"
	DATAGRAM_SYMBOL = "data_symbol"
	DATAGRAM_TAG    = "data_tag"
)

// Datagram
type Datagram struct {
	Event

	uid    string                 `json:"uuid"`
	fields map[string]interface{} `json:"fields"`
	tags   map[string]string      `json:"tags"`
	time   time.Time              `json:"time"`
}

// NewDataGram
func NewDataGram() *Datagram {
	return &Datagram{
		fields: make(map[string]interface{}),
		tags:   make(map[string]string),
	}
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
	d.tags[key] = val
}

// Tags
func (d *Datagram) Tags() map[string]string {
	res := make(map[string]string)
	for key, val := range d.tags {
		res[key] = val
	}
	return res
}

// Tag
func (d *Datagram) Tag() string {
	return d.tags[DATAGRAM_TAG]
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
	status int

	eventCh chan EventHandler
}

// NewDataGramMaster ...
func NewDataGramMaster(m map[string]string) *DataGramMaster {
	host := m["influxdb_host"]
	user := m["influxdb_user"]
	pwd := m["influxdb_pwd"]
	database := m["influxdb_database"]
	return &DataGramMaster{
		eventCh: make(chan EventHandler, 20),
		status:  GobackStop,
		influx:  NewInfluxdb(host, user, pwd, database),
	}
}

// AddPoint ...
func (m *DataGramMaster) AddPoint(datagram DataGramEvent) error {
	// set id here
	//datagram.SetId("test")

	point, err := client.NewPoint(
		datagram.Id(),
		datagram.Tags(),
		datagram.Val(),
		datagram.Time(),
	)
	if err != nil {
		return err
	}
	return m.influx.WritesPoints(point)
}

// Connect ...
func (m *DataGramMaster) Connect() error {
	return m.influx.Connect()
}

// Close ...
func (m *DataGramMaster) Close() error {
	return m.influx.Close()
}

// Start ...
func (m *DataGramMaster) Start() error {
	if m.status == GobackRun {
		return errors.New("already running")
	}
	m.status = GobackRun
	/*
		if err := m.Close(); err != nil{
			log.Error("influxdb close fail")
			return err
		}
	*/
	go m.Run()

	return nil
}

// Stop ...
func (m *DataGramMaster) Stop() error {
	if m.status == GobackStop || m.status == GobackPending {
		return errors.New("already stop or pending")

	}
	m.status = GobackPending
	var cmd Cmd
	cmd.SetCmd("stop")
	m.eventCh <- &cmd
	/*
		time.Sleep(time.Millisecond *500)
		if err := m.Close(); err != nil{
			log.Error("influxdb close fail")
			return err
		}
	*/

	return nil
}

// Run ...
func (m *DataGramMaster) Run() error {
	for {

		e := <-m.eventCh

		switch event := e.(type) {

		case DataGramEvent:
			log.Infof("DataGram Get dataGram event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())
			if err := m.AddPoint(event); err != nil {
				log.Error("add Point fail")
				return err
			}

		case CmdEvent:
			log.Infof("DataGram Get cmd event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())
			m.status = GobackStop
			if err := m.influx.Close(); err != nil {
				log.Error("close influxdb fail")
			}
			return nil
		}
	}
}

// Status ...
func (m *DataGramMaster) Status() int {
	return m.status
}

// AddDataGram ...
func (m *DataGramMaster) AddDataGram(datagram DataGramEvent) error {
	m.eventCh <- datagram
	return nil
}
