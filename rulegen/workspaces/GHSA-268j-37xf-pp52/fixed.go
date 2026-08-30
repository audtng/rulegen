package main

						Delete(deleteMilestone)
				}, reqRepoWriter())

				m.Patch("/issue-tracker", reqRepoAdmin(), bind(editIssueTrackerRequest{}), issueTracker)
				m.Patch("/wiki", reqRepoAdmin(), bind(editWikiRequest{}), wiki)
				m.Post("/mirror-sync", reqRepoAdmin(), mirrorSync)
				m.Get("/editorconfig/:filename", context.RepoRef(), getEditorconfig)
			}, repoAssignment())
		}, reqToken())
