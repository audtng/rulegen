package main


	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/nef/internal/sbi/processor"
	nef_util "github.com/free5gc/nef/internal/util"
	"github.com/free5gc/nef/pkg/app"
	"github.com/free5gc/nef/pkg/factory"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/metrics"

	s.router.Use(metrics.InboundMetrics())

	// Callback is always mounted – NEF acts as a notification receiver here,
	// so no inbound token check is required.
	callbackGroup := s.router.Group(factory.NefCallbackResUriPrefix)
	applyRoutes(callbackGroup, s.getCallbackRoutes())

	// All other route groups are mounted only when their service is declared in
	// ServiceList, and each group is protected by OAuth2 middleware.
	for _, service := range s.Config().ServiceList() {
		switch service.ServiceName {
		case factory.ServiceNefPfd:
			// nnef-pfdmanagement covers both the external AF-facing API
			// (/3gpp-pfd-management) and the SBI PFDF API (/nnef-pfdmanagement).
			authCheck := nef_util.NewRouterAuthorizationCheck(models.ServiceName_NNEF_PFDMANAGEMENT)

			pfdMngGroup := s.router.Group(factory.PfdMngResUriPrefix)
			pfdMngGroup.Use(func(c *gin.Context) {
				authCheck.Check(c, s.Context())
			})
			applyRoutes(pfdMngGroup, s.getPFDManagementRoutes())

			pfdFGroup := s.router.Group(factory.NefPfdMngResUriPrefix)
			pfdFGroup.Use(func(c *gin.Context) {
				authCheck.Check(c, s.Context())
			})
			applyRoutes(pfdFGroup, s.getPFDFRoutes())

		case factory.ServiceNefOam:
			authCheck := nef_util.NewRouterAuthorizationCheck(models.ServiceName(factory.ServiceNefOam))

			oamGroup := s.router.Group(factory.NefOamResUriPrefix)
			oamGroup.Use(func(c *gin.Context) {
				authCheck.Check(c, s.Context())
			})
			applyRoutes(oamGroup, s.getOamRoutes())

		case factory.ServiceTraffInflu:
			// 3gpp-traffic-influence is an AF-facing API (3GPP TS 29.522);
			authCheck := nef_util.NewRouterAuthorizationCheck(models.ServiceName_3GPP_TRAFFIC_INFLUENCE)
			tiGroup := s.router.Group(factory.TraffInfluResUriPrefix)
			tiGroup.Use(func(c *gin.Context) {
				authCheck.Check(c, s.Context())
			})
			applyRoutes(tiGroup, s.getTrafficInfluenceRoutes())
		}
	}

	s.router.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
package util

import (
	"net/http"

	nef_context "github.com/free5gc/nef/internal/context"
	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/openapi/models"
	"github.com/gin-gonic/gin"
)

// Metrics consts
const (
	METRICS_APP_PFDS_CREATION_ERR_MSG = "PFDs for all application were not created successfully"
)

type RouterAuthorizationCheck struct {
	serviceName models.ServiceName
}

func NewRouterAuthorizationCheck(serviceName models.ServiceName) *RouterAuthorizationCheck {
	return &RouterAuthorizationCheck{
		serviceName: serviceName,
	}
}

func (rac *RouterAuthorizationCheck) Check(c *gin.Context, nefCtx nef_context.NFContext) {
	token := c.Request.Header.Get("Authorization")
	if err := nefCtx.AuthorizationCheck(token, rac.serviceName); err != nil {
		logger.UtilLog.Debugf("RouterAuthorizationCheck::Check Unauthorized: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	logger.UtilLog.Debugf("RouterAuthorizationCheck::Check Authorized")
}
	Config() *factory.Config
}

// NFContext is the interface used by middleware to perform inbound OAuth2 token checks.
type NFContext interface {
	AuthorizationCheck(token string, serviceName models.ServiceName) error
}

var _ NFContext = &NefContext{}

type NefContext struct {
	nef

	return oauth.GetTokenCtx(models.NrfNfManagementNfType_NEF, targetNF,
		c.nfInstID, c.Config().NrfUri(), string(serviceName))
}

// AuthorizationCheck validates the inbound OAuth2 bearer token against serviceName.
// When OAuth2 is disabled it returns nil immediately (pass-through for dev/test).
func (c *NefContext) AuthorizationCheck(token string, serviceName models.ServiceName) error {
	if !c.OAuth2Required {
		logger.CtxLog.Debugf("NefContext::AuthorizationCheck: OAuth2 not required")
		return nil
	}
	logger.CtxLog.Debugf("NefContext::AuthorizationCheck: token[%s] serviceName[%s]", token, serviceName)
	return oauth.VerifyOAuth(token, string(serviceName), c.Config().NrfCertPem())
}
