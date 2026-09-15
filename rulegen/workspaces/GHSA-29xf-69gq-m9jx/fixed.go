package main

		return
	}

	// Only owners can change the password of another owner.
	if apiKey.UserID != user.ID && slices.Contains(user.RBACRoles, rbac.RoleOwner().String()) {
		actingUser, err := api.Database.GetUserByID(ctx, apiKey.UserID)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, codersdk.Response{
				Message: "Internal error fetching acting user.",
				Detail:  err.Error(),
			})
			return
		}
		if !slices.Contains(actingUser.RBACRoles, rbac.RoleOwner().String()) {
			httpapi.Write(ctx, rw, http.StatusBadRequest, codersdk.Response{
				Message: "Only owners can change the password of an owner.",
			})
			return
		}
	}

	if !httpapi.Read(ctx, rw, r, &params) {
		return
	}
