package model

// UserFixture returns user for test
func UserFixture() *User {
	return &User{
		Email:                 "user@example.org",
		EmailVerified:         true,
		EmailVerificationCode: "some_verification_code",
		Nickname:              "nick",
		FirstName:             "FirstName",
		LastName:              "LastName",
		Password:              "Password1",
	}
}
