package main

	var re = regexp.MustCompile(`^(.*/)([^?].*)?[?|.]*$`)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		matches := re.FindStringSubmatch(r.RequestURI)
		path := matches[2]

			}
			_, _ = w.Write([]byte(doc))
		case "":
			http.Redirect(w, r, h.Prefix+"index.html", http.StatusMovedPermanently)
		default:
			h.ServeHTTP(w, r)
		}
