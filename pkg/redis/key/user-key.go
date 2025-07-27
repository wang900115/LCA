package rediskey

const (
	// key: user_id
	REDIS_TABLE_USER = "user:"

	// key: user_id
	REDIS_SET_USER_CHANNELS = "user:channels:"

	// ttl short
	// key: user_id + token_id
	REDIS_STRING_ACCESS_TOKEN = "access-token:"

	// ttl long
	// key: user_id + token_id
	REDIS_STRING_REFRESH_TOKEN = "refresh-token:"

	// REDIS_LIST_BLACK_USER = "black-list:"
)

const (
	REDIS_FIELD_USER_USERNAME  = "username"
	REDIS_FIELD_USER_ROLE      = "role"
	REDIS_FIELD_USER_NICKNAME  = "nickname"
	REDIS_FIELD_USER_FIRSTNAME = "firstname"
	REDIS_FIELD_USER_LASTNAME  = "lastname"
	REDIS_FIELD_USER_BIRTH     = "birth"

	REDIS_FIELD_USER_CREATEDAT = "created_at"
	REDIS_FIELD_USER_UPDATEDAT = "updated_at"
)

const (
	REDIS_FIELD_TOKEN_ID       = "id"
	REDIS_FIELD_TOKEN_DEVICE   = "device"
	REDIS_FIELD_TOKEN_IP       = "ip"
	REDIS_FIELD_TOKEN_UA       = "ua"
	REDIS_FIELD_TOKEN_ISSUEDAT = "issued_at"
)
