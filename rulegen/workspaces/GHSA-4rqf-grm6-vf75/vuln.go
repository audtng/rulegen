package main

func SendOnDataChangeNotify(ueId string, notifyItems []models.NotifyItem) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.HttpLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

		if ueId == subscriptionDataSubscription.UeId {
			onDataChangeNotifyUrl := subscriptionDataSubscription.CallbackReference

			dataChangeReq := DataRepository.SubscriptionDataSubscriptionsOnDataChangePostRequest{}
			dataChangeReq.DataChangeNotify.UeId = ueId
			dataChangeReq.DataChangeNotify.OriginalCallbackReference = []string{
				subscriptionDataSubscription.OriginalCallbackReference,
			}
			dataChangeReq.DataChangeNotify.NotifyItems = notifyItems
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
	var pd *models.ProblemDetails = nil

	if !ok {
		pd = util.ProblemDetailsNotFound("USER_NOT_FOUND")
		logger.DataRepoLog.Errorf("RemoveAmfSubscriptionsInfoProcedure err: %s", pd.Detail)
	}

	UESubsData := value.(*udr_context.UESubsData)
	_, ok = UESubsData.EeSubscriptionCollection[subsId]

	if !ok {
		pd = util.ProblemDetailsNotFound("SUBSCRIPTION_NOT_FOUND")
		logger.DataRepoLog.Errorf("RemoveAmfSubscriptionsInfoProcedure err: %s", pd.Detail)
	}

	if UESubsData.EeSubscriptionCollection[subsId].AmfSubscriptionInfos == nil {
		pd = util.ProblemDetailsNotFound("AMFSUBSCRIPTION_NOT_FOUND")
	}

			c.JSON(http.StatusInternalServerError, problemDetails)
			return
		}
		for _, smData := range tmp {
			dnnConfigurations := smData.DnnConfigurations
			tmpDnnConfigurations := make(map[string]models.DnnConfiguration)
			for escapedDnn, dnnConf := range dnnConfigurations {
				dnn := util.UnescapeDnn(escapedDnn)
				tmpDnnConfigurations[dnn] = dnnConf
			}
			smData.DnnConfigurations = tmpDnnConfigurations
		}
		provisionedDataSets.SmData.IndividualSmSubsData = tmp
	}
