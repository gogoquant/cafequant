package gobacktest

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
)

/*
const (
	MyDB          = "test"
	username      = "admin"
	password      = ""
	MyMeasurement = "cpu_usage"
)


func main() {
	conn := connInflux()
	fmt.Println(conn)

	//insert
	WritesPoints(conn)

	//获取10条数据并展示
	qs := fmt.Sprintf("SELECT * FROM %s LIMIT %d", MyMeasurement, 10)
	res, err := QueryDB(conn, qs)
	if err != nil {
		log.Fatal(err)
	}

	for i, row := range res[0].Series[0].Values {
		t, err := time.Parse(time.RFC3339, row[0].(string))
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(reflect.TypeOf(row[1]))
		valu := row[2].(json.Number)
		log.Printf("[%2d] %s: %s\n", i, t.Format(time.Stamp), valu)
	}
}
*/

type InfluxdbHandler interface {
	Connect() error
	Close() error
	QueryDB(cmd string) (res []client.Result, err error)
	WritesPoints(point *client.Point) error
}

// NewInfluxdb ...
func NewInfluxdb(host, user, pwd, database string) InfluxdbHandler {
	return &Influxdb{
		host:     host,
		user:     user,
		pwd:      pwd,
		database: database,
	}
}

type Influxdb struct {
	host     string
	user     string
	pwd      string
	database string
	c        client.Client
}

func (d *Influxdb) Connect() error {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     d.host,
		Username: d.user,
		Password: d.pwd,
	})
	if err != nil {
		return err
	}
	d.c = cli
	return nil
}

func (d *Influxdb) Close() error {
	return d.c.Close()
}

//query
func (d *Influxdb) QueryDB(cmd string) (res []client.Result, err error) {
	cli := d.c
	q := client.Query{
		Command:  cmd,
		Database: d.database,
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

//Insert
func (d *Influxdb) WritesPoints(point *client.Point) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  d.database,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	bp.AddPoint(point)

	if err := d.c.Write(bp); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
