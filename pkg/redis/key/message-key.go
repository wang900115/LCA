package rediskey

const (
	// key: message_id
	REDIS_TABLE_MESSAGE = "message:"
)

const (
	REDIS_FIELD_MESSAGE_CHANNELID = "channel_id"
	REDIS_FIELD_MESSAGE_USERID    = "user_id"
	REDIS_FIELD_MESSAGE_MSGTYPE   = "msg_type"
	REDIS_FIELD_MESSAGE_STATUS    = "status"
	REDIS_FIELD_MESSAGE_REPLYTOID = "reply_to_id"
	REDIS_FIELD_MESSAGE_CONTENT   = "content"
	REDIS_FIELD_MESSAGE_ATTACHURL = "attach_url"

	REDIS_FIELD_MESSAGE_CREATEDAT = "created_at"
	REDIS_FIELD_MESSAGE_UPDATEDAT = "updated_at"
)
