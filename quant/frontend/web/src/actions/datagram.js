import * as actions from '../constants/actions';
import { Client } from 'hprose-js';

// List
function datagramListRequest() {
  return { type: actions.DATAGRAM_LIST_REQUEST };
}

function datagramListSuccess(list, col, mode) {
  return { type: actions.DATAGRAM_LIST_SUCCESS, list, col, mode };
}

function datagramListFailure(message) {
  return { type: actions.DATAGRAM_LIST_FAILURE, message };
}

export function DatagramList(traderId, mode) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(datagramListRequest());
    if (!cluster || !token) {
      dispatch(logListFailure('No authorization'));
      return;
    }

    console.log('call datagram show with id:' + traderId);
    const client = Client.create(`${cluster}/api`, { Datagram: ['List'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Datagram.List(traderId, mode, (resp) => {
      if (resp.success) {
        dispatch(datagramListSuccess(resp.data.list, resp.data.col, resp.data.mode));
      } else {
        dispatch(datagramListFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(datagramListFailure('Server error'));
      console.log('【Hprose】Log.List Error:', resp, err);
    });
  };
}

// Delete
function datagramDeleteRequest() {
  return { type: actions.DATAGRAM_DELETE_REQUEST };
}

function datagramDeleteSuccess() {
  return { type: actions.DATAGRAM_DELETE_SUCCESS};
}

function datagramDeleteFailure(message) {
  return { type: actions.DATAGRAM_DELETE_FAILURE, message };
}

export function DatagramDelete(traderId) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(datagramDeleteRequest());
    if (!cluster || !token) {
      dispatch(logListFailure('No authorization'));
      return;
    }

    console.log('call datagram delete with id:' + traderId);
    const client = Client.create(`${cluster}/api`, { Datagram: ['Delete'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Datagram.Delete(traderId, (resp) => {
      if (resp.success) {
        dispatch(datagramDeleteSuccess());
      } else {
        dispatch(datagramDeleteFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(datagramDeleteFailure('Server error'));
      console.log('【Hprose】Log.List Error:', resp, err);
    });
  };
}
