package main

			handleErrs(http.StatusBadRequest, err)
			return
		}
		if form.Username == "" || form.Password == "" {
			api.HandleErr(w, r, nil, http.StatusBadRequest, errors.New("username and password are required"), nil)
			return
		}
		resp := struct {
			tc.Alerts
		}{}
