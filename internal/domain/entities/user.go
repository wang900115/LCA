package entities

type User struct {
	ID       uint    `json:"id"`
	Username string  `json:"userName"`
	Password *string `json:"password"`
	NickName string  `json:"nickName"`
	FullName string  `json:"fullName"`
	LastName string  `json:"lastName"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Birth    int64   `json:"birth"`
	Status   *string `json:"status"`
}

type UserLogin struct {
	LastLogin  int64   `json:"lastLogin"`
	IPAddress  *string `json:"ipAddress"`
	DeviceInfo *string `json:"deviceInfo"`
}

type UserChannel struct {
	Role     string `json:"role"`
	LastJoin int64  `json:"lastJoin"`
}
