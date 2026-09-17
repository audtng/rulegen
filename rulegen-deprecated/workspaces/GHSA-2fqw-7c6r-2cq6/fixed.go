package main

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

