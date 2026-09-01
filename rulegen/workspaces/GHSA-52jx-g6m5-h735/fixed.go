package main

	"github.com/fleetdm/fleet/v4/server/fleet"
	rtvalidator "github.com/mattermost/xml-roundtrip-validator"
	dsig "github.com/russellhaering/goxmldsig"
)

type Validator interface {
		// If entire doc is signed, success, we're done.
		return validated, nil
	}
	// Some IdPs (like Google) do not sign the root, and only sign the Assertion.
	if err == dsig.ErrMissingSignature {
		if err := v.validateAssertionSignature(elt); err != nil {
			return nil, err
		}
		return elt, nil
	}
	return nil, err
}

var (
	errMissingAssertion              = errors.New("missing Assertion element under namespace urn:oasis:names:tc:SAML:2.0:assertion")
	errMultipleAssertions            = errors.New("multiple Assertions elements found")
	errAssertionWithInvalidNamespace = errors.New("Assertion with invalid namespace found")
)

// validateAssertionSignature validates that one "Assertion" child element exists under
// the "urn:oasis:names:tc:SAML:2.0:assertion" namespace and that it's signed by the IdP.
// It returns:
//   - errMissingAssertion if there is no "Assertion" child element under the given tree.
//   - errMultipleAssertions if there's more than one "Assertion" element under the given tree.
//   - errAssertionWithInvalidNamespace if an "Assertion" element has a namespace that's not
//     "urn:oasis:names:tc:SAML:2.0:assertion"
//   - an error if the signature of the one "Assertion" element is invalid.
func (v *validator) validateAssertionSignature(elt *etree.Element) error {
	var assertion *etree.Element
	for _, child := range elt.ChildElements() {
		if child.Tag == "Assertion" {
			if child.NamespaceURI() != "urn:oasis:names:tc:SAML:2.0:assertion" {
				return errAssertionWithInvalidNamespace
			}
			if assertion != nil {
				return errMultipleAssertions
			}
			assertion = child
		}
	}
	if assertion == nil {
		return errMissingAssertion
	}
	if _, err := v.context.Validate(assertion); err != nil {
		return fmt.Errorf("failed to validate assertion signature: %w", err)
	}
	return nil
}

const (

func ValidateAudiences(metadata Metadata, auth fleet.Auth, audiences ...string) error {
	validator, err := NewValidator(metadata, WithExpectedAudience(audiences...))
	if err != nil {
		return fmt.Errorf("create validator from metadata: %w", err)
	}
