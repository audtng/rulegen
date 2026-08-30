package main

import (
	"context"
	"emperror.dev/errors"
	"github.com/pterodactyl/wings/internal/database"
	"github.com/pterodactyl/wings/internal/models"
)

type activityCron struct{}

func (ac *activityCron) Run(ctx context.Context) error {
	var activities []models.Activity
	ids := make([]int, len(activities))
	for i, v := range activities {
		ids[i] = v.ID
	}

	// Delete all the activities that were sent to the Panel (or that were invalid).
	tx := database.Instance().WithContext(ctx).Where("id IN ?", ids).Delete(&models.Activity{})
	if tx.Error != nil {
		return errors.WithStack(tx.Error)
	}
	return nil
}
