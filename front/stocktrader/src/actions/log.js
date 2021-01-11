import * as actions from '../constants/actions';
import { Client } from 'hprose-js';

// List

function logListRequest() {
  return { type: actions.LOG_LIST_REQUEST };
}

function logListSuccess(total, list) {
  return { type: actions.LOG_LIST_SUCCESS, total, list };
}

function logListFailure(message) {
  return { type: actions.LOG_LIST_FAILURE, message };
}

// Status ...

function logStatusRequest() {
  return { type: actions.LOG_STATUS_REQUEST };
}

function logStatusSuccess(data) {
  return { type: actions.LOG_STATUS_SUCCESS, data };
}

function logStatusFailure(message) {
  return { type: actions.LOG_LSTATUS_FAILURE, message };
}

export function LogStatus(trader) {
  return (dispatch) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(logStatusRequest());
    if (!cluster || !token) {
      dispatch(logListFailure('No authorization'));
      return;
    }

    const client = Client.create(`${cluster}/api`, { Log: ['Status'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Log.Status(trader, (resp) => {
      if (resp.success) {
        dispatch(logStatusSuccess(resp.data));
      } else {
        dispatch(logStatusFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(logStatusFailure('Server error'));
      console.log('【Hprose】Log.Status Error:', resp, err);
    });
  };
}

export function LogList(trader, pagination, filters) {
  return (dispatch) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(logListRequest());
    if (!cluster || !token) {
      dispatch(logListFailure('No authorization'));
      return;
    }

    const client = Client.create(`${cluster}/api`, { Log: ['List'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Log.List(trader, pagination, filters, (resp) => {
      if (resp.success) {
        dispatch(logListSuccess(resp.data.total, resp.data.list));
      } else {
        dispatch(logListFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(logListFailure('Server error'));
      console.log('【Hprose】Log.List Error:', resp, err);
    });
  };
}
