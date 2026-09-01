package main

		To(n.User.Email).
		Subject(i18n.T(lang, "notifications.task.reminder.subject", n.Task.Title, n.Project.Title)).
		Greeting(i18n.T(lang, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(lang, "notifications.task.reminder.message", n.Task.Title, n.Project.Title)).
		Action(i18n.T(lang, "notifications.common.actions.open_task"), config.ServicePublicURL.GetString()+"tasks/"+strconv.FormatInt(n.Task.ID, 10)).
		Line(i18n.T(lang, "notifications.common.have_nice_day"))
}
			From(n.Doer.GetNameAndFromEmail()).
			Subject(i18n.T(lang, "notifications.task.assigned.subject_to_assignee", n.Task.Title, n.Task.GetFullIdentifier())).
			Greeting(i18n.T(lang, "notifications.greeting", n.Target.GetName())).
			Line(i18n.T(lang, "notifications.task.assigned.message_to_assignee", n.Doer.GetName(), n.Task.Title)).
			Action(i18n.T(lang, "notifications.common.actions.open_task"), n.Task.GetFrontendURL()).
			IncludeLinkToSettings(lang)
	}
			From(n.Doer.GetNameAndFromEmail()).
			Subject(i18n.T(lang, "notifications.task.assigned.subject_to_others_self", n.Task.Title, n.Task.GetFullIdentifier(), n.Doer.GetName())).
			Greeting(i18n.T(lang, "notifications.greeting", n.Target.GetName())).
			Line(i18n.T(lang, "notifications.task.assigned.message_to_others_self", n.Doer.GetName())).
			Action(i18n.T(lang, "notifications.common.actions.open_task"), n.Task.GetFrontendURL()).
			IncludeLinkToSettings(lang)
	}
		From(n.Doer.GetNameAndFromEmail()).
		Subject(i18n.T(lang, "notifications.task.assigned.subject_to_others", n.Task.Title, n.Task.GetFullIdentifier(), n.Assignee.GetName())).
		Greeting(i18n.T(lang, "notifications.greeting", n.Target.GetName())).
		Line(i18n.T(lang, "notifications.task.assigned.message_to_others", n.Doer.GetName(), n.Assignee.GetName())).
		Action(i18n.T(lang, "notifications.common.actions.open_task"), n.Task.GetFrontendURL()).
		IncludeLinkToSettings(lang)
}
func (n *TaskDeletedNotification) ToMail(lang string) *notifications.Mail {
	return notifications.NewMail().
		Subject(i18n.T(lang, "notifications.task.deleted.subject", n.Task.Title, n.Task.GetFullIdentifier())).
		Line(i18n.T(lang, "notifications.task.deleted.message", n.Doer.GetName(), n.Task.Title, n.Task.GetFullIdentifier()))
}

// ToDB returns the TaskDeletedNotification notification in a format which can be saved in the db
func (n *ProjectCreatedNotification) ToMail(lang string) *notifications.Mail {
	return notifications.NewMail().
		Subject(i18n.T(lang, "notifications.project.created", n.Doer.GetName(), n.Project.Title)).
		Line(i18n.T(lang, "notifications.project.created", n.Doer.GetName(), n.Project.Title)).
		Action(i18n.T(lang, "notifications.common.actions.open_project"), config.ServicePublicURL.GetString()+"projects/")
}

		Subject(i18n.T(lang, "notifications.team.member_added.subject", n.Doer.GetName(), n.Team.Name)).
		From(n.Doer.GetNameAndFromEmail()).
		Greeting(i18n.T(lang, "notifications.greeting", n.Member.GetName())).
		Line(i18n.T(lang, "notifications.team.member_added.message", n.Doer.GetName(), n.Team.Name)).
		Action(i18n.T(lang, "notifications.common.actions.open_team"), config.ServicePublicURL.GetString()+"teams/"+strconv.FormatInt(n.Team.ID, 10)+"/edit")
}

		IncludeLinkToSettings(lang).
		Subject(i18n.T(lang, "notifications.task.overdue.subject", n.Task.Title, n.Project.Title)).
		Greeting(i18n.T(lang, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(lang, "notifications.task.overdue.message", n.Task.Title, n.Project.Title, getOverdueSinceString(until, n.User.Language))).
		Action(i18n.T(lang, "notifications.common.actions.open_task"), config.ServicePublicURL.GetString()+"tasks/"+strconv.FormatInt(n.Task.ID, 10)).
		Line(i18n.T(lang, "notifications.common.have_nice_day"))
}
	overdueLine := ""
	for _, task := range sortedTasks {
		until := time.Until(task.DueDate).Round(1*time.Hour) * -1
		overdueLine += `* [` + task.Title + `](` + config.ServicePublicURL.GetString() + "tasks/" + strconv.FormatInt(task.ID, 10) + `) (` + n.Projects[task.ProjectID].Title + `), ` + i18n.T("notifications.task.overdue.overdue", getOverdueSinceString(until, n.User.Language)) + "\n"
	}

	return notifications.NewMail().
	return notifications.NewMail().
		Subject(i18n.T(lang, "notifications.api_token.expiring.week.subject", n.Token.Title)).
		Greeting(i18n.T(lang, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(lang, "notifications.api_token.expiring.week.message", n.Token.Title, n.Token.ExpiresAt.Format("2006-01-02"))).
		Action(i18n.T(lang, "notifications.api_token.expiring.action"), config.ServicePublicURL.GetString()+"user/settings/api-tokens").
		Line(i18n.T(lang, "notifications.common.have_nice_day"))
}
	return notifications.NewMail().
		Subject(i18n.T(lang, "notifications.api_token.expiring.day.subject", n.Token.Title)).
		Greeting(i18n.T(lang, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(lang, "notifications.api_token.expiring.day.message", n.Token.Title, n.Token.ExpiresAt.Format("2006-01-02"))).
		Action(i18n.T(lang, "notifications.api_token.expiring.action"), config.ServicePublicURL.GetString()+"user/settings/api-tokens").
		Line(i18n.T(lang, "notifications.common.have_nice_day"))
}
