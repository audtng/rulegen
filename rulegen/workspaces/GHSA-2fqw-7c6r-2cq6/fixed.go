package main

		metricMiddleware.Handler,
		middleware.RequestID,
		log.ConditionalLogger(s.cfg.Service.LogLevel, zap.L(), "router_agent"),
		auth.NewAgentAuthenticator(s.cfg.Service.Auth.AgentAuthenticationEnabled, s.store).Authenticator,
	)

	router.Use(
		middleware.Recoverer,
		oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapiOpts),
	store store.Store
}

func NewAgentAuthenticator(enabled bool, store store.Store) Authenticator {
	if enabled {
		zap.S().Named("auth").Info("agent authentication enabled")
		return newProductionAgentAuthenticator(store)
	}

	zap.S().Named("auth").Info("agent authentication disabled, using none authenticator")
	return NewNoneAgentAuthenticator()
}

func newProductionAgentAuthenticator(store store.Store) *AgentAuthenticator {
	return &AgentAuthenticator{store: store}
}

package auth_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"

			token := generateAgentToken("kid", "my_source", "GothamCity", privateKey)

			agentAuthenticator := auth.NewAgentAuthenticator(true, s).(*auth.AgentAuthenticator)

			agentJwt, err := agentAuthenticator.Authenticate(token)
			Expect(err).To(BeNil())

			token := generateAgentToken("missing-key-kid", "my_source", "GothamCity", privateKey)

			agentAuthenticator := auth.NewAgentAuthenticator(true, s).(*auth.AgentAuthenticator)

			_, err = agentAuthenticator.Authenticate(token)
			Expect(err).ToNot(BeNil())

			token := generateAgentToken("kid", "my_source", "GothamCity", signingKey)

			agentAuthenticator := auth.NewAgentAuthenticator(true, s).(*auth.AgentAuthenticator)

			_, err = agentAuthenticator.Authenticate(token)
			Expect(err).ToNot(BeNil())

			token := generateAgentToken("1234_kid", "my_source", "org_id", privateKey)

			agentAuthenticator := auth.NewAgentAuthenticator(true, s)
			h := &handler{}
			ts := httptest.NewServer(agentAuthenticator.Authenticator(h))
			defer ts.Close()
		})
	})

	Context("none authenticator middleware", func() {
		It("successfully extracts sourceId from request body", func() {
			innerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify AgentJWT is in context
				agentJWT := auth.MustHaveAgent(r.Context())
				Expect(agentJWT.SourceID).To(Equal("test-source-123"))

				// Verify body is still readable for downstream
				var body map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&body)
				Expect(err).To(BeNil(), "Body should be readable by downstream handler")
				Expect(body["sourceId"]).To(Equal("test-source-123"))

				w.WriteHeader(200)
			})

			noneAuthenticator := auth.NewNoneAgentAuthenticator()
			ts := httptest.NewServer(noneAuthenticator.Authenticator(innerHandler))
			defer ts.Close()

			body := `{"sourceId":"test-source-123"}`
			req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBufferString(body))
			Expect(err).To(BeNil())

			resp, rerr := http.DefaultClient.Do(req)
			Expect(rerr).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
		})
	})

})

func generateAgentToken(kid, sourceID, orgID string, signingKey *rsa.PrivateKey) string {
package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type NoneAgentAuthenticator struct{}

func NewNoneAgentAuthenticator() *NoneAgentAuthenticator {
	return &NoneAgentAuthenticator{}
}

func (n *NoneAgentAuthenticator) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read request body: %v", err), http.StatusBadRequest)
			return
		}
		_ = r.Body.Close()

		var req struct {
			SourceID string `json:"sourceId"`
		}

		if err := json.Unmarshal(bodyBytes, &req); err != nil {
			http.Error(w, fmt.Sprintf("missing source id from request body: %v", err), http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		agentJWT := AgentJWT{
			ExpireAt: time.Now().Add(defaultExpirationPeriod * time.Hour),
			IssueAt:  time.Now(),
			Issuer:   "none",
			OrgID:    "internal",
			SourceID: req.SourceID,
		}
		ctx := NewTokenContext(r.Context(), agentJWT)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

	v1alpha1 "github.com/kubev2v/migration-planner/api/v1alpha1/agent"
	agentServer "github.com/kubev2v/migration-planner/internal/api/server/agent"
	"github.com/kubev2v/migration-planner/internal/auth"
	apiMappers "github.com/kubev2v/migration-planner/internal/handlers/v1alpha1/mappers"
	"github.com/kubev2v/migration-planner/internal/handlers/validator"
	"github.com/kubev2v/migration-planner/internal/service"
		return agentServer.UpdateSourceInventory400JSONResponse{Message: "empty body"}, nil
	}

	agentJWT := auth.MustHaveAgent(ctx)
	if agentJWT.SourceID != request.Id.String() {
		return agentServer.UpdateSourceInventory403JSONResponse{
			Message: fmt.Sprintf("agent is not authorized to update source %s", request.Id),
		}, nil
	}

	data, err := json.Marshal(request.Body.Inventory)
	if err != nil {
		return agentServer.UpdateSourceInventory500JSONResponse{Message: err.Error()}, nil
		return agentServer.UpdateAgentStatus400JSONResponse{Message: err.Error()}, nil
	}

	agentJWT := auth.MustHaveAgent(ctx)
	if agentJWT.SourceID != request.Body.SourceId.String() {
		return agentServer.UpdateAgentStatus403JSONResponse{
			Message: fmt.Sprintf("agent is not authorized to update source %s", request.Body.SourceId),
		}, nil
	}

	_, created, err := h.srv.UpdateAgentStatus(ctx, mappers.AgentUpdateForm{
		ID:         request.Id,
		SourceID:   request.Body.SourceId,
			Expect(tx.Error).To(BeNil())
			agentID := uuid.New()

			agentJWT := auth.AgentJWT{
				OrgID:    "admin",
				SourceID: sourceID.String(),
			}
			ctx := auth.NewTokenContext(context.TODO(), agentJWT)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateAgentStatus(ctx, server.UpdateAgentStatusRequestObject{
			Expect(tx.Error).To(BeNil())
			agentID := uuid.New()

			agentJWT := auth.AgentJWT{
				OrgID:    "admin",
				SourceID: sourceID.String(),
			}
			ctx := auth.NewTokenContext(context.TODO(), agentJWT)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateAgentStatus(ctx, server.UpdateAgentStatusRequestObject{
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, agentID, "not-connected", "status-info-1", "cred_url-1", sourceID))
			Expect(tx.Error).To(BeNil())

			agentJWT := auth.AgentJWT{
				OrgID:    "admin",
				SourceID: sourceID.String(),
			}
			ctx := auth.NewTokenContext(context.TODO(), agentJWT)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateAgentStatus(ctx, server.UpdateAgentStatusRequestObject{
			tx := gormdb.Exec(fmt.Sprintf(insertSourceWithUsernameStm, sourceID, "admin", "admin"))
			Expect(tx.Error).To(BeNil())

			agentJWT := auth.AgentJWT{
				OrgID:    "wayne_enterprises",
				SourceID: sourceID,
			}
			ctx := auth.NewTokenContext(context.TODO(), agentJWT)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateAgentStatus(ctx, server.UpdateAgentStatusRequestObject{
			Expect(reflect.TypeOf(resp).String()).To(Equal(reflect.TypeOf(server.UpdateAgentStatus400JSONResponse{}).String()))
		})

		It("rejects update when JWT source_id does not match target source", func() {
			sourceID := uuid.New()
			differentSourceID := uuid.New()
			tx := gormdb.Exec(fmt.Sprintf(insertSourceWithUsernameStm, sourceID, "admin", "admin"))
			Expect(tx.Error).To(BeNil())
			agentID := uuid.New()

			// JWT has a different source_id than the one being updated
			agentJWT := auth.AgentJWT{
				OrgID:    "admin",
				SourceID: differentSourceID.String(),
			}
			ctx := auth.NewTokenContext(context.TODO(), agentJWT)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateAgentStatus(ctx, server.UpdateAgentStatusRequestObject{
				Id: agentID,
				Body: &apiAgent.UpdateAgentStatusJSONRequestBody{
					Status:        string(v1alpha1.AgentStatusWaitingForCredentials),
					StatusInfo:    "waiting-for-credentials",
					CredentialUrl: "http://agent.com",
					Version:       "version-1",
					SourceId:      sourceID,
				},
			})
			Expect(err).To(BeNil())
			Expect(reflect.TypeOf(resp).String()).To(Equal(reflect.TypeOf(server.UpdateAgentStatus403JSONResponse{}).String()))

			// Verify no agent was created
			count := -1
			tx = gormdb.Raw("SELECT COUNT(*) FROM agents;").Scan(&count)
			Expect(tx.Error).To(BeNil())
			Expect(count).To(Equal(0))
		})

		AfterEach(func() {
			gormdb.Exec("DELETE FROM agents;")
			gormdb.Exec("DELETE FROM sources;")
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, agentID, "not-connected", "status-info-1", "cred_url-1", sourceID))
			Expect(tx.Error).To(BeNil())

			agentJWT := auth.AgentJWT{
				OrgID:    "admin",
				SourceID: sourceID.String(),
			}
			ctx := auth.NewTokenContext(context.TODO(), agentJWT)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateSourceInventory(ctx, server.UpdateSourceInventoryRequestObject{
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, secondAgentID, "not-connected", "status-info-1", "cred_url-1", sourceID))
			Expect(tx.Error).To(BeNil())

			agentJWT := auth.AgentJWT{
				OrgID:    "admin",
				SourceID: sourceID.String(),
			}
			ctx := auth.NewTokenContext(context.TODO(), agentJWT)

			// first agent request
			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, uuid.New(), "not-connected", "status-info-1", "cred_url-1", secondSourceID))
			Expect(tx.Error).To(BeNil())

			agentJWT := auth.AgentJWT{
				OrgID:    "admin",
				SourceID: firstSourceID.String(),
			}
			ctx := auth.NewTokenContext(context.TODO(), agentJWT)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateSourceInventory(ctx, server.UpdateSourceInventoryRequestObject{
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, firstAgentID, "not-connected", "status-info-1", "cred_url-1", firstSourceID))
			Expect(tx.Error).To(BeNil())

			agentJWT := auth.AgentJWT{
				OrgID:    "admin",
				SourceID: firstSourceID.String(),
			}
			ctx := auth.NewTokenContext(context.TODO(), agentJWT)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateSourceInventory(ctx, server.UpdateSourceInventoryRequestObject{

		})

		It("rejects inventory update when JWT source_id does not match target source", func() {
			sourceID := uuid.New()
			differentSourceID := uuid.New()
			agentID := uuid.New()
			tx := gormdb.Exec(fmt.Sprintf(insertSourceWithUsernameStm, sourceID, "admin", "admin"))
			Expect(tx.Error).To(BeNil())
			tx = gormdb.Exec(fmt.Sprintf(insertAgentStm, agentID, "not-connected", "status-info-1", "cred_url-1", sourceID))
			Expect(tx.Error).To(BeNil())

			// JWT has a different source_id than the one being updated
			agentJWT := auth.AgentJWT{
				OrgID:    "admin",
				SourceID: differentSourceID.String(),
			}
			ctx := auth.NewTokenContext(context.TODO(), agentJWT)

			srv := handlers.NewAgentHandler(service.NewAgentService(s))
			resp, err := srv.UpdateSourceInventory(ctx, server.UpdateSourceInventoryRequestObject{
				Id: sourceID,
				Body: &apiAgent.SourceStatusUpdate{
					AgentId: agentID,
					Inventory: v1alpha1.Inventory{
						VcenterId: "vcenter",
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(reflect.TypeOf(resp).String()).To(Equal(reflect.TypeOf(server.UpdateSourceInventory403JSONResponse{}).String()))

			// Verify inventory was not updated
			source, err := s.Source().Get(ctx, sourceID)
			Expect(err).To(BeNil())
			Expect(source.Inventory).To(BeNil())
		})

		AfterEach(func() {
			gormdb.Exec("DELETE FROM agents;")
			gormdb.Exec("DELETE FROM sources;")
