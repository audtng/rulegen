package rules

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

type mockXRV struct{}

func (m mockXRV) Validate(r io.Reader) error {
	return nil
}

var xrv = mockXRV{}

type AuthnRequest struct {
	ID string `xml:"ID,attr"`
}

type Response struct {
	Assertion string `xml:"Assertion"`
}

type Assertion struct {
	Subject string `xml:"Subject"`
}

type EntityDescriptor struct {
	EntityID string `xml:"entityID,attr"`
}

// 1. Vulnerable: xml.Unmarshal of SAML AuthnRequest without roundtrip validation
func parseAuthnRequest(data []byte) (*AuthnRequest, error) {
	req := &AuthnRequest{}
	// ruleid: saml-xml-roundtrip-auth-bypass
	if err := xml.Unmarshal(data, req); err != nil {
		return nil, err
	}
	return req, nil
}

// 2. Safe: xml.Unmarshal with prior xrv.Validate in if-statement
func parseAuthnRequestSafe(data []byte) (*AuthnRequest, error) {
	if err := xrv.Validate(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	req := &AuthnRequest{}
	// ok: saml-xml-roundtrip-auth-bypass
	if err := xml.Unmarshal(data, req); err != nil {
		return nil, err
	}
	return req, nil
}

// 3. Vulnerable: xml.Unmarshal of SAML Response without roundtrip validation
func parseSAMLResponse(data []byte) (*Response, error) {
	var resp Response
	// ruleid: saml-xml-roundtrip-auth-bypass
	err := xml.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &resp, nil
}

// 4. Safe: xml.Unmarshal with prior xrv.Validate assigned to error
func parseSAMLResponseSafe(data []byte) (*Response, error) {
	err := xrv.Validate(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var resp Response
	// ok: saml-xml-roundtrip-auth-bypass
	err = xml.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &resp, nil
}

// 5. Vulnerable: SAML Decoder without roundtrip validation
func decodeSAMLAssertion(r io.Reader) (*Assertion, error) {
	decoder := xml.NewDecoder(r)
	var assertion Assertion
	// ruleid: saml-xml-roundtrip-auth-bypass
	if err := decoder.Decode(&assertion); err != nil {
		return nil, err
	}
	return &assertion, nil
}

// 6. Safe: SAML Decoder with prior roundtrip validation
func decodeSAMLAssertionSafe(buf []byte) (*Assertion, error) {
	if err := xrv.Validate(bytes.NewReader(buf)); err != nil {
		return nil, err
	}
	decoder := xml.NewDecoder(bytes.NewReader(buf))
	var assertion Assertion
	// ok: saml-xml-roundtrip-auth-bypass
	if err := decoder.Decode(&assertion); err != nil {
		return nil, err
	}
	return &assertion, nil
}

// 7. Vulnerable: Roundtrip validator called but error ignored
func validateBufferIgnored(data []byte) error {
	// ruleid: saml-xml-roundtrip-auth-bypass
	_ = xrv.Validate(bytes.NewReader(data))
	return nil
}

// 8. Safe: Roundtrip validator called with error checked
func validateBufferSafe(data []byte) error {
	// ok: saml-xml-roundtrip-auth-bypass
	if err := xrv.Validate(bytes.NewReader(data)); err != nil {
		return err
	}
	return nil
}

// 9. Vulnerable: SAML metadata unmarshaling without validation
func parseMetadataVuln(data []byte) (*EntityDescriptor, error) {
	var entity EntityDescriptor
	// ruleid: saml-xml-roundtrip-auth-bypass
	if err := xml.Unmarshal(data, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

// 10. Safe: SAML metadata unmarshaling with validation
func parseMetadataSafe(data []byte) (*EntityDescriptor, error) {
	if err := xrv.Validate(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	var entity EntityDescriptor
	// ok: saml-xml-roundtrip-auth-bypass
	if err := xml.Unmarshal(data, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
