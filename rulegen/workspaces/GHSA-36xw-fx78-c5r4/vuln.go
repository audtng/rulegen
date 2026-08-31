package main

	"io/ioutil"
	"net"
	"path/filepath"

	"github.com/containerd/console"
	"github.com/containerd/containerd/cmd/ctr/commands"
	s1 := filepath.Join(string(filepath.Separator), "containerd-shim", ns, id, "shim.sock")
	// this should not error, ctr always get a default ns
	ctx := namespaces.WithNamespace(gocontext.Background(), ns)
	s2, _ := shim.SocketAddress(ctx, id)

	for _, socket := range []string{s1, s2} {
		conn, err := net.Dial("unix", "\x00"+socket)
		if err == nil {
			client := ttrpc.NewClient(conn)

	if err != nil {
		return "", err
	}
	address, err := shim.SocketAddress(ctx, id)
	if err != nil {
		return "", err
	}
	socket, err := shim.NewSocket(address)
	if err != nil {
		return "", err
	}
	defer socket.Close()
	f, err := socket.File()
	if err != nil {
		return "", err
	}
	defer f.Close()

	cmd.ExtraFiles = append(cmd.ExtraFiles, f)

	}
	defer func() {
		if err != nil {
			cmd.Process.Kill()
		}
	}()
func (s *service) Shutdown(ctx context.Context, r *taskAPI.ShutdownRequest) (*ptypes.Empty, error) {
	s.cancel()
	close(s.events)
	return empty, nil
}

	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
		return nil, errors.Wrap(err, "failed to initialized platform behavior")
	}
	go s.forward(ctx, publisher)
	return s, nil
}


	containers map[string]*runc.Container

	cancel func()
}

func newCommand(ctx context.Context, id, containerdBinary, containerdAddress, containerdTTRPCAddress string) (*exec.Cmd, error) {
			break
		}
	}
	address, err := shim.SocketAddress(ctx, grouping)
	if err != nil {
		return "", err
	}
	socket, err := shim.NewSocket(address)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			if err := shim.WriteAddress("address", address); err != nil {
				return "", err
			}
			return address, nil
		}
		return "", err
	}
	defer socket.Close()
	f, err := socket.File()
	if err != nil {
		return "", err
	}
	defer f.Close()

	cmd.ExtraFiles = append(cmd.ExtraFiles, f)

	if err := cmd.Start(); err != nil {
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

	return empty, nil
}

	flag.BoolVar(&versionFlag, "v", false, "show the shim version and exit")
	flag.StringVar(&namespaceFlag, "namespace", "", "namespace that owns the shim")
	flag.StringVar(&idFlag, "id", "", "id of the task")
	flag.StringVar(&socketFlag, "socket", "", "abstract socket path to serve")
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
		defer l.Close()
		if err := server.Serve(ctx, l); err != nil &&
			!strings.Contains(err.Error(), "use of closed network connection") {
			logrus.WithError(err).Fatal("containerd-shim: ttrpc server failure")
		}
	}()
	return nil
}
		l, err = net.FileListener(os.NewFile(3, "socket"))
		path = "[inherited from parent]"
	} else {
		if len(path) > 106 {
			return nil, errors.Errorf("%q: unix socket path too long (> 106)", path)
		}
		l, err = net.Listen("unix", "\x00"+path)
	}
	if err != nil {
		return nil, err
	}
	logrus.WithField("socket", path).Debug("serving api on abstract socket")
	return l, nil
}

// ErrNoAddress is returned when the address file has no content
var ErrNoAddress = errors.New("no shim address")

// ReadAddress returns the shim's abstract socket address from the path
func ReadAddress(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
	"github.com/pkg/errors"
)

const shimBinaryFormat = "containerd-shim-%s-%s"

func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
	return nil
}

// SocketAddress returns an abstract socket address
func SocketAddress(ctx context.Context, id string) (string, error) {
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return "", err
	}
	d := sha256.Sum256([]byte(filepath.Join(ns, id)))
	return filepath.Join(string(filepath.Separator), "containerd-shim", fmt.Sprintf("%x.sock", d)), nil
}

// AnonDialer returns a dialer for an abstract socket
func AnonDialer(address string, timeout time.Duration) (net.Conn, error) {
	address = strings.TrimPrefix(address, "unix://")
	return dialer.Dialer("\x00"+address, timeout)
}

func AnonReconnectDialer(address string, timeout time.Duration) (net.Conn, error) {

// NewSocket returns a new socket
func NewSocket(address string) (*net.UnixListener, error) {
	if len(address) > 106 {
		return nil, errors.Errorf("%q: unix socket path too long (> 106)", address)
	}
	l, err := net.Listen("unix", "\x00"+address)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to listen to abstract unix socket %q", address)
	}
	return l.(*net.UnixListener), nil
}
		return c, nil
	}
}
cmd/containerd-shim/main_unix.go | 16 ++++--
runtime/v1/linux/bundle.go       | 15 ++++--
runtime/v1/shim/client/client.go | 92 ++++++++++++++++++++++++++++----
3 files changed, 105 insertions(+), 18 deletions(-)
func init() {
	flag.BoolVar(&debugFlag, "debug", false, "enable debug output in logs")
	flag.StringVar(&namespaceFlag, "namespace", "", "namespace that owns the shim")
	flag.StringVar(&socketFlag, "socket", "", "abstract socket path to serve")
	flag.StringVar(&addressFlag, "address", "", "grpc address back to main containerd")
	flag.StringVar(&workdirFlag, "workdir", "", "path used to storge large temporary data")
	flag.StringVar(&runtimeRootFlag, "runtime-root", process.RuncRoot, "root directory for the runtime")
		f.Close()
		path = "[inherited from parent]"
	} else {
		if len(path) > 106 {
			return errors.Errorf("%q: unix socket path too long (> 106)", path)
		}
		l, err = net.Listen("unix", "\x00"+path)
	}
	if err != nil {
		return err
	return func(b *bundle, ns string, ropts *runctypes.RuncOptions) (shim.Config, client.Opt) {
		config := b.shimConfig(ns, c, ropts)
		return config,
			client.WithStart(c.Shim, b.shimAddress(ns), daemonAddress, cgroup, c.ShimDebug, exitHandler)
	}
}


// Delete deletes the bundle from disk
func (b *bundle) Delete() error {
	err := atomicDelete(b.path)
	if err == nil {
		return atomicDelete(b.workDir)
	return filepath.Join(string(filepath.Separator), "containerd-shim", namespace, b.id, "shim.sock")
}

func (b *bundle) shimAddress(namespace string) string {
	d := sha256.Sum256([]byte(filepath.Join(namespace, b.id)))
	return filepath.Join(string(filepath.Separator), "containerd-shim", fmt.Sprintf("%x.sock", d))
}

func (b *bundle) loadAddress() (string, error) {
	return func(ctx context.Context, config shim.Config) (_ shimapi.ShimService, _ io.Closer, err error) {
		socket, err := newSocket(address)
		if err != nil {
			return nil, nil, err
		}
		defer socket.Close()
		f, err := socket.File()
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get fd for socket %s", address)
			if stderrLog != nil {
				stderrLog.Close()
			}
		}()
		log.G(ctx).WithFields(logrus.Fields{
			"pid":     cmd.Process.Pid,
	}
}

// setupOOMScore gets containerd's oom score and adds +1 to it
// to ensure a shim has a lower* score than the daemons
func setupOOMScore(shimPid int) error {
	return os.Rename(tempPath, path)
}

func newSocket(address string) (*net.UnixListener, error) {
	if len(address) > 106 {
		return nil, errors.Errorf("%q: unix socket path too long (> 106)", address)
	}
	l, err := net.Listen("unix", "\x00"+address)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to listen to abstract unix socket %q", address)
	}

	return l.(*net.UnixListener), nil
}

func connect(address string, d func(string, time.Duration) (net.Conn, error)) (net.Conn, error) {
	return d(address, 100*time.Second)
}

func annonDialer(address string, timeout time.Duration) (net.Conn, error) {
	address = strings.TrimPrefix(address, "unix://")
	return dialer.Dialer("\x00"+address, timeout)
}

// WithConnect connects to an existing shim
func WithConnect(address string, onClose func()) Opt {
	return func(ctx context.Context, config shim.Config) (shimapi.ShimService, io.Closer, error) {
		conn, err := connect(address, annonDialer)
		if err != nil {
			return nil, nil, err
		}
