package main

)

const (
	SMTPConfigProjectionTable = "projections.smtp_configs4"
	SMTPConfigTable           = SMTPConfigProjectionTable + "_" + smtpConfigSMTPTableSuffix
	SMTPConfigHTTPTable       = SMTPConfigProjectionTable + "_" + smtpConfigHTTPTableSuffix

		return nil, err
	}

	// Deal with old and unique SMTP settings (empty ID)
	id := e.ID
	description := e.Description
	state := domain.SMTPConfigStateInactive
	if e.ID == "" {
		id = e.Aggregate().ResourceOwner
		description = "generic"
		state = domain.SMTPConfigStateActive
	}
				handler.NewCol(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(SMTPConfigColumnAggregateID, e.Aggregate().ID),
				handler.NewCol(SMTPConfigColumnID, id),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(SMTPConfigColumnState, state),
				handler.NewCol(SMTPConfigColumnDescription, description),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigSMTPColumnID, id),
				handler.NewCol(SMTPConfigSMTPColumnTLS, e.TLS),
				handler.NewCol(SMTPConfigSMTPColumnSenderAddress, e.SenderAddress),
				handler.NewCol(SMTPConfigSMTPColumnSenderName, e.SenderName),
				handler.NewCol(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(SMTPConfigColumnAggregateID, e.Aggregate().ID),
				handler.NewCol(SMTPConfigColumnID, e.ID),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateInactive),
				handler.NewCol(SMTPConfigColumnDescription, e.Description),
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigSMTPColumnID, e.ID),
				handler.NewCol(SMTPConfigHTTPColumnEndpoint, e.Endpoint),
			},
			handler.WithTableSuffix(smtpConfigHTTPTableSuffix),
		stmts = append(stmts, handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, e.ID),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		))
		stmts = append(stmts, handler.AddUpdateStatement(
			smtpColumns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigHTTPColumnID, e.ID),
				handler.NewCond(SMTPConfigHTTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigHTTPTableSuffix),
		return nil, err
	}

	// Deal with old and unique SMTP settings (empty ID)
	id := e.ID
	if e.ID == "" {
		id = e.Aggregate().ResourceOwner
	}

	stmts := make([]func(eventstore.Event) handler.Exec, 0, 3)
	columns := []handler.Column{
		handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
		stmts = append(stmts, handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, id),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		))
	}

	httpColumns := make([]handler.Column, 0, 7)
	if e.TLS != nil {
		httpColumns = append(httpColumns, handler.NewCol(SMTPConfigSMTPColumnTLS, *e.TLS))
	}
	if e.FromAddress != nil {
		httpColumns = append(httpColumns, handler.NewCol(SMTPConfigSMTPColumnSenderAddress, *e.FromAddress))
	}
	if e.FromName != nil {
		httpColumns = append(httpColumns, handler.NewCol(SMTPConfigSMTPColumnSenderName, *e.FromName))
	}
	if e.ReplyToAddress != nil {
		httpColumns = append(httpColumns, handler.NewCol(SMTPConfigSMTPColumnReplyToAddress, *e.ReplyToAddress))
	}
	if e.Host != nil {
		httpColumns = append(httpColumns, handler.NewCol(SMTPConfigSMTPColumnHost, *e.Host))
	}
	if e.User != nil {
		httpColumns = append(httpColumns, handler.NewCol(SMTPConfigSMTPColumnUser, *e.User))
	}
	if e.Password != nil {
		httpColumns = append(httpColumns, handler.NewCol(SMTPConfigSMTPColumnPassword, *e.Password))
	}
	if len(httpColumns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			httpColumns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigSMTPColumnID, e.ID),
				handler.NewCond(SMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigSMTPTableSuffix),
		return nil, err
	}

	// Deal with old and unique SMTP settings (empty ID)
	id := e.ID
	if e.ID == "" {
		id = e.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
				handler.NewCol(SMTPConfigSMTPColumnPassword, e.Password),
			},
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, id),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigSMTPTableSuffix),
		),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, id),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
		return nil, err
	}

	// Deal with old and unique SMTP settings (empty ID)
	id := e.ID
	if e.ID == "" {
		id = e.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
				handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateInactive),
			},
			[]handler.Condition{
				handler.Not(handler.NewCond(SMTPConfigColumnID, e.ID)),
				handler.NewCond(SMTPConfigColumnState, domain.SMTPConfigStateActive),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
				handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateActive),
			},
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, id),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
		return nil, err
	}

	// Deal with old and unique SMTP settings (empty ID)
	id := e.ID
	if e.ID == "" {
		id = e.Aggregate().ResourceOwner
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnID, id),
			handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
		return nil, err
	}

	// Deal with old and unique SMTP settings (empty ID)
	id := e.ID
	if e.ID == "" {
		id = e.Aggregate().ResourceOwner
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnID, id),
			handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
package login

import (
	"net/http"

	"golang.org/x/text/language"
	OrgRegister        bool
}

func (l *Login) handleRegister(w http.ResponseWriter, r *http.Request) {
	data := new(registerFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
		l.renderError(w, r, authRequest, err)
		return
	}
	l.renderRegister(w, r, authRequest, data, nil)
}

func (l *Login) handleRegisterCheck(w http.ResponseWriter, r *http.Request) {
	data := new(registerFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
		l.renderError(w, r, authRequest, err)
		return
	}
	if data.Password != data.Password2 {
		err := zerrors.ThrowInvalidArgument(nil, "VIEW-KaGue", "Errors.User.Password.ConfirmationWrong")
		l.renderRegister(w, r, authRequest, data, err)
		return
	}

	resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()

	if authRequest != nil && authRequest.RequestedOrgID != "" && authRequest.RequestedOrgID != resourceOwner {
		resourceOwner = authRequest.RequestedOrgID
	}
	// For consistency with the external authentication flow,
	// the setMetadata() function is provided on the pre creation hook, for now,
	// like for the ExternalAuthentication flow.
		formData.Language = l.renderer.ReqLang(translator, r).String()
	}

	var resourceOwner string
	if authRequest != nil {
		resourceOwner = authRequest.RequestedOrgID
	}

	if resourceOwner == "" {
		resourceOwner = authz.GetInstance(r.Context()).DefaultOrganisationID()
	}

	data := registerData{
		baseData:         l.getBaseData(r, authRequest, translator, "RegistrationUser.Title", "RegistrationUser.Description", errID, errMessage),
		registerFormData: *formData,
			return CtxData{}, zerrors.ThrowUnauthenticated(errors.Join(err, sysTokenErr), "AUTH-7fs1e", "Errors.Token.Invalid")
		}
	}
	var projectID string
	var origins []string
	if clientID != "" {
		projectID, origins, err = t.ProjectIDAndOriginsByClientID(ctx, clientID)
		if err != nil {
			return CtxData{}, zerrors.ThrowPermissionDenied(err, "AUTH-GHpw2", "could not read projectid by clientid")
		}
		// We used to check origins for every token, but service users shouldn't be used publicly (native app / SPA).
		// Therefore, mostly won't send an origin and aren't able to configure them anyway.
		// For the current time we will only check origins for tokens issued to users through apps (code / implicit flow).
		if err := checkOrigin(ctx, origins); err != nil {
			return CtxData{}, err
		}
	}
	if orgID == "" && orgDomain == "" {
		orgID = resourceOwner
	}
	// System API calls don't have a resource owner
	if orgID != "" {
		orgID, err = t.ExistsOrg(ctx, orgID, orgDomain)
		if err != nil {
			return CtxData{}, zerrors.ThrowPermissionDenied(nil, "AUTH-Bs7Ds", "Organisation doesn't exist")
	}, nil
}

func SetCtxData(ctx context.Context, ctxData CtxData) context.Context {
	return context.WithValue(ctx, dataKey, ctxData)
}
