package main


	op := api.OrganizationPermissions{}

	if !organization.HasOrgOrUserVisible(ctx, o, ctx.Doer) {
		ctx.APIErrorNotFound("HasOrgOrUserVisible", nil)
		return
	}
