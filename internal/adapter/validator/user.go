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
