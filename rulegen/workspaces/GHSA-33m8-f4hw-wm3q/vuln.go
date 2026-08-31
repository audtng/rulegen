package main


import (
	"fmt"
)

// Role is the type of a role.
	if len(create.Username) < 4 {
		return fmt.Errorf("username is too short, minimum length is 4")
	}
	if len(create.Password) < 4 {
		return fmt.Errorf("password is too short, minimum length is 4")
	}

	return nil
}
	OpenID       *string
}

type UserFind struct {
	ID *int `json:"id"`

		if err := json.NewDecoder(c.Request().Body).Decode(userPatch); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Malformatted patch user request").SetInternal(err)
		}

		if userPatch.Email != nil && *userPatch.Email != "" && !common.ValidateEmail(*userPatch.Email) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid email format")
		}

		if userPatch.Password != nil && *userPatch.Password != "" {
