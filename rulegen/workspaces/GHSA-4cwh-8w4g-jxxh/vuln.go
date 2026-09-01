package main

		return
	}
	_, _ = uc.actionService.ActionRecordAdd(ctx, schema.ActionRecordTypeFindPass, ctx.ClientIP())
	code, err := uc.userService.RetrievePassWord(ctx, req)
	handler.HandleResponse(ctx, err, code)
}

// UseRePassWord godoc
