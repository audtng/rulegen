package main


func OriginChecker(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if *Cfg.ServiceSettings.AllowCorsFrom == "*" {
		return true
	}
	for _, allowed := range strings.Split(*Cfg.ServiceSettings.AllowCorsFrom, " ") {
		if allowed == origin {
			return true
		}
	}
	return false
}

func GetOriginChecker(r *http.Request) OriginCheckerProc {

func (cw *CorsWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(*utils.Cfg.ServiceSettings.AllowCorsFrom) > 0 {
		if utils.OriginChecker(r) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

			if r.Method == "OPTIONS" {
				w.Header().Set(
