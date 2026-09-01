package main

	Groups   []string `json:"groups"`
}

// session holds the user session information during the life of a request.
type session struct {
	Details
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fluxcd/pkg/runtime/cel"
)

const (

// Validate validates the AnonymousAuthenticationSpec configuration.
func (a *AnonymousAuthenticationSpec) Validate() error {
	if a.Username == "" && len(a.Groups) == 0 {
		return fmt.Errorf("at least one of 'username' or 'groups' must be set for Anonymous authentication")
	}
	return nil
}

			imp.Groups = []string{}
		}

		return &user.Details{
			Profile:       profile,
			Impersonation: imp,
