package main

	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	xrv "github.com/mattermost/xml-roundtrip-validator"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
	"github.com/russellhaering/goxmldsig/etreeutils"
		Response: string(decodedResponseXML),
	}

	// ensure that the response XML is well formed before we parse it
	if err := xrv.Validate(bytes.NewReader(decodedResponseXML)); err != nil {
		retErr.PrivateErr = fmt.Errorf("invalid xml: %s", err)
		return nil, retErr
	}

	// do some validation first before we decrypt
	resp := Response{}
	if err := xml.Unmarshal(decodedResponseXML, &resp); err != nil {
		retErr.PrivateErr = fmt.Errorf("cannot unmarshal response: %s", err)
		return nil, retErr
	}
		}
		retErr.Response = string(plaintextAssertion)

		// TODO(ross): add test case for this
		if err := xrv.Validate(bytes.NewReader(plaintextAssertion)); err != nil {
			retErr.PrivateErr = fmt.Errorf("plaintext response contains invalid XML: %s", err)
			return nil, retErr
		}

		doc = etree.NewDocument()
		if err := doc.ReadFromBytes(plaintextAssertion); err != nil {
			retErr.PrivateErr = fmt.Errorf("cannot parse plaintext response %v", err)
		}

		assertion = &Assertion{}
		// Note: plaintextAssertion is known to be safe to parse because
		// plaintextAssertion is unmodified from when xrv.Validate() was called above.
		if err := xml.Unmarshal(plaintextAssertion, assertion); err != nil {
			retErr.PrivateErr = err
			return nil, retErr
		return fmt.Errorf("unable to parse base64: %s", err)
	}

	// TODO(ross): add test case for this (SLO does not have tests right now)
	if err := xrv.Validate(bytes.NewReader(rawResponseBuf)); err != nil {
		return fmt.Errorf("response contains invalid XML: %s", err)
	}

	var resp LogoutResponse
	if err := xml.Unmarshal(rawResponseBuf, &resp); err != nil {
		return fmt.Errorf("cannot unmarshal response: %s", err)
	}
		return fmt.Errorf("unable to parse base64: %s", err)
	}

	gr, err := ioutil.ReadAll(flate.NewReader(bytes.NewBuffer(rawResponseBuf)))
	if err != nil {
		return err
	}

	if err := xrv.Validate(bytes.NewReader(gr)); err != nil {
		return err
	}

	decoder := xml.NewDecoder(bytes.NewReader(gr))

	var resp LogoutResponse

	}

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(bytes.NewReader(gr)); err != nil {
		return err
	}

package samlidp

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"

	xrv "github.com/mattermost/xml-roundtrip-validator"

	"github.com/crewjam/saml"
)
}

func getSPMetadata(r io.Reader) (spMetadata *saml.EntityDescriptor, err error) {
	var data []byte
	if data, err = ioutil.ReadAll(r); err != nil {
		return nil, err
	}

	spMetadata = &saml.EntityDescriptor{}
	if err := xrv.Validate(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(data, &spMetadata); err != nil {
		if err.Error() == "expected element type <EntityDescriptor> but have <EntitiesDescriptor>" {
			entities := &saml.EntitiesDescriptor{}
			if err := xml.Unmarshal(data, &entities); err != nil {
				return nil, err
			}

package samlsp

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"net/url"

	"github.com/crewjam/httperr"
	xrv "github.com/mattermost/xml-roundtrip-validator"

	"github.com/crewjam/saml"
)
// <EntityDescriptor>.
func ParseMetadata(data []byte) (*saml.EntityDescriptor, error) {
	entity := &saml.EntityDescriptor{}

	if err := xrv.Validate(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	err := xml.Unmarshal(data, entity)

	// this comparison is ugly, but it is how the error is generated in encoding/xml
