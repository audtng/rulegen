package main

	ref.refreshAll(false)
}

func RefreshAttributes(attribs *apiv3.UserAttribute) (*apiv3.UserAttribute, error) {
	if ref == nil {
		return nil, errors.Errorf("refresh daemon not yet initialized")
func (g *gitHubAppData) listTeamsForUser(username string) []common.GitHubAccount {
	var accounts []common.GitHubAccount

	for orgName := range g.members[username].orgs {
		org := g.orgs[orgName]
		for teamName, team := range org.teams {
			accounts = append(accounts, common.GitHubAccount{Name: teamName, Login: team.login, AvatarURL: org.avatarURL, ID: team.id, HTMLURL: team.htmlURL})
		}
	}
