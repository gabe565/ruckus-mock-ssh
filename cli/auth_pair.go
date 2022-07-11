package cli

type AuthPair struct {
	Username, Password string
}

func (a AuthPair) Validate(username, password string) bool {
	return username == a.Username && password == a.Password
}
