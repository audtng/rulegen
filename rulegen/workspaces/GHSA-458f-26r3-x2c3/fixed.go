package main

	"time"

	"github.com/gorilla/websocket"
	"github.com/microcosm-cc/bluemonday"
	"github.com/schollz/documentsimilarity"
	log "github.com/schollz/logger"
	"github.com/schollz/rwtxt/pkg/db"
	"github.com/schollz/rwtxt/pkg/utils"
)

var pbclean = bluemonday.UGCPolicy()

const DefaultBind = ":8152"

type RWTxt struct {
		return rwt.handleStatic(w, r)
	}

	fields := strings.Split(pbclean.Sanitize(r.URL.Path), "/")

	tr := NewTemplateRender(rwt)
	tr.Domain = "public"
