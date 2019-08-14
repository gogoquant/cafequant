export const BASE_URL = 'https://cnodejs.org/api/v1';

export const TopicsApi = {
  topic: '/topic/:id',
  topics: '/topics',
  topicsUpdate: '/topics/update',
};

export const UserApi = {
  user: '/user/:name',
  accessToken: '/accesstoken',
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
