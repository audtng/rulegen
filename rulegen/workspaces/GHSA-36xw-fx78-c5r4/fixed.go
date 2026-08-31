package main

	"io/ioutil"
	"net"
	"path/filepath"
	"strings"

	"github.com/containerd/console"
	"github.com/containerd/containerd/cmd/ctr/commands"
	s1 := filepath.Join(string(filepath.Separator), "containerd-shim", ns, id, "shim.sock")
	// this should not error, ctr always get a default ns
	ctx := namespaces.WithNamespace(gocontext.Background(), ns)
	s2, _ := shim.SocketAddress(ctx, context.GlobalString("address"), id)
	s2 = strings.TrimPrefix(s2, "unix://")

	for _, socket := range []string{s2, "\x00" + s1} {
		conn, err := net.Dial("unix", socket)
		if err == nil {
			client := ttrpc.NewClient(conn)

	if err != nil {
		return "", err
	}
	address, err := shim.SocketAddress(ctx, containerdAddress, id)
	if err != nil {
		return "", err
	}
	socket, err := shim.NewSocket(address)
	if err != nil {
		if !shim.SocketEaddrinuse(err) {
			return "", err
		}
		if err := shim.RemoveSocket(address); err != nil {
			return "", errors.Wrap(err, "remove already used socket")
		}
		if socket, err = shim.NewSocket(address); err != nil {
			return "", err
		}
	}
	f, err := socket.File()
	if err != nil {
		return "", err
	}

	cmd.ExtraFiles = append(cmd.ExtraFiles, f)

	}
	defer func() {
		if err != nil {
			_ = shim.RemoveSocket(address)
			cmd.Process.Kill()
		}
	}()
func (s *service) Shutdown(ctx context.Context, r *taskAPI.ShutdownRequest) (*ptypes.Empty, error) {
	s.cancel()
	close(s.events)
	if address, err := shim.ReadAddress("address"); err == nil {
		_ = shim.RemoveSocket(address)
	}
	return empty, nil
}

	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
		return nil, errors.Wrap(err, "failed to initialized platform behavior")
	}
	go s.forward(ctx, publisher)

	if address, err := shim.ReadAddress("address"); err == nil {
		s.shimAddress = address
	}
	return s, nil
}


	containers map[string]*runc.Container

	shimAddress string
	cancel      func()
}

func newCommand(ctx context.Context, id, containerdBinary, containerdAddress, containerdTTRPCAddress string) (*exec.Cmd, error) {
			break
		}
	}
	address, err := shim.SocketAddress(ctx, containerdAddress, grouping)
	if err != nil {
		return "", err
	}

	socket, err := shim.NewSocket(address)
	if err != nil {
		// the only time where this would happen is if there is a bug and the socket
		// was not cleaned up in the cleanup method of the shim or we are using the
		// grouping functionality where the new process should be run with the same
		// shim as an existing container
		if !shim.SocketEaddrinuse(err) {
			return "", errors.Wrap(err, "create new shim socket")
		}
		if shim.CanConnect(address) {
			if err := shim.WriteAddress("address", address); err != nil {
				return "", errors.Wrap(err, "write existing socket for shim")
			}
			return address, nil
		}
		if err := shim.RemoveSocket(address); err != nil {
			return "", errors.Wrap(err, "remove pre-existing socket")
		}
		if socket, err = shim.NewSocket(address); err != nil {
			return "", errors.Wrap(err, "try create new shim socket 2x")
		}
	}
	defer func() {
		if retErr != nil {
			socket.Close()
			_ = shim.RemoveSocket(address)
		}
	}()
	f, err := socket.File()
	if err != nil {
		return "", err
	}

	cmd.ExtraFiles = append(cmd.ExtraFiles, f)

	if err := cmd.Start(); err != nil {
		f.Close()
		return "", err
	}
	defer func() {
	if err != nil {
		return nil, err
	}
	runtime, err := runc.ReadRuntime(path)
	if err != nil {
		return nil, err
	if s.platform != nil {
		s.platform.Close()
	}
	if s.shimAddress != "" {
		_ = shim.RemoveSocket(s.shimAddress)
	}
	return empty, nil
}

	flag.BoolVar(&versionFlag, "v", false, "show the shim version and exit")
	flag.StringVar(&namespaceFlag, "namespace", "", "namespace that owns the shim")
	flag.StringVar(&idFlag, "id", "", "id of the task")
	flag.StringVar(&socketFlag, "socket", "", "socket path to serve")
	flag.StringVar(&bundlePath, "bundle", "", "path to the bundle if not workdir")

	flag.StringVar(&addressFlag, "address", "", "grpc address back to main containerd")
	ctx = context.WithValue(ctx, OptsKey{}, Opts{BundlePath: bundlePath, Debug: debugFlag})
	ctx = log.WithLogger(ctx, log.G(ctx).WithField("runtime", id))
	ctx, cancel := context.WithCancel(ctx)
	service, err := initFunc(ctx, idFlag, publisher, cancel)
	if err != nil {
		return err
		return err
	}
	go func() {
		if err := server.Serve(ctx, l); err != nil &&
			!strings.Contains(err.Error(), "use of closed network connection") {
			logrus.WithError(err).Fatal("containerd-shim: ttrpc server failure")
		}
		l.Close()
		if address, err := ReadAddress("address"); err == nil {
			_ = RemoveSocket(address)
		}

	}()
	return nil
}
		l, err = net.FileListener(os.NewFile(3, "socket"))
		path = "[inherited from parent]"
	} else {
		if len(path) > socketPathLimit {
			return nil, errors.Errorf("%q: unix socket path too long (> %d)", path, socketPathLimit)
		}
		l, err = net.Listen("unix", path)
	}
	if err != nil {
		return nil, err
	}
	logrus.WithField("socket", path).Debug("serving api on socket")
	return l, nil
}

// ErrNoAddress is returned when the address file has no content
var ErrNoAddress = errors.New("no shim address")

// ReadAddress returns the shim's socket address from the path
func ReadAddress(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
	"github.com/pkg/errors"
)

const (
	shimBinaryFormat = "containerd-shim-%s-%s"
	socketPathLimit  = 106
)

func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
	return nil
}

const socketRoot = "/run/containerd"

// SocketAddress returns a socket address
func SocketAddress(ctx context.Context, socketPath, id string) (string, error) {
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return "", err
	}
	d := sha256.Sum256([]byte(filepath.Join(socketPath, ns, id)))
	return fmt.Sprintf("unix://%s/%x", filepath.Join(socketRoot, "s"), d), nil
}

// AnonDialer returns a dialer for a socket
func AnonDialer(address string, timeout time.Duration) (net.Conn, error) {
	return dialer.Dialer(socket(address).path(), timeout)
}

func AnonReconnectDialer(address string, timeout time.Duration) (net.Conn, error) {

// NewSocket returns a new socket
func NewSocket(address string) (*net.UnixListener, error) {
	var (
		sock = socket(address)
		path = sock.path()
	)
	if !sock.isAbstract() {
		if err := os.MkdirAll(filepath.Dir(path), 0600); err != nil {
			return nil, errors.Wrapf(err, "%s", path)
		}
	}
	l, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(path, 0600); err != nil {
		os.Remove(sock.path())
		l.Close()
		return nil, err
	}
	return l.(*net.UnixListener), nil
}

const abstractSocketPrefix = "\x00"

type socket string

func (s socket) isAbstract() bool {
	return !strings.HasPrefix(string(s), "unix://")
}

func (s socket) path() string {
	path := strings.TrimPrefix(string(s), "unix://")
	// if there was no trim performed, we assume an abstract socket
	if len(path) == len(s) {
		path = abstractSocketPrefix + path
	}
	return path
}

// RemoveSocket removes the socket at the specified address if
// it exists on the filesystem
func RemoveSocket(address string) error {
	sock := socket(address)
	if !sock.isAbstract() {
		return os.Remove(sock.path())
	}
	return nil
}

// SocketEaddrinuse returns true if the provided error is caused by the
// EADDRINUSE error number
func SocketEaddrinuse(err error) bool {
	netErr, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	if netErr.Op != "listen" {
		return false
	}
	syscallErr, ok := netErr.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	errno, ok := syscallErr.Err.(syscall.Errno)
	if !ok {
		return false
	}
	return errno == syscall.EADDRINUSE
}

// CanConnect returns true if the socket provided at the address
// is accepting new connections
func CanConnect(address string) bool {
	conn, err := AnonDialer(address, 100*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
		return c, nil
	}
}

// RemoveSocket removes the socket at the specified address if
// it exists on the filesystem
func RemoveSocket(address string) error {
	return nil
}
cmd/containerd-shim/main_unix.go | 16 ++++--
runtime/v1/linux/bundle.go       | 15 ++++--
runtime/v1/shim/client/client.go | 92 ++++++++++++++++++++++++++++----
3 files changed, 105 insertions(+), 18 deletions(-)
func init() {
	flag.BoolVar(&debugFlag, "debug", false, "enable debug output in logs")
	flag.StringVar(&namespaceFlag, "namespace", "", "namespace that owns the shim")
	flag.StringVar(&socketFlag, "socket", "", "socket path to serve")
	flag.StringVar(&addressFlag, "address", "", "grpc address back to main containerd")
	flag.StringVar(&workdirFlag, "workdir", "", "path used to storge large temporary data")
	flag.StringVar(&runtimeRootFlag, "runtime-root", process.RuncRoot, "root directory for the runtime")
		f.Close()
		path = "[inherited from parent]"
	} else {
		const (
			abstractSocketPrefix = "\x00"
			socketPathLimit      = 106
		)
		p := strings.TrimPrefix(path, "unix://")
		if len(p) == len(path) {
			p = abstractSocketPrefix + p
		}
		if len(p) > socketPathLimit {
			return errors.Errorf("%q: unix socket path too long (> %d)", p, socketPathLimit)
		}
		l, err = net.Listen("unix", p)
	}
	if err != nil {
		return err
	return func(b *bundle, ns string, ropts *runctypes.RuncOptions) (shim.Config, client.Opt) {
		config := b.shimConfig(ns, c, ropts)
		return config,
			client.WithStart(c.Shim, b.shimAddress(ns, daemonAddress), daemonAddress, cgroup, c.ShimDebug, exitHandler)
	}
}


// Delete deletes the bundle from disk
func (b *bundle) Delete() error {
	address, _ := b.loadAddress()
	if address != "" {
		// we don't care about errors here
		client.RemoveSocket(address)
	}
	err := atomicDelete(b.path)
	if err == nil {
		return atomicDelete(b.workDir)
	return filepath.Join(string(filepath.Separator), "containerd-shim", namespace, b.id, "shim.sock")
}

const socketRoot = "/run/containerd"

func (b *bundle) shimAddress(namespace, socketPath string) string {
	d := sha256.Sum256([]byte(filepath.Join(socketPath, namespace, b.id)))
	return fmt.Sprintf("unix://%s/%x", filepath.Join(socketRoot, "s"), d)
}

func (b *bundle) loadAddress() (string, error) {
	return func(ctx context.Context, config shim.Config) (_ shimapi.ShimService, _ io.Closer, err error) {
		socket, err := newSocket(address)
		if err != nil {
			if !eaddrinuse(err) {
				return nil, nil, err
			}
			if err := RemoveSocket(address); err != nil {
				return nil, nil, errors.Wrap(err, "remove already used socket")
			}
			if socket, err = newSocket(address); err != nil {
				return nil, nil, err
			}
		}

		f, err := socket.File()
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get fd for socket %s", address)
			if stderrLog != nil {
				stderrLog.Close()
			}
			socket.Close()
			RemoveSocket(address)
		}()
		log.G(ctx).WithFields(logrus.Fields{
			"pid":     cmd.Process.Pid,
	}
}

func eaddrinuse(err error) bool {
	cause := errors.Cause(err)
	netErr, ok := cause.(*net.OpError)
	if !ok {
		return false
	}
	if netErr.Op != "listen" {
		return false
	}
	syscallErr, ok := netErr.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	errno, ok := syscallErr.Err.(syscall.Errno)
	if !ok {
		return false
	}
	return errno == syscall.EADDRINUSE
}

// setupOOMScore gets containerd's oom score and adds +1 to it
// to ensure a shim has a lower* score than the daemons
func setupOOMScore(shimPid int) error {
	return os.Rename(tempPath, path)
}

const (
	abstractSocketPrefix = "\x00"
	socketPathLimit      = 106
)

type socket string

func (s socket) isAbstract() bool {
	return !strings.HasPrefix(string(s), "unix://")
}

func (s socket) path() string {
	path := strings.TrimPrefix(string(s), "unix://")
	// if there was no trim performed, we assume an abstract socket
	if len(path) == len(s) {
		path = abstractSocketPrefix + path
	}
	return path
}

func newSocket(address string) (*net.UnixListener, error) {
	if len(address) > socketPathLimit {
		return nil, errors.Errorf("%q: unix socket path too long (> %d)", address, socketPathLimit)
	}
	var (
		sock = socket(address)
		path = sock.path()
	)
	if !sock.isAbstract() {
		if err := os.MkdirAll(filepath.Dir(path), 0600); err != nil {
			return nil, errors.Wrapf(err, "%s", path)
		}
	}
	l, err := net.Listen("unix", path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to listen to unix socket %q (abstract: %t)", address, sock.isAbstract())
	}
	if err := os.Chmod(path, 0600); err != nil {
		l.Close()
		return nil, err
	}

	return l.(*net.UnixListener), nil
}

// RemoveSocket removes the socket at the specified address if
// it exists on the filesystem
func RemoveSocket(address string) error {
	sock := socket(address)
	if !sock.isAbstract() {
		return os.Remove(sock.path())
	}
	return nil
}

func connect(address string, d func(string, time.Duration) (net.Conn, error)) (net.Conn, error) {
	return d(address, 100*time.Second)
}

func anonDialer(address string, timeout time.Duration) (net.Conn, error) {
	return dialer.Dialer(socket(address).path(), timeout)
}

// WithConnect connects to an existing shim
func WithConnect(address string, onClose func()) Opt {
	return func(ctx context.Context, config shim.Config) (shimapi.ShimService, io.Closer, error) {
		conn, err := connect(address, anonDialer)
		if err != nil {
			return nil, nil, err
		}
