package main

	defer resp.Body.Close()

	var response model.PostActionIntegrationResponse
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", model.NewAppError("DoPostActionWithCookie", "api.post.do_action.action_integration.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}
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
