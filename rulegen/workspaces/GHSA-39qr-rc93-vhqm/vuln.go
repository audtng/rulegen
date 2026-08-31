package main

	"net/http"
)

type apiRequest struct {
	Servers []string `json:"servers"`
	Type    string   `json:"type"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	var request apiRequest
	var response apiResponse
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response = apiResponse{
			Error: err.Error(),
		}
}

func webHandlerTelegramBot(w http.ResponseWriter, r *http.Request) {
	// Parse only needed fields of incoming JSON body
	var err error
	var request tgWebhookRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		println(err.Error())
		return
	}
