package model

// UserFixture returns user for test
func UserFixture() *User {
	return &User{
		Email:     "user@example.org",
		Nick:      "nick",
		FirstName: "FirstName",
		LastName:  "LastName",
		Password:  "Password1",
	}
}
