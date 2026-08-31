package main

	return ok
}

var worldNets = map[identity.NumericIdentity][]*net.IPNet{
	identity.ReservedIdentityWorld: {
		{IP: net.IPv4zero, Mask: net.CIDRMask(0, net.IPv4len*8)},
					// identity is removed.
					newEntry.AddDependent(newKeyCpy)
				}
			} else if (newKey.Identity == k.Identity ||
				identityIsSupersetOf(newKey.Identity, k.Identity, identities)) &&
				(newKey.PortProtoIsBroader(k) || newKey.PortProtoIsEqual(k)) {
				// If the new-entry is a superset (or equal) of the iterated-allow-entry and
				// the new-entry has a broader (or equal) port-protocol then we
				// should delete the iterated-allow-entry
				ms.deleteKeyWithChanges(k, nil, changes)
			}
			return true
		})
				return true
			}

			if (newKey.Identity == k.Identity ||
				identityIsSupersetOf(k.Identity, newKey.Identity, identities)) &&
				k.DestPort == 0 && k.Nexthdr == 0 &&
				!v.HasDependent(newKey) {
				// If this iterated-deny-entry is a supserset (or equal) of the new-entry and
				// the iterated-deny-entry is an L3-only policy then we
				// should not insert the new entry (as long as it is not one
				// of the special L4-only denies we created to cover the special
				// case of a superset-allow with a more specific port-protocol).
				//
				// NOTE: This condition could be broader to reject more deny entries,
				// but there *may* be performance tradeoffs.
				bailed = true
				return false
			} else if (newKey.Identity == k.Identity ||
				identityIsSupersetOf(newKey.Identity, k.Identity, identities)) &&
				newKey.DestPort == 0 && newKey.Nexthdr == 0 &&
				!newEntry.HasDependent(k) {
				// If this iterated-deny-entry is a subset (or equal) of the new-entry and
				// the new-entry is an L3-only policy then we
				// should delete the iterated-deny-entry (as long as it is not one
				// of the special L4-only denies we created to cover the special
				// case of a superset-allow with a more specific port-protocol).
				//
				// NOTE: This condition could be broader to reject more deny entries,
				// but there *may* be performance tradeoffs.
				ms.deleteKeyWithChanges(k, nil, changes)
			}
			return true
	} else {
		// NOTE: We do not delete redundant allow entries.
		bailed := false
		ms.ForEachDeny(func(k Key, v MapStateEntry) bool {
			// Protocols and traffic directions that don't match ensure that the policies
			// do not interact in anyway.
					// identity is removed.
					ms.addDependentOnEntry(k, v, denyKeyCpy, changes)
				}
			} else if (k.Identity == newKey.Identity ||
				identityIsSupersetOf(k.Identity, newKey.Identity, identities)) &&
				(k.PortProtoIsBroader(newKey) || k.PortProtoIsEqual(newKey)) &&
				!v.HasDependent(newKey) {
				// If the iterated-deny-entry is a superset (or equal) of the new-entry and has a
				// broader (or equal) port-protocol than the new-entry then the new
				// entry should not be inserted.
				bailed = true
				return false
			}

			return true
		})

		if !bailed {
			ms.authPreferredInsert(newKey, newEntry, features, changes)
		}
