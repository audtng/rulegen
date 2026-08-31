package main

	return ok
}

// HasSameOwners returns true if both MapStateEntries
// have the same owners as one another.
// MapStateEntries are stored by value, so we do not check for nil pointers here.
func (e *MapStateEntry) HasSameOwners(bEntry *MapStateEntry) bool {
	if len(e.owners) != len(bEntry.owners) {
		return false
	}
	for owner := range e.owners {
		if _, ok := bEntry.owners[owner]; !ok {
			return false
		}
	}
	return true
}

var worldNets = map[identity.NumericIdentity][]*net.IPNet{
	identity.ReservedIdentityWorld: {
		{IP: net.IPv4zero, Mask: net.CIDRMask(0, net.IPv4len*8)},
					// identity is removed.
					newEntry.AddDependent(newKeyCpy)
				}
			} else if newKey.PortProtoIsBroader(k) || newKey.PortProtoIsEqual(k) {
				// If newKey has a broader (or equal) port-protocol then we should
				// either delete the iterated-allow-entry (if the identity is the
				// same or the newKey is L3 wildcard), or change it to a deny entry
				// if the newKey's identity is a superset of the iterated identity
				// (e.g., newKey has a wider CIDR (say 10/8 covering the iterated
				// identity of more specific CIDR (say 10.1.1.1). Note that the
				// security identities assigned to these CIDRs have no numerical
				// relation to each other (e.g, they could be any numbers X and Y)
				// and the datapath does an exact match on them.
				if newKey.Identity == 0 || newKey.Identity == k.Identity {
					ms.deleteKeyWithChanges(k, nil, changes)
				} else if identityIsSupersetOf(newKey.Identity, k.Identity, identities) {
					// When newKey.Identity is not ANY and is different from the
					// subset key, but still a superset (e.g., in CIDR sense) we
					// must keep the subset key and make it a deny instead.
					l3l4DenyEntry := NewMapStateEntry(newKey, newEntry.DerivedFromRules, 0, "", 0, true, DefaultAuthType, AuthTypeDisabled)
					ms.addKeyWithChanges(k, l3l4DenyEntry, changes)
					newEntry.AddDependent(k)
				}
			} else if identityIsSupersetOf(newKey.Identity, k.Identity, identities) {
				// k.PortProtoIsBroader(newKey) // due to if statements above

				// Deny takes precedence for the port/proto of the newKey
				// for each allow with broader port/proto and narrower ID.

				// If newKey is a superset of the iterated allow key and newKey has
				// a less specific port-protocol than the iterated allow key then an
				// additional deny entry with port/proto of newKey and with the
				// identity of the iterated allow key must be added.
				denyKeyCpy := newKey
				denyKeyCpy.Identity = k.Identity
				l3l4DenyEntry := NewMapStateEntry(newKey, newEntry.DerivedFromRules, 0, "", 0, true, DefaultAuthType, AuthTypeDisabled)
				ms.addKeyWithChanges(denyKeyCpy, l3l4DenyEntry, changes)
				newEntry.AddDependent(denyKeyCpy)
			}
			return true
		})
				return true
			}

			// A narrower of two deny keys is redundant in the datapath only if
			// the broader ID is 0, or the IDs are the same. This is because the
			// ID will be assigned from the ipcache and datapath has no notion
			// of one ID being related to another (e.g., in a CIDR sense).
			if (k.Identity == 0 || k.Identity == newKey.Identity) &&
				(k.PortProtoIsEqual(newKey) || k.PortProtoIsBroader(newKey)) {
				// If this iterated-deny-entry is an deny-all-L3 or has the same ID
				// as the new-entry and the iterated-deny-entry has a broader (or
				// equal) port-protocol it will match all the packets the newKey
				// would, given that we do not allow more specific allow rules to be
				// inserted.

				// Identical key needs to be added if owners are different (to merge
				// them). This has no effect on the datapath policy map but is
				// needed for internal bookkeeping.
				if k != newKey || v.HasSameOwners(&newEntry) {
					bailed = true
					return false
				}
			} else if (newKey.Identity == 0 || newKey.Identity == k.Identity) &&
				(newKey.PortProtoIsEqual(k) || newKey.PortProtoIsBroader(k)) &&
				!newEntry.HasDependent(k) {
				// If this iterated-deny-entry is a subset (or equal) of the
				// new-entry and the new-entry has a broader (or equal)
				// port-protocol the newKey will match all the packets the iterated
				// key would, given that there are no more specific or L4-only allow
				// entries. We removed the more specific allow rules in the loop
				// above, and added more specific deny rules if there was an L4-only
				// allow rule. We use 'HasDependant' to figure out that 'k' must
				// remain to take precedence over the L4-only allow key.

				// Identical key would have been captured in the block above, so we
				// do not need to check for it here.
				ms.deleteKeyWithChanges(k, nil, changes)
			}
			return true
	} else {
		// NOTE: We do not delete redundant allow entries.
		bailed := false
		changeToDeny := false
		ms.ForEachDeny(func(k Key, v MapStateEntry) bool {
			// Protocols and traffic directions that don't match ensure that the policies
			// do not interact in anyway.
					// identity is removed.
					ms.addDependentOnEntry(k, v, denyKeyCpy, changes)
				}
			} else if k.PortProtoIsBroader(newKey) || k.PortProtoIsEqual(newKey) {
				if k.Identity == 0 || k.Identity == newKey.Identity {
					// If the iterated-deny-entry is a datapath superset (or
					// equal) of the new-entry and has a broader (or equal)
					// port-protocol than the new-entry then the new entry
					// should not be inserted.
					bailed = true
					return false
				} else if identityIsSupersetOf(k.Identity, newKey.Identity, identities) {
					// if newKey is not bailed due to being covered in the
					// datapath by a deny entry, but is covered by a deny entry
					// in the CIDR sense, we must change this allow entry to a
					// deny entry so that the covering deny policy is honored
					// also for this ID in the datapath.
					changeToDeny = true
				}
			} else { // newKey.PortProtoIsBroader(k)
				if identityIsSupersetOf(k.Identity, newKey.Identity, identities) {
					// If the new-entry is a subset of the iterated-deny-entry
					// and the new-entry has a less specific port-protocol than
					// the iterated-deny-entry then an additional copy of the
					// iterated-deny-entry with the identity of the new-entry
					// must be added.
					denyKeyCpy := k
					denyKeyCpy.Identity = newKey.Identity
					l3l4DenyEntry := NewMapStateEntry(k, v.DerivedFromRules, 0, "", 0, true, DefaultAuthType, AuthTypeDisabled)
					ms.addKeyWithChanges(denyKeyCpy, l3l4DenyEntry, changes)
					// L3-only entries can be deleted incrementally so we need
					// to track their effects on other entries so that those
					// effects can be reverted when the identity is removed.
					ms.addDependentOnEntry(k, v, denyKeyCpy, changes)
				}
			}

			return true
		})

		if changeToDeny {
			newEntry.IsDeny = true
			newEntry.ProxyPort = 0
			newEntry.Listener = ""
			newEntry.priority = 0
			newEntry.hasAuthType = DefaultAuthType
			newEntry.AuthType = AuthTypeDisabled
		}
		if !bailed {
			ms.authPreferredInsert(newKey, newEntry, features, changes)
		}
