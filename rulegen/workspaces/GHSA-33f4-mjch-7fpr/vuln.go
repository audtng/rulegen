package main

package reviewbot

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v59/github"
	"github.com/rs/zerolog/log"
)

// re-requested reviews should remove last review
// - fire event

type PullRequestInfo struct {
	owner          string
	repo           string
	user           string
	installationId int64
	headSHA        string
	number         int
}

func runPRCheck(config Config, pr PullRequestInfo) error {
	tr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID, pr.installationId, config.GitHub.PrivateKeyPath)
	if err != nil {
		log.Error().Interface("pr", pr).Err(err).Msg("Could not read key")
		return err
	}

	client := github.NewClient(&http.Client{Transport: tr})
	ctx := context.Background()

	// TODO: get repo-level overwrites, if available
	minReviewsRequired := config.MinReviewsRequired

	// List of approvers to verify
	var approvalCandidates = map[string]bool{
		// Add PR Creator as someone to check
		pr.user: true,
	}

	optListReviews := &github.ListOptions{PerPage: 100}

	// Check reviews
	for {
		reviews, resp, err := client.PullRequests.ListReviews(ctx, pr.owner, pr.repo, pr.number, optListReviews)
		if err != nil {
			log.Error().Interface("pr", pr).Err(err).Msg("Could not list reviews")
			return err
		}

		for _, review := range reviews {
			login := review.GetUser().GetLogin()
			association := review.GetAuthorAssociation()
			state := review.GetState()

			// Ignore accounts without association with the repo and comments
			if association == "NONE" || state == "COMMENTED" {
				continue
			}

			log.Debug().Interface("pr", pr).Str("login", login).Str("association", association).Str("state", state).Msg("Found a review candidate")

			if state == "APPROVED" {
				approvalCandidates[login] = true
			} else {
				delete(approvalCandidates, login)
			}
		}

		if resp.NextPage == 0 {
			break
		}

		optListReviews.Page = resp.NextPage
	}

	// Points for approval
	var points uint64 = 0

	for login := range approvalCandidates {
		permissionLevel, _, err := client.Repositories.GetPermissionLevel(ctx, pr.owner, pr.repo, login)
		if err != nil {
			return err
		}

		permission := permissionLevel.GetPermission()
		isAuthorized := permission == "admin" || permission == "write"

		log.Debug().Interface("pr", pr).Str("login", login).Str("permission", permission).Msg("Approver Authorization")

		if isAuthorized {
			points++

			if points == minReviewsRequired {
				// no need to waste resources - we have enough authorized approvers
				break
			}
		}
	}

	log.Info().Interface("pr", pr).Uint64("points", points).Msg("Check's State")

	statusComplete := "completed"
	titlePrefix := "⭐️ Allstar Pull Request Review Bot - "
	text := fmt.Sprintf("PR has %d authorized approvals, %d required", points, minReviewsRequired)
	timestamp := github.Timestamp{
		Time: time.Now(),
	}

	check := github.CreateCheckRunOptions{
		Name:        "Allstar Review Bot",
		Status:      &statusComplete,
		CompletedAt: &timestamp,
		Output: &github.CheckRunOutput{
			Text: &text,
		},
		HeadSHA: pr.headSHA,
	}

	if points >= minReviewsRequired {
		conclusion := "success"
		title := titlePrefix + conclusion
		summary := "Pull request has enough authorized approvals"

		check.Conclusion = &conclusion
		check.Output.Title = &title
		check.Output.Summary = &summary
	} else {
		conclusion := "failure"
		title := titlePrefix + conclusion

		delta := minReviewsRequired - points
		deltaMessage := fmt.Sprintf("need %d more approval(s)", delta)

		summary := "Pull request does not have enough authorized approvals - " + deltaMessage

		check.Conclusion = &conclusion
		check.Output.Title = &title
		check.Output.Summary = &summary
	}

	checkRun, _, err := client.Checks.CreateCheckRun(ctx, pr.owner, pr.repo, check)
	if err != nil {
		return err
	}

	log.Info().Interface("pr", pr).Interface("Check Run", checkRun).Msg("Created Check Run")

	return nil
}
// Copyright 2022 Allstar Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/ossf/allstar/pkg/reviewbot"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultAppID = 169668
const defaultMinReviewsRequired = 2
const defaultPort = 8080
const defaultSecretToken = "FooBar"

func main() {
	setupLog()

	config := reviewbot.Config{}

	if err := determineConfig(&config); err != nil {
		log.Fatal().Err(err).Msg("Error determining configuration")
	}

	if err := reviewbot.HandleWebhooks(&config); err != nil {
		log.Fatal().Err(err).Msg("Error listening to webhooks")
	}
}

func determineConfigFromEnv(config *reviewbot.Config) error {
	if envPort, ok := os.LookupEnv("PORT"); ok {
		port, err := strconv.ParseUint(envPort, 10, 16)

		if err != nil {
			return err
		}

		config.Port = port
	}

	if envAppId, ok := os.LookupEnv("APP_ID"); ok {
		appId, err := strconv.ParseInt(envAppId, 10, 64)

		if err != nil {
			return err
		}

		config.GitHub.AppId = appId
	}

	if envPrivateKeyPath, ok := os.LookupEnv("PRIVATE_KEY_PATH"); ok {
		config.GitHub.PrivateKeyPath = envPrivateKeyPath
	}

	if envSecretToken, ok := os.LookupEnv("SECRET_TOKEN"); ok {
		config.GitHub.SecretToken = envSecretToken
	}

	return nil
}

func determineConfigFromFlags(config *reviewbot.Config) error {
	flagAppID := flag.Int64("app-id", defaultAppID, "A GitHub App Id")
	flagPrivateKeyPath := flag.String("private-key-path", "", "A path to a GitHub Private Key")
	flagSecretToken := flag.String("secret-token", defaultSecretToken, "GitHub Private Key")
	flagMinReviewsRequired := flag.Uint64("min-reviews-required", defaultMinReviewsRequired, "The global minimum number of reviews required")
	flagPort := flag.Uint64("port", defaultPort, "A port to listen on")

	flag.Parse()

	if *flagAppID != defaultAppID {
		config.GitHub.AppId = *flagAppID
	}

	if *flagPrivateKeyPath != "" {
		config.GitHub.PrivateKeyPath = *flagPrivateKeyPath
	}

	if *flagSecretToken != defaultSecretToken {
		config.GitHub.PrivateKeyPath = *flagSecretToken
	}

	if *flagMinReviewsRequired != defaultMinReviewsRequired {
		config.MinReviewsRequired = *flagMinReviewsRequired
	}

	if *flagPort != defaultPort {
		config.Port = *flagPort
	}

	return nil
}

func determineConfig(config *reviewbot.Config) error {
	// Set defaults
	config.GitHub.AppId = defaultAppID
	config.GitHub.SecretToken = defaultSecretToken
	config.MinReviewsRequired = defaultMinReviewsRequired
	config.Port = defaultPort

	// Determine from environment variables
	if err := determineConfigFromEnv(config); err != nil {
		return err
	}

	// Determine from flags
	if err := determineConfigFromFlags(config); err != nil {
		return err
	}

	return nil
}

func setupLog() {
	// Match expected values in GCP
	zerolog.LevelFieldName = "severity"
	zerolog.LevelTraceValue = "DEFAULT"
	zerolog.LevelDebugValue = "DEBUG"
	zerolog.LevelInfoValue = "INFO"
	zerolog.LevelWarnValue = "WARNING"
	zerolog.LevelErrorValue = "ERROR"
	zerolog.LevelFatalValue = "CRITICAL"
	zerolog.LevelPanicValue = "CRITICAL"
}
package reviewbot

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v59/github"
	"github.com/rs/zerolog/log"
)

const secretToken = "FooBar"
const appID = 169668

type Config struct {
	// Configuration for GitHub
	// TODO: future: option to get the below values from Secret Manager. See: `setKeySecret` in pkg/config/operator/operator.go
	GitHub struct {
		// The GitHub App's id.
		// See: https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#authenticating-as-a-github-app
		AppId int64

		// Path to private key
		PrivateKeyPath string

		// See https://docs.github.com/en/developers/webhooks-and-events/webhooks/securing-your-webhooks
		SecretToken string
	}

	// The global minimum reviews required for approval
	MinReviewsRequired uint64

	// Port to listen on
	Port uint64
}

type WebookHandler struct {
	config Config
}

// Handle GitHub Webhooks for Review Bot.
//
// Example:
//
//	config := Config{...}
//	reviewbot.HandleWebhooks(&config)
func HandleWebhooks(config *Config) error {
	w := WebookHandler{*config}

	http.HandleFunc("/", w.HandleRoot)

	address := fmt.Sprintf(":%d", config.Port)

	return http.ListenAndServe(address, nil)
}

// Handle the root path
func (h *WebookHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	// Validate payload
	payload, err := github.ValidatePayload(r, []byte(secretToken))
	if err != nil {
		log.Error().Interface("payload", payload).Err(err).Msg("Got an invalid payload")
		w.WriteHeader(400)
		if _, err := fmt.Fprintln(w, "Got an invalid payload"); err != nil {
			log.Error().Err(err).Msg("Failed to write http response")
		}
		return
	}

	// Parse the webhook to a GitHub event type
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Error().Interface("event", event).Err(err).Msg("Failed to parse the webhook payload")
		w.WriteHeader(400)
		if _, err := fmt.Fprintln(w, "Failed to parse the webhook payload"); err != nil {
			log.Error().Err(err).Msg("Failed to write http response")
		}
		return
	}

	var pr PullRequestInfo

	// Extract relevant PR information from event, if is a PR-related event
	switch event := event.(type) {
	case *github.PullRequestEvent:
		pr = PullRequestInfo{
			owner:          event.GetRepo().GetOwner().GetLogin(),
			repo:           event.GetRepo().GetName(),
			user:           event.GetPullRequest().GetUser().GetLogin(),
			installationId: event.GetInstallation().GetID(),
			headSHA:        event.PullRequest.GetHead().GetSHA(),
			number:         event.GetPullRequest().GetNumber(),
		}
	default:
		log.Warn().Interface("event", event).Msg("Unknown event")
		w.WriteHeader(400)
		if _, err := fmt.Fprintln(w, "Unknown GitHub Event"); err != nil {
			log.Error().Err(err).Msg("Failed to write http response")
		}
		return
	}

	log.Info().Interface("pr", pr).Msg("Handling Pull Request Review Event")

	// Run PR Check
	err = runPRCheck(h.config, pr)
	if err != nil {
		log.Error().Interface("pr", pr).Err(err).Msg("Error handling webhook")
		w.WriteHeader(500)
		if _, err := fmt.Fprintln(w, "Error handling webhook"); err != nil {
			log.Error().Err(err).Msg("Failed to write http response")
		}
		return
	}
}
