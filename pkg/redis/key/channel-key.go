package rediskey

const (
	// key: channel_id
	REDIS_TABLE_CHANNEL = "channel:"

	// key: channel_id values: user_id
	REDIS_SET_CHANNEL_USER = "channel-user:"

	// key: channel_id values: message_id
	REDIS_SET_CHANNEL_MESSAGE = "channel-message:"

	// key: channel_id + user_id values: message_id
	REDIS_LIST_CHANNEL_USER_MESSAGE = "channel-user-message:"
)

const (
	REDIS_FIELD_CHANNEL_NAME = "name"
	REDIS_FIELD_CHANNEL_TYPE = "type"

	REDIS_FIELD_CHANNEL_CREATEDAT = "created_at"
	REDIS_FIELD_CHANNEL_UPDATEDAT = "updated_at"
)
