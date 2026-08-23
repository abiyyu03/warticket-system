package author

type (
	RegisterRequest struct {
		Name     string
		Email    string
		Password string
	}

	LoginRequest struct {
		Email    string
		Password string
	}

	RefreshRequest struct {
		RefreshToken string
	}

	// AuthResponse dikembalikan login & refresh.
	AuthResponse struct {
		TokenType    string `json:"token_type"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"` // detik, umur access token
		Role         string `json:"role"`
	}
)
