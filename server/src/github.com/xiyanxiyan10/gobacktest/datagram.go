package gobacktest

import (
	"errors"
	"github.com/influxdata/influxdb/client/v2"
	"gopkg.in/logger.v1"
	"time"
)

const (
	DATAGRAM_VAL    = "val"
	DATAGRAM_TAG    = "tag"
	DATAGRAM_SYMBOL = "symbol"
)

// Datagram
type Datagram struct {
	Event

	uid    string                 `json:"uuid"`
	fields map[string]interface{} `json:"fields"`
	tags   map[string]string      `json:"tags"`
}

// DatagramInfo
type DatagramInfo struct {
	Symbol interface{}
	Timestamp interface{}
	Fields map[string]interface{}
}

// NewDataGram
func NewDataGram() DataGramEvent {
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


// SetVal
func (d *Datagram) SetFields(f map[string]interface{}) {
	d.fields = f
}

// Fields
func (d *Datagram) Fields() map[string]interface{} {
	vals := make(map[string]interface{})
	for key, val  := range d.fields{
		vals[key]=val
	}
	return vals
}

// Tags
func (d *Datagram) Tags() map[string]string {
	vals := make(map[string]string)
	for key, val  := range d.tags{
		vals[key]=val
	}
	return vals
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

// QueryDB
func (m *DataGramMaster) QueryDB(cmd string) (infos []DatagramInfo, table []string, err error){
	data, err := m.influx.QueryDB(cmd)
	if err != nil{
		return
	}
	if len(data) == 0{
		return
	}
	if len(data[0].Series) == 0{
		return
	}
	keyMap := make(map[string]int)
	for key, val := range data[0].Series[0].Columns{
		keyMap[val] = key
	}
	for key, _ := range keyMap{
		if key == "time" || key == "symbol"{
			continue
		}
		table = append(table, key)
	}

	for _, val := range data[0].Series[0].Values{
		var info DatagramInfo
		info.Fields = make(map[string]interface{})
		for name, idx := range keyMap {
			if name == "time" {
				info.Timestamp = val[idx]
				continue
			}
			if name == "symbol"{
				info.Symbol = val[idx]
				continue
			}
			info.Fields[name] = val[idx]
		}
		infos = append(infos, info)
	}
	return
}

// AddPoint ...
func (m *DataGramMaster) AddPoint(datagram DataGramEvent) error {

	fields := datagram.Fields()
	tags := datagram.Tags()
	fields[DATAGRAM_SYMBOL] = datagram.Symbol()
	tags[DATAGRAM_SYMBOL] = datagram.Symbol()

	point, err := client.NewPoint(
		datagram.Id(),
		tags,
		fields,
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

	time.Sleep(time.Millisecond * 500)
	if err := m.Close(); err != nil {
		log.Error("influxdb close fail")
		return err
	}

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
