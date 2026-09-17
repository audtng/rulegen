package main

import (
	"bytes"
	"context"
	"crypto/md5" //nolint:gosec // Windows MDM Auth uses MD5
	"crypto/x509"
	"database/sql"
	"encoding/base64"
// The Application Provisioning configuration is used for bootstrapping a device with an OMA DM account
// The paramenters here maps to the W7 application CSP
// https://learn.microsoft.com/en-us/windows/client-management/mdm/w7-application-csp
func NewApplicationProvisioningData(mdmEndpoint string, username string, secret string) mdm_types.Characteristic {
	provDoc := newCharacteristic("APPLICATION", []mdm_types.Param{
		// The PROVIDER-ID parameter specifies the server identifier for a management server used in the current management session
		newParm("PROVIDER-ID", syncml.DocProvisioningAppProviderID, ""),
			newParm("AAUTHLEVEL", "APPSRV", ""),
			// DIGEST - Specifies that the SyncML DM 'syncml:auth-md5' authentication type.
			newParm("AAUTHTYPE", "DIGEST", ""),
			newParm("AAUTHNAME", username, ""),
			newParm("AAUTHSECRET", secret, ""),
			newParm("AAUTHDATA", "nonce", ""), // We don't care about setting the first round nonce, as when the device checks in we will prompt the credentials and pass a new nonce.
		}, nil),
	})

	}

	// Getting the device provisioning information in the form of a WapProvisioningDoc
	deviceProvisioning, credentialsHash, err := svc.getDeviceProvisioningInformation(ctx, secTokenMsg)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "device provisioning information")
	}
	//
	// This method also creates the relevant enrollment activity as it has
	// access to the device information.
	err = svc.storeWindowsMDMEnrolledDevice(ctx, userID, hostUUID, secTokenMsg, credentialsHash)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "enrolled device information cannot be stored")
	}
	}

	// Checking if the incoming request is trusted
	requestAuthState, err := svc.isTrustedRequest(ctx, reqSyncML, reqCerts)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "management request is not trusted")
	}

	// Token is authorized
	svc.authz.SkipAuthorization(ctx)

	if requestAuthState == RequestAuthStateRekey {
		// If we signalled to rekey the device, we short-circuit into a rekey flow.
		return svc.rekeyWindowsDevice(ctx, reqSyncML)
	}

	// Getting the management response message
	resSyncMLmsg, err := svc.getManagementResponse(ctx, reqSyncML, requestAuthState)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "management response message")
	}

	return resSyncMLmsg, nil
}

	return htmlBuf.String(), nil
}

type requestAuthState int

const (
	RequestAuthStateUntrusted requestAuthState = iota
	RequestAuthStateUnauthorized
	RequestAuthStateChallenge
	RequestAuthStateRekey
	RequestAuthStateTrusted
)

// isTrustedRequest checks if the incoming request was sent from MDM enrolled device
// It returns a boolean if we should challenge the device and an error if the request/device calling is not trusted
func (svc *Service) isTrustedRequest(ctx context.Context, reqSyncML *fleet.SyncML, reqCerts []*x509.Certificate) (requestAuthState, error) {
	if reqSyncML == nil {
		return RequestAuthStateUntrusted, fleet.NewInvalidArgumentError("syncml req message", "message is not present")
	}

	// Checking if calling request is coming from an already MDM enrolled device
	deviceID, err := reqSyncML.GetSource()
	if err != nil || deviceID == "" {
		return RequestAuthStateUntrusted, fmt.Errorf("invalid SyncML message %w", err)
	}

	enrolledDevice, err := svc.ds.MDMWindowsGetEnrolledDeviceWithDeviceID(ctx, deviceID)
	if err != nil || enrolledDevice == nil {
		return RequestAuthStateUntrusted, errors.New("device was not MDM enrolled")
	}

	// Check if TLS certs contains device ID on its common name
	if len(reqCerts) > 0 {
		for _, reqCert := range reqCerts {
			if strings.Contains(reqCert.Subject.CommonName, deviceID) {
				return RequestAuthStateTrusted, nil
			}
		}
	}

	if !enrolledDevice.CredentialsAcknowledged && enrolledDevice.CredentialsHash == nil {
		// Device has not gotten new credentials, rekey the device only once
		return RequestAuthStateRekey, nil
	}

	if reqSyncML.SyncHdr.Cred == nil {
		// No certs, but no credentials present - challenge the device
		return RequestAuthStateChallenge, nil
	}

	// Extract the last nonce used to generate the credentials hash
	nonce, err := svc.keyValueStore.Get(ctx, fleet.WindowsMDMAuthNoncePrefix+deviceID)
	if err != nil {
		return RequestAuthStateUntrusted, ctxerr.Wrap(ctx, err, "get device nonce from kv store")
	}

	if nonce == nil || *nonce == "" {
		// Challenge the device if nonce is missing, which will send a new nonce and store it
		return RequestAuthStateChallenge, nil
	}

	// Credentials are present, validate it
	credFormat := reqSyncML.SyncHdr.Cred.Meta.Format
	credType := reqSyncML.SyncHdr.Cred.Meta.Type
	credData := reqSyncML.SyncHdr.Cred.Data

	if credFormat == nil || credType == nil || credFormat.Content == nil || credType.Content == nil {
		return RequestAuthStateUntrusted, errors.New("SyncML credentials format or type is missing")
	}

	if *credFormat.Content != syncml.AuthB64Format || *credType.Content != syncml.AuthMD5 {
		return RequestAuthStateUntrusted, errors.New("SyncML credentials format or type is invalid")
	}

	// MD5 auth digest, which includes (username:password):nonce
	// Where username:password is hashed and b64 encoded, and then further hased with the nonce and finally b64 encoded for transport
	// https://www.openmobilealliance.org/release/DM/V1_2_1-20080617-A/OMA-TS-DM_Security-V1_2_1-20080617-A.pdf Chaper (5.3)
	receivedDigestHash, err := base64.StdEncoding.DecodeString(credData)
	if err != nil {
		return RequestAuthStateUntrusted, ctxerr.Wrap(ctx, err, "decode SyncML credentials data")
	}

	encodedCredentialsHash := base64.StdEncoding.EncodeToString(*enrolledDevice.CredentialsHash)
	expectedDigest := fmt.Sprintf("%s:%s", encodedCredentialsHash, *nonce)
	expectedDigestHash := md5.Sum([]byte(expectedDigest)) //nolint:gosec // Windows MDM Auth uses MD5

	if !bytes.Equal(receivedDigestHash, expectedDigestHash[:]) {
		// Credentials do not match what we expect
		return RequestAuthStateUnauthorized, nil
	}

	// We verified the username, password and nonce match what we expect, so we can ack the rekeyed credentials
	if !enrolledDevice.CredentialsAcknowledged {
		err = svc.ds.MDMWindowsAcknowledgeEnrolledDeviceCredentials(ctx, enrolledDevice.MDMDeviceID)
		if err != nil {
			return RequestAuthStateUntrusted, ctxerr.Wrap(ctx, err, "mark device credentials as acknowledged")
		}
	}

	return RequestAuthStateTrusted, nil
}

func (svc *Service) rekeyWindowsDevice(ctx context.Context, reqSyncML *fleet.SyncML) (*fleet.SyncML, error) {
	if reqSyncML == nil {
		return nil, fleet.NewInvalidArgumentError("syncml req message", "message is not present")
	}

	// Getting the device ID from the SyncML message
	deviceID, err := reqSyncML.GetSource()
	if err != nil || deviceID == "" {
		return nil, fmt.Errorf("invalid SyncML message %w", err)
	}

	enrolledDevice, err := svc.ds.MDMWindowsGetEnrolledDeviceWithDeviceID(ctx, deviceID)
	if err != nil || enrolledDevice == nil {
		return nil, errors.New("device was not MDM enrolled")
	}

	username := deviceID
	password := uuid.NewString()
	credentialsHash := md5.Sum(fmt.Appendf([]byte{}, "%s:%s", username, password)) //nolint:gosec // Windows MDM Auth uses MD5

	// Store the new credentials hash and mark that the device has not acknowledged them yet
	err = svc.ds.MDMWindowsUpdateEnrolledDeviceCredentials(ctx, enrolledDevice.MDMDeviceID, credentialsHash[:])
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "update enrolled device credentials")
	}

	// Queue two replace commands to update the credentials on the device
	accountUid := "0x0000800cABF8CA7D87C0FA21BB21B896AE2C8468CA07710326415228ACF2392EA1327077"
	usernameReplace := newSyncMLCmdText("Replace", fmt.Sprintf("./SyncML/DMAcc/%s/AppAuth/CLCRED/AAuthName", accountUid), username)
	usernameReplace.CmdID = mdm_types.CmdID{
		IncludeFleetComment: true,
		Value:               "rekey-credentials-username",
	}
	passwordReplace := newSyncMLCmdText("Replace", fmt.Sprintf("./SyncML/DMAcc/%s/AppAuth/CLCRED/AAuthSecret", accountUid), password)
	passwordReplace.CmdID = mdm_types.CmdID{
		IncludeFleetComment: true,
		Value:               "rekey-credentials-password",
	}

	// Get the incoming MessageID
	reqMessageID, err := reqSyncML.GetMessageID()
	if err != nil {
		return nil, fmt.Errorf("get incoming msg: %w", err)
	}

	// We only create a response here, and does not persist it to the DB to avoid saving the secrets in plain text
	return svc.createResponseSyncML(ctx, reqSyncML, []*mdm_types.SyncMLCmd{
		NewSyncMLCmdStatus(reqMessageID, "0", syncml.SyncMLHdrName, syncml.CmdStatusOK), // We need to ack the incoming message to send rekey commands
		usernameReplace,
		passwordReplace,
	})
}

// isFleetdPresentOnDevice checks if the device requires Fleetd to be deployed

// processIncomingMDMCmds process the incoming message from the device
// It will return the list of operations that need to be sent to the device
func (svc *Service) processIncomingMDMCmds(ctx context.Context, deviceID string, reqMsg *fleet.SyncML, requestAuthState requestAuthState) ([]*fleet.SyncMLCmd, error) {
	var responseCmds []*fleet.SyncMLCmd

	saveResponse := func(topLevelExists []string) error {
		enrichedSyncML := fleet.NewEnrichedSyncML(reqMsg)
		if enrichedSyncML.HasCommands() {
			if err := svc.ds.MDMWindowsSaveResponse(ctx, deviceID, enrichedSyncML, topLevelExists); err != nil {
				return fmt.Errorf("store incoming msgs: %w", err)
			}
		}
		return nil
	}

	// Get the incoming MessageID
	reqMessageID, err := reqMsg.GetMessageID()
	if err != nil {
		return nil, fmt.Errorf("get incoming msg: %w", err)
	}

	if requestAuthState == RequestAuthStateChallenge || requestAuthState == RequestAuthStateUnauthorized {
		nonce := uuid.NewString() // using UUID as nonce since it has 122 bits of entropy
		base64Nonce := base64.StdEncoding.EncodeToString([]byte(nonce))
		err := svc.keyValueStore.Set(ctx, fleet.WindowsMDMAuthNoncePrefix+deviceID, nonce, 5*time.Minute)
		if err != nil {
			return nil, ctxerr.Wrap(ctx, err, "store device nonce in kv store")
		}

		status := syncml.CmdStatusAuthenticationRequired
		if requestAuthState == RequestAuthStateUnauthorized {
			status = syncml.CmdStatusInvalidCredentials
		}

		ackMsg := NewSyncMLCmdStatus(reqMessageID, "0", syncml.SyncMLHdrName, status)
		ackMsg.Chal = &fleet.SyncMLChallenge{
			Meta: fleet.ChallengeMeta{
				NextNonce: fleet.MetaAttr{
					XMLNS:   syncml.SyncMLMetaNamespace,
					Content: &base64Nonce,
				},
				Meta: fleet.Meta{
					Type: &fleet.MetaAttr{
						XMLNS:   syncml.SyncMLMetaNamespace,
						Content: ptr.String(syncml.AuthMD5),
					},
					Format: &fleet.MetaAttr{
						XMLNS:   syncml.SyncMLMetaNamespace,
						Content: ptr.String(syncml.AuthB64Format),
					},
				},
			},
		}

		responseCmds = append(responseCmds, ackMsg)
		err = saveResponse([]string{})
		if err != nil {
			return nil, err
		}
		return responseCmds, nil
	}

	if requestAuthState != RequestAuthStateTrusted {
		return nil, errors.New("untrusted request cannot be processed")
	}

	// Acknowledge the message header
	// msgref is always 0 for the header
	if err = reqMsg.IsValidHeader(); err == nil {
		// We always return 200 here, we could also return 212 to indicate we don't need the credentials every time,
		// but our logic is built around it being present on each request
		ackMsg := NewSyncMLCmdStatus(reqMessageID, "0", syncml.SyncMLHdrName, syncml.CmdStatusOK)
		responseCmds = append(responseCmds, ackMsg)
	}
		return nil, err
	}

	err = saveResponse(topLevelExists)
	if err != nil {
		return nil, err
	}

	return responseCmds, nil
}

// getManagementResponse returns a valid SyncML response message
func (svc *Service) getManagementResponse(ctx context.Context, reqMsg *fleet.SyncML, requestAuthState requestAuthState) (*mdm_types.SyncML, error) {
	if reqMsg == nil {
		return nil, fleet.NewInvalidArgumentError("syncml req message", "message is not present")
	}
	}

	// Process the incoming MDM protocol commands and get the response MDM protocol commands
	resIncomingCmds, err := svc.processIncomingMDMCmds(ctx, deviceID, reqMsg, requestAuthState)
	if err != nil {
		return nil, fmt.Errorf("message processing error %w", err)
	}

	resPendingCmds := []*mdm_types.SyncMLCmd{}

	if requestAuthState == RequestAuthStateTrusted {
		// Process the pending operations and get the MDM response protocol commands
		pendingCmds, err := svc.getPendingMDMCmds(ctx, deviceID)
		if err != nil {
			return nil, fmt.Errorf("message processing error %w", err)
		}
		resPendingCmds = pendingCmds
	}

	// Create the response SyncML message
// This information is used to configure the device management client
// See section 2.2.9.1 for more details on the XML provision schema used here
// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-mde2/35e1aca6-1b8a-48ba-bbc0-23af5d46907a
func (svc *Service) getDeviceProvisioningInformation(ctx context.Context, secTokenMsg *fleet.RequestSecurityToken) (string, []byte, error) {
	reqDeviceID, err := GetContextItem(secTokenMsg, syncml.ReqSecTokenContextItemDeviceID)
	if err != nil {
		return "", nil, err
	}

	// Getting the HW DeviceID from the RequestSecurityToken msg
	reqHWDeviceID, err := GetContextItem(secTokenMsg, syncml.ReqSecTokenContextItemHWDevID)
	if err != nil {
		return "", nil, err
	}

	// Getting the EnrollmentType information from the RequestSecurityToken msg
	reqEnrollType, err := GetContextItem(secTokenMsg, syncml.ReqSecTokenContextItemEnrollmentType)
	if err != nil {
		return "", nil, err
	}

	// Getting the BinarySecurityToken from the RequestSecurityToken msg
	binSecurityTokenData, err := secTokenMsg.GetBinarySecurityTokenData()
	if err != nil {
		return "", nil, err
	}

	// Getting the BinarySecurityToken type from the RequestSecurityToken msg
	binSecurityTokenType, err := secTokenMsg.GetBinarySecurityTokenType()
	if err != nil {
		return "", nil, err
	}

	// Getting the client CSR request from the device
	clientCSR, err := microsoft_mdm.GetClientCSR(binSecurityTokenData, binSecurityTokenType)
	if err != nil {
		return "", nil, err
	}

	// Getting the signed, DER-encoded certificate bytes and its uppercased, hex-endcoded SHA1 fingerprint
	rawSignedCertDER, rawSignedCertFingerprint, err := svc.SignMDMMicrosoftClientCSR(ctx, reqHWDeviceID, clientCSR)
	if err != nil {
		return "", nil, err
	}

	// Preparing client certificate and identity certificate information to be sent to the Windows MDM Enrollment Client
	// Preparing the provisioning information that includes the location of the Device Management Service (DMS)
	appCfg, err := svc.ds.AppConfig(ctx)
	if err != nil {
		return "", nil, err
	}

	// Getting the MS-MDM management URL to provision the device
	urlManagementEndpoint, err := microsoft_mdm.ResolveWindowsMDMManagement(appCfg.ServerSettings.ServerURL)
	if err != nil {
		return "", nil, err
	}

	// generate username and password for device management service
	username := reqDeviceID
	password := uuid.NewString()
	credentialsHash := md5.Sum(fmt.Appendf(nil, "%s:%s", username, password)) //nolint:gosec // Windows MDM Auth uses MD5

	// Preparing the Application Provisioning information
	appConfigProvisioningData := NewApplicationProvisioningData(urlManagementEndpoint, username, password)

	// Preparing the DM Client Provisioning information
	appDMClientProvisioningData := NewDMClientProvisioningData()
	provDoc := NewProvisioningDoc(certStoreProvisioningData, appConfigProvisioningData, appDMClientProvisioningData)
	encodedProvDoc, err := provDoc.GetEncodedB64Representation()
	if err != nil {
		return "", nil, err
	}

	return encodedProvDoc, credentialsHash[:], nil
}

// storeWindowsMDMEnrolledDevice stores the device information to the list of MDM enrolled devices
func (svc *Service) storeWindowsMDMEnrolledDevice(ctx context.Context, userID string, hostUUID string, secTokenMsg *fleet.RequestSecurityToken, credentialsHash []byte) error {
	const (
		error_tag = "windows MDM enrolled storage: "
	)

	// Getting the Windows Enrolled Device Information
	enrolledDevice := &fleet.MDMWindowsEnrolledDevice{
		MDMDeviceID:             reqDeviceID,
		MDMHardwareID:           reqHWDevID,
		MDMDeviceState:          microsoft_mdm.MDMDeviceStateEnrolled,
		MDMDeviceType:           reqDeviceType,
		MDMDeviceName:           reqDeviceName,
		MDMEnrollType:           reqEnrollType,
		MDMEnrollUserID:         userID, // This could be Host UUID or UPN email
		MDMEnrollProtoVersion:   reqEnrollVersion,
		MDMEnrollClientVersion:  reqAppVersion,
		MDMNotInOOBE:            reqNotInOOBE,
		HostUUID:                hostUUID,
		CredentialsHash:         &credentialsHash,
		CredentialsAcknowledged: true,
	}

	if err := svc.ds.MDMWindowsInsertEnrolledDevice(ctx, enrolledDevice); err != nil {

import (
	"bytes"
	"crypto/md5" //nolint:gosec // Windows MDM Auth uses MD5
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	jwtSigningKey *rsa.PrivateKey
	// jwtSigningKeyID is the ID to report in the header for the signing key
	jwtSigningKeyID string

	username string
	password string
	nonce    string
}

// This is a test-only enrollment type to force erroneous behavior.
		fmt.Println(string(rawXMLReq))
	}

	sendRequest := func(req []byte) (*fleet.SyncML, error) {
		managementResp, err := c.request(microsoft_mdm.MDE2ManagementPath, req)
		if err != nil {
			return nil, err
		}

		rawXMLResp, err := io.ReadAll(managementResp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response body: %w", err)
		}

		if c.debug {
			fmt.Println("=============== management response ================")
			fmt.Println(string(rawXMLResp))
		}

		var syncML fleet.SyncML
		if err := xml.Unmarshal(rawXMLResp, &syncML); err != nil {
			return nil, fmt.Errorf("unmarshalling response body: %w", err)
		}

		return &syncML, nil
	}

	syncML, err := sendRequest(rawXMLReq)
	if err != nil {
		return nil, err
	}

	if username, password := c.isRekeyRequest(syncML); username != "" && password != "" {
		c.username = username
		c.password = password

		// We rekeyed, so we need to resend the original request
		syncML, err = sendRequest(rawXMLReq)
		if err != nil {
			return nil, err
		}
	}

	if shouldAuth, nonce := c.shouldAuth(syncML); shouldAuth {
		var reqSyncML fleet.SyncML
		if err := xml.Unmarshal(rawXMLReq, &reqSyncML); err != nil {
			return nil, fmt.Errorf("unmarshalling request body for auth: %w", err)
		}

		extractedNonce, _ := base64.StdEncoding.DecodeString(*nonce)
		c.nonce = string(extractedNonce)
		reqSyncML.SyncHdr.Cred = c.getCredHDR()

		// resend the request but now with credentials
		rawXMLReq, err = xml.MarshalIndent(reqSyncML, "", "\t")
		if err != nil {
			return nil, fmt.Errorf("serializing XML req with auth: %w", err)
		}

		if c.debug {
			fmt.Println("=============== management request with auth ================")
			fmt.Println(string(rawXMLReq))
		}

		syncML, err = sendRequest(rawXMLReq)
		if err != nil {
			return nil, err
		}
	}

	c.lastManagementResp = syncML

	cmds := make(map[string]fleet.ProtoCmdOperation)
	for _, p := range c.lastManagementResp.GetOrderedCmds() {
	return cmds, nil
}

func (c *TestWindowsMDMClient) isRekeyRequest(req *fleet.SyncML) (username string, password string) {
	for _, cmd := range req.GetOrderedCmds() {
		if cmd.Verb == fleet.CmdReplace && strings.Contains(cmd.Cmd.GetTargetURI(), "AAuthName") {
			username = cmd.Cmd.GetTargetData()
		} else if cmd.Verb == fleet.CmdReplace && strings.Contains(cmd.Cmd.GetTargetURI(), "AAuthSecret") {
			password = cmd.Cmd.GetTargetData()
		}
	}
	return
}

func (c *TestWindowsMDMClient) shouldAuth(req *fleet.SyncML) (bool, *string) {
	for _, cmd := range req.GetOrderedCmds() {
		if cmd.Verb == fleet.CmdStatus && cmd.Cmd.Chal != nil {
			return true, cmd.Cmd.Chal.Meta.NextNonce.Content
		}
	}
	return false, nil
}

func (c *TestWindowsMDMClient) SendResponse() (map[string]fleet.ProtoCmdOperation, error) {
	// Get SessionID
	sessionID, err := c.lastManagementResp.GetSessionID()
		Target: &fleet.LocURI{
			LocURI: ptr.String(c.fleetServerURL + microsoft_mdm.MDE2ManagementPath),
		},
		Cred: c.getCredHDR(),
	}

	// iterate over mocked responses and append them to the SyncML message
	return c.doManagementReq(xmlReq)
}

func (c *TestWindowsMDMClient) getCredHDR() *fleet.CredHdr {
	return &fleet.CredHdr{
		Meta: fleet.Meta{
			Type: &fleet.MetaAttr{
				XMLNS:   syncml.SyncMLMetaNamespace,
				Content: ptr.String(syncml.AuthMD5),
			},
			Format: &fleet.MetaAttr{
				XMLNS:   syncml.SyncMLMetaNamespace,
				Content: ptr.String(syncml.AuthB64Format),
			},
		},
		Data: c.hashedCredentials(),
	}
}

func (c *TestWindowsMDMClient) hashedCredentials() string {
	credentials := fmt.Sprintf("%s:%s", c.username, c.password)
	credentialsHash := md5.Sum([]byte(credentials)) //nolint:gosec // Windows MDM Auth uses MD5
	credentialsWithNonce := fmt.Sprintf("%s:%s", base64.StdEncoding.EncodeToString(credentialsHash[:]), c.nonce)
	digestHash := md5.Sum([]byte(credentialsWithNonce)) //nolint:gosec // Windows MDM Auth uses MD5
	return base64.StdEncoding.EncodeToString(digestHash[:])
}

// AppendResponse sets a response for a specific command UUID.
func (c *TestWindowsMDMClient) AppendResponse(op fleet.SyncMLCmd) {
	c.queuedCommandResponses[op.CmdID.Value] = op
		return fmt.Errorf("enroll request returned SOAP fault: %s", string(body))
	}

	var soapResponse fleet.SoapResponse
	if err := xml.Unmarshal(body, &soapResponse); err != nil {
		return fmt.Errorf("unmarshalling enroll response body: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(soapResponse.Body.RequestSecurityTokenResponseCollection.RequestSecurityTokenResponse.RequestedSecurityToken.BinarySecurityToken.Content)
	if err != nil {
		return fmt.Errorf("decoding enroll response binary security token: %w", err)
	}

	// strip xml header
	decoded = bytes.TrimPrefix(decoded, []byte(xml.Header))
	var provDoc fleet.WapProvisioningDoc
	if err := xml.Unmarshal(decoded, &provDoc); err != nil {
		return fmt.Errorf("unmarshalling enroll response provisioning doc: %w", err)
	}

Outer:
	for _, char := range provDoc.Characteristics {
		if char.Type != "APPLICATION" {
			continue
		}

		for _, appChar := range char.Characteristics {
			if appChar.Type != "APPAUTH" {
				continue
			}
			username := ""
			password := ""
			for _, appAuthParam := range appChar.Params {
				if appAuthParam.Name == "AAUTHNAME" {
					username = appAuthParam.Value
				}
				if appAuthParam.Name == "AAUTHSECRET" {
					password = appAuthParam.Value
				}
			}

			// We can do this since only the client credentials characteristic has both username and password, the other one only has password.
			if username != "" && password != "" {
				c.username = username
				c.password = password
				break Outer
			}
		}
	}

	return nil
}


	return binarySecToken, tokenValueType, nil
}

func (c *TestWindowsMDMClient) Unenroll() error {
	unenrollRequest := []byte(`
			 <SyncML xmlns="SYNCML:SYNCML1.2">
			<SyncHdr>
				<VerDTD>1.2</VerDTD>
				<VerProto>DM/1.2</VerProto>
				<SessionID>2</SessionID>
				<MsgID>1</MsgID>
				<Target>
				<LocURI>` + c.fleetServerURL + microsoft_mdm.MDE2ManagementPath + `</LocURI>
				</Target>
				<Source>
				<LocURI>` + c.DeviceID + `</LocURI>
				</Source>
			</SyncHdr>
			<SyncBody>
				<Alert>
				<CmdID>4</CmdID>
				<Data>1226</Data>
				<Item>
					<Meta>
					<Type xmlns="syncml:metinf">com.microsoft:mdm.unenrollment.userrequest</Type>
					<Format xmlns="syncml:metinf">int</Format>
					</Meta>
					<Data>1</Data>
				</Item>
				</Alert>
				<Final/>
			</SyncBody>
			</SyncML>`)

	_, err := c.doManagementReq(unenrollRequest)
	return err
}

const (
	WINDOWS_SCEP_LOC_URI_PART = "/Vendor/MSFT/ClientCertificateInstall/SCEP"
	WindowsMDMAuthNoncePrefix = "mwenonce:"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// SoapResponse is the Soap Envelope Response type for MS-MDE2 responses from the server
// This envelope XML message is composed by a mandatory SOAP envelope, a SOAP header, and a SOAP body
type SoapResponse struct {
	XMLName xml.Name       `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	XMLNSS  string         `xml:"xmlns:s,attr"`
	XMLNSA  string         `xml:"xmlns:a,attr"`
	XMLNSU  *string        `xml:"xmlns:u,attr,omitempty"`
	Header  ResponseHeader `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Body    BodyResponse   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

// SoapRequest is the Soap Envelope Request type for MS-MDE2 responses to the server
/// Contains the information of the enrolled Windows host

type MDMWindowsEnrolledDevice struct {
	ID                      uint      `db:"id"`
	HostUUID                string    `db:"host_uuid"`
	MDMDeviceID             string    `db:"mdm_device_id"`
	MDMHardwareID           string    `db:"mdm_hardware_id"`
	MDMDeviceState          string    `db:"device_state"`
	MDMDeviceType           string    `db:"device_type"`
	MDMDeviceName           string    `db:"device_name"`
	MDMEnrollType           string    `db:"enroll_type"`
	MDMEnrollUserID         string    `db:"enroll_user_id"`
	MDMEnrollProtoVersion   string    `db:"enroll_proto_version"`
	MDMEnrollClientVersion  string    `db:"enroll_client_version"`
	MDMNotInOOBE            bool      `db:"not_in_oobe"`
	CredentialsHash         *[]byte   `db:"credentials_hash"`
	CredentialsAcknowledged bool      `db:"credentials_acknowledged"`
	CreatedAt               time.Time `db:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"`
}

func (e MDMWindowsEnrolledDevice) AuthzType() string {
	Target    *LocURI  `xml:"Target,omitempty"`
	Source    *LocURI  `xml:"Source,omitempty"`
	Meta      *MetaHdr `xml:"Meta,omitempty"`
	Cred      *CredHdr `xml:"Cred,omitempty"`
}

type MetaHdr struct {
	MaxMsgSize *string `xml:"MaxMsgSize,omitempty"`
}

type CredHdr struct {
	Meta Meta   `xml:"Meta"`
	Data string `xml:"Data"`
}

// ProtoCmds contains a slice of SyncML protocol commands
type ProtoCmds []SyncMLCmd


// Protocol Command
type SyncMLCmd struct {
	XMLName xml.Name         `xml:",omitempty"`
	CmdID   CmdID            `xml:"CmdID"`
	MsgRef  *string          `xml:"MsgRef,omitempty"`
	CmdRef  *string          `xml:"CmdRef,omitempty"`
	Cmd     *string          `xml:"Cmd,omitempty"`
	Data    *string          `xml:"Data,omitempty"`
	Items   []CmdItem        `xml:"Item,omitempty"`
	Chal    *SyncMLChallenge `xml:"Chal,omitempty"`

	// ReplaceCommands is a catch-all for any nested <Replace> commands,
	// which can be found under <Atomic> elements.
	ExecCommands []SyncMLCmd `xml:"Exec,omitempty"`
}

type SyncMLChallenge struct {
	Meta ChallengeMeta `xml:"Meta"`
}

type ChallengeMeta struct {
	Meta
	NextNonce MetaAttr `xml:"NextNonce,omitempty"`
}

// ParseWindowsMDMCommand parses the raw XML as a single Windows MDM command.
// A single <Exec> command is accepted as input.
func ParseWindowsMDMCommand(rawXMLCmd []byte) (*SyncMLCmd, error) {
