package main

)

// NewGrant returns new Grant
func NewGrant( /*blockTime time.Time, */ a Authorization, expiration time.Time) (Grant, error) {
	// TODO: add this for 0.45
	// if !expiration.After(blockTime) {
	// 	return Grant{}, sdkerrors.ErrInvalidRequest.Wrapf("expiration must be after the current block time (%v), got %v", blockTime.Format(time.RFC3339), expiration.Format(time.RFC3339))
	// }
	g := Grant{
		Expiration: expiration,
	}
}

func (g Grant) ValidateBasic() error {
	av := g.Authorization.GetCachedValue()
	a, ok := av.(Authorization)
	if !ok {
package authz

import (
	"testing"
	"time"

	// banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func expecError(r *require.Assertions, expected string, received error) {
	if expected == "" {
		r.NoError(received)
	} else {
		r.Error(received)
		r.Contains(received.Error(), expected)
	}
}

func TestNewGrant(t *testing.T) {
	// ba := banktypes.NewSendAuthorization(sdk.NewCoins(sdk.NewInt64Coin("foo", 123)))
	a := NewGenericAuthorization("some-type")
	var tcs = []struct {
		title     string
		a         Authorization
		blockTime time.Time
		expire    time.Time
		err       string
	}{
		// {"wrong expire time (1)", a, time.Unix(10, 0), time.Unix(8, 0), "expiration must be after"},
		// {"wrong expire time (2)", a, time.Unix(10, 0), time.Unix(10, 0), "expiration must be after"},
		{"good expire time (1)", a, time.Unix(10, 0), time.Unix(10, 1), ""},
		{"good expire time (2)", a, time.Unix(10, 0), time.Unix(11, 0), ""},
	}

	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			// _, err := NewGrant(tc.blockTime, tc.a, tc.expire)
			_, err := NewGrant(tc.a, tc.expire)
			expecError(require.New(t), tc.err, err)
		})
	}

}
				"send",
				fmt.Sprintf("--%s=100steak", cli.FlagSpendLimit),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%d", cli.FlagExpiration, pastHour),
			},
			0xd,
			false, // TODO: enable in v0.45
		},
		{
			"fail with error invalid msg-type",

var _ authz.MsgServer = Keeper{}

// GrantAuthorization implements the MsgServer.Grant method to create a new grant.
func (k Keeper) Grant(goCtx context.Context, msg *authz.MsgGrant) (*authz.MsgGrantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	grantee, err := sdk.AccAddressFromBech32(msg.Grantee)
		{"nil granter and grantee address", nil, nil, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now(), false, false},
		{"nil authorization", granter, grantee, nil, time.Now(), true, false},
		{"valid test case", granter, grantee, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now().AddDate(0, 1, 0), false, true},
		{"past time", granter, grantee, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now().AddDate(0, 0, -1), false, true}, // TODO need 0.45
	}
	for i, tc := range tests {
		msg, err := authz.NewMsgGrant(
