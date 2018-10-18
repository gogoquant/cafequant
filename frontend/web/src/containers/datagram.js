import React from 'react';
import { connect } from 'react-redux';

import {Button, notification} from 'antd';
import {browserHistory} from 'react-router';
import {ResetError} from '../actions';
import {DatagramList} from '../actions/datagram';
import {
  Chart,
  Geom,
  Axis,
  Tooltip,
  Legend,
} from 'bizcharts';
import DataSet from '@antv/data-set';

class Datagram extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      messageErrorKey: '',
    };
    this.reload = this.reload.bind(this);
    this.handleCancel = this.handleCancel.bind(this);
  }

  componentWillReceiveProps(nextProps) {

    const { dispatch } = this.props;
    const { messageErrorKey } = this.state;
    const { datagram } = nextProps;

    // let id = this.props.params.id;
    // browserHistory.push('/datagram/' + id);

    if (!messageErrorKey && datagram.message) {
      this.setState({
        messageErrorKey: 'logError',
      });
      notification['error']({
        key: 'logError',
        message: 'Error',
        description: String(datagram.message),
        onClose: () => {
          if (this.state.messageErrorKey) {
            this.setState({ messageErrorKey: '' });
          }
          dispatch(ResetError());
        },
      });
    }
  }

  componentWillMount() {
    this.reload();
  }

  componentWillUnmount() {
    notification.destroy();
  }

  reload() {
    const { dispatch } = this.props;
    let id = this.props.params.id;
    dispatch(DatagramList(id));
  }

  handleCancel() {
    browserHistory.push('/algorithm');
  }

  render() {
    const { datagram } = this.props;
    console.log(datagram);

    let datas = [];
    let items = datagram.list;
    let col = datagram.col;
    let mode = this.props.params.mode;

    for (let i = 0;i < items.length;i++) {
      let data_map = Object.create(null);
      for (let j = 0;j < col.length;j++) {
        let key = col[j];
        if (key === 'time') {
          let date = new Date(items[i].fields[key]);
          if (mode === 'hour') {
            // set year key but used  as hour
            data_map['year'] = date.getHours();
          } if (mode === 'minute') {
            data_map['time'] = date.getMonth().toString() + '/' + date.getDay().toString() + ' ' + date.getHours().toString() + ':' + date.getMinutes().toString() ;
          } else {
            let mTime = date.toUTCString();
            data_map[key] = mTime;
          }
          continue;
        }
        let value = items[i].fields[key];
        let float_value = parseFloat(value).toFixed(2);
        console.log('convert ' + value + ' to ' + float_value);
        data_map[key] = float_value;
        // data_map[key] = i;
      }
      datas.push(data_map);
    }

    const ds = new DataSet();
    const dv = ds.createView().source(datas);

    let new_col = [];
    for (let i = 0;i < col.length;i++) {
      if (col[i] !== 'time') {
        new_col.push(col[i]);
      }
    }

    dv.transform({
      type: 'fold',
      fields: new_col,
      // 展开字段集
      key: 'symbol',
      // key字段
      value: 'amount' // value字段
    });
    console.log(dv);
    const cols = {
      hour: {
        range: [0, 1439]
      }
    };

    return (
      <div>
        <div className="table-operations">
          <Button type="primary" onClick={this.reload}>Reload</Button>
          <Button type="ghost" onClick={this.handleCancel}>Back</Button>
        </div>
        <div id="mountNode">
          <div>
            <Chart height={600} data={dv} scale={cols} forceFit>
              <Legend />
              <Axis name="time" />
              <Axis
                name="amount"
                label={{
                  formatter: val => `( ${val})`
                }}
              />
              <Tooltip
                crosshairs={{
                  type: 'y'
                }}
              />
              <Geom
                type="line"
                position="time*amount"
                size={2}
                color={'symbol'}
              />
              <Geom
                type="point"
                position="time*amount"
                size={4}
                shape={'circle'}
                color={'symbol'}
                style={{
                  stroke: '#fff',
                  lineWidth: 1
                }}
              />
            </Chart>
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  datagram: state.datagram,
});

export default connect(mapStateToProps)(Datagram);
