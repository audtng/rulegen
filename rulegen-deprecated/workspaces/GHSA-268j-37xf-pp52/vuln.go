package main

						Delete(deleteMilestone)
				}, reqRepoWriter())

				m.Patch("/issue-tracker", reqRepoWriter(), bind(editIssueTrackerRequest{}), issueTracker)
				m.Patch("/wiki", reqRepoWriter(), bind(editWikiRequest{}), wiki)
				m.Post("/mirror-sync", reqRepoWriter(), mirrorSync)
				m.Get("/editorconfig/:filename", context.RepoRef(), getEditorconfig)
			}, repoAssignment())
		}, reqToken())
