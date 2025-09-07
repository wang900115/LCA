package validator

type UserCreateRequest struct {
	Username string  `json:"userName" binding:"required"`
	Password *string `json:"password" binding:"required"`
	NickName string  `json:"nickName" binding:"required"`
	FullName string  `json:"fullName" binding:"required"`
	LastName string  `json:"lastName" binding:"required"`
	Email    string  `json:"email" binding:"required"`
	Phone    string  `json:"phone" binding:"required"`
	Birth    int64   `json:"birth" binding:"required"`
	Status   *string `json:"status" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    int64  `json:"login" binding:"required"`
}

type UserDeleteRequest struct {
	Password string `json:"password" binding:"required"`
}

type UserJoinRequest struct {
	ChannelID uint  `json:"channel" binding:"required"`
	JoinTime  int64 `json:"joinTime" binding:"required"`
}

type UserCommentRequest struct {
	Content string `json:"content" binding:"required"`
	// CommentTime int64  `json:"commentTime" binding:"required"`
}

type UserEditeRequest struct {
	NewContent string `json:"content" binding:"required"`
}

type UserRegainRequest struct {
	MessageID uint `json:"message" binding:"required"`
}
