package main

/*
Package addrmgr implements concurrency safe Bitcoin address manager.

Address Manager Overview

In order maintain the peer-to-peer Bitcoin network, there needs to be a source
of addresses to connect to as nodes come and go.  The Bitcoin protocol provides
import (
	"fmt"

	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/btcutil"
)

// maybeAcceptBlock potentially accepts a block into the block chain and, if
// from the block being located.
//
// For example, assume a block chain with a side chain as depicted below:
// 	genesis -> 1 -> 2 -> ... -> 15 -> 16  -> 17  -> 18
// 	                              \-> 16a -> 17a
//
// The block locator for block 17a would be the hashes of blocks:
// [17a 16a 15 14 13 12 11 10 9 8 7 6 4 genesis]
	// return sequence lock values of -1 indicating that this transaction
	// can be included within a block at any given height or time.
	mTx := tx.MsgTx()
	sequenceLockActive := mTx.Version >= 2 && csvSoftforkActive
	if !sequenceLockActive || IsCoinBase(tx) {
		return sequenceLock, nil
	}
// LockTimeToSequence converts the passed relative locktime to a sequence
// number in accordance to BIP-68.
// See: https://github.com/bitcoin/bips/blob/master/bip-0068.mediawiki
//  * (Compatibility)
func LockTimeToSequence(isSeconds bool, locktime uint32) uint32 {
	// If we're expressing the relative lock time in blocks, then the
	// corresponding sequence number is simply the desired input age.
// a reorganization to become the main chain).
//
// The flags modify the behavior of this function as follows:
//  - BFFastAdd: Avoids several expensive transaction validation operations.
//    This is useful when using checkpoints.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) connectBestChain(node *blockNode, block *btcutil.Block, flags BehaviorFlags) (bool, error) {
// isCurrent returns whether or not the chain believes it is current.  Several
// factors are used to guess, but the key factors that allow the chain to
// believe it is current are:
//  - Latest block height is after the latest checkpoint (if enabled)
//  - Latest block has a timestamp newer than 24 hours ago
//
// This function MUST be called with the chain state lock held (for reads).
func (b *BlockChain) isCurrent() bool {
// IsCurrent returns whether or not the chain believes it is current.  Several
// factors are used to guess, but the key factors that allow the chain to
// believe it is current are:
//  - Latest block height is after the latest checkpoint (if enabled)
//  - Latest block has a timestamp newer than 24 hours ago
//
// This function is safe for concurrent access.
func (b *BlockChain) IsCurrent() bool {
//
// In addition, there are two special cases:
//
// - When no locators are provided, the stop hash is treated as a request for
//   that block, so it will either return the node associated with the stop hash
//   if it is known, or nil if it is unknown
// - When locators are provided, but none of them are known, nodes starting
//   after the genesis block will be returned
//
// This is primarily a helper function for the locateBlocks and locateHeaders
// functions.
//
// In addition, there are two special cases:
//
// - When no locators are provided, the stop hash is treated as a request for
//   that block, so it will either return the stop hash itself if it is known,
//   or nil if it is unknown
// - When locators are provided, but none of them are known, hashes starting
//   after the genesis block will be returned
//
// This function is safe for concurrent access.
func (b *BlockChain) LocateBlocks(locator BlockLocator, hashStop *chainhash.Hash, maxHashes uint32) []chainhash.Hash {
//
// In addition, there are two special cases:
//
// - When no locators are provided, the stop hash is treated as a request for
//   that header, so it will either return the header for the stop hash itself
//   if it is known, or nil if it is unknown
// - When locators are provided, but none of them are known, headers starting
//   after the genesis block will be returned
//
// This function is safe for concurrent access.
func (b *BlockChain) LocateHeaders(locator BlockLocator, hashStop *chainhash.Hash) []wire.BlockHeader {
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// TestHaveBlock tests the HaveBlock API to ensure proper functionality.
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (
// for comparing chains.
//
// For example, assume a block chain with a side chain as depicted below:
//   genesis -> 1 -> 2 -> 3 -> 4  -> 5 ->  6  -> 7  -> 8
//                         \-> 4a -> 5a -> 6a
//
// The chain view for the branch ending in 6a consists of:
//   genesis -> 1 -> 2 -> 3 -> 4a -> 5a -> 6a
type chainView struct {
	mtx   sync.Mutex
	nodes []*blockNode
// view.
//
// For example, assume a block chain with a side chain as depicted below:
//   genesis -> 1 -> 2 -> 3 -> 4  -> 5 ->  6  -> 7  -> 8
//                         \-> 4a -> 5a -> 6a
//
// Further, assume the view is for the longer chain depicted above.  That is to
// say it consists of:
//   genesis -> 1 -> 2 -> 3 -> 4 -> 5 -> 6 -> 7 -> 8
//
// Invoking this function with block node 5 would return block node 6 while
// invoking it with block node 5a would return nil since that node is not part
// the chain view.  It will return nil if there is no common block.
//
// For example, assume a block chain with a side chain as depicted below:
//   genesis -> 1 -> 2 -> ... -> 5 -> 6  -> 7  -> 8
//                                \-> 6a -> 7a
//
// Further, assume the view is for the longer chain depicted above.  That is to
// say it consists of:
//   genesis -> 1 -> 2 -> ... -> 5 -> 6 -> 7 -> 8.
//
// Invoking this function with block node 7a would return block node 5 while
// invoking it with block node 7 would return itself since it is already part of
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/btcutil"
)

// CheckpointConfirmations is the number of blocks before the end of the current
// checkpoint candidate.
//
// The factors used to determine a good checkpoint are:
//  - The block must be in the main chain
//  - The block must be at least 'CheckpointConfirmations' blocks prior to the
//    current end of the main chain
//  - The timestamps for the blocks before and after the checkpoint must have
//    timestamps which are also before and after the checkpoint, respectively
//    (due to the median time allowance this is not always the case)
//  - The block must not contain any strange transaction such as those with
//    nonstandard scripts
//
// The intent is that candidates are reviewed by a developer to make the final
// decision and then manually added to the list of checkpoints for a network.
// Like IEEE754 floating point, there are three basic components: the sign,
// the exponent, and the mantissa.  They are broken out as follows:
//
//	* the most significant 8 bits represent the unsigned base 256 exponent
// 	* bit 23 (the 24th bit) represents the sign bit
//	* the least significant 23 bits represent the mantissa
//
//	-------------------------------------------------
//	|   Exponent     |    Sign    |    Mantissa     |
//	-------------------------------------------------
//
// The formula to calculate N is:
// 	N = (-1^sign) * mantissa * 256^(exponent-3)
//
// This compact form is only used in bitcoin to encode unsigned 256-bit numbers
// which represent difficulty targets, thus there really is not a need for a
such as orphan blocks which need their parents requested and newly connected
main chain blocks which might result in wallet updates.

Bitcoin Chain Processing Overview

Before a block is allowed into the block chain, it must go through an intensive
series of validation rules.  The following list serves as a general outline of
those rules to provide some intuition into what is going on under the hood, but
is by no means exhaustive:

 - Reject duplicate blocks
 - Perform a series of sanity checks on the block and its transactions such as
   verifying proof of work, timestamps, number and character of transactions,
   transaction amounts, script complexity, and merkle root calculations
 - Compare the block against predetermined checkpoints for expected timestamps
   and difficulty based on elapsed time since the checkpoint
 - Save the most recent orphan blocks for a limited time in case their parent
   blocks become available
 - Stop processing if the block is an orphan as the rest of the processing
   depends on the block's position within the block chain
 - Perform a series of more thorough checks that depend on the block's position
   within the block chain such as verifying block difficulties adhere to
   difficulty retarget rules, timestamps are after the median of the last
   several blocks, all transactions are finalized, checkpoint blocks match, and
   block versions are in line with the previous blocks
 - Determine how the block fits into the chain and perform different actions
   accordingly in order to ensure any side chains which have higher difficulty
   than the main chain become the new main chain
 - When a block is being connected to the main chain (either through
   reorganization of a side chain to the main chain or just extending the
   main chain), perform further checks on the block's transactions such as
   verifying transaction duplicates, script complexity for the combination of
   connected scripts, coinbase maturity, double spends, and connected
   transaction values
 - Run the transaction scripts to verify the spender is allowed to spend the
   coins
 - Insert the block into the block database

Errors

Errors returned by this package are either the raw errors provided by underlying
calls or of type blockchain.RuleError.  This allows the caller to differentiate
determine the specific rule violation by examining the ErrorCode field of the
type asserted blockchain.RuleError.

Bitcoin Improvement Proposals

This package includes spec changes outlined by the following BIPs:

		BIP0016 (https://en.bitcoin.it/wiki/BIP_0016)
		BIP0030 (https://en.bitcoin.it/wiki/BIP_0030)
		BIP0034 (https://en.bitcoin.it/wiki/BIP_0034)
*/
package blockchain
	"path/filepath"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/database"
	_ "github.com/btcsuite/btcd/database/ffldb"
	"github.com/btcsuite/btcd/btcutil"
)

// This example demonstrates how to create a new chain instance and use

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/blockchain/fullblocktests"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	_ "github.com/btcsuite/btcd/database/ffldb"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (
// - A coinbase that pays the required subsidy to an OP_TRUE script
// - When a spendable output is provided:
//   - A transaction that spends from the provided output the following outputs:
//     - One that pays the inputs amount minus 1 atom to an OP_TRUE script
//     - One that contains an OP_RETURN output with a random uint64 in order to
//       ensure the transaction has a unique hash
//
// Additionally, if one or more munge functions are specified, they will be
// invoked with the block prior to solving it.  This provides callers with the
	"sync"
	"time"

	"github.com/btcsuite/btclog"
	"github.com/btcsuite/btcd/btcutil"
)

// blockProgressLogger provides periodic logging for other services in order

// newBlockProgressLogger returns a new block progress logger.
// The progress message is templated as follows:
//  {progressAction} {numProcessed} {blocks|block} in the last {timePeriod}
//  ({numTxs}, height {lastBlockHeight}, {lastBlockTimeStamp})
func newBlockProgressLogger(progressMessage string, logger btclog.Logger) *blockProgressLogger {
	return &blockProgressLogger{
		lastBlockLogTime: time.Now(),
	"errors"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/gcs"
	"github.com/btcsuite/btcd/btcutil/gcs/builder"
)

const (
	"errors"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/btcutil"
)

var (
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

var (
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (
	"fmt"
	"math"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/btcutil"
)

const (
//
// The above stored as a linear array is as follows:
//
// 	[h1 h2 h3 h4 h12 h34 root]
//
// As the above shows, the merkle root is always the last element in the array.
//
// Notification defines notification that is sent to the caller via the callback
// function provided during the call to New and consists of a notification type
// as well as associated data that depends on the type as follows:
// 	- NTBlockAccepted:     *btcutil.Block
// 	- NTBlockConnected:    *btcutil.Block
// 	- NTBlockDisconnected: *btcutil.Block
type Notification struct {
	Type NotificationType
	Data interface{}
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/btcutil"
)

// BehaviorFlags is a bitmask defining tweaks to the normal behavior when
//
// The legacy format is as follows:
//
//   <version><height><header code><unspentness bitmap>[<compressed txouts>,...]
//
//   Field                Type     Size
//   version              VLQ      variable
//   block height         VLQ      variable
//   header code          VLQ      variable
//   unspentness bitmap   []byte   variable
//   compressed txouts
//     compressed amount  VLQ      variable
//     compressed script  []byte   variable
//
// The serialized header code format is:
//   bit 0 - containing transaction is a coinbase
//   bit 1 - output zero is unspent
//   bit 2 - output one is unspent
//   bits 3-x - number of bytes in unspentness bitmap.  When both bits 1 and 2
//     are unset, it encodes N-1 since there must be at least one unspent
//     output.
//
// The rationale for the header code scheme is as follows:
//   - Transactions which only pay to a single output and a change output are
// From tx in main blockchain:
// Blk 1, 0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098
//
//    010103320496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52
//    <><><><------------------------------------------------------------------>
//     | | \--------\                               |
//     | height     |                      compressed txout 0
//  version    header code
//
//  - version: 1
//  - height: 1
//  - header code: 0x03 (coinbase, output zero unspent, 0 bytes of unspentness)
//  - unspentness: Nothing since it is zero bytes
//  - compressed txout 0:
//    - 0x32: VLQ-encoded compressed amount for 5000000000 (50 BTC)
//    - 0x04: special script type pay-to-pubkey
//    - 0x96...52: x-coordinate of the pubkey
//
// Example 2:
// From tx in main blockchain:
// Blk 113931, 4a16969aa4764dd7507fc1de7f0baa4850a246de90c45e59a3207f9a26b5036f
//
//    0185f90b0a011200e2ccd6ec7c6e2e581349c77e067385fa8236bf8a800900b8025be1b3efc63b0ad48e7f9f10e87544528d58
//    <><----><><><------------------------------------------><-------------------------------------------->
//     |    |  | \-------------------\            |                            |
//  version |  \--------\       unspentness       |                    compressed txout 2
//        height     header code          compressed txout 0
//
//  - version: 1
//  - height: 113931
//  - header code: 0x0a (output zero unspent, 1 byte in unspentness bitmap)
//  - unspentness: [0x01] (bit 0 is set, so output 0+2 = 2 is unspent)
//    NOTE: It's +2 since the first two outputs are encoded in the header code
//  - compressed txout 0:
//    - 0x12: VLQ-encoded compressed amount for 20000000 (0.2 BTC)
//    - 0x00: special script type pay-to-pubkey-hash
//    - 0xe2...8a: pubkey hash
//  - compressed txout 2:
//    - 0x8009: VLQ-encoded compressed amount for 15000000 (0.15 BTC)
//    - 0x00: special script type pay-to-pubkey-hash
//    - 0xb8...58: pubkey hash
//
// Example 3:
// From tx in main blockchain:
// Blk 338156, 1b02d1c8cfef60a189017b9a420c682cf4a0028175f2f563209e4ff61c8c3620
//
//    0193d06c100000108ba5b9e763011dd46a006572d820e448e12d2bbb38640bc718e6
//    <><----><><----><-------------------------------------------------->
//     |    |  |   \-----------------\            |
//  version |  \--------\       unspentness       |
//        height     header code          compressed txout 22
//
//  - version: 1
//  - height: 338156
//  - header code: 0x10 (2+1 = 3 bytes in unspentness bitmap)
//    NOTE: It's +1 since neither bit 1 nor 2 are set, so N-1 is encoded.
//  - unspentness: [0x00 0x00 0x10] (bit 20 is set, so output 20+2 = 22 is unspent)
//    NOTE: It's +2 since the first two outputs are encoded in the header code
//  - compressed txout 22:
//    - 0x8ba5b9e763: VLQ-encoded compressed amount for 366875659 (3.66875659 BTC)
//    - 0x01: special script type pay-to-script-hash
//    - 0x1d...e6: script hash
func deserializeUtxoEntryV0(serialized []byte) (map[uint32]*UtxoEntry, error) {
	// Deserialize the version.
	//
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// TestSequenceLocksActive tests the SequenceLockActive function to ensure it
	x, y string
}

//TODO: add more test vectors
var s256BaseMultTests = []baseMultTest{
	{
		"AA5E28D6A97A2479A65527F7290311A3624D4CC0FA1578598EE3C2613BF99522",
	},
}

//TODO: test different curves as well?
func TestBaseMult(t *testing.T) {
	s256 := S256()
	for i, e := range s256BaseMultTests {
// arithmetic over the secp256k1 group order. This means all arithmetic is
// performed modulo:
//
//   0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141
//
// It only implements the arithmetic needed for elliptic curve operations,
// however, the operations that are not implemented can typically be worked
/*
Package btcjson provides primitives for working with the bitcoin JSON-RPC API.

Overview

When communicating via the JSON-RPC protocol, all of the commands need to be
marshalled to and from the the wire in the appropriate format.  This package
In addition, it also provides some additional features such as custom command
registration, command categorization, and reflection-based help generation.

JSON-RPC Protocol Overview

This information is not necessary in order to use this package, but it does
provide some intuition into what the marshalling and unmarshalling that is
string in the result field with a null error for certain commands.  However,
for the most part, the error field will be set as described on failure.

Marshalling and Unmarshalling

Based upon the discussion above, it should be easy to see how the types of this
package map into the required parts of the protocol

  - Request Objects (type Request)
    - Commands (type <Foo>Cmd)
    - Notifications (type <Foo>Ntfn)
  - Response Objects (type Response)
    - Result (type <Foo>Result)

To simplify the marshalling of the requests and responses, the MarshalCmd and
MarshalResponse functions are provided.  They return the raw bytes ready to be
sent across the wire.

Unmarshalling a received Request object is a two step process:
  1) Unmarshal the raw bytes into a Request struct instance via json.Unmarshal
  2) Use UnmarshalCmd on the Result field of the unmarshalled Request to create
     a concrete command or notification instance with all struct fields set
     accordingly

This approach is used since it provides the caller with access to the additional
fields in the request that are not part of the command such as the ID.

Unmarshalling a received Response object is also a two step process:
  1) Unmarhsal the raw bytes into a Response struct instance via json.Unmarshal
  2) Depending on the ID, unmarshal the Result field of the unmarshalled
     Response to create a concrete type instance

As above, this approach is used since it provides the caller with access to the
fields in the response such as the ID and Error.

Command Creation

This package provides two approaches for creating a new command.  This first,
and preferred, method is to use one of the New<Foo>Cmd functions.  This allows
actually executed.  However, it is quite useful for user-supplied commands
that are intentionally dynamic.

Custom Command Registration

The command handling of this package is built around the concept of registered
commands.  This is true for the wide variety of commands already provided by the
A list of all registered methods can be obtained with the RegisteredCmdMethods
function.

Command Inspection

All registered commands are registered with flags that identify information such
as whether the command applies to a chain server, wallet server, or is a
with the MethodUsageFlags flags, and the method can be obtained with the
CmdMethod function.

Help Generation

To facilitate providing consistent help to users of the RPC server, this package
exposes the GenerateHelp and function which uses reflection on registered
In addition, the MethodUsageText function is provided to generate consistent
one-line usage for registered commands and notifications using reflection.

Errors

There are 2 distinct type of errors supported by this package:

// an error will use the key in place of the description.
//
// The following outlines the required keys:
//   "<method>--synopsis"             Synopsis for the command
//   "<method>-<lowerfieldname>"      Description for each command argument
//   "<typename>-<lowerfieldname>"    Description for each object field
//   "<method>--condition<#>"         Description for each result condition
//   "<method>--result<#>"            Description for each primitive result num
//
// Notice that the "special" keys synopsis, condition<#>, and result<#> are
// preceded by a double dash to ensure they don't conflict with field names.
// For example, consider the 'help' command itself.  There are two possible
// returns depending on the provided parameters.  So, the help would be
// generated by calling the function as follows:
//   GenerateHelp("help", descs, (*string)(nil), (*string)(nil)).
//
// The following keys would then be required in the provided descriptions map:
//
//   "help--synopsis":   "Returns a list of all commands or help for ...."
//   "help-command":     "The command to retrieve help for",
//   "help--condition0": "no command provided"
//   "help--condition1": "command specified"
//   "help--result0":    "List of commands"
//   "help--result1":    "Help for specified command"
func GenerateHelp(method string, descs map[string]string, resultTypes ...interface{}) (string, error) {
	// Look up details about the provided method and error out if not
	// registered.
//
// Descriptors are typically ranged when specified in the form of generic HD
// chain paths.
//   Example of a ranged descriptor: pkh(tpub.../*)
//
// The value can be an int to specify the end of the range, or the range
// itself, as []int{begin, end}.
// Reference: https://bitcoincore.org/en/doc/0.20.0/rpc/wallet/getaddressinfo
//
// The GetAddressInfoResult has three segments:
//   1. General information about the address.
//   2. Metadata (Timestamp, HDKeyPath, HDSeedID) and wallet fields
//      (IsMine, IsWatchOnly).
//   3. Information about the embedded address in case of P2SH or P2WSH.
//      Same structure as (1).
type GetAddressInfoResult struct {
	embeddedAddressInfo
	IsMine      bool                 `json:"ismine"`
// (%LOCALAPPDATA%) that is used by default.
//
// Example results:
//  dir := AppDataDir("myapp", false)
//   POSIX (Linux/BSD): ~/.myapp
//   Mac OS: $HOME/Library/Application Support/Myapp
//   Windows: %LOCALAPPDATA%\Myapp
//   Plan 9: $home/myapp
func AppDataDir(appName string, roaming bool) string {
	return appDataDir(runtime.GOOS, appName, roaming)
}
Package base58 provides an API for working with modified base58 and Base58Check
encodings.

Modified Base58 Encoding

Standard base58 encoding is similar to standard base64 encoding except, as the
name implies, it uses a 58 character alphabet which results in an alphanumeric
0, O, I, and l characters that look the same in many fonts and are therefore
hard to humans to distinguish.

Base58Check Encoding Scheme

The Base58Check encoding scheme is primarily used for Bitcoin addresses at the
time of this writing, however it can be used to generically encode arbitrary
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//+build ignore

package main

	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/davecgh/go-spew/spew"
)

	"math/rand"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil/bloom"
)

// This example demonstrates how to create a new bloom filter, add a transaction
	"math"
	"sync"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// ln2Squared is simply the square of the natural log of 2.
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bloom"
)

// TestFilterLarge ensures a maximum sized filter can be created.

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// merkleBlock is used to house intermediate information needed to generate a
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bloom"
)

func TestMerkleBlock3(t *testing.T) {
	"errors"
	"sort"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// Coin represents a spendable transaction outpoint
}

// TotalValueAge returns the total value * number of confirmations
//  of the coins in the set.
func (cs *CoinSet) TotalValueAge() (valueAge int64) {
	return cs.totalValueAge
}
// input priority over the threshold, but no guarantees will be made as to
// minimality of the selection.  The selection below is almost certainly
// suboptimal.
//
type MinPriorityCoinSelector struct {
	MaxInputs              int
	MinChangeAmount        btcutil.Amount
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/coinset"
)

type TestCoin struct {
/*
Package btcutil provides bitcoin-specific convenience functions and types.

Block Overview

A Block defines a bitcoin block that provides easier and more efficient
manipulation of raw wire protocol blocks.  It also memoizes hashes for the
block and its transactions on their first access so subsequent accesses don't
have to repeat the relatively expensive hashing operations.

Tx Overview

A Tx defines a bitcoin transaction that provides more efficient manipulation of
raw wire protocol transactions.  It memoizes the hash for the transaction on its
first access so subsequent accesses don't have to repeat the relatively
expensive hashing operations.

Address Overview

The Address interface provides an abstraction for a Bitcoin address.  While the
most common type is a pay-to-pubkey-hash, Bitcoin already supports others and
	"fmt"
	"math"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil/gcs"
)

const (
/*
Package gcs provides an API for building and using a Golomb-coded set filter.

Golomb-Coded Set

A Golomb-coded set is a probabilistic data structure used similarly to a Bloom
filter. A filter uses constant-size overhead plus on average n+2 bits per
item added to the filter, where 2^-n is the desired false positive (collision)
probability.

GCS use in Bitcoin

GCS filters are a proposed mechanism for storing and transmitting per-block
filters in Bitcoin. The usage is intended to be the inverse of Bloom filters:
// described in:
// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
//
//  * v * N  >> log_2(N)
//
// In our case, using 64-bit integers, log_2 is 64. As most processors don't
// support 128-bit arithmetic natively, we'll be super portable and unfold the
Package hdkeychain provides an API for bitcoin hierarchical deterministic
extended keys (BIP0032).

Overview

The ability to implement hierarchical deterministic wallets depends on the
ability to create and derive hierarchical deterministic extended keys.
functions.  Each extended key can either be a private or public extended key
which itself is capable of deriving a child extended key.

Determining the Extended Key Type

Whether an extended key is a private or public extended key can be determined
with the IsPrivate function.

Transaction Signing Keys and Payment Addresses

In order to create and sign transactions, or provide others with addresses to
send funds to, the underlying key and address material must be accessible.  This
package provides the ECPubKey, ECPrivKey, and Address functions for this
purpose.

The Master Node

As previously mentioned, the extended keys are hierarchical meaning they are
used to form a tree.  The root of that tree is called the master node and this
random seed.  The GenerateSeed function is provided as a convenient way to
create a random seed for use with the NewMaster function.

Deriving Children

Once you have created a tree root (or have deserialized an extended key as
discussed later), the child extended keys can be derived by using the Derive
Derive function.  This provides the ability to cascade the keys into a tree and
hence generate the hierarchical deterministic key chains.

Normal vs Hardened Derived Extended Keys

A private extended key can be used to derive both hardened and non-hardened
(normal) child private and public extended keys.  A public extended key can only
account level in the tree. This way, a leak of an account-specific (or below)
private key never risks compromising the master or other accounts."

Neutering a Private Extended Key

A private extended key can be converted to a new instance of the corresponding
public extended key with the Neuter function.  The original extended key is not
modified.  A public extended key is still capable of deriving non-hardened child
public extended keys.

Serializing and Deserializing Extended Keys

Extended keys are serialized and deserialized with the String and
NewKeyFromString functions.  The serialized key is a Base58-encoded string which
looks like the following:
	public key:   xpub68Gmy5EdvgibQVfPdqkBBCHxA5htiqg55crXYuXoQRKfDBFA1WEjWgP6LHhwBZeNK1VTsfTFUHCdrfp1bgwQ9xv5ski8PX9rL2dZXvgGDnw
	private key:  xprv9uHRZZhk6KAJC1avXpDAp4MDc3sQKNxDiPvvkX8Br5ngLNv1TxvUxt4cV1rGL5hj6KCesnDYUhd7oWgT11eZG7XnxHrnYeSvkzY7d2bhkJ7

Network

Extended keys are much like normal Bitcoin addresses in that they have version
bytes which tie them to a specific network.  The SetNet and IsForNet functions
import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
)

// This example demonstrates how to generate a cryptographically random seed
// on the SLIP132 standard (serializable to yprv/ypub, zprv/zpub, etc.).
//
// References:
//   [SLIP132]: SLIP-0132 - Registered HD version bytes for BIP-0032
//   https://github.com/satoshilabs/slips/blob/master/slip-0132.md
func (k *ExtendedKey) CloneWithVersion(version []byte) (*ExtendedKey, error) {
	if len(version) != 4 {
		// TODO: The semantically correct error to return here is
// extended keys.
//
// The following tool was used for generating the tests:
//   https://jlopp.github.io/xpub-converter
func TestCloneWithVersion(t *testing.T) {
	tests := []struct {
		name    string
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// +build !appengine

package btcutil
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// +build appengine

package btcutil
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/davecgh/go-spew/spew"
)

/*
Package txsort provides the transaction sorting according to BIP 69.

Overview

BIP 69 defines a standard lexicographical sort order of transaction inputs and
outputs.  This is useful to standardize transactions for faster multi-party
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil/txsort"
)

// TestSort ensures the transaction sorting works according to the BIP.
// The WIF string must be a base58-encoded string of the following byte
// sequence:
//
//  * 1 byte to identify the network, must be 0x80 for mainnet or 0xef for
//    either testnet3 or the regression test network
//  * 32 bytes of a binary-encoded, big-endian, zero-padded private key
//  * Optional 1 byte (equal to 0x01) if the address being imported or exported
//    was created by taking the RIPEMD160 after SHA256 hash of a serialized
//    compressed (33-byte) public key
//  * 4 bytes of checksum, must equal the first four bytes of the double SHA256
//    of every byte before the checksum in this sequence
//
// If the base58-decoded byte sequence does not match this, DecodeWIF will
// return a non-nil error.  ErrMalformedPrivateKey is returned when the WIF
// When a network parameter is needed, it may then be looked up through this
// variable (either directly, or hidden in a library call).
//
//  package main
//
//  import (
//          "flag"
//          "fmt"
//          "log"
//
//          "github.com/btcsuite/btcd/btcutil"
//          "github.com/btcsuite/btcd/chaincfg"
//  )
//
//  var testnet = flag.Bool("testnet", false, "operate on the testnet Bitcoin network")
//
//  // By default (without -testnet), use mainnet.
//  var chainParams = &chaincfg.MainNetParams
//
//  func main() {
//          flag.Parse()
//
//          // Modify active network parameters if operating on testnet.
//          if *testnet {
//                  chainParams = &chaincfg.TestNet3Params
//          }
//
//          // later...
//
//          // Create and print new payment address, specific to the active network.
//          pubKeyHash := make([]byte, 20)
//          addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, chainParams)
//          if err != nil {
//                  log.Fatal(err)
//          }
//          fmt.Println(addr)
//  }
//
// If an application does not use one of the three standard Bitcoin networks,
// a new Params struct may be created which defines the parameters for the
// ErrInvalidHDKeyID error will be returned.
//
// Reference:
//   SLIP-0132 : Registered HD version bytes for BIP-0032
//   https://github.com/satoshilabs/slips/blob/master/slip-0132.md
func RegisterHDKeyID(hdPublicKeyID []byte, hdPrivateKeyID []byte) error {
	if len(hdPublicKeyID) != 4 || len(hdPrivateKeyID) != 4 {
		return ErrInvalidHDKeyID
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/database"
	_ "github.com/btcsuite/btcd/database/ffldb"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	flags "github.com/jessevdk/go-flags"
)


	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/blockchain/indexers"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

var zeroHash = chainhash.Hash{}
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/btcutil"
	flags "github.com/jessevdk/go-flags"
)

// line options.
//
// The configuration proceeds as follows:
// 	1) Start with a default config with sane settings
// 	2) Pre-parse the command line to check for an alternative config file
// 	3) Load configuration file overwriting defaults with any specified options
// 	4) Parse CLI options and overwrite/add any specified options
//
// The above results in functioning properly without any config settings
// while still allowing the user to override settings with config files and
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/database"
	_ "github.com/btcsuite/btcd/database/ffldb"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	flags "github.com/jessevdk/go-flags"
)

	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/connmgr"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/go-socks/socks"
	flags "github.com/jessevdk/go-flags"
)
// line options.
//
// The configuration proceeds as follows:
// 	1) Start with a default config with sane settings
// 	2) Pre-parse the command line to check for an alternative config file
// 	3) Load configuration file overwriting defaults with any specified options
// 	4) Parse CLI options and overwrite/add any specified options
//
// The above results in btcd functioning properly without any config settings
// while still allowing the user to override settings with config files and
/*
Package connmgr implements a generic Bitcoin network connection manager.

Connection Manager Overview

Connection Manager handles all the general connection concerns such as
maintaining a set number of outbound connections, sourcing peers, banning,
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/database"
	_ "github.com/btcsuite/btcd/database/ffldb"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

var (
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// importCmd defines the configuration options for the insecureimport command.
/*
Package database provides a block and metadata storage database.

Overview

As of Feb 2016, there are over 400,000 blocks in the Bitcoin block chain and
and over 112 million transactions (which turns out to be over 60GB of data).

A quick overview of the features database provides are as follows:

 - Key/value metadata store
 - Bitcoin block storage
 - Efficient retrieval of block headers and regions (transactions, scripts, etc)
 - Read-only and read-write transactions with both manual and managed modes
 - Nested buckets
 - Supports registration of backend databases
 - Comprehensive test coverage

Database

The main entry point is the DB interface.  It exposes functionality for
transactional-based access and storage of metadata and block data.  It is
functions provide a managed transaction.  These are described in more detail
below.

Transactions

The Tx interface provides facilities for rolling back or committing changes that
took place while the transaction was active.  It also provides the root metadata
bucket under which all keys, values, and nested buckets are stored.  A
transaction can either be read-only or read-write and managed or unmanaged.

Managed versus Unmanaged Transactions

A managed transaction is one where the caller provides a function to execute
within the context of the transaction and the commit or rollback is handled
open for long periods of time can have several adverse effects, so it is
recommended that managed transactions are used instead.

Buckets

The Bucket interface provides the ability to manipulate key/value pairs and
nested buckets as well as iterate through them.
buckets.  The ForEach function allows the caller to provide a function to be
called with each key/value pair and nested bucket in the current bucket.

Metadata Bucket

As discussed above, all of the functions which are used to manipulate key/value
pairs and nested buckets exist on the Bucket interface.  The root metadata
same time as the database.  Use the Metadata function on the Tx interface
to retrieve it.

Nested Buckets

The CreateBucket and CreateBucketIfNotExists functions on the Bucket interface
provide the ability to create an arbitrary number of nested buckets.  It is
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/btcutil"
)

// BenchmarkBlockHeader benchmarks how long it takes to load the mainnet genesis
// were partially written.
//
// There are effectively two scenarios to consider here:
//   1) Transient write failures from which recovery is possible
//   2) More permanent failures such as hard disk death and/or removal
//
// In either case, the write cursor will be repositioned to the old block file
// offset regardless of any other errors that occur while attempting to undo
	"sort"
	"sync"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/database/internal/treap"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/comparer"
	ldberrors "github.com/syndtr/goleveldb/leveldb/errors"
for the metadata, flat files for block storage, and checksums in key areas to
ensure data integrity.

Usage

This package is a driver to the database package and provides the database type
of "ffldb".  The parameters the Open and Create functions take are the
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/database/ffldb"
	"github.com/btcsuite/btcd/btcutil"
)

// dbType is the database type name for this driver.
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

var (
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/syndtr/goleveldb/leveldb"
	ldberrors "github.com/syndtr/goleveldb/leveldb/errors"
)
package database

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/btcutil"
)

// Cursor represents a cursor over key/value pairs and nested buckets of a
// unexpected keys and/or values.
//
// For example:
//   iter := t.Iterator(nil, nil)
//   for iter.Next() {
//   	if someCondition {
//   		t.Delete(iter.Key())
//   		iter.ForceReseek()
//   	}
//   }
func (t *Mutable) Iterator(startKey, limitKey []byte) *Iterator {
	iter := &Iterator{
		t:        t,
this location.

Usage:
  btcd [OPTIONS]

Application Options:
      --addcheckpoint=        Add a custom checkpoint.  Format:
                              '<height>:<hash>'
  -a, --addpeer=              Add a peer to connect with at startup
      --addrindex             Maintain a full address-based transaction index
                              which makes the searchrawtransactions RPC
                              available
      --banduration=          How long to ban misbehaving peers.  Valid time
                              units are {s, m, h}.  Minimum 1 second (default:
                              24h0m0s)
      --banthreshold=         Maximum allowed ban score before disconnecting
                              and banning misbehaving peers. (default: 100)
      --blockmaxsize=         Maximum block size in bytes to be used when
                              creating a block (default: 750000)
      --blockminsize=         Mininum block size in bytes to be used when
                              creating a block
      --blockmaxweight=       Maximum block weight to be used when creating a
                              block (default: 3000000)
      --blockminweight=       Mininum block weight to be used when creating a
                              block
      --blockprioritysize=    Size in bytes for high-priority/low-fee
                              transactions when creating a block (default:
                              50000)
      --blocksonly            Do not accept transactions from remote peers.
  -C, --configfile=           Path to configuration file
      --connect=              Connect only to the specified peers at startup
      --cpuprofile=           Write CPU profile to the specified file
  -b, --datadir=              Directory to store data
      --dbtype=               Database backend to use for the Block Chain
                              (default: ffldb)
  -d, --debuglevel=           Logging level for all subsystems {trace, debug,
                              info, warn, error, critical} -- You may also
                              specify
                              <subsystem>=<level>,<subsystem2>=<level>,... to
                              set the log level for individual subsystems --
                              Use show to list available subsystems (default:
                              info)
      --dropaddrindex         Deletes the address-based transaction index from
                              the database on start up and then exits.
      --dropcfindex           Deletes the index used for committed filtering
                              (CF) support from the database on start up and
                              then exits.
      --droptxindex           Deletes the hash-based transaction index from the
                              database on start up and then exits.
      --externalip=           Add an ip to the list of local addresses we claim
                              to listen on to peers
      --generate              Generate (mine) bitcoins using the CPU
      --limitfreerelay=       Limit relay of transactions with no transaction
                              fee to the given amount in thousands of bytes per
                              minute (default: 15)
      --listen=               Add an interface/port to listen for connections
                              (default all interfaces port: 8333, testnet:
                              18333, signet: 38333)
      --logdir=               Directory to log output
      --maxorphantx=          Max number of orphan transactions to keep in
                              memory (default: 100)
      --maxpeers=             Max number of inbound and outbound peers
                              (default: 125)
      --miningaddr=           Add the specified payment address to the list of
                              addresses to use for generated blocks -- At least
                              one address is required if the generate option is
                              set
      --minrelaytxfee=        The minimum transaction fee in BTC/kB to be
                              considered a non-zero fee. (default: 1e-05)
      --nobanning             Disable banning of misbehaving peers
      --nocfilters            Disable committed filtering (CF) support
      --nocheckpoints         Disable built-in checkpoints.  Don't do this
                              unless you know what you're doing.
      --nodnsseed             Disable DNS seeding for peers
      --nolisten              Disable listening for incoming connections --
                              NOTE: Listening is automatically disabled if the
                              --connect or --proxy options are used without
                              also specifying listen interfaces via --listen
      --noonion               Disable connecting to tor hidden services
      --nopeerbloomfilters    Disable bloom filtering support
      --norelaypriority       Do not require free or low-fee transactions to
                              have high priority for relaying
      --norpc                 Disable built-in RPC server -- NOTE: The RPC
                              server is disabled by default if no
                              rpcuser/rpcpass or rpclimituser/rpclimitpass is
                              specified
      --notls                 Disable TLS for the RPC server -- NOTE: This is
                              only allowed if the RPC server is bound to
                              localhost
      --onion=                Connect to tor hidden services via SOCKS5 proxy
                              (eg. 127.0.0.1:9050)
      --onionpass=            Password for onion proxy server
      --onionuser=            Username for onion proxy server
      --profile=              Enable HTTP profiling on given port -- NOTE port
                              must be between 1024 and 65536
      --proxy=                Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)
      --proxypass=            Password for proxy server
      --proxyuser=            Username for proxy server
      --regtest               Use the regression test network
      --rejectnonstd          Reject non-standard transactions regardless of
                              the default settings for the active network.
      --relaynonstd           Relay non-standard transactions regardless of the
                              default settings for the active network.
      --rpccert=              File containing the certificate file
      --rpckey=               File containing the certificate key
      --rpclimitpass=         Password for limited RPC connections
      --rpclimituser=         Username for limited RPC connections
      --rpclisten=            Add an interface/port to listen for RPC
                              connections (default port: 8334, testnet: 18334)
      --rpcmaxclients=        Max number of RPC clients for standard
                              connections (default: 10)
      --rpcmaxconcurrentreqs= Max number of concurrent RPC requests that may be
                              processed concurrently (default: 20)
      --rpcmaxwebsockets=     Max number of RPC websocket connections (default:
                              25)
      --rpcquirks             Mirror some JSON-RPC quirks of Bitcoin Core --
                              NOTE: Discouraged unless interoperability issues
                              need to be worked around
  -P, --rpcpass=              Password for RPC connections
  -u, --rpcuser=              Username for RPC connections
      --sigcachemaxsize=      The maximum number of entries in the signature
                              verification cache (default: 100000)
      --simnet                Use the simulation test network
      --testnet               Use the test network
      --torisolation          Enable Tor stream isolation by randomizing user
                              credentials for each connection.
      --trickleinterval=      Minimum time between attempts to send new
                              inventory to a connected peer (default: 10s)
      --txindex               Maintain a full hash-based transaction index
                              which makes all transactions available via the
                              getrawtransaction RPC
      --uacomment=            Comment to add to the user agent -- See BIP 14
                              for more information.
      --upnp                  Use UPnP to map our listening port outside of NAT
  -V, --version               Display version information and exit
      --whitelist=            Add an IP network or IP that will not be banned.
                              (eg. 192.168.1.0/24 or ::1)

Help Options:
  -h, --help           Show this help message

*/
package main
// - Assert the chain height is 0 and the state is ThresholdDefined
// - Generate 1 fewer blocks than needed to reach the first state transition
//   - Assert chain height is expected and state is still ThresholdDefined
// - Generate 1 more block to reach the first state transition
//   - Assert chain height is expected and state moved to ThresholdStarted
// - Generate enough blocks to reach the next state transition window, but only
//   signal support in 1 fewer than the required number to achieve
//   ThresholdLockedIn
//   - Assert chain height is expected and state is still ThresholdStarted
// - Generate enough blocks to reach the next state transition window with only
//   the exact number of blocks required to achieve locked in status signalling
//   support.
//   - Assert chain height is expected and state moved to ThresholdLockedIn
// - Generate 1 fewer blocks than needed to reach the next state transition
//   - Assert chain height is expected and state is still ThresholdLockedIn
// - Generate 1 more block to reach the next state transition
//   - Assert chain height is expected and state moved to ThresholdActive
func TestBIP0009(t *testing.T) {
	t.Parallel()
// Overview:
// - Generate block 1
//   - Assert bit is NOT set (ThresholdDefined)
// - Generate enough blocks to reach first state transition
//   - Assert bit is NOT set for block prior to state transition
//   - Assert bit is set for block at state transition (ThresholdStarted)
// - Generate enough blocks to reach second state transition
//   - Assert bit is set for block at state transition (ThresholdLockedIn)
// - Generate enough blocks to reach third state transition
//   - Assert bit is set for block prior to state transition (ThresholdLockedIn)
//   - Assert bit is NOT set for block at state transition (ThresholdActive)

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/integration/rpctest"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (
// them.
//
// Overview:
//  - Pre soft-fork:
//    - Transactions with non-final lock-times from the PoV of MTP should be
//      rejected from the mempool.
//    - Transactions within non-final MTP based lock-times should be accepted
//      in valid blocks.
//
//  - Post soft-fork:
//    - Transactions with non-final lock-times from the PoV of MTP should be
//      rejected from the mempool and when found within otherwise valid blocks.
//    - Transactions with final lock-times from the PoV of MTP should be
//      accepted to the mempool and mined in future block.
func TestBIP0113Activation(t *testing.T) {
	t.Parallel()

// 112 and BIP 68 rule-set after the activation of the CSV-package soft-fork.
//
// Overview:
//  - Pre soft-fork:
//    - A transaction spending a CSV output validly should be rejected from the
//    mempool, but accepted in a valid generated block including the
//    transaction.
//  - Post soft-fork:
//    - See the cases exercised within the table driven tests towards the end
//    of this test.
func TestBIP0068AndBIP0112Activation(t *testing.T) {
	t.Parallel()

	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mining"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// solveBlock attempts to find a nonce which makes the passed block header hash
	"runtime"
	"time"

	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/btcutil"
)

// nodeConfig contains all the args, and data required to launch a btcd process
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

func testSendOutputs(r *Harness, t *testing.T) {
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// +build !windows,!plan9

package limits
to other peers on the network and notifying the mining process that new
transactions are available.

Feature Overview

The following is a quick overview of the major features.  It is not intended to
be an exhaustive list.

 - Maintain a pool of fully validated transactions
   - Reject non-fully-spent duplicate transactions
   - Reject coinbase transactions
   - Reject double spends (both from the chain and other transactions in pool)
   - Reject invalid transactions according to the network consensus rules
   - Full script execution and validation with signature cache support
   - Individual transaction query support
 - Orphan transaction support (transactions that spend from unknown outputs)
   - Configurable limits (see transaction acceptance policy)
   - Automatic addition of orphan transactions that are no longer orphans as new
     transactions are added to the pool
   - Individual orphan transaction query support
 - Configurable transaction acceptance policy
   - Option to accept or reject standard transactions
   - Option to accept or reject transactions based on priority calculations
   - Rate limiting of low-fee and free transactions
   - Non-zero fee threshold
   - Max signature operations per transaction
   - Max orphan transaction size
   - Max number of orphan transactions allowed
 - Additional metadata tracking for each transaction
   - Timestamp when the transaction was added to the pool
   - Most recent block height when the transaction was added to the pool
   - The fee the transaction pays
   - The starting priority for the transaction
 - Manual control of transaction removal
   - Recursive removal of all dependent transactions

Errors

Errors returned by this package are either the raw errors provided by underlying
calls or of type mempool.RuleError.  Since there are two classes of rules
	"strings"
	"sync"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mining"
	"github.com/btcsuite/btcd/btcutil"
)

// TODO incorporate Alex Morcos' modifications to Gavin's initial model
	"math/rand"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mining"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// newTestFeeEstimator creates a feeEstimator with some different parameters
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mining"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (
//
// Given the above, a block generated by this function is of the following form:
//
//   -----------------------------------  --  --
//  |      Coinbase Transaction         |   |   |
//  |-----------------------------------|   |   |
//  |                                   |   |   | ----- policy.BlockPrioritySize
//  |   High-priority Transactions      |   |   |
//  |                                   |   |   |
//  |-----------------------------------|   | --
//  |                                   |   |
//  |                                   |   |
//  |                                   |   |--- policy.BlockMaxSize
//  |  Transactions prioritized by fee  |   |
//  |  until <= policy.TxMinFreeFee     |   |
//  |                                   |   |
//  |                                   |   |
//  |                                   |   |
//  |-----------------------------------|   |
//  |  Low-fee/Non high-priority (free) |   |
//  |  transactions (while block size   |   |
//  |  <= policy.BlockMinSize)          |   |
//   -----------------------------------  --
func (g *BlkTmplGenerator) NewBlockTemplate(payToAddress btcutil.Address) (*BlockTemplate, error) {
	// Extend the most recently known best block.
	best := g.chain.BestSnapshot()

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (
	"testing"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// newHashFromStr converts the passed big-endian hex string into a
	"sync"
	"time"

	"github.com/btcsuite/btclog"
	"github.com/btcsuite/btcd/btcutil"
)

// blockProgressLogger provides periodic logging for other services in order

// newBlockProgressLogger returns a new block progress logger.
// The progress message is templated as follows:
//  {progressAction} {numProcessed} {blocks|block} in the last {timePeriod}
//  ({numTxs}, height {lastBlockHeight}, {lastBlockTimeStamp})
func newBlockProgressLogger(progressMessage string, logger btclog.Logger) *blockProgressLogger {
	return &blockProgressLogger{
		lastBlockLogTime: time.Now(),

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// PeerNotifier exposes methods to notify peers of status changes to
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/mempool"
	peerpkg "github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (
Package peer provides a common base for creating and managing Bitcoin network
peers.

Overview

This package builds upon the wire package, which provides the fundamental
primitives necessary to speak the bitcoin wire protocol, in order to simplify

A quick overview of the major features peer provides are as follows:

 - Provides a basic concurrent safe bitcoin peer for handling bitcoin
   communications via the peer-to-peer protocol
 - Full duplex reading and writing of bitcoin protocol messages
 - Automatic handling of the initial handshake process including protocol
   version negotiation
 - Asynchronous message queuing of outbound messages with optional channel for
   notification when the message is actually sent
 - Flexible peer configuration
   - Caller is responsible for creating outgoing connections and listening for
     incoming connections so they have flexibility to establish connections as
     they see fit (proxies, etc)
   - User agent name and version
   - Bitcoin network
   - Service support signalling (full nodes, bloom filters, etc)
   - Maximum supported protocol version
   - Ability to register callbacks for handling bitcoin protocol messages
 - Inventory message batching and send trickling with known inventory detection
   and avoidance
 - Automatic periodic keep-alive pinging and pong responses
 - Random nonce generation and self connection detection
 - Proper handling of bloom filter related commands when the caller does not
   specify the related flag to signal support
   - Disconnects the peer when the protocol version is high enough
   - Does not invoke the related callbacks for older protocol versions
 - Snapshottable peer statistics such as the total number of bytes read and
   written, the remote address, user agent, and negotiated protocol version
 - Helper functions pushing addresses, getblocks, getheaders, and reject
   messages
   - These could all be sent manually via the standard message output function,
     but the helpers provide additional nice functionality such as duplicate
     filtering and address randomization
 - Ability to wait for shutdown/disconnect
 - Comprehensive test coverage

Peer Configuration

All peer configuration is handled with the Config struct.  This allows the
caller to specify things such as the user agent name and version, the bitcoin
messages are received.  See the documentation for each field of the Config
struct for more details.

Inbound and Outbound Peers

A peer can either be inbound or outbound.  The caller is responsible for
establishing the connection to remote peers and listening for incoming peers.
WaitForDisconnect can be used to block until peer disconnection and resource
cleanup has completed.

Callbacks

In order to do anything useful with a peer, it is necessary to react to bitcoin
messages.  This is accomplished by creating an instance of the MessageListeners
callback handlers.  This provides a clean method for accessing that state when
callbacks are invoked.

Queuing Messages and Inventory

The QueueMessage function provides the fundamental means to send messages to the
remote peer.  As the name implies, this employs a non-blocking queue.  A done
intelligent known remote peer inventory detection and avoidance through the use
of a most-recently used algorithm.

Message Sending Helper Functions

In addition to the bare QueueMessage function previously described, the
PushAddrMsg, PushGetBlocksMsg, PushGetHeadersMsg, and PushRejectMsg functions
optionally provides a flag to cause it to block until the message is actually
sent.

Peer Statistics

A snapshot of the current peer statistics can be obtained with the StatsSnapshot
function.  This includes statistics such as the total number of bytes read and
written, the remote address, user agent, and negotiated protocol version.

Logging

This package provides extensive logging capabilities through the UseLogger
function which allows a btclog.Logger to be specified.  For example, logging at
logging at the trace level provides full dumps of parsed messages as well as the
raw message bytes using a format similar to hexdump -C.

Bitcoin Improvement Proposals

This package supports all BIPS supported by the wire package.
(https://pkg.go.dev/github.com/btcsuite/btcd/wire#hdr-Bitcoin_Improvement_Proposals)
// peer. The events should occur in the following order, otherwise an error is
// returned:
//
//   1. Remote peer sends their version.
//   2. We send our version.
//   3. We send sendaddrv2 if their version is >= 70016.
//   4. We send our verack.
//   5. Wait until sendaddrv2 or verack is received. Unknown messages are
//      skipped as it could be wtxidrelay or a different message in the future
//      that btcd does not implement but bitcoind does.
//   6. If remote peer sent sendaddrv2 above, wait until receipt of verack.
func (p *Peer) negotiateInboundProtocol() error {
	if err := p.readRemoteVersionMsg(); err != nil {
		return err
// peer. The events should occur in the following order, otherwise an error is
// returned:
//
//   1. We send our version.
//   2. Remote peer sends their version.
//   3. We send sendaddrv2 if their version is >= 70016.
//   4. We send our verack.
//   5. We wait to receive sendaddrv2 or verack, skipping unknown messages as
//      in the inbound case.
//   6. If sendaddrv2 was received, wait for receipt of verack.
func (p *Peer) negotiateOutboundProtocol() error {
	if err := p.writeLocalVersionMsg(); err != nil {
		return err
	"sync/atomic"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/netsync"
	"github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// rpcPeer provides a peer for use with the RPC server and implements the
/*
Package rpcclient implements a websocket-enabled Bitcoin JSON-RPC client.

Overview

This client provides a robust and easy to use client for interfacing with a
Bitcoin RPC server that uses a btcd/bitcoin core compatible Bitcoin JSON-RPC
fall back to HTTP POST and disable TLS to support talking with inferior bitcoin
core style RPC servers.

Websockets vs HTTP POST

In HTTP POST-based JSON-RPC, every request creates a new HTTP connection,
issues the call, waits for the response, and closes the connection.  This adds
call.  In addition, the websocket interface provides other nice features such as
the ability to register for asynchronous notifications of various events.

Synchronous vs Asynchronous API

The client provides both a synchronous (blocking) and asynchronous API.

immediately if it has already arrived, or block until it has.  This is useful
since it provides the caller with greater control over concurrency.

Notifications

The first important part of notifications is to realize that they will only
work when connected via websockets.  This should intuitively make sense
if you want to be notified when funds are received by a set of addresses, you
register the addresses via the NotifyReceived (or NotifyReceivedAsync) function.

Notification Handlers

Notifications are exposed by the client through the use of callback handlers
which are setup via a NotificationHandlers instance that is specified by the
returns, but the callback would be waiting for a response.   Thus, any
additional RPCs must be issued an a completely decoupled manner.

Automatic Reconnection

By default, when running in websockets mode, this client will automatically
keep trying to reconnect to the RPC server should the connection be lost.  There
need chain-related RPCs.  In contrast, btcwallet provides pass through treatment
for chain-related RPCs, so it supports them in addition to wallet-related RPCs.

Errors

There are 3 categories of errors that will be returned throughout this package:

detected by type asserting the error in a *btcjson.RPCError.  For example, to
detect if a command is unimplemented by the remote RPC server:

  amount, err := client.GetBalance("")
  if err != nil {
  	if jerr, ok := err.(*btcjson.RPCError); ok {
  		switch jerr.Code {
  		case btcjson.ErrRPCUnimplemented:
  			// Handle not implemented error

  		// Handle other specific errors you care about
		}
  	}

  	// Log or otherwise handle the error knowing it was not one returned
  	// from the remote RPC server.
  }

Example Usage

The following full-blown client examples are in the examples directory:

 - bitcoincorehttp
   Connects to a bitcoin core RPC server using HTTP POST mode with TLS disabled
   and gets the current block count
 - btcdwebsockets
   Connects to a btcd RPC server using TLS-secured websockets, registers for
   block connected and block disconnected notifications, and gets the current
   block count
 - btcwalletwebsockets
   Connects to a btcwallet RPC server using TLS-secured websockets, registers
   for notifications about changes to account balances, and gets a list of
   unspent transaction outputs (utxos) the wallet can sign
*/
package rpcclient
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

func main() {
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/davecgh/go-spew/spew"
)

	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

// FutureDebugLevelResult is a future promise to deliver the result of a
// specification.
//
// The levelspec can be either a debug level or of the form:
// 	<subsystem>=<level>,<subsystem2>=<level2>,...
//
// Additionally, the special keyword 'show' can be used to get a list of the
// available subsystems.
	"errors"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/btcutil"
)

// FutureGenerateResult is a future promise to deliver the result of a
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

var (
	"encoding/json"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (
//
// Optional parameters can be specified using functional-options pattern. The
// following functions are available:
//   * WithCreateWalletDisablePrivateKeys
//   * WithCreateWalletBlank
//   * WithCreateWalletPassphrase
//   * WithCreateWalletAvoidReuse
func (c *Client) CreateWallet(name string, opts ...CreateWalletOpt) (*btcjson.CreateWalletResult, error) {
	return c.CreateWalletAsync(name, opts...).Receive()
}

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/websocket"
	"golang.org/x/crypto/ripemd160"
)
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package main
This package provides data structures and functions to parse and execute
bitcoin transaction scripts.

Script Overview

Bitcoin transaction scripts are written in a stack-base, FORTH-like language.

One benefit of using a scripting language is added flexibility in specifying
what conditions must be met in order to spend bitcoins.

Errors

Errors returned by this package are of type txscript.Error.  This allows the
caller to programmatically determine the specific error by examining the

// Error identifies a script-related error.  It is used to indicate three
// classes of errors:
// 1) Script execution failures due to violating one of the many requirements
//    imposed by the script engine or evaluating to false
// 2) Improper API usage by callers
// 3) Internal consistency check failures
//
// The caller can use type assertions on the returned errors to access the
// ErrorCode field to ascertain the specific reason for the error.  As an

	// Transaction version numbers not high enough to trigger CSV rules must
	// fail.
	if vm.tx.Version < 2 {
		str := fmt.Sprintf("invalid transaction version: %d",
			vm.tx.Version)
		return scriptError(ErrUnsatisfiedLockTime, str)
// Bytes returns the number serialized as a little endian with a sign bit.
//
// Example encodings:
//       127 -> [0x7f]
//      -127 -> [0xff]
//       128 -> [0x80 0x00]
//      -128 -> [0x80 0x80]
//       129 -> [0x81 0x00]
//      -129 -> [0x81 0x80]
//       256 -> [0x00 0x01]
//      -256 -> [0x00 0x81]
//     32767 -> [0xff 0x7f]
//    -32767 -> [0xff 0xff]
//     32768 -> [0x00 0x80 0x00]
//    -32768 -> [0x00 0x80 0x80]
func (n scriptNum) Bytes() []byte {
	// Zero encodes as an empty byte slice.
	if n == 0 {
with the specifics of message handling such as what to do when a message is
received.  This provides the caller with a high level of flexibility.

Bitcoin Message Overview

The bitcoin protocol consists of exchanging messages between peers.  Each
message is preceded by a header which identifies information about it such as
wire using bitcoin encoding are handled so the caller doesn't have to concern
themselves with the specifics.

Message Interaction

The following provides a quick summary of how the bitcoin messages are intended
to interact with one another.  As stated above, these interactions are not
	  in BIP0031.  The BIP0031Version constant can be used to detect a recent
	  enough protocol version for this purpose (version > BIP0031Version).

Common Parameters

There are several common parameters that arise when using this package to read
and write bitcoin messages.  The following sections provide a quick overview of
these parameters so the next sections can build on them.

Protocol Version

The protocol version should be negotiated with the remote peer at a higher
level than this package via the version (MsgVersion) message exchange, however,
for all outbound connections before a potentially lower protocol version is
negotiated.

Bitcoin Network

The bitcoin network is a magic number which is used to identify the start of a
message and which bitcoin network the message applies to.  This package provides
	wire.TestNet3 (Test network version 3)
	wire.SimNet   (Simulation test network)

Determining Message Type

As discussed in the bitcoin message overview section, this package reads
and writes bitcoin messages using a generic interface named Message.  In
		fmt.Printf("Number of tx in block: %v", msg.Header.TxnCount)
	}

Reading Messages

In order to unmarshall bitcoin messages from the wire, use the ReadMessage
function.  It accepts any io.Reader, but typically this will be a net.Conn to
		// Log and handle the error
	}

Writing Messages

In order to marshall bitcoin messages to the wire, use the WriteMessage
function.  It accepts any io.Writer, but typically this will be a net.Conn to
		// Log and handle the error
	}

Errors

Errors returned by this package are either the raw errors provided by underlying
calls to read/write from streams such as io.EOF, io.ErrUnexpectedEOF, and
differentiate between general IO errors and malformed messages through type
assertions.

Bitcoin Improvement Proposals

This package includes spec changes outlined by the following BIPs:

