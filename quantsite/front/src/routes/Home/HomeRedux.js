import { Map } from 'immutable';
import { createAction, handleActions } from 'redux-actions';
import { TopicsService } from 'services';

const LOAD_TOPICS = 'routes/Home/LOAD_TOPICS';
const LOAD_TOPICS_SUCCESS = 'routes/Home/LOAD_TOPICS_SUCCESS';
const LOAD_TOPICS_ERROR = 'routes/Home/LOAD_TOPICS_ERROR';

const LOAD_MORE_TOPICS = 'routes/Home/LOAD_MORE_TOPICS';
const LOAD_MORE_TOPICS_SUCCESS = 'routes/Home/LOAD_MORE_TOPICS_SUCCESS';
const LOAD_MORE_TOPICS_ERROR = 'routes/Home/LOAD_MORE_TOPICS_ERROR';

const LOAD_TOPICS_FINISH = 'routes/Home/LOAD_TOPICS_FINISH';

const SAVE_SCROLL_HEIGHT = 'routes/Home/SAVE_SCROLL_HEIGHT';

const PAGE_SIZE = 10;

const initialState = Map({
  loading: false,
  loadingMore: false,
  hasMore: true,
  error: false,
  topicsData: [],
  scrollHeight: 0,
});

const reducer = handleActions(
  {
    [LOAD_TOPICS]: state =>
      state
        .set('loading', true)
        .set('hasMore', true)
        .set('error', false),

    [LOAD_TOPICS_SUCCESS]: (state, { payload }) =>
      state.set('loading', false).set('topicsData', payload),

    [LOAD_TOPICS_ERROR]: state =>
      state.set('loading', false).set('error', true),

    [LOAD_MORE_TOPICS]: state =>
      state.set('loadingMore', true).set('error', false),

    [LOAD_MORE_TOPICS_SUCCESS]: (state, { payload }) =>
      state
        .set('loadingMore', false)
        .update('topicsData', list => list.concat(payload)),

    [LOAD_MORE_TOPICS_ERROR]: state =>
      state.set('loadingMore', false).set('error', true),

    [LOAD_TOPICS_FINISH]: state => state.set('hasMore', false),

    [SAVE_SCROLL_HEIGHT]: (state, { payload }) =>
      state.set('scrollHeight', payload),
  },
  initialState,
);

export const loadTopics = createAction(LOAD_TOPICS);
export const loadTopicsSuccess = createAction(LOAD_TOPICS_SUCCESS);
export const loadTopicsError = createAction(LOAD_TOPICS_ERROR);

export const loadMoreTopics = createAction(LOAD_MORE_TOPICS);
export const loadMoreTopicsSuccess = createAction(LOAD_MORE_TOPICS_SUCCESS);
export const loadMoreTopicsError = createAction(LOAD_MORE_TOPICS_ERROR);

export const loadTopicsFinish = createAction(LOAD_TOPICS_FINISH);

export const saveScrollHeight = createAction(SAVE_SCROLL_HEIGHT);

export const getTopicsData = (tab, page, callback) => async dispatch => {
  dispatch(loadTopics());
  try {
    const {
      data: { data },
    } = await TopicsService.getTopics(tab, page, PAGE_SIZE);
    dispatch(loadTopicsSuccess(data));
    if (data.length < PAGE_SIZE) {
      dispatch(loadTopicsFinish());
    }
    callback();
  } catch (e) {
    dispatch(loadTopicsError());
  }
};

export const getMoreTopicsData = (tab, page, callback) => async dispatch => {
  dispatch(loadMoreTopics());
  try {
    const {
      data: { data },
    } = await TopicsService.getTopics(tab, page, PAGE_SIZE);
    dispatch(loadMoreTopicsSuccess(data));
    if (data.length < PAGE_SIZE) {
      dispatch(loadTopicsFinish());
    }
    callback();
  } catch (e) {
    dispatch(loadMoreTopicsError());
  }
};

export default reducer;
