package method

import "gorm.io/gen"

// Querier is the custom query API implemented by gorm/gen.
type Querier interface {
	// SELECT * FROM @@table WHERE username = @username
	GetUserByUsername(username string) (*gen.T, error)

	// SELECT * FROM @@table WHERE phone = @phone
	GetUserByPhone(phone string) (*gen.T, error)

	// SELECT * FROM @@table WHERE email = @email
	GetUserByEmail(email string) (*gen.T, error)
}
