package main

	"net/http"
)

const maxRequestBodySize = 100 * 1024 // 100KB

type apiRequest struct {
	Servers []string `json:"servers"`
	Type    string   `json:"type"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	var request apiRequest
	var response apiResponse
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		response = apiResponse{
			Error: err.Error(),
		}
}

func webHandlerTelegramBot(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	var err error
	var request tgWebhookRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		println(err.Error())
		return
	}
