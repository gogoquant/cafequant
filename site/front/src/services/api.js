// request api base url
 export const BASE_URL = 'http://127.0.0.1:8080';
// export const BASE_URL = '/proxy';

// cors support
export const CORS = true;

export const TopicsApi = {
  topic: '/topic/:id',
  topics: '/topics',
  topicsUpdate: '/topics/update',
};

export const UserApi = {
  user: '/user/:name',
  accessToken: '/accesstoken',
  register: '/api/user/register',
};

export const ReplyApi = {
  reply: '/topic/:topic_id/replies',
  upReply: '/reply/:reply_id/ups',
};

export const TopicCollectApi = {
  userCollect: '/topic_collect/:name',
  collectTopic: '/topic_collect/collect',
  cancelCollect: '/topic_collect/de_collect',
};
