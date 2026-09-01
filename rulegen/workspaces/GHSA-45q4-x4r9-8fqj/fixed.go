package main

// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package notifications

import "strings"

// markdownSpecialChars is the full CommonMark §2.4 backslash-escapable set.
// '<' is included to neutralize autolinks (`<https://evil.com>`), which
// bluemonday's UGC policy would otherwise render as clickable <a> tags.
// '\' is handled separately first so inserted backslashes are not re-escaped.
const markdownSpecialChars = "`*_{}[]()<>#+-.!|~"

// EscapeMarkdown escapes every CommonMark-special character in s. Fixes
// GHSA-45q4-x4r9-8fqj (Markdown injection in notification emails).
func EscapeMarkdown(s string) string {
	// Backslash first so inserted backslashes are not double-escaped.
	s = strings.ReplaceAll(s, `\`, `\\`)
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r < 128 && strings.ContainsRune(markdownSpecialChars, r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}
		To(n.User.Email).
		Subject(i18n.T(lang, "notifications.task.reminder.subject", n.Task.Title, n.Project.Title)).
		Greeting(i18n.T(lang, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(lang, "notifications.task.reminder.message", notifications.EscapeMarkdown(n.Task.Title), notifications.EscapeMarkdown(n.Project.Title))).
		Action(i18n.T(lang, "notifications.common.actions.open_task"), config.ServicePublicURL.GetString()+"tasks/"+strconv.FormatInt(n.Task.ID, 10)).
		Line(i18n.T(lang, "notifications.common.have_nice_day"))
}
			From(n.Doer.GetNameAndFromEmail()).
			Subject(i18n.T(lang, "notifications.task.assigned.subject_to_assignee", n.Task.Title, n.Task.GetFullIdentifier())).
			Greeting(i18n.T(lang, "notifications.greeting", n.Target.GetName())).
			Line(i18n.T(lang, "notifications.task.assigned.message_to_assignee", notifications.EscapeMarkdown(n.Doer.GetName()), notifications.EscapeMarkdown(n.Task.Title))).
			Action(i18n.T(lang, "notifications.common.actions.open_task"), n.Task.GetFrontendURL()).
			IncludeLinkToSettings(lang)
	}
			From(n.Doer.GetNameAndFromEmail()).
			Subject(i18n.T(lang, "notifications.task.assigned.subject_to_others_self", n.Task.Title, n.Task.GetFullIdentifier(), n.Doer.GetName())).
			Greeting(i18n.T(lang, "notifications.greeting", n.Target.GetName())).
			Line(i18n.T(lang, "notifications.task.assigned.message_to_others_self", notifications.EscapeMarkdown(n.Doer.GetName()))).
			Action(i18n.T(lang, "notifications.common.actions.open_task"), n.Task.GetFrontendURL()).
			IncludeLinkToSettings(lang)
	}
		From(n.Doer.GetNameAndFromEmail()).
		Subject(i18n.T(lang, "notifications.task.assigned.subject_to_others", n.Task.Title, n.Task.GetFullIdentifier(), n.Assignee.GetName())).
		Greeting(i18n.T(lang, "notifications.greeting", n.Target.GetName())).
		Line(i18n.T(lang, "notifications.task.assigned.message_to_others", notifications.EscapeMarkdown(n.Doer.GetName()), notifications.EscapeMarkdown(n.Assignee.GetName()))).
		Action(i18n.T(lang, "notifications.common.actions.open_task"), n.Task.GetFrontendURL()).
		IncludeLinkToSettings(lang)
}
func (n *TaskDeletedNotification) ToMail(lang string) *notifications.Mail {
	return notifications.NewMail().
		Subject(i18n.T(lang, "notifications.task.deleted.subject", n.Task.Title, n.Task.GetFullIdentifier())).
		Line(i18n.T(lang, "notifications.task.deleted.message", notifications.EscapeMarkdown(n.Doer.GetName()), notifications.EscapeMarkdown(n.Task.Title), notifications.EscapeMarkdown(n.Task.GetFullIdentifier())))
}

// ToDB returns the TaskDeletedNotification notification in a format which can be saved in the db
func (n *ProjectCreatedNotification) ToMail(lang string) *notifications.Mail {
	return notifications.NewMail().
		Subject(i18n.T(lang, "notifications.project.created", n.Doer.GetName(), n.Project.Title)).
		Line(i18n.T(lang, "notifications.project.created", notifications.EscapeMarkdown(n.Doer.GetName()), notifications.EscapeMarkdown(n.Project.Title))).
		Action(i18n.T(lang, "notifications.common.actions.open_project"), config.ServicePublicURL.GetString()+"projects/")
}

		Subject(i18n.T(lang, "notifications.team.member_added.subject", n.Doer.GetName(), n.Team.Name)).
		From(n.Doer.GetNameAndFromEmail()).
		Greeting(i18n.T(lang, "notifications.greeting", n.Member.GetName())).
		Line(i18n.T(lang, "notifications.team.member_added.message", notifications.EscapeMarkdown(n.Doer.GetName()), notifications.EscapeMarkdown(n.Team.Name))).
		Action(i18n.T(lang, "notifications.common.actions.open_team"), config.ServicePublicURL.GetString()+"teams/"+strconv.FormatInt(n.Team.ID, 10)+"/edit")
}

		IncludeLinkToSettings(lang).
		Subject(i18n.T(lang, "notifications.task.overdue.subject", n.Task.Title, n.Project.Title)).
		Greeting(i18n.T(lang, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(lang, "notifications.task.overdue.message", notifications.EscapeMarkdown(n.Task.Title), notifications.EscapeMarkdown(n.Project.Title), getOverdueSinceString(until, n.User.Language))).
		Action(i18n.T(lang, "notifications.common.actions.open_task"), config.ServicePublicURL.GetString()+"tasks/"+strconv.FormatInt(n.Task.ID, 10)).
		Line(i18n.T(lang, "notifications.common.have_nice_day"))
}
	overdueLine := ""
	for _, task := range sortedTasks {
		until := time.Until(task.DueDate).Round(1*time.Hour) * -1
		overdueLine += `* [` + notifications.EscapeMarkdown(task.Title) + `](` + config.ServicePublicURL.GetString() + "tasks/" + strconv.FormatInt(task.ID, 10) + `) (` + notifications.EscapeMarkdown(n.Projects[task.ProjectID].Title) + `), ` + i18n.T("notifications.task.overdue.overdue", getOverdueSinceString(until, n.User.Language)) + "\n"
	}

	return notifications.NewMail().
	return notifications.NewMail().
		Subject(i18n.T(lang, "notifications.api_token.expiring.week.subject", n.Token.Title)).
		Greeting(i18n.T(lang, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(lang, "notifications.api_token.expiring.week.message", notifications.EscapeMarkdown(n.Token.Title), n.Token.ExpiresAt.Format("2006-01-02"))).
		Action(i18n.T(lang, "notifications.api_token.expiring.action"), config.ServicePublicURL.GetString()+"user/settings/api-tokens").
		Line(i18n.T(lang, "notifications.common.have_nice_day"))
}
	return notifications.NewMail().
		Subject(i18n.T(lang, "notifications.api_token.expiring.day.subject", n.Token.Title)).
		Greeting(i18n.T(lang, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(lang, "notifications.api_token.expiring.day.message", notifications.EscapeMarkdown(n.Token.Title), n.Token.ExpiresAt.Format("2006-01-02"))).
		Action(i18n.T(lang, "notifications.api_token.expiring.action"), config.ServicePublicURL.GetString()+"user/settings/api-tokens").
		Line(i18n.T(lang, "notifications.common.have_nice_day"))
}
