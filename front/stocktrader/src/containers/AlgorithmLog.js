import { ResetError } from '../actions';
import { LogList, LogStatus } from '../actions/log';
import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { Button, Table, Tag, notification, Switch, Collapse } from 'antd';
const Panel = Collapse.Panel;

function PrefixInteger(num, m) {
  return (num + Array(m).join(0)).slice(0, m);
}

class Log extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      sync: false,
      syncTime: 5000,
      messageErrorKey: '',
      pagination: {
        pageSize: 12,
        current: 1,
        total: 0
      },
      filters: {}
    };
    this.reload = this.reload.bind(this);
    this.loadStatus = this.loadStatus.bind(this);
    this.handleSync = this.handleSync.bind(this);
    this.handleDiagram = this.handleDiagram.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { dispatch } = this.props;
    const { messageErrorKey, pagination } = this.state;
    const { trader, log } = nextProps;

    if (!trader.cache.name) {
      browserHistory.push('/algorithm');
    }

    if (!messageErrorKey && log.message) {
      this.setState({
        messageErrorKey: 'logError'
      });
      notification['error']({
        key: 'logError',
        message: 'Error',
        description: String(log.message),
        onClose: () => {
          if (this.state.messageErrorKey) {
            this.setState({ messageErrorKey: '' });
          }
          dispatch(ResetError());
        }
      });
    }
    pagination.total = log.total;
    this.setState({ pagination });
  }

  componentWillMount() {
    this.filters = {};
    this.reload();
  }

  componentWillUnmount() {
    notification.destroy();
    this.syncTimer && clearTimeout(this.syncTimer);
    this.syncTimer = null;
  }

  reload() {
    const { pagination } = this.state;
    const { trader, dispatch } = this.props;

    dispatch(LogList(trader.cache, pagination, this.filters));
  }

  loadStatus() {
    const { trader, dispatch } = this.props;
    dispatch(LogStatus(trader.cache));
  }

  handleTableChange(newPagination, filters) {
    const { pagination } = this.state;

    pagination.current = newPagination.current;
    this.filters = filters;
    this.setState({ pagination });
    this.reload();
  }

  handleCancel() {
    browserHistory.push('/algorithm');
  }

  handleDiagram() {
    const { trader } = this.props;
    const cluster = localStorage.getItem('cluster');
    window.location.href = cluster + '/' + String(trader.cache.id) + '.html';
  }

  handleSync(sync) {
    if (sync === this.state.sync) {
      return;
    }
    if (sync === true) {
      this.state.pagination.current = 0;
      this.syncTimer = setInterval(() => this.reload(), this.state.syncTime);
    } else {
      this.syncTimer && clearTimeout(this.syncTimer);
      this.syncTimer = null;
    }
    this.state.sync = sync;
  }

  render() {
    const { pagination } = this.state;
    const { log } = this.props;
    const colors = {
      INFO: '#A9A9A9',
      ERROR: '#F50F50',
      PROFIT: '#4682B4',
      CANCEL: '#5F9EA0'
    };
    const columns = [
      {
        width: 160,
        title: 'Time',
        dataIndex: 'time',
        render: v =>
          v.toLocaleString() +
          '.' +
          PrefixInteger(v.getMilliseconds(), 3).toString()
      },
      {
        width: 100,
        title: 'Exchange',
        dataIndex: 'exchangeType',
        render: v => <Tag color={v === 'global' ? '' : '#00BFFF'}>{v}</Tag>
      },
      {
        width: 100,
        title: 'Type',
        dataIndex: 'type',
        render: v => <Tag color={colors[v] || '#00BFFF'}>{v}</Tag>
      },
      {
        title: 'Price',
        dataIndex: 'price',
        width: 100
      },
      {
        width: 100,
        title: 'Amount',
        dataIndex: 'amount'
      },
      {
        title: 'Message',
        dataIndex: 'message'
      }
    ];
    console.log(log);
    return (
      <div>
        <div className="table-operations">
          <Button type="primary" onClick={this.reload}>
            Reload
          </Button>
          <Button type="ghost" onClick={this.handleCancel}>
            Back
          </Button>
          <Button type="ghost" onClick={this.handleDiagram}>
            Diagram
          </Button>
          <Switch
            checkedChildren="sync"
            unCheckedChildren=""
            onChange={this.handleSync}
          />
        </div>
        <Collapse onChange={this.loadStatus}>
          <Panel header="status" key="1" extra={log.data}>
            <p>{log.data}</p>
          </Panel>
        </Collapse>
        <Table
          rowKey="id"
          columns={columns}
          dataSource={log.list}
          pagination={pagination}
          loading={log.loading}
          onChange={this.handleTableChange}
        />
      </div>
    );
  }
}

const mapStateToProps = state => ({
  user: state.user,
  trader: state.trader,
  log: state.log
});

export default connect(mapStateToProps)(Log);
