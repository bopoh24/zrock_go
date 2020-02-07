package model

// UserFixture returns user for test
func UserFixture() *User {
	return &User{
		Email:     "user@example.org",
		Nickname:  "nick",
		FirstName: "FirstName",
		LastName:  "LastName",
		Password:  "Password1",
	}
}
