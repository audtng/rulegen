package main

	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
	MaxTriggerDepth            = 10
)

var validTrackingIDPattern = regexp.MustCompile(`^[a-fA-F0-9\-]+$`)

func isValidTrackingID(id string) bool {
	const MaxTrackingIDLength = 36

	return id != "" && len(id) <= MaxTrackingIDLength && validTrackingIDPattern.MatchString(id)
}

var (
	metricActionsRequested = promauto.NewCounter(prometheus.CounterOpts{
		Name: "olivetin_actions_requested_count",
	}

	_, isDuplicate := e.GetLog(req.TrackingID)
	if isDuplicate || !isValidTrackingID(req.TrackingID) {
		req.TrackingID = uuid.NewString()
	}

