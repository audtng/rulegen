package main


import (
	"fmt"

	"github.com/usememos/memos/common"
)

// Role is the type of a role.
	if len(create.Username) < 4 {
		return fmt.Errorf("username is too short, minimum length is 4")
	}
	if len(create.Username) > 32 {
		return fmt.Errorf("username is too long, maximum length is 32")
	}
	if len(create.Password) < 4 {
		return fmt.Errorf("password is too short, minimum length is 4")
	}
	if len(create.Nickname) > 64 {
		return fmt.Errorf("nickname is too long, maximum length is 64")
	}
	if create.Email != "" {
		if len(create.Email) > 256 {
			return fmt.Errorf("email is too long, maximum length is 256")
		}
		if common.ValidateEmail(create.Email) {
			return fmt.Errorf("invalid email format")
		}
	}

	return nil
}
	OpenID       *string
}

func (patch UserPatch) Validate() error {
	if patch.Username != nil && len(*patch.Username) < 4 {
		return fmt.Errorf("username is too short, minimum length is 4")
	}
	if patch.Username != nil && len(*patch.Username) > 32 {
		return fmt.Errorf("username is too long, maximum length is 32")
	}
	if patch.Password != nil && len(*patch.Password) < 4 {
		return fmt.Errorf("password is too short, minimum length is 4")
	}
	if patch.Nickname != nil && len(*patch.Nickname) > 64 {
		return fmt.Errorf("nickname is too long, maximum length is 64")
	}
	if patch.Email != nil {
		if len(*patch.Email) > 256 {
			return fmt.Errorf("email is too long, maximum length is 256")
		}
		if common.ValidateEmail(*patch.Email) {
			return fmt.Errorf("invalid email format")
		}
	}

	return nil
}

type UserFind struct {
	ID *int `json:"id"`

		if err := json.NewDecoder(c.Request().Body).Decode(userPatch); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Malformatted patch user request").SetInternal(err)
		}
		if err := userPatch.Validate(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid user patch format.").SetInternal(err)
		}

		if userPatch.Password != nil && *userPatch.Password != "" {
