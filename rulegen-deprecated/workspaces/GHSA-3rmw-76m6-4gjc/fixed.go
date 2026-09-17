package main

)

const (
	SMTPConfigProjectionTable = "projections.smtp_configs5"
	SMTPConfigTable           = SMTPConfigProjectionTable + "_" + smtpConfigSMTPTableSuffix
	SMTPConfigHTTPTable       = SMTPConfigProjectionTable + "_" + smtpConfigHTTPTableSuffix

		return nil, err
	}

	description := e.Description
	state := domain.SMTPConfigStateInactive
	if e.ID == "" {
		description = "generic"
		state = domain.SMTPConfigStateActive
	}
				handler.NewCol(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(SMTPConfigColumnAggregateID, e.Aggregate().ID),
				handler.NewCol(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(SMTPConfigColumnState, state),
				handler.NewCol(SMTPConfigColumnDescription, description),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigSMTPColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCol(SMTPConfigSMTPColumnTLS, e.TLS),
				handler.NewCol(SMTPConfigSMTPColumnSenderAddress, e.SenderAddress),
				handler.NewCol(SMTPConfigSMTPColumnSenderName, e.SenderName),
				handler.NewCol(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(SMTPConfigColumnAggregateID, e.Aggregate().ID),
				handler.NewCol(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateInactive),
				handler.NewCol(SMTPConfigColumnDescription, e.Description),
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigHTTPColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigHTTPColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCol(SMTPConfigHTTPColumnEndpoint, e.Endpoint),
			},
			handler.WithTableSuffix(smtpConfigHTTPTableSuffix),
		stmts = append(stmts, handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		))
		stmts = append(stmts, handler.AddUpdateStatement(
			smtpColumns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigHTTPColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigHTTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigHTTPTableSuffix),
		return nil, err
	}

	stmts := make([]func(eventstore.Event) handler.Exec, 0, 3)
	columns := []handler.Column{
		handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
		stmts = append(stmts, handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		))
	}

	smtpColumns := make([]handler.Column, 0, 7)
	if e.TLS != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnTLS, *e.TLS))
	}
	if e.FromAddress != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnSenderAddress, *e.FromAddress))
	}
	if e.FromName != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnSenderName, *e.FromName))
	}
	if e.ReplyToAddress != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnReplyToAddress, *e.ReplyToAddress))
	}
	if e.Host != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnHost, *e.Host))
	}
	if e.User != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnUser, *e.User))
	}
	if e.Password != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnPassword, *e.Password))
	}
	if len(smtpColumns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			smtpColumns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigSMTPColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigSMTPTableSuffix),
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
				handler.NewCol(SMTPConfigSMTPColumnPassword, e.Password),
			},
			[]handler.Condition{
				handler.NewCond(SMTPConfigSMTPColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigSMTPTableSuffix),
		),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
				handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateInactive),
			},
			[]handler.Condition{
				handler.Not(handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate()))),
				handler.NewCond(SMTPConfigColumnState, domain.SMTPConfigStateActive),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
				handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateActive),
			},
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
		return nil, err
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
			handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
		return nil, err
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
			handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func getSMTPConfigID(id string, aggregate *eventstore.Aggregate) string {
	if id != "" {
		return id
	}
	// Deal with old and unique SMTP settings (empty ID)
	return aggregate.ResourceOwner
}
package login

import (
	"context"
	"net/http"

	"golang.org/x/text/language"
	OrgRegister        bool
}

func determineResourceOwner(ctx context.Context, authRequest *domain.AuthRequest) string {
	if authRequest != nil && authRequest.RequestedOrgID != "" {
		return authRequest.RequestedOrgID
	}
	return authz.GetInstance(ctx).DefaultOrganisationID()
}

func (l *Login) handleRegister(w http.ResponseWriter, r *http.Request) {
	data := new(registerFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
		l.renderError(w, r, authRequest, err)
		return
	}
	if err := l.checkRegistrationAllowed(r, determineResourceOwner(r.Context(), authRequest), authRequest); err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}
	l.renderRegister(w, r, authRequest, data, nil)
}

func (l *Login) checkRegistrationAllowed(r *http.Request, orgID string, authReq *domain.AuthRequest) error {
	if authReq != nil {
		if registrationAllowed(authReq) {
			return nil
		}
		return zerrors.ThrowPreconditionFailed(nil, "VIEW-RRGRXz4kGw", "Errors.Org.LoginPolicy.RegistrationNotAllowed")
	}
	loginPolicy, err := l.getLoginPolicy(r, orgID)
	if err != nil {
		return err
	}
	if loginPolicy.AllowRegister && loginPolicy.AllowUsernamePassword {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "VIEW-Vq3bduAacD", "Errors.Org.LoginPolicy.RegistrationNotAllowed")
}

func (l *Login) handleRegisterCheck(w http.ResponseWriter, r *http.Request) {
	data := new(registerFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
		l.renderError(w, r, authRequest, err)
		return
	}
	resourceOwner := determineResourceOwner(r.Context(), authRequest)
	if err := l.checkRegistrationAllowed(r, resourceOwner, authRequest); err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}
	if data.Password != data.Password2 {
		err := zerrors.ThrowInvalidArgument(nil, "VIEW-KaGue", "Errors.User.Password.ConfirmationWrong")
		l.renderRegister(w, r, authRequest, data, err)
		return
	}
	// For consistency with the external authentication flow,
	// the setMetadata() function is provided on the pre creation hook, for now,
	// like for the ExternalAuthentication flow.
		formData.Language = l.renderer.ReqLang(translator, r).String()
	}

	resourceOwner := determineResourceOwner(r.Context(), authRequest)
	data := registerData{
		baseData:         l.getBaseData(r, authRequest, translator, "RegistrationUser.Title", "RegistrationUser.Description", errID, errMessage),
		registerFormData: *formData,
			return CtxData{}, zerrors.ThrowUnauthenticated(errors.Join(err, sysTokenErr), "AUTH-7fs1e", "Errors.Token.Invalid")
		}
	}
	projectID, err := projectIDAndCheckOriginForClientID(ctx, clientID, t)
	if err != nil {
		return CtxData{}, err
	}
	if orgID == "" && orgDomain == "" {
		orgID = resourceOwner
	}
	// System API calls don't have a resource owner
	if orgID != "" || orgDomain != "" {
		orgID, err = t.ExistsOrg(ctx, orgID, orgDomain)
		if err != nil {
			return CtxData{}, zerrors.ThrowPermissionDenied(nil, "AUTH-Bs7Ds", "Organisation doesn't exist")
	}, nil
}

func projectIDAndCheckOriginForClientID(ctx context.Context, clientID string, t APITokenVerifier) (string, error) {
	if clientID == "" {
		return "", nil
	}
	projectID, origins, err := t.ProjectIDAndOriginsByClientID(ctx, clientID)
	logging.WithFields("clientID", clientID).OnError(err).Debug("could not check projectID and origin of clientID (might be service account)")

	// We used to check origins for every token, but service users shouldn't be used publicly (native app / SPA).
	// Therefore, mostly won't send an origin and aren't able to configure them anyway.
	// For the current time we will only check origins for tokens issued to users through apps (code / implicit flow).
	if projectID == "" {
		return "", nil
	}
	return projectID, checkOrigin(ctx, origins)
}

func SetCtxData(ctx context.Context, ctxData CtxData) context.Context {
	return context.WithValue(ctx, dataKey, ctxData)
}
