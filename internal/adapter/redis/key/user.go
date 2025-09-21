package rediskey

// key: user_id --> login status
const (
	REDIS_USER_LOGIN_TABLE            = "user-login:"
	REDIS_USER_LOGIN_FIELD_IPADDRESS  = "ip_address"
	REDIS_USER_LOGIN_FIELD_DEVICEINFO = "device_info"
	REDIS_USER_LOGIN_FIELD_LASTLOGIN  = "last_login"
)

// key: user_id:channel_id --> channel status
const (
	REDIS_USER_CHANNEL_TABLE          = "user-channel:"
	REDIS_USER_CHANNEL_FIELD_ROLE     = "role"
	REDIS_USER_CHANNEL_FIELD_LASTJOIN = "last_join"
)

// key: user_id:event_id --> event status
const (
	REDIS_USER_EVENT_TABLE               = "user-event:"
	REDIS_USER_EVENT_FIELD_LASTPARTICATE = "particate"
	REDIS_USER_EVENT_FIELD_ROLE          = "role"
)
