package main

	originStorage Storage
	dirtyStorage  Storage

	address common.Address

	// flags
		account.CodeHash = emptyCodeHash
	}
	return &stateObject{
		db:            db,
		address:       address,
		account:       account,
		originStorage: make(Storage),
		dirtyStorage:  make(Storage),
	}
}

				return errorsmod.Wrap(err, "failed to set account")
			}
			for _, key := range obj.dirtyStorage.SortedKeys() {
				value := obj.dirtyStorage[key]
				// Skip noop changes, persist actual changes
				if value == obj.originStorage[key] {
					continue
				}
				s.keeper.SetState(s.ctx, obj.Address(), key, value.Bytes())
			}
		}
	}

// Run executes the precompiled contract bank query methods defined in the ABI.
func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) (bz []byte, err error) {
	ctx, _, method, initialGas, args, err := p.RunSetup(evm, contract, readOnly, p.IsTransaction)
	if err != nil {
		return nil, err
	}
	// It avoids panics and returns the out of gas error so the EVM can continue gracefully.
	defer cmn.HandleGasError(ctx, contract, initialGas, &err)()

	switch method.Name {
	// Bank queries
	case BalancesMethod:
