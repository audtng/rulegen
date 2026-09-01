package main

	Groups   []string `json:"groups"`
}

// SanitizeAndValidate sanitizes and validates the user impersonation details.
func (imp *Impersonation) SanitizeAndValidate() error {
	imp.Username = strings.TrimSpace(imp.Username)
	for i, g := range imp.Groups {
		imp.Groups[i] = strings.TrimSpace(g)
	}
	if imp.Groups == nil {
		imp.Groups = []string{}
	}
	slices.Sort(imp.Groups)
	if imp.Username == "" && len(imp.Groups) == 0 {
		return fmt.Errorf("at least one of 'username' or 'groups' must be set for user impersonation")
	}
	for i, g := range imp.Groups {
		if g == "" {
			return fmt.Errorf("group[%d] is an empty string", i)
		}
	}
	return nil
}

// session holds the user session information during the life of a request.
type session struct {
	Details
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fluxcd/pkg/runtime/cel"

	"github.com/controlplaneio-fluxcd/flux-operator/internal/web/user"
)

const (

// Validate validates the AnonymousAuthenticationSpec configuration.
func (a *AnonymousAuthenticationSpec) Validate() error {
	imp := (user.Impersonation)(*a)
	if err := imp.SanitizeAndValidate(); err != nil {
		return fmt.Errorf("invalid anonymous authentication impersonation: %w", err)
	}
	*a = (AnonymousAuthenticationSpec)(imp)
	return nil
}

			imp.Groups = []string{}
		}

		// Sanitize and validate the extracted impersonation.
		if err := imp.SanitizeAndValidate(); err != nil {
			return nil, fmt.Errorf("impersonation validation failed: %w", err)
		}

		return &user.Details{
			Profile:       profile,
			Impersonation: imp,
