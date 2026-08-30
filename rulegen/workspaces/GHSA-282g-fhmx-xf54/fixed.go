package main

				}
			},
		},
		{
			name: "verified email",
			testCase: func(runId, userId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Email: &user.SetHumanEmail{
										Email: integration.Email(),
										Verification: &user.SetHumanEmail_IsVerified{
											IsVerified: true,
										},
									},
								},
							},
						},
					},
					wantErr: false,
				}
			},
		},
		{
			name: "verified email (self-management), error",
			testCase: func(runId, userId string) testCase {
				Instance.SetUserPassword(LoginCTX, userId, integration.UserPassword, false)
				_, token, _, _ := Instance.CreatePasswordSession(t, LoginCTX, userId, integration.UserPassword)
				return testCase{
					args: args{
						integration.WithAuthorizationToken(CTX, token),
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Email: &user.SetHumanEmail{
										Email: integration.Email(),
										Verification: &user.SetHumanEmail_IsVerified{
											IsVerified: true,
										},
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "password not complexity conform",
			testCase: func(runId, userId string) testCase {
	}

	if human.Changed() {
		// Changing metadata or setting email, resp. phone to verified is only allowed with user write permissions, but not for self-management.
		requireWritePermission := metadataChanged || (human.Email != nil && human.Email.Verified) || (human.Phone != nil && human.Phone.Verified)
		if err := c.checkPermissionUpdateUser(ctx, existingHuman.ResourceOwner, existingHuman.AggregateID, !requireWritePermission); err != nil {
			return err
		}
	}
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
				},
			},
		},
		{
			name: "change human email verified (self-management), not allowed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
				tarpit:          expectTarpit(0),
			},
			args: args{
				ctx:   authz.NewMockContext("instance1", "org1", "user1"),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address:  "changed@example.com",
						Verified: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change human email verified, ok",
			fields: fields{
				},
			},
		},
		{
			name: "change human phone verified (self-management), not allowed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
				tarpit:          expectTarpit(0),
			},
			args: args{
				ctx:   authz.NewMockContext("instance1", "org1", "user1"),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Number:   "+41791234567",
						Verified: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change human phone verified, ok",
			fields: fields{
