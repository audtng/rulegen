package main

		totalCoins = totalCoins.Add(period.Amount...)
	}

	if s.BankKeeper.BlockedAddr(to) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	if acc := s.AccountKeeper.GetAccount(ctx, to); acc != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "account %s already exists", msg.ToAddress)
	}
			expErr:    true,
			expErrMsg: "already exists",
		},
		"create for blocked account": {
			preRun: func() {
				s.bankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), fooCoin).Return(nil)
				s.bankKeeper.EXPECT().BlockedAddr(to1Addr).Return(true)
			},
			input: vestingtypes.NewMsgCreateVestingAccount(
				fromAddr,
				to1Addr,
				sdk.Coins{fooCoin},
				time.Now().Unix(),
				true,
			),
			expErr:    true,
			expErrMsg: "not allowed to receive funds",
		},
		"create a valid delayed vesting account": {
			preRun: func() {
				s.bankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), fooCoin).Return(nil)
			expErr:    true,
			expErrMsg: "already exists",
		},
		"create for blocked account": {
			preRun: func() {
				toAcc := s.accountKeeper.NewAccountWithAddress(s.ctx, to1Addr)
				s.bankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), fooCoin).Return(nil)
				s.bankKeeper.EXPECT().BlockedAddr(to1Addr).Return(true)
				s.accountKeeper.SetAccount(s.ctx, toAcc)
			},
			input: vestingtypes.NewMsgCreatePermanentLockedAccount(
				fromAddr,
				to1Addr,
				sdk.Coins{fooCoin},
			),
			expErr:    true,
			expErrMsg: "not allowed to receive funds",
		},

		"create a valid permanent locked account": {
			preRun: func() {
				s.bankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), fooCoin).Return(nil)
		{
			name: "create for existing account",
			preRun: func() {
				s.bankKeeper.EXPECT().BlockedAddr(to1Addr).Return(false)
				toAcc := s.accountKeeper.NewAccountWithAddress(s.ctx, to1Addr)
				s.accountKeeper.SetAccount(s.ctx, toAcc)
			},
			expErr:    true,
			expErrMsg: "already exists",
		},
		{
			name: "create for blocked address",
			preRun: func() {
				s.bankKeeper.EXPECT().BlockedAddr(to2Addr).Return(true)
			},
			input: vestingtypes.NewMsgCreatePeriodicVestingAccount(
				fromAddr,
				to2Addr,
				time.Now().Unix(),
				[]vestingtypes.Period{
					{
						Length: 10,
						Amount: sdk.NewCoins(periodCoin),
					},
					{
						Length: 20,
						Amount: sdk.NewCoins(fooCoin),
					},
				},
			),
			expErr:    true,
			expErrMsg: "not allowed to receive funds",
		},
		{
			name: "create a valid periodic vesting account",
			preRun: func() {
				s.bankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), periodCoin.Add(fooCoin)).Return(nil)
				s.bankKeeper.EXPECT().BlockedAddr(to2Addr).Return(false)
				s.bankKeeper.EXPECT().SendCoins(gomock.Any(), fromAddr, to2Addr, gomock.Any()).Return(nil)
			},
			input: vestingtypes.NewMsgCreatePeriodicVestingAccount(
