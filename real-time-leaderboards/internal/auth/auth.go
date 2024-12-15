package auth

// User sends username and password
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User received access and refresh token back, if authentication is succesful
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}