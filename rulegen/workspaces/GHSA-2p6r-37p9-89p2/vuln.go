package main

)

// NewGrant returns new Grant
func NewGrant(a Authorization, expiration time.Time) (Grant, error) {
	g := Grant{
		Expiration: expiration,
	}
}

func (g Grant) ValidateBasic() error {
	if g.Expiration.Unix() < time.Now().Unix() {
		return sdkerrors.Wrap(ErrInvalidExpirationTime, "Time can't be in the past")
	}

	av := g.Authorization.GetCachedValue()
	a, ok := av.(Authorization)
	if !ok {
				"send",
				fmt.Sprintf("--%s=100steak", cli.FlagSpendLimit),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
				fmt.Sprintf("--%s=%d", cli.FlagExpiration, pastHour),
			},
			0,
			true,
		},
		{
			"fail with error invalid msg-type",

var _ authz.MsgServer = Keeper{}

// GrantAuthorization implements the MsgServer.Grant method.
func (k Keeper) Grant(goCtx context.Context, msg *authz.MsgGrant) (*authz.MsgGrantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	grantee, err := sdk.AccAddressFromBech32(msg.Grantee)
		{"nil granter and grantee address", nil, nil, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now(), false, false},
		{"nil authorization", granter, grantee, nil, time.Now(), true, false},
		{"valid test case", granter, grantee, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now().AddDate(0, 1, 0), false, true},
		{"past time", granter, grantee, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now().AddDate(0, 0, -1), false, false},
	}
	for i, tc := range tests {
		msg, err := authz.NewMsgGrant(
