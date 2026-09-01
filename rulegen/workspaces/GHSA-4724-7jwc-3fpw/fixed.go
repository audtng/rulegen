package main

	return &createSnapshotResponse, nil
}

func createOriginalDashboardURL(appURL string, cmd *dashboardsnapshots.CreateDashboardSnapshotCommand) (string, error) {
	dashUID := cmd.Dashboard.Get("uid").MustString("")
	if ok := util.IsValidShortUID(dashUID); !ok {
		return "", fmt.Errorf("invalid dashboard UID")
	}

	return fmt.Sprintf("/d/%v", dashUID), nil
}

// swagger:route POST /snapshots snapshots createDashboardSnapshot
//
// When creating a snapshot using the API, you have to provide the full dashboard payload including the snapshot data. This endpoint is designed for the Grafana UI.
		cmd.Name = "Unnamed snapshot"
	}

	var snapshotUrl string
	cmd.ExternalUrl = ""
	cmd.OrgId = c.OrgID
	cmd.UserId = c.UserID
	originalDashboardURL, err := createOriginalDashboardURL(hs.Cfg.AppURL, &cmd)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Invalid app URL", err)
	}

	if cmd.External {
		if !setting.ExternalEnabled {
			return nil
		}

		snapshotUrl = response.Url
		cmd.Key = response.Key
		cmd.DeleteKey = response.DeleteKey
		cmd.ExternalUrl = response.Url

		metrics.MApiDashboardSnapshotExternal.Inc()
	} else {
		cmd.Dashboard.SetPath([]string{"snapshot", "originalUrl"}, originalDashboardURL)

		if cmd.Key == "" {
			var err error
			cmd.Key, err = util.GetRandomString(32)
			}
		}

		snapshotUrl = setting.ToAbsUrl("dashboard/snapshot/" + cmd.Key)

		metrics.MApiDashboardSnapshotCreate.Inc()
	}
	c.JSON(http.StatusOK, util.DynMap{
		"key":       cmd.Key,
		"deleteKey": cmd.DeleteKey,
		"url":       snapshotUrl,
		"deleteUrl": setting.ToAbsUrl("api/snapshots-delete/" + cmd.DeleteKey),
		"id":        cmd.Result.Id,
	})
