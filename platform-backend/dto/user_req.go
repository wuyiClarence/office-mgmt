package dto

type UserCreateReq struct {
	UserName        string  `json:"user_name" binding:"required"`
	UserDisplayName string  `json:"user_display_name" binding:"required"`
	Email           string  `json:"email"`
	PhoneNumber     string  `json:"phone_number"`
	Password        string  `json:"password"`
	RoleIDS         []int64 `json:"role_ids"`
}

type UserDeleteReq struct {
	UserIDs []int64 `json:"user_ids" binding:"required"`
}

type UserUpdateInfoReq struct {
	UserID          int64   `json:"user_id"`
	UserDisplayName string  `json:"user_display_name"`
	Email           string  `json:"email"`
	PhoneNumber     string  `json:"phone_number"`
	Password        string  `json:"password"`
	RoleIDS         []int64 `json:"role_ids"`
}

type UserResetPasswordReq struct {
	UserID   int64  `json:"user_id" binding:"required"`
	Password string `json:"password"`
}

type UserLoginApiReq struct {
	UserName string `form:"username" json:"username" binding:"required"`
	// Use DES-DCB encryption
	Password string `form:"password" json:"password" binding:"required"`
}

type UserLoginApiRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserUpdatePasswordApiReq struct {
	// Use DES-DCB encryption
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type UserRefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserRefreshTokenRes struct {
	AccessToken string `json:"access_token"`
}

type UserInfoRes struct {
	Name        string            `json:"name"`
	UserId      int64             `json:"user_id"`
	Permissions PermissionListRes `json:"permissions"`
}

type UserRoleInfo struct {
	RoleId   int64  `json:"role_id" gorm:"column:role_id"`
	RoleName string `json:"role_name" gorm:"column:role_name"`
}

type UserListItem struct {
	UserID                int64           `json:"user_id"`
	UserName              string          `json:"user_name"`
	UserDisplayName       string          `json:"user_display_name"`
	Email                 string          `json:"email"`
	PhoneNumber           string          `json:"phone_number"`
	CreateUserName        string          `json:"create_user_name"`
	CreateUserDisplayName string          `json:"create_user_display_name"`
	CreatedAt             string          `json:"created_at"`
	RoleInfos             []*UserRoleInfo `json:"role_infos"`
	Permissions           []Permission    `json:"permissions" gorm:"-"`
}
type UserListRes struct {
	Total int64           `json:"total"`
	List  []*UserListItem `json:"list"`
}
type UserListReq struct {
	ListReq
}
