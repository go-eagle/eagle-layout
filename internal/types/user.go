package types

// User include user base info and user profile
type User struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	LoginAt   int64  `json:"login_at"` // login time for last times
	Status    int32  `json:"status"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Gender    int32  `json:"gender"`
	Birthday  string `json:"birthday"`
	Bio       string `json:"bio"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
