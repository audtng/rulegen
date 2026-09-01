package main

	"github.com/fleetdm/fleet/v4/server/fleet"
	rtvalidator "github.com/mattermost/xml-roundtrip-validator"
	dsig "github.com/russellhaering/goxmldsig"
	"github.com/russellhaering/goxmldsig/etreeutils"
)

type Validator interface {
		// If entire doc is signed, success, we're done.
		return validated, nil
	}
	if err == dsig.ErrMissingSignature {
		// If entire document is not signed find signed assertions, remove assertions
		// that are not signed.
		err = v.validateAssertionSignature(elt)
		if err != nil {
			return nil, err
		}
		return elt, nil
	}

	return nil, err
}

func (v *validator) validateAssertionSignature(elt *etree.Element) error {
	validateAssertion := func(ctx etreeutils.NSContext, unverified *etree.Element) error {
		if unverified.Parent() != elt {
			return fmt.Errorf("assertion with unexpected parent: %s", unverified.Parent().Tag)
		}
		// Remove assertions that are not signed.
		detached, err := etreeutils.NSDetatch(ctx, unverified)
		if err != nil {
			return err
		}
		signed, err := v.context.Validate(detached)
		if err != nil {
			return err
		}
		elt.RemoveChild(unverified)
		elt.AddChild(signed)
		return nil
	}
	return etreeutils.NSFindIterate(elt, "urn:oasis:names:tc:SAML:2.0:assertion", "Assertion", validateAssertion)
}

const (

func ValidateAudiences(metadata Metadata, auth fleet.Auth, audiences ...string) error {
	validator, err := NewValidator(metadata, WithExpectedAudience(audiences...))

	if err != nil {
		return fmt.Errorf("create validator from metadata: %w", err)
	}
