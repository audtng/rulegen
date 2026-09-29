package rules

type Profile struct {
	Email string
	EmailVerified bool
}

func testStruct(email string) Profile {
	return Profile{
		Email: email,
		EmailVerified: email != "",
	}
}

func testAssign(p *Profile, email string) {
	p.EmailVerified = email != ""
}
