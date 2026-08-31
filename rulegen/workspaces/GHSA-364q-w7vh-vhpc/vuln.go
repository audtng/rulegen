package main

	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
	MaxTriggerDepth            = 10
)

var (
	metricActionsRequested = promauto.NewCounter(prometheus.CounterOpts{
		Name: "olivetin_actions_requested_count",
	}

	_, isDuplicate := e.GetLog(req.TrackingID)

	if isDuplicate || req.TrackingID == "" {
		req.TrackingID = uuid.NewString()
	}

