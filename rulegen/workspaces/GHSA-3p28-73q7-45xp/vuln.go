package main


	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/nef/internal/sbi/processor"
	"github.com/free5gc/nef/pkg/app"
	"github.com/free5gc/nef/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/metrics"

	s.router.Use(metrics.InboundMetrics())

	endpoints := s.getTrafficInfluenceRoutes()
	group := s.router.Group(factory.TraffInfluResUriPrefix)
	applyRoutes(group, endpoints)

	endpoints = s.getPFDManagementRoutes()
	group = s.router.Group(factory.PfdMngResUriPrefix)
	applyRoutes(group, endpoints)

	endpoints = s.getPFDFRoutes()
	group = s.router.Group(factory.NefPfdMngResUriPrefix)
	applyRoutes(group, endpoints)

	endpoints = s.getOamRoutes()
	group = s.router.Group(factory.NefOamResUriPrefix)
	applyRoutes(group, endpoints)

	endpoints = s.getCallbackRoutes()
	group = s.router.Group(factory.NefCallbackResUriPrefix)
	applyRoutes(group, endpoints)

	s.router.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
package util

// Metrics consts
const (
	METRICS_APP_PFDS_CREATION_ERR_MSG = "PFDs for all application were not created successfully"
)
	Config() *factory.Config
}

type NefContext struct {
	nef

	return oauth.GetTokenCtx(models.NrfNfManagementNfType_NEF, targetNF,
		c.nfInstID, c.Config().NrfUri(), string(serviceName))
}
