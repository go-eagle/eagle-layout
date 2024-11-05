package method

import "gorm.io/gen"

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE username = @username
	GetUserByUsername(username string) (*gen.T, error)

	// SELECT * FROM @@table WHERE phone = @phone
	GetUserByPhone(phone string) (*gen.T, error)

	// SELECT * FROM @@table WHERE email = @email
	GetUserByEmail(email string) (*gen.T, error)
}
