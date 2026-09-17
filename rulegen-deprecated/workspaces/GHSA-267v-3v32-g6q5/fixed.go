package main


import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"

	"github.com/beevik/etree"
// SOAPBinding is the official URN for the SOAP binding (transport)
const SOAPBinding = "urn:oasis:names:tc:SAML:2.0:bindings:SOAP"

// SOAPBindingV1 is the URN for the SOAP binding in SAML 1.0
const SOAPBindingV1 = "urn:oasis:names:tc:SAML:1.0:bindings:SOAP-binding"

// EntitiesDescriptor represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-metadata-2.0-os.pdf §2.3.1
	ResponseLocation string `xml:"ResponseLocation,attr,omitempty"`
}

func checkEndpointLocation(binding string, location string) (string, error) {
	// Within the SAML standard, the complex type EndpointType describes a
	// SAML protocol binding endpoint at which a SAML entity can be sent
	// protocol messages. In particular, the location of an endpoint type is
	// defined as follows in the Metadata for the OASIS Security Assertion
	// Markup Language (SAML) V2.0 - 2.2.2 Complex Type EndpointType:
	//
	//   Location [Required] A required URI attribute that specifies the
	//   location of the endpoint. The allowable syntax of this URI depends
	//   on the protocol binding.
	switch binding {
	case HTTPPostBinding,
		HTTPRedirectBinding,
		HTTPArtifactBinding,
		SOAPBinding,
		SOAPBindingV1:
		locationURL, err := url.Parse(location)
		if err != nil {
			return "", fmt.Errorf("invalid url %q: %w", location, err)
		}
		switch locationURL.Scheme {
		case "http", "https":
		// ok
		default:
			return "", fmt.Errorf("invalid url scheme %q for binding %q",
				locationURL.Scheme, binding)
		}
	default:
		// We don't know what form location should take, but the protocol
		// requires that we validate its syntax.
		//
		// In practice, lots of metadata contains random bindings, for example
		// "urn:mace:shibboleth:1.0:profiles:AuthnRequest" from our own test suite.
		//
		// We can't fail, but we also can't allow a location parameter whose syntax we
		// cannot verify. The least-bad course of action here is to set location to
		// and empty string, and hope the caller doesn't care need it.
		location = ""
	}

	return location, nil
}

// UnmarshalXML implements xml.Unmarshaler
func (m *Endpoint) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias Endpoint
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := d.DecodeElement(aux, &start); err != nil {
		return err
	}

	var err error
	m.Location, err = checkEndpointLocation(m.Binding, m.Location)
	if err != nil {
		return err
	}
	if m.ResponseLocation != "" {
		m.ResponseLocation, err = checkEndpointLocation(m.Binding, m.ResponseLocation)
		if err != nil {
			return err
		}
	}

	return nil
}

// IndexedEndpoint represents the SAML IndexedEndpointType object.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-metadata-2.0-os.pdf §2.2.3
	IsDefault        *bool   `xml:"isDefault,attr"`
}

// UnmarshalXML implements xml.Unmarshaler
func (m *IndexedEndpoint) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias IndexedEndpoint
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := d.DecodeElement(aux, &start); err != nil {
		return err
	}

	var err error
	m.Location, err = checkEndpointLocation(m.Binding, m.Location)
	if err != nil {
		return err
	}
	if m.ResponseLocation != nil {
		responseLocation, err := checkEndpointLocation(m.Binding, *m.ResponseLocation)
		if err != nil {
			return err
		}
		if responseLocation != "" {
			m.ResponseLocation = &responseLocation
		} else {
			m.ResponseLocation = nil
		}
	}

	return nil
}

// SSODescriptor represents the SAML complex type SSODescriptor
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-metadata-2.0-os.pdf §2.4.2
