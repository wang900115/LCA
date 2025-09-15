package websocketevent

const (
	EVENT_USER_JOIN    = "join"
	EVENT_USER_LEAVE   = "leave"
	EVENT_USER_COMMENT = "speak"
	EVENT_USER_EDIT    = "edit"
	EVENT_USER_DELETE  = "delete"
)

const (
	EVENT_SYSTEM_FIX            = "fix"
	EVENT_SYSTEM_GLOBAL         = "global"
	EVENT_SYSTEM_LOCAL          = "local"
	EVENT_SYSTEM_CHANNEL_FIX    = "channel-fix"
	EVENT_SYSTEM_CHANNEL_DELETE = "channel-delete"
)
