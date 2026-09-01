package main

		return
	}
	_, _ = uc.actionService.ActionRecordAdd(ctx, schema.ActionRecordTypeFindPass, ctx.ClientIP())
	_, err := uc.userService.RetrievePassWord(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UseRePassWord godoc
