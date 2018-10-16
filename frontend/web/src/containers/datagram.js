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
    browserHistory.push('/datagram/' + id);
  }

  render() {
    const { datagram } = this.props;
    console.log(datagram);

    let datas = [];
    let items = datagram.list;
    let col = datagram.col;
    for (let i = 0;i < items.length;i++) {
      let data = new Map();
      for (let j = 0;j < col.length;j++) {
        let key = col[j];
        data[key] = items[i].fields[key];
      }
      datas.push(data);
    }

    const ds = new DataSet();
    const dv = ds.createView().source(datas);
    dv.transform({
      type: 'fold',
      fields: col,
      // 展开字段集
      key: 'symbol',
      // key字段
      value: 'amount' // value字段
    });
    console.log(dv);
    const cols = {
      timestamp: {
        range: [0, 100]
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
            <Chart height={400} data={dv} scale={cols} forceFit>
              <Legend />
              <Axis name="month" />
              <Axis
                name="temperature"
                label={{
                  formatter: val => `${val}°C`
                }}
              />
              <Tooltip
                crosshairs={{
                  type: 'y'
                }}
              />
              <Geom
                type="line"
                position="month*temperature"
                size={2}
                color={'city'}
              />
              <Geom
                type="point"
                position="month*temperature"
                size={4}
                shape={'circle'}
                color={'city'}
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
