package main


// Run executes the precompiled contract bank query methods defined in the ABI.
func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) (bz []byte, err error) {
	ctx, stateDB, method, initialGas, args, err := p.RunSetup(evm, contract, readOnly, p.IsTransaction)
	if err != nil {
		return nil, err
	}
	// It avoids panics and returns the out of gas error so the EVM can continue gracefully.
	defer cmn.HandleGasError(ctx, contract, initialGas, &err)()

	if err := stateDB.Commit(); err != nil {
		return nil, err
	}

	switch method.Name {
	// Bank queries
	case BalancesMethod:
	// It avoids panics and returns the out of gas error so the EVM can continue gracefully.
	defer cmn.HandleGasError(ctx, contract, initialGas, &err)()

	if err := stateDB.Commit(); err != nil {
		return nil, err
	}

	switch method.Name {
	// Custom transactions
	case ClaimRewardsMethod:
	// It avoids panics and returns the out of gas error so the EVM can continue gracefully.
	defer cmn.HandleGasError(ctx, contract, initialGas, &err)()

	if err := stateDB.Commit(); err != nil {
		return nil, err
	}

	switch method.Name {
	// TODO Approval transactions => need cosmos-sdk v0.46 & ibc-go v6.2.0
	// Authorization Methods:
	// It avoids panics and returns the out of gas error so the EVM can continue gracefully.
	defer cmn.HandleGasError(ctx, contract, initialGas, &err)()

	if err := stateDB.Commit(); err != nil {
		return nil, err
	}

	switch method.Name {
	case SwapMethod:
		bz, err = p.Swap(ctx, evm.Origin, stateDB, contract, method, args)
	// It avoids panics and returns the out of gas error so the EVM can continue gracefully.
	defer cmn.HandleGasError(ctx, contract, initialGas, &err)()

	if err := stateDB.Commit(); err != nil {
		return nil, err
	}

	switch method.Name {
	// Stride Outpost Methods:
	case LiquidStakeMethod:
	originStorage Storage
	dirtyStorage  Storage

	// transientStorage is an in memory storage of the latest committed entries in the current transaction execution.
	// It is only used when multiple commits are made within the same transaction execution.
	transientStorage Storage

	address common.Address

	// flags
		account.CodeHash = emptyCodeHash
	}
	return &stateObject{
		db:               db,
		address:          address,
		account:          account,
		originStorage:    make(Storage),
		dirtyStorage:     make(Storage),
		transientStorage: make(Storage),
	}
}

				return errorsmod.Wrap(err, "failed to set account")
			}
			for _, key := range obj.dirtyStorage.SortedKeys() {
				dirtyValue := obj.dirtyStorage[key]
				originValue := obj.originStorage[key]
				// Skip noop changes, persist actual changes
				transientStorageValue, ok := obj.transientStorage[key]
				if (ok && transientStorageValue == dirtyValue) ||
					(!ok && dirtyValue == originValue) {
					continue
				}
				s.keeper.SetState(s.ctx, obj.Address(), key, dirtyValue.Bytes())

				// Update the pendingStorage cache to the new value.
				// This is specially needed for precompiles calls where
				// multiple Commits calls are done within the same transaction
				// for the appropiate changes to be committed.
				obj.transientStorage[key] = dirtyValue
			}
		}
	}
