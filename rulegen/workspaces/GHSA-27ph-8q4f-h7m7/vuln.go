package main

		SuppFeat:          util.StringToPtr(request.SuppFeat),
	}

	if exists {
		// Update existing subscription
		bsfContext.BsfSelf.UpdateSubscription(subId, subscription)

		// Return updated subscription
		response := models.BsfSubscriptionResp{
			Events:            subscription.Events,
			NotifUri:          subscription.NotifUri,
			NotifCorreId:      subscription.NotifCorreId,
			Supi:              subscription.Supi,
			Gpsi:              util.PtrToString(subscription.Gpsi),
			SnssaiDnnPairs:    subscription.SnssaiDnnPairs,
			AddSnssaiDnnPairs: subscription.AddSnssaiDnnPairs,
			SuppFeat:          util.PtrToString(subscription.SuppFeat),
		}
		c.JSON(http.StatusOK, response)
	} else {
		// Create new subscription with given ID
		subscription.SubId = subId
		bsfContext.BsfSelf.Subscriptions[subId] = subscription

		c.Status(http.StatusNoContent)
	}
}

// DeleteIndividualSubcription handles DELETE /subscriptions/{subId}
internal/sbi/processor/subscriptions.go | 2 +-
1 file changed, 1 insertion(+), 1 deletion(-)
		SuppFeat:          util.StringToPtr(request.SuppFeat),
	}

	if (!exists) {
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
