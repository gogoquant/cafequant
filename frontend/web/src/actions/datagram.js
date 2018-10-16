import * as actions from '../constants/actions';
import { Client } from 'hprose-js';

// List
function datagramListRequest() {
  return { type: actions.DATAGRAM_LIST_REQUEST };
}

function datagramListSuccess(list, col) {
  return { type: actions.DATAGRAM_LIST_SUCCESS, list, col };
}

function datagramListFailure(message) {
  return { type: actions.DATAGRAM_LIST_FAILURE, message };
}

export function DatagramList(traderId) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(datagramListRequest());
    if (!cluster || !token) {
      dispatch(logListFailure('No authorization'));
      return;
    }

    console.log('call datagram with id:' + traderId);
    const client = Client.create(`${cluster}/api`, { Datagram: ['List'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Datagram.List(traderId, (resp) => {
      if (resp.success) {
        dispatch(datagramListSuccess(resp.data.list, resp.data.col));
      } else {
        dispatch(datagramListFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(datagramListFailure('Server error'));
      console.log('【Hprose】Log.List Error:', resp, err);
    });
  };
}
