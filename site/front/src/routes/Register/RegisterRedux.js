import { Map } from 'immutable';
import { createAction, handleActions } from 'redux-actions';
import { UserService } from 'services';
import { toastUtils, userUtils } from 'utils';

const POST_REGISTER_TOKEN = 'routes/Register/POST_REGISTER_TOKEN';
const POST_REGISTER_TOKEN_SUCCESS =
  'routes/Register/POST_REGISTER_TOKEN_SUCCESS';
const POST_REGISTER_TOKEN_ERROR = 'routes/Register/POST_REGISTER_TOKEN_ERROR';

const initialState = Map({
  loading: false,
  userData: {},
  error: false,
});

const reducer = handleActions(
  {
    [POST_REGISTER_TOKEN]: state =>
      state
        .set('loading', true)
        .set('error', false)
        .set('userData', {}),

    [POST_REGISTER_TOKEN_SUCCESS]: (state, { payload }) =>
      state.set('loading', false).set('userData', payload),

    [POST_REGISTER_TOKEN_ERROR]: state =>
      state.set('loading', false).set('error', true),
  },
  initialState,
);

export const postAccessToken = createAction(POST_REGISTER_TOKEN);
export const postAccessTokenSuccess = createAction(POST_REGISTER_TOKEN_SUCCESS);
export const postAccessTokenError = createAction(POST_REGISTER_TOKEN_ERROR);

export const register = (email, passwd, name, callback) => async dispatch => {
  dispatch(postAccessToken());
  try {
    data = await UserService.register(email, passwd, name);
    toastUtils.success('注册成功');
    dispatch(postAccessTokenSuccess(email));
    callback();
  } catch (e) {
    dispatch(postAccessTokenError());
  }
};

export default reducer;
