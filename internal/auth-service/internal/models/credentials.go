package models

type LoginCredentials struct {
	Login    string
	Password string
	Context  InputContext
}

type LoginResult struct {
	User    User
	Session Session
	Device  Device
	Tokens  TokenPair
}

type RefreshTokenResult struct {
	Tokens  TokenPair
	User    User
	Session Session
}
