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

	// SQLite has a limitation of how many parameters we can specify in a single
	// query, so we need to delete the activies in chunks of 32,000 instead of
	// all at once.
	i := 0
	idsLen := len(ids)
	for i < idsLen {
		start := i
		end := min(i+32000, idsLen)
		batchSize := end - start

		tx := database.Instance().WithContext(ctx).Where("id IN ?", ids[start:end]).Delete(&models.Activity{})
		if tx.Error != nil {
			return errors.WithStack(tx.Error)
		}

		i += batchSize
	}
	return nil
}
