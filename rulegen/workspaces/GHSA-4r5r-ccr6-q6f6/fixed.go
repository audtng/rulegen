package main

			return
		}

		if !v.CanPerformActions() || v.User.GlobalRole == nil || *v.User.GlobalRole != fleet.RoleAdmin {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
