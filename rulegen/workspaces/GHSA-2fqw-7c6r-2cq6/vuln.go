package main

		metricMiddleware.Handler,
		middleware.RequestID,
		log.ConditionalLogger(s.cfg.Service.LogLevel, zap.L(), "router_agent"),
	)

	zap.S().Infow("agent authentication", "enabled", s.cfg.Service.Auth.AgentAuthenticationEnabled)
	if s.cfg.Service.Auth.AgentAuthenticationEnabled {
		router.Use(
			auth.NewAgentAuthenticator(s.store).Authenticator,
		)
	}

	router.Use(
		middleware.Recoverer,
		oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapiOpts),
	store store.Store
}

func NewAgentAuthenticator(store store.Store) *AgentAuthenticator {
	return &AgentAuthenticator{store: store}
}

package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"

			token := generateAgentToken("kid", "my_source", "GothamCity", privateKey)

			agentAuthenticator := auth.NewAgentAuthenticator(s)

			agentJwt, err := agentAuthenticator.Authenticate(token)
			Expect(err).To(BeNil())

			token := generateAgentToken("missing-key-kid", "my_source", "GothamCity", privateKey)

			agentAuthenticator := auth.NewAgentAuthenticator(s)

			_, err = agentAuthenticator.Authenticate(token)
			Expect(err).ToNot(BeNil())

			token := generateAgentToken("kid", "my_source", "GothamCity", signingKey)

			agentAuthenticator := auth.NewAgentAuthenticator(s)

			_, err = agentAuthenticator.Authenticate(token)
			Expect(err).ToNot(BeNil())

			token := generateAgentToken("1234_kid", "my_source", "org_id", privateKey)

			agentAuthenticator := auth.NewAgentAuthenticator(s)
			h := &handler{}
			ts := httptest.NewServer(agentAuthenticator.Authenticator(h))
			defer ts.Close()
		})
	})

})

func generateAgentToken(kid, sourceID, orgID string, signingKey *rsa.PrivateKey) string {

	v1alpha1 "github.com/kubev2v/migration-planner/api/v1alpha1/agent"
	agentServer "github.com/kubev2v/migration-planner/internal/api/server/agent"
	apiMappers "github.com/kubev2v/migration-planner/internal/handlers/v1alpha1/mappers"
	"github.com/kubev2v/migration-planner/internal/handlers/validator"
	"github.com/kubev2v/migration-planner/internal/service"
		return agentServer.UpdateSourceInventory400JSONResponse{Message: "empty body"}, nil
	}

	data, err := json.Marshal(request.Body.Inventory)
	if err != nil {
		return agentServer.UpdateSourceInventory500JSONResponse{Message: err.Error()}, nil
		return agentServer.UpdateAgentStatus400JSONResponse{Message: err.Error()}, nil
	}

	_, created, err := h.srv.UpdateAgentStatus(ctx, mappers.AgentUpdateForm{
		ID:         request.Id,
		SourceID:   request.Body.SourceId,
			Expect(tx.Error).To(BeNil())
			agentID := uuid.New()

			user := auth.User{
				Username:     "admin",
				Organization: "admin",
			}
			ctx := auth.NewTokenContext(context.TODO(), user)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateAgentStatus(ctx, server.UpdateAgentStatusRequestObject{
			Expect(tx.Error).To(BeNil())
			agentID := uuid.New()

			user := auth.User{
				Username:     "admin",
				Organization: "admin",
			}
			ctx := auth.NewTokenContext(context.TODO(), user)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateAgentStatus(ctx, server.UpdateAgentStatusRequestObject{
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, agentID, "not-connected", "status-info-1", "cred_url-1", sourceID))
			Expect(tx.Error).To(BeNil())

			user := auth.User{
				Username:     "admin",
				Organization: "admin",
			}
			ctx := auth.NewTokenContext(context.TODO(), user)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateAgentStatus(ctx, server.UpdateAgentStatusRequestObject{
			tx := gormdb.Exec(fmt.Sprintf(insertSourceWithUsernameStm, sourceID, "admin", "admin"))
			Expect(tx.Error).To(BeNil())

			user := auth.User{
				Username:     "batman",
				Organization: "wayne_enterprises",
			}
			ctx := auth.NewTokenContext(context.TODO(), user)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateAgentStatus(ctx, server.UpdateAgentStatusRequestObject{
			Expect(reflect.TypeOf(resp).String()).To(Equal(reflect.TypeOf(server.UpdateAgentStatus400JSONResponse{}).String()))
		})

		AfterEach(func() {
			gormdb.Exec("DELETE FROM agents;")
			gormdb.Exec("DELETE FROM sources;")
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, agentID, "not-connected", "status-info-1", "cred_url-1", sourceID))
			Expect(tx.Error).To(BeNil())

			user := auth.User{
				Username:     "admin",
				Organization: "admin",
			}
			ctx := auth.NewTokenContext(context.TODO(), user)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateSourceInventory(ctx, server.UpdateSourceInventoryRequestObject{
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, secondAgentID, "not-connected", "status-info-1", "cred_url-1", sourceID))
			Expect(tx.Error).To(BeNil())

			user := auth.User{
				Username:     "admin",
				Organization: "admin",
			}
			ctx := auth.NewTokenContext(context.TODO(), user)

			// first agent request
			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, uuid.New(), "not-connected", "status-info-1", "cred_url-1", secondSourceID))
			Expect(tx.Error).To(BeNil())

			user := auth.User{
				Username:     "admin",
				Organization: "admin",
			}
			ctx := auth.NewTokenContext(context.TODO(), user)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateSourceInventory(ctx, server.UpdateSourceInventoryRequestObject{
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, firstAgentID, "not-connected", "status-info-1", "cred_url-1", firstSourceID))
			Expect(tx.Error).To(BeNil())

			user := auth.User{
				Username:     "admin",
				Organization: "admin",
			}
			ctx := auth.NewTokenContext(context.TODO(), user)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateSourceInventory(ctx, server.UpdateSourceInventoryRequestObject{

		})

		AfterEach(func() {
			gormdb.Exec("DELETE FROM agents;")
			gormdb.Exec("DELETE FROM sources;")
