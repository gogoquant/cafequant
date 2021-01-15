import * as actions from '../constants/actions';
import assign from 'lodash/assign';

const LOG_INIT = {
  loading: false,
  total: 0,
  list: [],
  message: ''
};

function log(state = LOG_INIT, action) {
  switch (action.type) {
    case actions.RESET_ERROR:
      return assign({}, state, {
        loading: false,
        message: ''
      });
    case actions.LOG_LIST_REQUEST:
      return assign({}, state, {
        loading: true
      });
    case actions.LOG_LIST_SUCCESS:
      return assign({}, state, {
        loading: false,
        total: action.total,
        list: action.list,
        data: state.data
      });
    case actions.LOG_LIST_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message
      });
    case actions.LOG_STATUS_REQUEST:
      return assign({}, state, {
        loading: true
      });
    case actions.LOG_STATUS_SUCCESS:
      return assign({}, state, {
        loading: false,
        total: state.total,
        list: state.list,
        data: action.data
      });
    case actions.LOG_STATUS_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message
      });
    default:
      return state;
  }
}

export default log;
