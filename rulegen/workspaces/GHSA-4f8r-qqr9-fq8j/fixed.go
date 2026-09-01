package main

package main

import (
	"crypto"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rdimitrov/ngo-tuf/metadata"
	"github.com/rdimitrov/ngo-tuf/repo"
	"github.com/sigstore/sigstore/pkg/signature"
	"golang.org/x/crypto/ed25519"
)

// A TUF repository example using the low-level TUF Metadata API.

// The example code in this file demonstrates how to *manually* create and
// maintain repository metadata using the low-level Metadata API.
// Contents:
//  * creation of top-level metadata
//  * target file handling
//  * consistent snapshots
//  * key management
//  * top-level delegation and signing thresholds
//  * metadata verification
//  * target delegation
//  * in-band and out-of-band metadata signing
//  * writing and reading metadata files
//  * root key rotation

// NOTE: Metadata files will be written to a 'tmp*'-directory in CWD.

func main() {
	// Create top-level metadata
	// =========================
	// Every TUF repository has at least four roles, i.e. the top-level roles
	// 'targets', 'snapshot', 'timestamp' and 'root'. Below we will discuss their
	// purpose, show how to create the corresponding metadata, and how to use them
	// to provide integrity, consistency and freshness for the files TUF aims to
	// protect, i.e. target files.

	// Define containers for metadata objects and cryptographic keys created below. This
	// allows us to sign and write metadata in a batch more easily.
	roles := repo.New()
	keys := map[string]ed25519.PrivateKey{}

	// Targets (integrity)
	// -------------------
	// The targets role guarantees integrity for the files that TUF aims to protect,
	// i.e. target files. It does so by listing the relevant target files, along
	// with their hash and length.
	targets := metadata.Targets(helperExpireIn(7))
	roles.SetTargets("targets", targets)

	// For the purpose of this example we use the top-level targets role to protect
	// the integrity of this very example script. The metadata entry contains the
	// hash and length of this file at the local path. In addition, it specifies the
	// 'target path', which a client uses to locate the target file relative to a
	// configured mirror base URL.
	// 	   |----base URL---||-----target path-----|
	// e.g. tuf-examples.org/examples/basic_repo.py
	targetPath, localPath := helperGetPathForTarget("basic_repo.go")
	targetFileInfo, err := metadata.TargetFile().FromFile(targetPath, localPath)
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "generating target file info failed", err))
	}
	targets.Signed.Targets[targetPath] = *targetFileInfo
	// Snapshot (consistency)
	// ----------------------
	// The snapshot role guarantees consistency of the entire repository. It does so
	// by listing all available targets metadata files at their latest version. This
	// becomes relevant, when there are multiple targets metadata files in a
	// repository and we want to protect the client against mix-and-match attacks.
	snapshot := metadata.Snapshot(helperExpireIn(7))
	roles.SetSnapshot(snapshot)
	// Timestamp (freshness)
	// ---------------------
	// The timestamp role guarantees freshness of the repository metadata. It does
	// so by listing the latest snapshot (which in turn lists all the latest
	// targets) metadata. A short expiration interval requires the repository to
	// regularly issue new timestamp metadata and thus protects the client against
	// freeze attacks.
	// Note that snapshot and timestamp use the same generic wireline metadata
	// format.
	timestamp := metadata.Timestamp(helperExpireIn(1))
	roles.SetTimestamp(timestamp)

	// Root (root of trust)
	// --------------------
	// The root role serves as root of trust for all top-level roles, including
	// itself. It does so by mapping cryptographic keys to roles, i.e. the keys that
	// are authorized to sign any top-level role metadata, and signing thresholds,
	// i.e. how many authorized keys are required for a given role (see 'roles'
	// field). This is called top-level delegation.

	// In addition, root provides all public keys to verify these signatures (see
	// 'keys' field), and a configuration parameter that describes whether a
	// repository uses consistent snapshots (see section 'Persist metadata' below
	// for more details).

	// Create root metadata object
	root := metadata.Root(helperExpireIn(365))
	roles.SetRoot(root)

	// For this example, we generate one private key of type 'ed25519' for each top-level role
	for _, name := range []string{"targets", "snapshot", "timestamp", "root"} {
		_, private, err := ed25519.GenerateKey(nil)
		if err != nil {
			panic(fmt.Sprintln("basic_repo.go:", "key generation failed", err))
		}
		keys[name] = private
		key, err := metadata.KeyFromPublicKey(private.Public())
		if err != nil {
			panic(fmt.Sprintln("basic_repo.go:", "key conversion failed", err))
		}
		err = roles.Root().Signed.AddKey(key, name)
		if err != nil {
			panic(fmt.Sprintln("basic_repo.go:", "adding key to root failed", err))
		}
	}
	// NOTE: We only need the public part to populate root, so it is possible to use
	// out-of-band mechanisms to generate key pairs and only expose the public part
	// to whoever maintains the root role. As a matter of fact, the very purpose of
	// signature thresholds is to avoid having private keys all in one place.

	// Signature thresholds
	// --------------------
	// Given the importance of the root role, it is highly recommended to require a
	// threshold of multiple keys to sign root metadata. For this example we
	// generate another root key (you can pretend it's out-of-band) and increase the
	// required signature threshold.
	_, anotherRootKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "key generation failed", err))
	}
	// TODO: Extend the example to showcase a mixture of keys, i.e.
	// anotherRootKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	// anotherRootKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	anotherKey, err := metadata.KeyFromPublicKey(anotherRootKey.Public())
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "key conversion failed", err))
	}
	err = roles.Root().Signed.AddKey(anotherKey, "root")
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "adding another key to root failed", err))
	}
	roles.Root().Signed.Roles["root"].Threshold = 2

	// Sign top-level metadata (in-band)
	// =================================
	// In this example we have access to all top-level signing keys, so we can use
	// them to create and add a signature for each role metadata.
	for _, name := range []string{"targets", "snapshot", "timestamp", "root"} {
		key := keys[name]
		signer, err := signature.LoadSigner(key, crypto.Hash(0))
		if err != nil {
			panic(fmt.Sprintln("basic_repo.go:", "loading a signer failed", err))
		}
		switch name {
		case "targets":
			_, err = roles.Targets("targets").Sign(signer)
		case "snapshot":
			_, err = roles.Snapshot().Sign(signer)
		case "timestamp":
			_, err = roles.Timestamp().Sign(signer)
		case "root":
			_, err = roles.Root().Sign(signer)
		}
		if err != nil {
			panic(fmt.Sprintln("basic_repo.go:", "metadata signing failed", err))
		}
	}

	// Persist metadata (consistent snapshot)
	// ======================================
	// It is time to publish the first set of metadata for a client to safely
	// download the target file that we have registered for this example repository.

	// For the purpose of this example we will follow the consistent snapshot naming
	// convention for all metadata. This means that each metadata file, must be
	// prefixed with its version number, except for timestamp. The naming convention
	// also affects the target files, but we don't cover this in the example. See
	// the TUF specification for more details:
	// https://theupdateframework.github.io/specification/latest/#writing-consistent-snapshots

	// Also note that the TUF specification does not mandate a wireline format. In
	// this demo we use a non-compact JSON format and store all metadata in
	// temporary directory at CWD for review.
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "getting cwd failed", err))
	}
	tmpDir, err := os.MkdirTemp(cwd, "tmp")
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "creating a temporary folder failed", err))
	}

	for _, name := range []string{"targets", "snapshot", "timestamp", "root"} {
		switch name {
		case "targets":
			filename := fmt.Sprintf("%d.%s.json", roles.Targets("targets").Signed.Version, name)
			err = roles.Targets("targets").ToFile(filepath.Join(tmpDir, filename), true)
		case "snapshot":
			filename := fmt.Sprintf("%d.%s.json", roles.Snapshot().Signed.Version, name)
			err = roles.Snapshot().ToFile(filepath.Join(tmpDir, filename), true)
		case "timestamp":
			filename := fmt.Sprintf("%s.json", name)
			err = roles.Timestamp().ToFile(filepath.Join(tmpDir, filename), true)
		case "root":
			filename := fmt.Sprintf("%d.%s.json", roles.Root().Signed.Version, name)
			err = roles.Root().ToFile(filepath.Join(tmpDir, filename), true)
		}
		if err != nil {
			panic(fmt.Sprintln("basic_repo.go:", "saving metadata to file failed", err))
		}
	}

	// Threshold signing (out-of-band)
	// ===============================
	// As mentioned above, using signature thresholds usually entails that not all
	// signing keys for a given role are in the same place. Let's briefly pretend
	// this is the case for the second root key we registered above, and we are now
	// on that key owner's computer. All the owner has to do is read the metadata
	// file, sign it, and write it back to the same file, and this can be repeated
	// until the threshold is satisfied.
	outofbandRoot, err := metadata.Root().FromFile(filepath.Join(tmpDir, "1.root.json"))
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "loading root metadata from file failed", err))
	}
	outofbandSigner, err := signature.LoadSigner(anotherRootKey, crypto.Hash(0))
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "loading a signer failed", err))
	}
	_, err = outofbandRoot.Sign(outofbandSigner)
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "signing root failed", err))
	}
	err = outofbandRoot.ToFile(filepath.Join(tmpDir, "1.root.json"), true)
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "saving root metadata to file failed", err))
	}

	// Verify that metadata is signed correctly
	// ====================================
	// Verify root
	err = outofbandRoot.VerifyDelegate("root", outofbandRoot)
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "verifying root metadata failed", err))
	}

	// Verify targets
	err = outofbandRoot.VerifyDelegate("targets", roles.Targets("targets"))
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "verifying targets metadata failed", err))
	}

	// Verify snapshot
	err = outofbandRoot.VerifyDelegate("snapshot", roles.Snapshot())
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "verifying snapshot metadata failed", err))
	}

	// Verify timestamp
	err = outofbandRoot.VerifyDelegate("timestamp", roles.Timestamp())
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "verifying timestamp metadata failed", err))
	}

	// Targets delegation
	// ==================
	// Similar to how the root role delegates responsibilities about integrity,
	// consistency and freshness to the corresponding top-level roles, a targets
	// role may further delegate its responsibility for target files (or a subset
	// thereof) to other targets roles. This allows creation of a granular trust
	// hierarchy, and further reduces the impact of a single role compromise.

	// In this example the top-level targets role trusts a new "go-scripts"
	// targets role to provide integrity for any target file that ends with ".go".
	delegateeName := "go-scripts"
	_, delegateePrivateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "key generation failed", err))
	}
	keys[delegateeName] = delegateePrivateKey

	// Delegatee
	// ---------
	// Create a new targets role, akin to how we created top-level targets above, and
	// add target file info from above according to the delegatee's responsibility.
	delegatee := metadata.Targets(helperExpireIn(7))
	delegatee.Signed.Targets[targetPath] = *targetFileInfo
	roles.SetTargets(delegateeName, delegatee)

	// Delegator
	// ---------
	// Akin to top-level delegation, the delegator expresses its trust in the
	// delegatee by authorizing a threshold of cryptographic keys to provide
	// signatures for the delegatee metadata. It also provides the corresponding
	// public key store.
	// The delegation info defined by the delegator further requires the provision
	// of a unique delegatee name and constraints about the target files the
	// delegatee is responsible for, e.g. a list of path patterns. For details about
	// all configuration parameters see
	// https://theupdateframework.github.io/specification/latest/#delegations
	delegateeKey, _ := metadata.KeyFromPublicKey(delegateePrivateKey.Public())
	roles.Targets("targets").Signed.Delegations = &metadata.Delegations{
		Keys: map[string]*metadata.Key{
			delegateeKey.ID(): delegateeKey,
		},
		Roles: []metadata.DelegatedRole{
			{
				Name:        delegateeName,
				KeyIDs:      []string{delegateeKey.ID()},
				Threshold:   1,
				Terminating: true,
				Paths:       []string{"*.go"},
			},
		},
	}

	// Remove target file info from top-level targets (delegatee is now responsible)
	delete(roles.Targets("targets").Signed.Targets, targetPath)

	// Increase expiry (delegators should be less volatile)
	roles.Targets("targets").Signed.Expires = helperExpireIn(365)

	// Snapshot + Timestamp + Sign + Persist
	// -------------------------------------
	// In order to publish a new consistent set of metadata, we need to update
	// dependent roles (snapshot, timestamp) accordingly, bumping versions of all
	// changed metadata.

	// Bump targets version
	roles.Targets("targets").Signed.Version += 1

	// Update snapshot to account for changed and new targets(delegatee) metadata
	roles.Snapshot().Signed.Meta["targets.json"] = *metadata.MetaFile(roles.Targets("targets").Signed.Version)
	roles.Snapshot().Signed.Meta[delegateeName] = *metadata.MetaFile(1)
	roles.Snapshot().Signed.Version += 1

	// Update timestamp to account for changed snapshot metadata
	roles.Timestamp().Signed.Meta["snapshot.json"] = *metadata.MetaFile(roles.Snapshot().Signed.Version)
	roles.Timestamp().Signed.Version += 1

	// Sign and write metadata for all changed roles, i.e. all but root
	for _, name := range []string{"targets", "snapshot", "timestamp", delegateeName} {
		key := keys[name]
		signer, err := signature.LoadSigner(key, crypto.Hash(0))
		if err != nil {
			panic(fmt.Sprintln("basic_repo.go:", "loading a signer failed", err))
		}
		switch name {
		case "targets":
			roles.Targets("targets").ClearSignatures()
			_, err = roles.Targets("targets").Sign(signer)
			if err != nil {
				panic(fmt.Sprintln("basic_repo.go:", "signing metadata failed", err))
			}
			filename := fmt.Sprintf("%d.%s.json", roles.Targets("targets").Signed.Version, name)
			err = roles.Targets("targets").ToFile(filepath.Join(tmpDir, filename), true)
		case "snapshot":
			roles.Snapshot().ClearSignatures()
			_, err = roles.Snapshot().Sign(signer)
			if err != nil {
				panic(fmt.Sprintln("basic_repo.go:", "signing metadata failed", err))
			}
			filename := fmt.Sprintf("%d.%s.json", roles.Snapshot().Signed.Version, name)
			err = roles.Snapshot().ToFile(filepath.Join(tmpDir, filename), true)
		case "timestamp":
			roles.Timestamp().ClearSignatures()
			_, err = roles.Timestamp().Sign(signer)
			if err != nil {
				panic(fmt.Sprintln("basic_repo.go:", "signing metadata failed", err))
			}
			filename := fmt.Sprintf("%s.json", name)
			err = roles.Timestamp().ToFile(filepath.Join(tmpDir, filename), true)
		case delegateeName:
			roles.Targets(delegateeName).ClearSignatures()
			_, err = roles.Targets(delegateeName).Sign(signer)
			if err != nil {
				panic(fmt.Sprintln("basic_repo.go:", "signing metadata failed", err))
			}
			filename := fmt.Sprintf("%d.%s.json", roles.Targets(delegateeName).Signed.Version, name)
			err = roles.Targets(delegateeName).ToFile(filepath.Join(tmpDir, filename), true)
		}
		if err != nil {
			panic(fmt.Sprintln("basic_repo.go:", "saving metadata to file failed", err))
		}
	}

	// Root key rotation (recover from a compromise / key loss)
	// ========================================================
	// TUF makes it easy to recover from a key compromise in-band. Given the trust
	// hierarchy through top-level and targets delegation you can easily
	// replace compromised or lost keys for any role using the delegating role, even
	// for the root role.
	// However, since root authorizes its own keys, it always has to be signed with
	// both the threshold of keys from the previous version and the threshold of
	// keys from the new version. This establishes a trusted line of continuity.

	// In this example we will replace a root key, and sign a new version of root
	// with the threshold of old and new keys. Since one of the previous root keys
	// remains in place, it can be used to count towards the old and new threshold.
	_, newRootKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "key generation failed", err))
	}
	oldRootKey, err := metadata.KeyFromPublicKey(keys["root"].Public())
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "key conversion failed", err))
	}
	err = roles.Root().Signed.RevokeKey(oldRootKey.ID(), "root")
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "revoking key failed", err))
	}
	// Add new key for root
	newRootKeyTUF, err := metadata.KeyFromPublicKey(newRootKey.Public())
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "key conversion failed", err))
	}
	err = roles.Root().Signed.AddKey(newRootKeyTUF, "root")
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "adding key to root failed", err))
	}
	roles.Root().Signed.Version += 1
	roles.Root().ClearSignatures()

	// Sign root
	for _, k := range []ed25519.PrivateKey{keys["root"], anotherRootKey, newRootKey} {
		signer, err := signature.LoadSigner(k, crypto.Hash(0))
		if err != nil {
			panic(fmt.Sprintln("basic_repo.go:", "loading a signer failed", err))
		}
		_, err = roles.Root().Sign(signer)
		if err != nil {
			panic(fmt.Sprintln("basic_repo.go:", "signing root failed", err))
		}
	}
	filename := fmt.Sprintf("%d.%s.json", roles.Root().Signed.Version, "root")
	err = roles.Root().ToFile(filepath.Join(tmpDir, filename), true)
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "saving root to file failed", err))
	}
	fmt.Println("Done! Metadata files location:", tmpDir)
}

// helperExpireIn returns time offset by days
func helperExpireIn(days int) time.Time {
	return time.Now().AddDate(0, 0, days).UTC()
}

// helperGetPathForTarget returns local and target paths for target
func helperGetPathForTarget(name string) (string, string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintln("basic_repo.go:", "getting cwd failed", err))
	}
	_, dir := filepath.Split(cwd)
	return filepath.Join(dir, name), filepath.Join(cwd, name)
}
package metadata

import (
	"bytes"
	"crypto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/secure-systems-lab/go-securesystemslib/cjson"
	"github.com/sigstore/sigstore/pkg/signature"
)

// Root create new metadata instance of type Root
func Root(expires ...time.Time) *Metadata[RootType] {
	// expire now if there's nothing set
	if len(expires) == 0 {
		expires = []time.Time{time.Now().UTC()}
	}
	roles := map[string]*Role{}
	for _, r := range []string{ROOT, SNAPSHOT, TARGETS, TIMESTAMP} {
		roles[r] = &Role{
			KeyIDs:    []string{},
			Threshold: 1,
		}
	}
	return &Metadata[RootType]{
		Signed: RootType{
			Type:               "root",
			SpecVersion:        SPECIFICATION_VERSION,
			Version:            1,
			Expires:            expires[0],
			Keys:               map[string]*Key{},
			Roles:              roles,
			ConsistentSnapshot: false,
		},
		Signatures: []Signature{},
	}
}

// Snapshot create new metadata instance of type Snapshot
func Snapshot(expires ...time.Time) *Metadata[SnapshotType] {
	// expire now if there's nothing set
	if len(expires) == 0 {
		expires = []time.Time{time.Now().UTC()}
	}
	return &Metadata[SnapshotType]{
		Signed: SnapshotType{
			Type:        "snapshot",
			SpecVersion: SPECIFICATION_VERSION,
			Version:     1,
			Expires:     expires[0],
			Meta: map[string]MetaFiles{
				"targets.json": {
					Version: 1,
				},
			},
		},
		Signatures: []Signature{},
	}
}

// Timestamp create new metadata instance of type Timestamp
func Timestamp(expires ...time.Time) *Metadata[TimestampType] {
	// expire now if there's nothing set
	if len(expires) == 0 {
		expires = []time.Time{time.Now().UTC()}
	}
	return &Metadata[TimestampType]{
		Signed: TimestampType{
			Type:        "timestamp",
			SpecVersion: SPECIFICATION_VERSION,
			Version:     1,
			Expires:     expires[0],
			Meta: map[string]MetaFiles{
				"snapshot.json": {
					Version: 1,
				},
			},
		},
		Signatures: []Signature{},
	}
}

// Targets create new metadata instance of type Targets
func Targets(expires ...time.Time) *Metadata[TargetsType] {
	// expire now if there's nothing set
	if len(expires) == 0 {
		expires = []time.Time{time.Now().UTC()}
	}
	return &Metadata[TargetsType]{
		Signed: TargetsType{
			Type:        "targets",
			SpecVersion: SPECIFICATION_VERSION,
			Version:     1,
			Expires:     expires[0],
			Targets:     map[string]TargetFiles{},
			Delegations: &Delegations{
				Keys:  map[string]*Key{},
				Roles: []DelegatedRole{},
			},
		},
		Signatures: []Signature{},
	}
}

// TargetFile create new metadata instance of type TargetFiles
func TargetFile() *TargetFiles {
	return &TargetFiles{
		Length: 0,
		Hashes: Hashes{},
	}
}

// MetaFile create new metadata instance of type MetaFile
func MetaFile(version int64) *MetaFiles {
	return &MetaFiles{
		Length:  0,
		Hashes:  Hashes{},
		Version: version,
	}
}

// FromFile load metadata from file
func (meta *Metadata[T]) FromFile(name string) (*Metadata[T], error) {
	m, err := fromFile[T](name)
	if err != nil {
		return nil, fmt.Errorf("error generating metadata from bytes - %s", name)
	}
	*meta = *m
	return meta, nil
}

// FromBytes deserialize metadata from bytes
func (meta *Metadata[T]) FromBytes(bytes []byte) (*Metadata[T], error) {
	m, err := fromBytes[T](bytes)
	if err != nil {
		return nil, err
	}
	*meta = *m
	return meta, nil
}

// ToBytes serialize metadata to bytes
func (meta *Metadata[T]) ToBytes(pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(*meta, "", "\t")
	}
	return json.Marshal(*meta)
}

// ToFile save metadata to file
func (meta *Metadata[T]) ToFile(name string, pretty bool) error {
	bytes, err := meta.ToBytes(pretty)
	if err != nil {
		return fmt.Errorf("failed serializing metadata")
	}
	return ioutil.WriteFile(name, bytes, 0644)
}

// Sign create signature over Signed and assign it to Signatures
func (meta *Metadata[T]) Sign(signer signature.Signer) (*Signature, error) {
	// encode the Signed part to canonical JSON so signatures are consistent
	payload, err := cjson.EncodeCanonical(meta.Signed)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Signed in canonical format during Sign()")
	}
	// sign the Signed part
	sb, err := signer.SignMessage(bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to Sign(), returned signature should not be nil")
	}
	// get the signer's PublicKey
	publ, err := signer.PublicKey()
	if err != nil {
		return nil, err
	}
	// convert to TUF Key type to get keyID
	key, err := KeyFromPublicKey(publ)
	if err != nil {
		return nil, err
	}
	// build signature
	sig := &Signature{
		KeyID:     key.ID(),
		Signature: sb,
	}
	// update the Signatures part
	meta.Signatures = append(meta.Signatures, *sig)
	// return the new signature
	return sig, nil
}

// VerifyDelegate verifies that “delegated_metadata“ is signed with the required
// threshold of keys for the delegated role “delegated_role“
func (meta *Metadata[T]) VerifyDelegate(delegated_role string, delegated_metadata any) error {
	var keys map[string]*Key
	var roleKeyIDs []string
	var roleThreshold int
	var sign Signature
	var payload []byte
	signing_keys := map[string]bool{}
	i := any(meta)
	// collect keys, keyIDs and threshold based on delegator type
	switch i := i.(type) {
	case *Metadata[RootType]:
		keys = i.Signed.Keys
		if role, ok := (*i).Signed.Roles[delegated_role]; ok {
			roleKeyIDs = role.KeyIDs
			roleThreshold = role.Threshold
		} else {
			return fmt.Errorf("no delegation found for %s", delegated_role)
		}
	case *Metadata[TargetsType]:
		keys = i.Signed.Delegations.Keys
		for _, v := range i.Signed.Delegations.Roles {
			if v.Name == delegated_role {
				roleKeyIDs = v.KeyIDs
				roleThreshold = v.Threshold
				break
			}
		}
	default:
		return fmt.Errorf("call is valid only on delegator metadata (root or targets)")
	}
	// if there are no keyIDs for that role it means there's no delegation found
	if len(roleKeyIDs) == 0 {
		fmt.Println("no delegation found for", delegated_role)
		return fmt.Errorf("no delegation found for %s", delegated_role)
	}
	// loop through each role keyID
	for _, v := range roleKeyIDs {
		// convert to a PublicKey type
		key, err := keys[v].ToPublicKey()
		if err != nil {
			fmt.Println("failed to generate crypto.PublicKey from Key")
			return err
		}
		// load a verifier based on that key
		verifier, err := signature.LoadVerifier(key, crypto.Hash(0))
		if err != nil {
			fmt.Println("failed to load verifier")
			return err
		}
		// collect the signature for that key and build the payload we'll verify
		// based on the Signed part of the delegated metadata
		switch d := delegated_metadata.(type) {
		case *Metadata[RootType]:
			for _, s := range d.Signatures {
				if s.KeyID == v {
					sign = s
				}
			}
			payload, err = cjson.EncodeCanonical(d.Signed)
			if err != nil {
				fmt.Println("failed to encode Signed in canonical format during verify")
			}
		case *Metadata[SnapshotType]:
			for _, s := range d.Signatures {
				if s.KeyID == v {
					sign = s
				}
			}
			payload, err = cjson.EncodeCanonical(d.Signed)
			if err != nil {
				fmt.Println("failed to encode Signed in canonical format during verify")
			}
		case *Metadata[TimestampType]:
			for _, s := range d.Signatures {
				if s.KeyID == v {
					sign = s
				}
			}
			payload, err = cjson.EncodeCanonical(d.Signed)
			if err != nil {
				fmt.Println("failed to encode Signed in canonical format during verify")
			}
		case *Metadata[TargetsType]:
			for _, s := range d.Signatures {
				if s.KeyID == v {
					sign = s
				}
			}
			payload, err = cjson.EncodeCanonical(d.Signed)
			if err != nil {
				fmt.Println("failed to encode Signed in canonical format during verify")
			}
		default:
			fmt.Println("unknown delegated metadata type")
		}
		// verify if the signature for that payload corresponds to the given key
		if err := verifier.VerifySignature(bytes.NewReader(sign.Signature), bytes.NewReader(payload)); err == nil {
			// save the verified keyID only if there's no err value
			signing_keys[v] = true
		}
	}
	// check if the amount of valid signatures is enough
	if len(signing_keys) < roleThreshold {
		return fmt.Errorf("signature verification failed, not enough signatures")
	}
	return nil
}

// IsExpired returns true if metadata is expired.
// It checks if referenceTime is after Signed.Expires
func (signed *RootType) IsExpired(referenceTime time.Time) bool {
	return referenceTime.After(signed.Expires)
}

// IsExpired returns true if metadata is expired.
// It checks if referenceTime is after Signed.Expires
func (signed *SnapshotType) IsExpired(referenceTime time.Time) bool {
	return referenceTime.After(signed.Expires)
}

// IsExpired returns true if metadata is expired.
// It checks if referenceTime is after Signed.Expires
func (signed *TimestampType) IsExpired(referenceTime time.Time) bool {
	return referenceTime.After(signed.Expires)
}

// IsExpired returns true if metadata is expired.
// It checks if referenceTime is after Signed.Expires
func (signed *TargetsType) IsExpired(referenceTime time.Time) bool {
	return referenceTime.After(signed.Expires)
}

// VerifyLengthHashes checks whether the data matches its corresponding
// length and hashes
func (f *MetaFiles) VerifyLengthHashes(data []byte) error {
	err := verifyHashes(data, f.Hashes)
	if err != nil {
		return err
	}
	err = verifyLength(data, f.Length)
	if err != nil {
		return err
	}
	return nil
}

// FromFile generates TargetFiles from file
func (t *TargetFiles) FromFile(targetPath, localPath string) (*TargetFiles, error) {
	return &TargetFiles{}, nil
}

// ClearSignatures clears the Signatures
func (meta *Metadata[T]) ClearSignatures() {
	meta.Signatures = []Signature{}
}
package metadata

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/secure-systems-lab/go-securesystemslib/cjson"
	"golang.org/x/exp/slices"
)

const (
	// MaxJSONKeySize defines the maximum length of a JSON payload.
	MaxJSONKeySize = 512 * 1024 // 512Kb
	KeyIDLength    = sha256.Size * 2

	KeyTypeEd25519           KeyType = "ed25519"
	KeyTypeECDSA_SHA2_P256   KeyType = "ecdsa-sha2-nistp256"
	KeyTypeRSASSA_PSS_SHA256 KeyType = "rsa"

	KeySchemeEd25519           KeyScheme = "ed25519"
	KeySchemeECDSA_SHA2_P256   KeyScheme = "ecdsa-sha2-nistp256"
	KeySchemeRSASSA_PSS_SHA256 KeyScheme = "rsassa-pss-sha256"
)

type helperED25519 struct {
	PublicKey HexBytes `json:"public"`
}
type helperRSAECDSA struct {
	PublicKey crypto.PublicKey `json:"public"`
}

// ToPublicKey generate crypto.PublicKey from metadata type Key
func (k *Key) ToPublicKey() (crypto.PublicKey, error) {
	switch k.Type {
	case KeyTypeRSASSA_PSS_SHA256:
		return k.toPublicKeyRSA()
	case KeyTypeECDSA_SHA2_P256:
		return k.toPublicKeyECDSA()
	case KeyTypeEd25519:
		return k.toPublicKeyED25519()
	}
	return nil, fmt.Errorf("unsupported public key type")
}

// KeyFromPublicKey generate metadata type Key from crypto.PublicKey
func KeyFromPublicKey(k crypto.PublicKey) (*Key, error) {
	var b []byte
	var err error
	key := &Key{}
	switch k := k.(type) {
	case *rsa.PublicKey:
		key.Type = KeyTypeRSASSA_PSS_SHA256
		key.Scheme = KeySchemeRSASSA_PSS_SHA256
		// pemKey, err := cryptoutils.MarshalPublicKeyToPEM(k)
		s := &helperRSAECDSA{
			PublicKey: k,
			// PublicKey: string(pemKey),
		}
		b, err = json.Marshal(s)
		if err != nil {
			return nil, err
		}
	case *ecdsa.PublicKey:
		key.Type = KeyTypeECDSA_SHA2_P256
		key.Scheme = KeySchemeECDSA_SHA2_P256
		// pemKey, err := cryptoutils.MarshalPublicKeyToPEM(k)
		s := &helperRSAECDSA{
			PublicKey: k,
			// PublicKey: string(pemKey),
		}
		b, err = json.Marshal(s)
		if err != nil {
			return nil, err
		}
	case ed25519.PublicKey:
		key.Type = KeyTypeEd25519
		key.Scheme = KeySchemeEd25519
		s := &helperED25519{
			PublicKey: []byte(k),
		}
		b, err = json.Marshal(s)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported public key type")
	}
	key.Value = b
	return key, nil
}

// AddKey adds new signing key for delegated role "role"
// keyID: Identifier of the key to be added for “role“.
// key: Signing key to be added for “role“.
// role: Name of the role, for which “key“ is added.
func (signed *RootType) AddKey(key *Key, role string) error {
	// verify role is present
	if _, ok := signed.Roles[role]; !ok {
		return fmt.Errorf("Role %s doesn't exist", role)
	}
	// add keyID to role
	if !slices.Contains(signed.Roles[role].KeyIDs, key.ID()) {
		signed.Roles[role].KeyIDs = append(signed.Roles[role].KeyIDs, key.ID())
	}
	// update Keys
	signed.Keys[key.ID()] = key
	return nil
}

// RevokeKey revoke key from “role“ and updates the Keys store.
// keyID: Identifier of the key to be removed for “role“.
// role: Name of the role, for which a signing key is removed.
func (signed *RootType) RevokeKey(keyID, role string) error {
	// verify role is present
	if _, ok := signed.Roles[role]; !ok {
		return fmt.Errorf("Role %s doesn't exist", role)
	}
	// verify keyID is present for given role
	if !slices.Contains(signed.Roles[role].KeyIDs, keyID) {
		return fmt.Errorf("Key with id %s is not used by %s", keyID, role)
	}
	// remove keyID from role
	filteredKeyIDs := []string{}
	for _, k := range signed.Roles[role].KeyIDs {
		if k != keyID {
			filteredKeyIDs = append(filteredKeyIDs, k)
		}
	}
	// overwrite the old keyID slice
	signed.Roles[role].KeyIDs = filteredKeyIDs
	// check if keyID is used by other roles too
	for _, r := range signed.Roles {
		if slices.Contains(r.KeyIDs, keyID) {
			return nil
		}
	}
	// delete the keyID from Keys if it's not used anywhere else
	delete(signed.Keys, keyID)
	return nil
}

// AddKey adds new signing key for delegated role "role"
// key: Signing key to be added for “role“.
// role: Name of the role, for which “key“ is added.
func (signed *TargetsType) AddKey(key *Key, role string) error {
	// check if Delegations are even present
	if signed.Delegations == nil {
		return fmt.Errorf("delegated role %s doesn't exist", role)
	}
	// loop through all delegated roles
	for i, d := range signed.Delegations.Roles {
		// if role is found
		if d.Name == role {
			// add key if keyID is not already part of keyIDs for that role
			if !slices.Contains(d.KeyIDs, key.ID()) {
				signed.Delegations.Roles[i].KeyIDs = append(signed.Delegations.Roles[i].KeyIDs, key.ID())
				signed.Delegations.Keys[key.ID()] = key
				return nil
			}
			return fmt.Errorf("delegated role %s already has keyID %s", role, key.ID())
		}
	}
	return fmt.Errorf("delegated role %s doesn't exist", role)
}

// RevokeKey revokes key from delegated role "role" and updates the delegations key store
// keyID: Identifier of the key to be removed for “role“.
// role: Name of the role, for which a signing key is removed.
func (signed *TargetsType) RevokeKey(keyID string, role string) error {
	// check if Delegations are even present
	if signed.Delegations == nil {
		return fmt.Errorf("delegated role %s doesn't exist", role)
	}
	// loop through all delegated roles
	for i, d := range signed.Delegations.Roles {
		// if role is found
		if d.Name == role {
			// check if keyID is present in keyIDs for that role
			if !slices.Contains(d.KeyIDs, keyID) {
				return fmt.Errorf("Key with id %s is not used by %s", keyID, role)
			}
			// remove keyID from role
			filteredKeyIDs := []string{}
			for _, k := range signed.Delegations.Roles[i].KeyIDs {
				if k != keyID {
					filteredKeyIDs = append(filteredKeyIDs, k)
				}
			}
			// overwrite the old keyID slice
			signed.Delegations.Roles[i].KeyIDs = filteredKeyIDs
			break
		}
	}
	// check if keyID is used by other roles too
	for _, r := range signed.Delegations.Roles {
		if slices.Contains(r.KeyIDs, keyID) {
			return nil
		}
	}
	// delete the keyID from Keys if it's not used anywhere else
	delete(signed.Delegations.Keys, keyID)
	return nil
}

func (k *Key) toPublicKeyED25519() (crypto.PublicKey, error) {
	// Prepare decoder limited to 512Kb
	dec := json.NewDecoder(io.LimitReader(bytes.NewReader(k.Value), MaxJSONKeySize))
	s := &helperED25519{}
	// Unmarshal key value
	if err := dec.Decode(s); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, fmt.Errorf("the public key is truncated or too large: %w", err)
		}
		return nil, err
	}
	if n := len(s.PublicKey); n != ed25519.PublicKeySize {
		return nil, fmt.Errorf("unexpected public key length for ed25519 key, expected %d, got %d", ed25519.PublicKeySize, n)
	}
	ed25519Key := ed25519.PublicKey(s.PublicKey)
	if _, err := x509.MarshalPKIXPublicKey(ed25519Key); err != nil {
		return nil, fmt.Errorf("marshalling to PKIX key: invalid public key")
	}
	return ed25519Key, nil
}

func (k *Key) toPublicKeyECDSA() (crypto.PublicKey, error) {
	// Prepare decoder limited to 512Kb
	dec := json.NewDecoder(io.LimitReader(bytes.NewReader(k.Value), MaxJSONKeySize))
	s := &helperRSAECDSA{}
	// Unmarshal key value
	if err := dec.Decode(s); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, fmt.Errorf("the public key is truncated or too large: %w", err)
		}
		return nil, err
	}
	ecdsaKey, ok := s.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key")
	}

	if _, err := x509.MarshalPKIXPublicKey(ecdsaKey); err != nil {
		return nil, fmt.Errorf("marshalling to PKIX key: invalid public key")
	}
	return ecdsaKey, nil
}

func (k *Key) toPublicKeyRSA() (crypto.PublicKey, error) {
	// Prepare decoder limited to 512Kb
	dec := json.NewDecoder(io.LimitReader(bytes.NewReader(k.Value), MaxJSONKeySize))
	s := &helperRSAECDSA{}
	// Unmarshal key value
	if err := dec.Decode(s); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, fmt.Errorf("the public key is truncated or too large: %w", err)
		}
		return nil, err
	}
	rsaKey, ok := s.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key")
	}

	if _, err := x509.MarshalPKIXPublicKey(rsaKey); err != nil {
		return nil, fmt.Errorf("marshalling to PKIX key: invalid public key")
	}
	return rsaKey, nil
}

// ID returns the keyID value for the given Key
func (k *Key) ID() string {
	k.idOnce.Do(func() {
		data, err := cjson.EncodeCanonical(k)
		if err != nil {
			panic(fmt.Errorf("tuf: error creating key ID: %w", err))
		}
		digest := sha256.Sum256(data)
		k.id = hex.EncodeToString(digest[:])
	})
	return k.id
}
