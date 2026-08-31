package main

	defer resp.Body.Close()

	var response model.PostActionIntegrationResponse
	limitedReader := io.LimitReader(resp.Body, MaxIntegrationResponseSize)
	respBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", model.NewAppError("DoPostActionWithCookie", "api.post.do_action.action_integration.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}
	})
}

// infiniteReader generates unlimited data for testing response size limits
type infiniteReader struct{}

func (r infiniteReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 'a'
	}
	return len(p), nil
}

// MM-67074: TestPostActionResponseSizeLimit verifies that DoPostActionWithCookie
// properly limits response sizes to prevent OOM attacks
func TestPostActionResponseSizeLimit(t *testing.T) {
	mainHelper.Parallel(t)
	th := Setup(t).InitBasic(t)

	channel := th.BasicChannel
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.ServiceSettings.AllowedUntrustedInternalConnections = "localhost,127.0.0.1"
	})

	t.Run("large valid JSON response is truncated", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Send response larger than MaxIntegrationResponseSize (1MB)
			// Response starts as valid JSON but becomes truncated
			_, _ = io.Copy(w, io.MultiReader(
				strings.NewReader(`{"update":{"message":"`),
				infiniteReader{},
				strings.NewReader(`"}}`),
			))
		}))
		defer server.Close()

		interactivePost := model.Post{
			Message:       "Interactive post",
			ChannelId:     channel.Id,
			PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
			UserId:        th.BasicUser.Id,
			Props: model.StringInterface{
				model.PostPropsAttachments: []*model.SlackAttachment{
					{
						Text: "hello",
						Actions: []*model.PostAction{
							{
								Type: model.PostActionTypeButton,
								Name: "action",
								Integration: &model.PostActionIntegration{
									URL: server.URL,
								},
							},
						},
					},
				},
			},
		}

		post, err := th.App.CreatePostAsUser(th.Context, &interactivePost, "", true)
		require.Nil(t, err)
		attachments, ok := post.GetProp(model.PostPropsAttachments).([]*model.SlackAttachment)
		require.True(t, ok)

		// Should return error due to truncated JSON, but NOT crash or OOM
		_, err = th.App.DoPostActionWithCookie(th.Context, post.Id,
			attachments[0].Actions[0].Id, th.BasicUser.Id, "", nil)
		require.NotNil(t, err)
		// Truncated JSON causes unmarshal error
		assert.Equal(t, "api.post.do_action.action_integration.app_error", err.Id)
	})

	t.Run("large invalid response is truncated", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Send infinite non-JSON data
			_, _ = io.Copy(w, infiniteReader{})
		}))
		defer server.Close()

		interactivePost := model.Post{
			Message:       "Interactive post",
			ChannelId:     channel.Id,
			PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
			UserId:        th.BasicUser.Id,
			Props: model.StringInterface{
				model.PostPropsAttachments: []*model.SlackAttachment{
					{
						Text: "hello",
						Actions: []*model.PostAction{
							{
								Type: model.PostActionTypeButton,
								Name: "action",
								Integration: &model.PostActionIntegration{
									URL: server.URL,
								},
							},
						},
					},
				},
			},
		}

		post, err := th.App.CreatePostAsUser(th.Context, &interactivePost, "", true)
		require.Nil(t, err)
		attachments, ok := post.GetProp(model.PostPropsAttachments).([]*model.SlackAttachment)
		require.True(t, ok)

		// Should return error due to invalid JSON, but NOT crash or OOM
		_, err = th.App.DoPostActionWithCookie(th.Context, post.Id,
			attachments[0].Actions[0].Id, th.BasicUser.Id, "", nil)
		require.NotNil(t, err)
		assert.Equal(t, "api.post.do_action.action_integration.app_error", err.Id)
	})
}

func TestPostAction(t *testing.T) {
	mainHelper.Parallel(t)
	testCases := []struct {
			func (p *MyPlugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
				var request model.SubmitDialogRequest
				json.NewDecoder(r.Body).Decode(&request)

				response := &model.LookupDialogResponse{
					Items: []model.DialogSelectOption{
						{Text: "Plugin Option 1", Value: "plugin_value1"},
