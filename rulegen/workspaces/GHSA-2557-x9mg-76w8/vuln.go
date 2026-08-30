package main

		totalCoins = totalCoins.Add(period.Amount...)
	}

	if acc := s.AccountKeeper.GetAccount(ctx, to); acc != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "account %s already exists", msg.ToAddress)
	}
			expErr:    true,
			expErrMsg: "already exists",
		},
		"create a valid delayed vesting account": {
			preRun: func() {
				s.bankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), fooCoin).Return(nil)
			expErr:    true,
			expErrMsg: "already exists",
		},
		"create a valid permanent locked account": {
			preRun: func() {
				s.bankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), fooCoin).Return(nil)
		{
			name: "create for existing account",
			preRun: func() {
				toAcc := s.accountKeeper.NewAccountWithAddress(s.ctx, to1Addr)
				s.accountKeeper.SetAccount(s.ctx, toAcc)
			},
			expErr:    true,
			expErrMsg: "already exists",
		},
		{
			name: "create a valid periodic vesting account",
			preRun: func() {
				s.bankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), periodCoin.Add(fooCoin)).Return(nil)
				s.bankKeeper.EXPECT().SendCoins(gomock.Any(), fromAddr, to2Addr, gomock.Any()).Return(nil)
			},
			input: vestingtypes.NewMsgCreatePeriodicVestingAccount(
