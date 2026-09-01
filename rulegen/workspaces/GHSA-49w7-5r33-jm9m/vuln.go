package main

	var re = regexp.MustCompile(`^(.*/)([^?].*)?[?|.]*$`)

	return func(w http.ResponseWriter, r *http.Request) {
		matches := re.FindStringSubmatch(r.RequestURI)
		path := matches[2]

			}
			_, _ = w.Write([]byte(doc))
		case "":
			http.Redirect(w, r, h.Prefix+"index.html", 301)
		default:
			h.ServeHTTP(w, r)
		}
