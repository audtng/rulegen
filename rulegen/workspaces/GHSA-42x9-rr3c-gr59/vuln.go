package main


func OriginChecker(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	return *Cfg.ServiceSettings.AllowCorsFrom == "*" || strings.Contains(*Cfg.ServiceSettings.AllowCorsFrom, origin)
}

func GetOriginChecker(r *http.Request) OriginCheckerProc {

func (cw *CorsWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(*utils.Cfg.ServiceSettings.AllowCorsFrom) > 0 {
		origin := r.Header.Get("Origin")
		if *utils.Cfg.ServiceSettings.AllowCorsFrom == "*" || strings.Contains(*utils.Cfg.ServiceSettings.AllowCorsFrom, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)

			if r.Method == "OPTIONS" {
				w.Header().Set(
