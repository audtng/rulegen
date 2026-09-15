package main

// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/auth"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
				ctx.HTML(200, "user/auth/activate")
				return
			}
			if ctx.IsSigned && auth.IsAPIPath(ctx.Req.URL.Path) && ctx.IsBasicAuth {
				twofa, err := models.GetTwoFactorByUID(ctx.User.ID)
				if err != nil {
					if models.IsErrTwoFactorNotEnrolled(err) {
						return // No 2FA enrollment for this user
					}
					ctx.Error(500)
					return
				}
				otpHeader := ctx.Req.Header.Get("X-Gitea-OTP")
				ok, err := twofa.ValidateTOTP(otpHeader)
				if err != nil {
					ctx.Error(500)
					return
				}
				if !ok {
					ctx.JSON(403, map[string]string{
						"message": "Only signed in user is allowed to call APIs.",
					})
					return
				}
			}
		}

		// Redirect to log in page if auto-signin info is provided and has not signed in.
	}
}

// CheckForOTP validateds OTP
func (ctx *APIContext) CheckForOTP() {
	otpHeader := ctx.Req.Header.Get("X-Gitea-OTP")
	twofa, err := models.GetTwoFactorByUID(ctx.Context.User.ID)
	if err != nil {
		if models.IsErrTwoFactorNotEnrolled(err) {
			return // No 2FA enrollment for this user
		}
		ctx.Context.Error(500)
		return
	}
	ok, err := twofa.ValidateTOTP(otpHeader)
	if err != nil {
		ctx.Context.Error(500)
		return
	}
	if !ok {
		ctx.Context.Error(401)
		return
	}
}

// APIContexter returns apicontext as macaron middleware
func APIContexter() macaron.Handler {
	return func(c *Context) {
		if true == ctx.Data["IsApiToken"] {
			return
		}
		if ctx.Context.IsBasicAuth {
			ctx.CheckForOTP()
			return
		}
		if ctx.IsSigned {
			ctx.RequireCSRF()
			return
}

func reqBasicAuth() macaron.Handler {
	return func(ctx *context.APIContext) {
		if !ctx.Context.IsBasicAuth {
			ctx.Context.Error(401)
			return
		}
		ctx.CheckForOTP()
	}
}

