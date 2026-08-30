package main

	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
		defer ensureAdditionalGids(s)
		setProcess(s)
		s.Process.User.AdditionalGids = nil

		// For LCOW it's a bit harder to confirm that the user actually exists on the host as a rootfs isn't
		// mounted on the host and shared into the guest, but rather the rootfs is constructed entirely in the
		switch len(parts) {
		case 1:
			v, err := strconv.Atoi(parts[0])
			if err != nil {
				// if we cannot parse as a uint they try to see if it is a username
				return WithUsername(userstr)(ctx, client, c, s)
			}
			return WithUserID(uint32(v))(ctx, client, c, s)
			)
			var uid, gid uint32
			v, err := strconv.Atoi(parts[0])
			if err != nil {
				username = parts[0]
			} else {
				uid = uint32(v)
			}
			if v, err = strconv.Atoi(parts[1]); err != nil {
				groupname = parts[1]
			} else {
				gid = uint32(v)
	"golang.org/x/sys/unix"
)

//nolint:gosec
func TestWithUserID(t *testing.T) {
	t.Parallel()
