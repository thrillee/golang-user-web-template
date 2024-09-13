package oauth

type OAuthUser struct {
	Email    string
	Phone    string
	Fullname string
	ImageURL string
}

type OAuth interface {
	HandleOAuth(code string) (*OAuthUser, error)
}
