package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
// The Application Provisioning configuration is used for bootstrapping a device with an OMA DM account
// The paramenters here maps to the W7 application CSP
// https://learn.microsoft.com/en-us/windows/client-management/mdm/w7-application-csp
func NewApplicationProvisioningData(mdmEndpoint string) mdm_types.Characteristic {
	provDoc := newCharacteristic("APPLICATION", []mdm_types.Param{
		// The PROVIDER-ID parameter specifies the server identifier for a management server used in the current management session
		newParm("PROVIDER-ID", syncml.DocProvisioningAppProviderID, ""),
			newParm("AAUTHLEVEL", "APPSRV", ""),
			// DIGEST - Specifies that the SyncML DM 'syncml:auth-md5' authentication type.
			newParm("AAUTHTYPE", "DIGEST", ""),
			newParm("AAUTHNAME", "dummy", ""),
			newParm("AAUTHSECRET", "dummy", ""),
			newParm("AAUTHDATA", "nonce", ""),
		}, nil),
	})

	}

	// Getting the device provisioning information in the form of a WapProvisioningDoc
	deviceProvisioning, err := svc.getDeviceProvisioningInformation(ctx, secTokenMsg)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "device provisioning information")
	}
	//
	// This method also creates the relevant enrollment activity as it has
	// access to the device information.
	err = svc.storeWindowsMDMEnrolledDevice(ctx, userID, hostUUID, secTokenMsg)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "enrolled device information cannot be stored")
	}
	}

	// Checking if the incoming request is trusted
	err := svc.isTrustedRequest(ctx, reqSyncML, reqCerts)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "management request is not trusted")
	}

	// Getting the management response message
	resSyncMLmsg, err := svc.getManagementResponse(ctx, reqSyncML)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "management response message")
	}

	// Token is authorized
	svc.authz.SkipAuthorization(ctx)

	return resSyncMLmsg, nil
}

	return htmlBuf.String(), nil
}

// isTrustedRequest checks if the incoming request was sent from MDM enrolled device
func (svc *Service) isTrustedRequest(ctx context.Context, reqSyncML *fleet.SyncML, reqCerts []*x509.Certificate) error {
	if reqSyncML == nil {
		return fleet.NewInvalidArgumentError("syncml req message", "message is not present")
	}

	// Checking if calling request is coming from an already MDM enrolled device
	deviceID, err := reqSyncML.GetSource()
	if err != nil || deviceID == "" {
		return fmt.Errorf("invalid SyncML message %w", err)
	}

	enrolledDevice, err := svc.ds.MDMWindowsGetEnrolledDeviceWithDeviceID(ctx, deviceID)
	if err != nil || enrolledDevice == nil {
		return errors.New("device was not MDM enrolled")
	}

	// Check if TLS certs contains device ID on its common name
	if len(reqCerts) > 0 {
		for _, reqCert := range reqCerts {
			if strings.Contains(reqCert.Subject.CommonName, deviceID) {
				return nil
			}
		}
	}

	// TODO: Latest version of the MDM client stack don't populate TLS.PeerCertificates array
	// This is a temporary workaround to allow the management request to proceed
	// Transport-level security should be replaced for Application-level security
	// Transport-level security is defined in the MS-MDM spec in section 1.3.1
	// On the other hand, Application-level security is defined here
	// https://www.openmobilealliance.org/release/DM/V1_2_1-20080617-A/OMA-TS-DM_Security-V1_2_1-20080617-A.pdf
	// The initial values for Application-level security configuration are defined in the
	// WAP Profile blob that is sent to the device during the enrollment process. Example below
	//	<characteristic type="APPAUTH">
	//		<parm name="AAUTHLEVEL" value="CLIENT"/>
	//		<parm name="AAUTHTYPE" value="DIGEST"/>
	//		<parm name="AAUTHSECRET" value="2jsidqgffx"/>
	//		<parm name="AAUTHDATA" value="aGVsbG8gd29ybGQ="/>
	//	</characteristic>
	//	<characteristic type="APPAUTH">
	//		<parm name="AAUTHLEVEL" value="APPSRV"/>
	//		<parm name="AAUTHTYPE" value="DIGEST"/>
	//		<parm name="AAUTHNAME" value="43f8bf591b8557346021"/>
	//		<parm name="AAUTHSECRET" value="crbr3w2cab"/>
	//		<parm name="AAUTHDATA" value="aGVsbG8gd29ybGQ="/>
	//	</characteristic>

	if len(reqCerts) == 0 {
		return nil
	}

	return errors.New("calling device is not trusted")
}

// isFleetdPresentOnDevice checks if the device requires Fleetd to be deployed

// processIncomingMDMCmds process the incoming message from the device
// It will return the list of operations that need to be sent to the device
func (svc *Service) processIncomingMDMCmds(ctx context.Context, deviceID string, reqMsg *fleet.SyncML) ([]*fleet.SyncMLCmd, error) {
	var responseCmds []*fleet.SyncMLCmd

	// Get the incoming MessageID
	reqMessageID, err := reqMsg.GetMessageID()
	if err != nil {
		return nil, fmt.Errorf("get incoming msg: %w", err)
	}

	// Acknowledge the message header
	// msgref is always 0 for the header
	if err = reqMsg.IsValidHeader(); err == nil {
		ackMsg := NewSyncMLCmdStatus(reqMessageID, "0", syncml.SyncMLHdrName, syncml.CmdStatusOK)
		responseCmds = append(responseCmds, ackMsg)
	}
		return nil, err
	}

	enrichedSyncML := fleet.NewEnrichedSyncML(reqMsg)
	if enrichedSyncML.HasCommands() {
		if err := svc.ds.MDMWindowsSaveResponse(ctx, deviceID, enrichedSyncML, topLevelExists); err != nil {
			return nil, fmt.Errorf("store incoming msgs: %w", err)
		}
	}

	return responseCmds, nil
}

// getManagementResponse returns a valid SyncML response message
func (svc *Service) getManagementResponse(ctx context.Context, reqMsg *fleet.SyncML) (*mdm_types.SyncML, error) {
	if reqMsg == nil {
		return nil, fleet.NewInvalidArgumentError("syncml req message", "message is not present")
	}
	}

	// Process the incoming MDM protocol commands and get the response MDM protocol commands
	resIncomingCmds, err := svc.processIncomingMDMCmds(ctx, deviceID, reqMsg)
	if err != nil {
		return nil, fmt.Errorf("message processing error %w", err)
	}

	// Process the pending operations and get the MDM response protocol commands
	resPendingCmds, err := svc.getPendingMDMCmds(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("message processing error %w", err)
	}

	// Create the response SyncML message
// This information is used to configure the device management client
// See section 2.2.9.1 for more details on the XML provision schema used here
// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-mde2/35e1aca6-1b8a-48ba-bbc0-23af5d46907a
func (svc *Service) getDeviceProvisioningInformation(ctx context.Context, secTokenMsg *fleet.RequestSecurityToken) (string, error) {
	// Getting the HW DeviceID from the RequestSecurityToken msg
	reqHWDeviceID, err := GetContextItem(secTokenMsg, syncml.ReqSecTokenContextItemHWDevID)
	if err != nil {
		return "", err
	}

	// Getting the EnrollmentType information from the RequestSecurityToken msg
	reqEnrollType, err := GetContextItem(secTokenMsg, syncml.ReqSecTokenContextItemEnrollmentType)
	if err != nil {
		return "", err
	}

	// Getting the BinarySecurityToken from the RequestSecurityToken msg
	binSecurityTokenData, err := secTokenMsg.GetBinarySecurityTokenData()
	if err != nil {
		return "", err
	}

	// Getting the BinarySecurityToken type from the RequestSecurityToken msg
	binSecurityTokenType, err := secTokenMsg.GetBinarySecurityTokenType()
	if err != nil {
		return "", err
	}

	// Getting the client CSR request from the device
	clientCSR, err := microsoft_mdm.GetClientCSR(binSecurityTokenData, binSecurityTokenType)
	if err != nil {
		return "", err
	}

	// Getting the signed, DER-encoded certificate bytes and its uppercased, hex-endcoded SHA1 fingerprint
	rawSignedCertDER, rawSignedCertFingerprint, err := svc.SignMDMMicrosoftClientCSR(ctx, reqHWDeviceID, clientCSR)
	if err != nil {
		return "", err
	}

	// Preparing client certificate and identity certificate information to be sent to the Windows MDM Enrollment Client
	// Preparing the provisioning information that includes the location of the Device Management Service (DMS)
	appCfg, err := svc.ds.AppConfig(ctx)
	if err != nil {
		return "", err
	}

	// Getting the MS-MDM management URL to provision the device
	urlManagementEndpoint, err := microsoft_mdm.ResolveWindowsMDMManagement(appCfg.ServerSettings.ServerURL)
	if err != nil {
		return "", err
	}

	// Preparing the Application Provisioning information
	appConfigProvisioningData := NewApplicationProvisioningData(urlManagementEndpoint)

	// Preparing the DM Client Provisioning information
	appDMClientProvisioningData := NewDMClientProvisioningData()
	provDoc := NewProvisioningDoc(certStoreProvisioningData, appConfigProvisioningData, appDMClientProvisioningData)
	encodedProvDoc, err := provDoc.GetEncodedB64Representation()
	if err != nil {
		return "", err
	}

	return encodedProvDoc, nil
}

// storeWindowsMDMEnrolledDevice stores the device information to the list of MDM enrolled devices
func (svc *Service) storeWindowsMDMEnrolledDevice(ctx context.Context, userID string, hostUUID string, secTokenMsg *fleet.RequestSecurityToken) error {
	const (
		error_tag = "windows MDM enrolled storage: "
	)

	// Getting the Windows Enrolled Device Information
	enrolledDevice := &fleet.MDMWindowsEnrolledDevice{
		MDMDeviceID:            reqDeviceID,
		MDMHardwareID:          reqHWDevID,
		MDMDeviceState:         microsoft_mdm.MDMDeviceStateEnrolled,
		MDMDeviceType:          reqDeviceType,
		MDMDeviceName:          reqDeviceName,
		MDMEnrollType:          reqEnrollType,
		MDMEnrollUserID:        userID, // This could be Host UUID or UPN email
		MDMEnrollProtoVersion:  reqEnrollVersion,
		MDMEnrollClientVersion: reqAppVersion,
		MDMNotInOOBE:           reqNotInOOBE,
		HostUUID:               hostUUID,
	}

	if err := svc.ds.MDMWindowsInsertEnrolledDevice(ctx, enrolledDevice); err != nil {

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	jwtSigningKey *rsa.PrivateKey
	// jwtSigningKeyID is the ID to report in the header for the signing key
	jwtSigningKeyID string
}

// This is a test-only enrollment type to force erroneous behavior.
		fmt.Println(string(rawXMLReq))
	}

	// TODO: this request works because we're allowing devices without
	// certificates to communicate with the server. We will need to include the
	// certificate we generated during enrollment when we fix that.
	managementResp, err := c.request(microsoft_mdm.MDE2ManagementPath, rawXMLReq)
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
	c.lastManagementResp = &syncML

	cmds := make(map[string]fleet.ProtoCmdOperation)
	for _, p := range c.lastManagementResp.GetOrderedCmds() {
	return cmds, nil
}

func (c *TestWindowsMDMClient) SendResponse() (map[string]fleet.ProtoCmdOperation, error) {
	// Get SessionID
	sessionID, err := c.lastManagementResp.GetSessionID()
		Target: &fleet.LocURI{
			LocURI: ptr.String(c.fleetServerURL + microsoft_mdm.MDE2ManagementPath),
		},
	}

	// iterate over mocked responses and append them to the SyncML message
	return c.doManagementReq(xmlReq)
}

// AppendResponse sets a response for a specific command UUID.
func (c *TestWindowsMDMClient) AppendResponse(op fleet.SyncMLCmd) {
	c.queuedCommandResponses[op.CmdID.Value] = op
		return fmt.Errorf("enroll request returned SOAP fault: %s", string(body))
	}

	return nil
}


	return binarySecToken, tokenValueType, nil
}

const (
	WINDOWS_SCEP_LOC_URI_PART = "/Vendor/MSFT/ClientCertificateInstall/SCEP"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// SoapResponse is the Soap Envelope Response type for MS-MDE2 responses from the server
// This envelope XML message is composed by a mandatory SOAP envelope, a SOAP header, and a SOAP body
type SoapResponse struct {
	XMLName xml.Name       `xml:"s:Envelope"`
	XMLNSS  string         `xml:"xmlns:s,attr"`
	XMLNSA  string         `xml:"xmlns:a,attr"`
	XMLNSU  *string        `xml:"xmlns:u,attr,omitempty"`
	Header  ResponseHeader `xml:"s:Header"`
	Body    BodyResponse   `xml:"s:Body"`
}

// SoapRequest is the Soap Envelope Request type for MS-MDE2 responses to the server
/// Contains the information of the enrolled Windows host

type MDMWindowsEnrolledDevice struct {
	ID                     uint      `db:"id"`
	HostUUID               string    `db:"host_uuid"`
	MDMDeviceID            string    `db:"mdm_device_id"`
	MDMHardwareID          string    `db:"mdm_hardware_id"`
	MDMDeviceState         string    `db:"device_state"`
	MDMDeviceType          string    `db:"device_type"`
	MDMDeviceName          string    `db:"device_name"`
	MDMEnrollType          string    `db:"enroll_type"`
	MDMEnrollUserID        string    `db:"enroll_user_id"`
	MDMEnrollProtoVersion  string    `db:"enroll_proto_version"`
	MDMEnrollClientVersion string    `db:"enroll_client_version"`
	MDMNotInOOBE           bool      `db:"not_in_oobe"`
	CreatedAt              time.Time `db:"created_at"`
	UpdatedAt              time.Time `db:"updated_at"`
}

func (e MDMWindowsEnrolledDevice) AuthzType() string {
	Target    *LocURI  `xml:"Target,omitempty"`
	Source    *LocURI  `xml:"Source,omitempty"`
	Meta      *MetaHdr `xml:"Meta,omitempty"`
}

type MetaHdr struct {
	MaxMsgSize *string `xml:"MaxMsgSize,omitempty"`
}

// ProtoCmds contains a slice of SyncML protocol commands
type ProtoCmds []SyncMLCmd


// Protocol Command
type SyncMLCmd struct {
	XMLName xml.Name  `xml:",omitempty"`
	CmdID   CmdID     `xml:"CmdID"`
	MsgRef  *string   `xml:"MsgRef,omitempty"`
	CmdRef  *string   `xml:"CmdRef,omitempty"`
	Cmd     *string   `xml:"Cmd,omitempty"`
	Data    *string   `xml:"Data,omitempty"`
	Items   []CmdItem `xml:"Item,omitempty"`

	// ReplaceCommands is a catch-all for any nested <Replace> commands,
	// which can be found under <Atomic> elements.
	ExecCommands []SyncMLCmd `xml:"Exec,omitempty"`
}

// ParseWindowsMDMCommand parses the raw XML as a single Windows MDM command.
// A single <Exec> command is accepted as input.
func ParseWindowsMDMCommand(rawXMLCmd []byte) (*SyncMLCmd, error) {
