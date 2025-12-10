package user

type UserFormatter struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Occupation     string `json:"occupation"`
	Email          string `json:"email"`
	AvatarFileName string `json:"avatar_file_name"`
	Role           string `json:"role"`
	Token          string `json:"token"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:             user.ID,
		Name:           user.Name,
		Occupation:     user.Occupation,
		Email:          user.Email,
		AvatarFileName: user.AvatarFileName,
		Role:           user.Role,
		Token:          token,
	}

	return formatter
}
