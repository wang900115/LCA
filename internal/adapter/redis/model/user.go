package redismodel

const (
	REDIS_USER_CHANNEL_TABLE = "user-channel:"
	REDIS_USER_LOGIN_TABLE   = "user-login:"
)

const (
	REDIS_USER_CHANNEL_FIELD_ID       = "id"
	REDIS_USER_CHANNEL_FIELD_ROLE     = "role"
	REDIS_USER_CHANNEL_FIELD_LASTJOIN = "last_join"
)

const (
	REDIS_USER_LOGIN_FIELD_IPADDRESS  = "ip_address"
	REDIS_USER_LOGIN_FIELD_DEVICEINFO = "device_info"
	REDIS_USER_LOGIN_FIELD_LASTLOGIN  = "last_login"
)
