RECOMMEND_WEIBO = "weibo"
RECOMMEND_HOTUSER = "hotuser"
RECOMMEND_SAMEFOLLOW = "samefollow"
RECOMMEND_SAMEINTEREST = "sameinterest"

RECOMMEND_RANKS = {
    RECOMMEND_WEIBO: 4,
    RECOMMEND_SAMEFOLLOW: 3,
    RECOMMEND_SAMEINTEREST: 2,
    RECOMMEND_HOTUSER: 1
}


def _compare_recommend_follow(rf1, rf2):
    rank1 = RECOMMEND_RANKS[rf1.get('recommend_type', RECOMMEND_HOTUSER)]
    rank2 = RECOMMEND_RANKS[rf2.get('recommend_type', RECOMMEND_HOTUSER)]
    type_flag = -cmp(rank1, rank2)
    if type_flag:
        return type_flag
    else:
        return -cmp(
            rf1.get('recommend_type_count', 0),
            rf2.get('recommend_type_count', 0))
