package main


	op := api.OrganizationPermissions{}

	if !organization.HasOrgOrUserVisible(ctx, o, ctx.ContextUser) {
		ctx.APIErrorNotFound("HasOrgOrUserVisible", nil)
		return
	}
