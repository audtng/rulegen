package main

func SendOnDataChangeNotify(ueId string, notifyItems []models.NotifyItem) {
	defer func() {
		if p := recover(); p != nil {
			logger.HttpLog.Errorf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

		if ueId == subscriptionDataSubscription.UeId {
			onDataChangeNotifyUrl := subscriptionDataSubscription.CallbackReference

			dataChangeReq := DataRepository.SubscriptionDataSubscriptionsOnDataChangePostRequest{
				DataChangeNotify: &models.DataChangeNotify{
					UeId: ueId,
					OriginalCallbackReference: []string{
						subscriptionDataSubscription.OriginalCallbackReference,
					},
					NotifyItems: notifyItems,
				},
			}
			rsp, err := client.SubsToNotifyCollectionApi.SubscriptionDataSubscriptionsOnDataChangePost(
				context.TODO(), onDataChangeNotifyUrl, &dataChangeReq)

				policyDataChangeNotification,
			},
		}
		rsp, err := client.PolicyDataSubscriptionsCollectionApi.
			CreateIndividualPolicyDataSubscriptionPolicyDataChangeNotificationPost(context.TODO(),
				policyDataChangeNotificationUrl, &req)
			req := DataRepository.CreateIndividualInfluenceDataSubscriptionTrafficInfluenceDataChangeNotificationPostRequest{
				RequestBody: []interface{}{trafficInfluDataNotif},
			}
			rsp, err := client.InfluenceDataSubscriptionsCollectionApi.
				CreateIndividualInfluenceDataSubscriptionTrafficInfluenceDataChangeNotificationPost(
					context.TODO(), influenceDataChangeNotificationUrl, &req)
			req := DataRepository.CreateIndividualInfluenceDataSubscriptionTrafficInfluenceDataChangeNotificationPostRequest{
				RequestBody: []interface{}{trafficInfluDataNotif},
			}
			rsp, err := client.InfluenceDataSubscriptionsCollectionApi.
				CreateIndividualInfluenceDataSubscriptionTrafficInfluenceDataChangeNotificationPost(
					context.TODO(), influenceDataChangeNotificationUrl, &req)
func (p *Processor) RemoveAmfSubscriptionsInfoProcedure(c *gin.Context, subsId string, ueId string) {
	udrSelf := udr_context.GetSelf()
	value, ok := udrSelf.UESubsCollection.Load(ueId)
	var pd *models.ProblemDetails

	if !ok {
		pd = util.ProblemDetailsNotFound("USER_NOT_FOUND")
		logger.DataRepoLog.Errorf("RemoveAmfSubscriptionsInfoProcedure err: %s", pd.Detail)
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	UESubsData := value.(*udr_context.UESubsData)
	eeSub, ok := UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		pd = util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
		logger.DataRepoLog.Errorf("RemoveAmfSubscriptionsInfoProcedure err: %s", pd.Detail)
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	if eeSub == nil || eeSub.AmfSubscriptionInfos == nil {
		pd = util.ProblemDetailsNotFound("AMFSUBSCRIPTION_NOT_FOUND")
	}

			c.JSON(http.StatusInternalServerError, problemDetails)
			return
		}
		for i := range tmp {
			dnnConfigurations := tmp[i].DnnConfigurations
			tmpDnnConfigurations := make(map[string]models.DnnConfiguration)
			for escapedDnn, dnnConf := range dnnConfigurations {
				dnn := util.UnescapeDnn(escapedDnn)
				tmpDnnConfigurations[dnn] = dnnConf
			}
			tmp[i].DnnConfigurations = tmpDnnConfigurations
		}
		if provisionedDataSets.SmData == nil {
			provisionedDataSets.SmData = &models.SmSubsData{}
		}
		provisionedDataSets.SmData.IndividualSmSubsData = tmp
	}
