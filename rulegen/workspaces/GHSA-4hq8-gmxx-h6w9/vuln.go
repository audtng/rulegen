package main

	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
	"github.com/russellhaering/goxmldsig/etreeutils"
		Response: string(decodedResponseXML),
	}

	// do some validation first before we decrypt
	resp := Response{}
	if err := xml.Unmarshal([]byte(decodedResponseXML), &resp); err != nil {
		retErr.PrivateErr = fmt.Errorf("cannot unmarshal response: %s", err)
		return nil, retErr
	}
		}
		retErr.Response = string(plaintextAssertion)

		doc = etree.NewDocument()
		if err := doc.ReadFromBytes(plaintextAssertion); err != nil {
			retErr.PrivateErr = fmt.Errorf("cannot parse plaintext response %v", err)
		}

		assertion = &Assertion{}
		if err := xml.Unmarshal(plaintextAssertion, assertion); err != nil {
			retErr.PrivateErr = err
			return nil, retErr
		return fmt.Errorf("unable to parse base64: %s", err)
	}

	var resp LogoutResponse

	if err := xml.Unmarshal(rawResponseBuf, &resp); err != nil {
		return fmt.Errorf("cannot unmarshal response: %s", err)
	}
		return fmt.Errorf("unable to parse base64: %s", err)
	}

	gr := flate.NewReader(bytes.NewBuffer(rawResponseBuf))

	decoder := xml.NewDecoder(gr)

	var resp LogoutResponse

	}

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(gr); err != nil {
		return err
	}

package samlidp

import (
	"errors"
	"io/ioutil"

	"encoding/xml"

	"io"

	"github.com/crewjam/saml"
)
}

func getSPMetadata(r io.Reader) (spMetadata *saml.EntityDescriptor, err error) {
	var bytes []byte

	if bytes, err = ioutil.ReadAll(r); err != nil {
		return nil, err
	}

	spMetadata = &saml.EntityDescriptor{}

	if err := xml.Unmarshal(bytes, &spMetadata); err != nil {
		if err.Error() == "expected element type <EntityDescriptor> but have <EntitiesDescriptor>" {
			entities := &saml.EntitiesDescriptor{}

			if err := xml.Unmarshal(bytes, &entities); err != nil {
				return nil, err
			}

package samlsp

import (
	"context"
	"encoding/xml"
	"errors"
	"net/url"

	"github.com/crewjam/httperr"

	"github.com/crewjam/saml"
)
// <EntityDescriptor>.
func ParseMetadata(data []byte) (*saml.EntityDescriptor, error) {
	entity := &saml.EntityDescriptor{}
	err := xml.Unmarshal(data, entity)

	// this comparison is ugly, but it is how the error is generated in encoding/xml
