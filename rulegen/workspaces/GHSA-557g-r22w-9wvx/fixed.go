package main

	"strings"
	"syscall"

	"github.com/sylabs/singularity/internal/pkg/util/user"
)

const (
	OciSubDir = "oci"
	// SingSubDir represents directory where Singularity instance files are stored
	SingSubDir = "sing"
	// LogSubDir represents directory where Singularity instance log files are stored
	LogSubDir = "logs"
)

const (
	instancePath    = ".singularity/instances"
	authorizedChars = `^[a-zA-Z0-9._-]+$`
	prognameFormat  = "Singularity instance: %s [%s]"
)

// File represents an instance file storing instance information
type File struct {
	Path   string `json:"-"`
	Pid    int    `json:"pid"`
	PPid   int    `json:"ppid"`
	Name   string `json:"name"`
	User   string `json:"user"`
	Image  string `json:"image"`
	Config []byte `json:"config"`
	UserNs bool   `json:"userns"`
}

// ProcName returns processus name based on instance name
}

// getPath returns the path where searching for instance files
func getPath(username string, subDir string) (string, error) {
	path := ""
	var pw *user.User
	var err error
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		return path, err
	}

	path = filepath.Join(pw.Dir, instancePath, subDir, hostname, pw.Name)
	return path, nil
}

// GetDir returns directory where instances file will be stored
func GetDir(name string, subDir string) (string, error) {
	if err := CheckName(name); err != nil {
		return "", err
	}
	path, err := getPath("", subDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(path, name), nil
}

// Get returns the instance file corresponding to instance name
func Get(name string, subDir string) (*File, error) {
	if err := CheckName(name); err != nil {

// Add creates an instance file for a named instance in a privileged
// or unprivileged path
func Add(name string, subDir string) (*File, error) {
	if err := CheckName(name); err != nil {
		return nil, err
	}
	if err == nil {
		return nil, fmt.Errorf("instance %s already exists", name)
	}
	i := &File{Name: name}
	i.Path, err = getPath("", subDir)
	if err != nil {
		return nil, err
	}
// List returns instance files matching username and/or name pattern
func List(username string, name string, subDir string) ([]*File, error) {
	list := make([]*File, 0)

	path, err := getPath(username, subDir)
	if err != nil {
		return nil, err
	}
	pattern := filepath.Join(path, name, name+".json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		r, err := os.Open(file)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(r)
		r.Close()
		if err != nil {
			return nil, err
		}
		f := &File{Path: file}
		if err := json.Unmarshal(b, f); err != nil {
			return nil, err
		}
		list = append(list, f)
	}

	return list, nil
}

// Delete deletes instance file
func (i *File) Delete() error {
	return os.RemoveAll(filepath.Dir(i.Path))
}

// Update stores instance information in associated instance file
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(i.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|syscall.O_NOFOLLOW, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(b); err != nil {
		return fmt.Errorf("failed to write instance file %s: %s", i.Path, err)
	}

	return file.Sync()
}

// SetLogFile replaces stdout/stderr streams and redirect content
// to log file
func SetLogFile(name string, uid int, subDir string) (*os.File, *os.File, error) {
	path, err := getPath("", subDir)
	if err != nil {
		return nil, nil, err
	}
		return nil, nil, err
	}

	stderr, err := os.OpenFile(stderrPath, os.O_RDWR|os.O_CREATE|os.O_APPEND|syscall.O_NOFOLLOW, 0644)
	if err != nil {
		return nil, nil, err
	}

	stdout, err := os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND|syscall.O_NOFOLLOW, 0644)
	if err != nil {
		return nil, nil, err
	}
package singularity

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/sylabs/singularity/pkg/util/fs/proc"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sylabs/singularity/internal/pkg/buildcfg"
	"github.com/sylabs/singularity/pkg/util/capabilities"
)

var nsProcName = map[specs.LinuxNamespaceType]string{
	specs.PIDNamespace:     "pid",
	specs.UTSNamespace:     "uts",
	specs.IPCNamespace:     "ipc",
	specs.MountNamespace:   "mnt",
	specs.CgroupNamespace:  "cgroup",
	specs.NetworkNamespace: "net",
	specs.UserNamespace:    "user",
}

// prepareUserCaps is responsible for checking that user's requested
// capabilities are authorized
func (e *EngineOperations) prepareUserCaps() error {
		return err
	}

	// basic checks:
	// 1. a user must not use SUID workflow to join an instance
	//    started with user namespace
	// 2. a user must use SUID workflow to join an instance
	//    started without user namespace
	if starterConfig.GetIsSUID() && file.UserNs {
		return fmt.Errorf("joining user namespace with SUID workflow is not allowed")
	} else if !starterConfig.GetIsSUID() && !file.UserNs {
		return fmt.Errorf("a setuid installation is required to join this instance")
	}

	// Pid and PPid are stored in instance file and can be controlled
	// by users, just check for cool values
	if file.Pid <= 1 || file.PPid <= 1 {
		return fmt.Errorf("bad instance process ID found")
	}

	// instance configuration holding configuration read
	// from instance file
	instanceEngineConfig := singularityConfig.NewConfig()

	// extract engine configuration from instance file, the whole content
	// of this file can't be trusted
	instanceConfig := &config.Common{
		EngineConfig: instanceEngineConfig,
	}
		return err
	}

	// configuration may be altered, be sure to not panic
	if instanceEngineConfig.OciConfig.Linux == nil {
		instanceEngineConfig.OciConfig.Linux = &specs.Linux{}
	}

	// go into /proc/<pid> directory to open namespaces inodes
	// relative to current working directory while joining
	// namespaces within C starter code as changing directory
	// here also affects starter process thanks to CLONE_FS.
	// Additionally it would prevent TOCTOU races and symlink
	// usage.
	// And if instance process exits during checks or while
	// entering in namespace, we would get a "no such process"
	// error because current working directory would point to a
	// deleted inode: "/proc/self/cwd -> /proc/<pid> (deleted)"
	path := filepath.Join("/proc", strconv.Itoa(file.Pid))
	if err := mainthread.Chdir(path); err != nil {
		return err
	}

	uid := os.Getuid()
	gid := os.Getgid()

	// enforce checks while joining an instance process with SUID workflow
	// since instance file is stored in user home directory, we can't trust
	// its content when using SUID workflow
	if !file.UserNs && uid != 0 {
		// check if instance is running with user namespace enabled
		// by reading /proc/pid/uid_map
		_, hid, err := proc.ReadIDMap("uid_map")

		// if the error returned is "no such file or directory" it means
		// that user namespaces are not supported, just skip this check
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to read user namespace mapping: %s", err)
		} else if err == nil && hid > 0 {
			// a host uid greater than 0 means user namespace is in use for this process
			return fmt.Errorf("trying to join an instance running with user namespace enabled")
		}

		// read "/proc/pid/root" link of instance process must return a permission denied error
		if _, err := mainthread.Readlink("root"); !os.IsPermission(err) {
			return fmt.Errorf("trying to join a wrong instance process")
		}
		// "/proc/pid/task" directory must be owned by user UID/GID
		fi, err := os.Stat("task")
		if err != nil {
			return fmt.Errorf("error while getting information for instance task directory: %s", err)
		}
		st := fi.Sys().(*syscall.Stat_t)
		if st.Uid != uint32(uid) || st.Gid != uint32(gid) {
			return fmt.Errorf("instance process owned by %d:%d instead of %d:%d", st.Uid, st.Gid, uid, gid)
		}

		ppid := -1

		// read "/proc/pid/status" to check if instance process
		// is neither orphaned or faked
		f, err := os.Open("status")
		if err != nil {
			return fmt.Errorf("could not open status: %s", err)
		}

		for s := bufio.NewScanner(f); s.Scan(); {
			if n, _ := fmt.Sscanf(s.Text(), "PPid:\t%d", &ppid); n == 1 {
				break
			}
		}
		f.Close()

		// check that Ppid/Pid read from instance file are "somewhat" valid
		// processes
		if ppid <= 1 || ppid != file.PPid {
			return fmt.Errorf("orphaned (or faked) instance process")
		}

		// read "/proc/ppid/root" link of parent instance process must return
		// a permission denied error.
		// Also we don't use absolute path because we want to return an error
		// if current working directory is deleted.
		path := filepath.Join("..", strconv.Itoa(file.PPid), "root")
		if _, err := mainthread.Readlink(path); !os.IsPermission(err) {
			return fmt.Errorf("trying to join a wrong instance process")
		}
		// "/proc/ppid/task" directory must be owned by user UID/GID
		path = filepath.Join("..", strconv.Itoa(file.PPid), "task")
		fi, err = os.Stat(path)
		if err != nil {
			return fmt.Errorf("error while getting information for parent task directory: %s", err)
		}
		st = fi.Sys().(*syscall.Stat_t)
		if st.Uid != uint32(uid) || st.Gid != uint32(gid) {
			return fmt.Errorf("parent instance process owned by %d:%d instead of %d:%d", st.Uid, st.Gid, uid, gid)
		}
	}

	// get starter binary in use
	dest, err := mainthread.Readlink("/proc/self/exe")
	if err != nil {
		return fmt.Errorf("failed to read /proc/self/exe link: %s", err)
	}
	// should be either starter-suid or starter
	exe := filepath.Base(dest)
	path = filepath.Join("..", strconv.Itoa(file.PPid), "comm")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %s", path, err)
	}
	// check that the right starter binary is used according
	// to namespace configuration and joined instance
	if exe != strings.Trim(string(b), "\n") {
		return fmt.Errorf("%s not found in %s, wrong instance process", exe, path)
	}

	// tell starter that we are joining an instance
	starterConfig.SetJoinMount(true)

	// update namespaces path relative to /proc/<pid>
	// since starter process is in /proc/<pid> directory
	for i := range instanceEngineConfig.OciConfig.Linux.Namespaces {
		// ignore unknown namespaces
		t := instanceEngineConfig.OciConfig.Linux.Namespaces[i].Type
		if _, ok := nsProcName[t]; !ok {
			continue
		}
		// set namespace relative path
		instanceEngineConfig.OciConfig.Linux.Namespaces[i].Path = filepath.Join("ns", nsProcName[t])
	}

	// store namespace paths in starter configuration that will
	// be passed via a shared memory area and used by starter C code
	// once this process exit
	if err := starterConfig.SetNsPathFromSpec(instanceEngineConfig.OciConfig.Linux.Namespaces); err != nil {
		return err
	}

	// duplicate instance capabilities
		e.EngineConfig.OciConfig.Process.Capabilities.Ambient = instanceEngineConfig.OciConfig.Process.Capabilities.Ambient
	}

	// check if user is authorized to set those capabilities and remove
	// unauthorized capabilities from current set according to capability
	// configuration file
	if uid == 0 {
		if err := e.prepareRootCaps(); err != nil {
			return err
		}
		}
	}

	// restore apparmor profile or apply a new one if provided
	param := security.GetParam(e.EngineConfig.GetSecurity(), "apparmor")
	if param != "" {
		sylog.Debugf("Applying Apparmor profile %s", param)
		e.EngineConfig.OciConfig.SetProcessApparmorProfile(instanceEngineConfig.OciConfig.Process.ApparmorProfile)
	}

	// restore selinux context or apply a new one if provided
	param = security.GetParam(e.EngineConfig.GetSecurity(), "selinux")
	if param != "" {
		sylog.Debugf("Applying SELinux context %s", param)
		e.EngineConfig.OciConfig.SetProcessSelinuxLabel(instanceEngineConfig.OciConfig.Process.SelinuxLabel)
	}

	// restore seccomp filter or apply a new one if provided
	param = security.GetParam(e.EngineConfig.GetSecurity(), "seccomp")
	if param != "" {
		sylog.Debugf("Applying seccomp rule from %s", param)
			return err
		}
	} else {
		if e.EngineConfig.OciConfig.Linux == nil {
			e.EngineConfig.OciConfig.Linux = &specs.Linux{}
		}
		e.EngineConfig.OciConfig.Linux.Seccomp = instanceEngineConfig.OciConfig.Linux.Seccomp
	}

	// only root user can set this value based on instance file
	// and always set to true for normal users or if instance file
	// returned a wrong configuration
	if uid == 0 && instanceEngineConfig.OciConfig.Process != nil {
		e.EngineConfig.OciConfig.Process.NoNewPrivileges = instanceEngineConfig.OciConfig.Process.NoNewPrivileges
	} else {
		e.EngineConfig.OciConfig.Process.NoNewPrivileges = true
	}

	return nil
}

	"github.com/sylabs/singularity/internal/pkg/security"

	"github.com/sylabs/singularity/internal/pkg/util/user"

	specs "github.com/opencontainers/runtime-spec/specs-go"

	if engine.EngineConfig.GetInstance() {
		uid := os.Getuid()
		name := engine.CommonConfig.ContainerID

		if err := os.Chdir("/"); err != nil {
			return fmt.Errorf("failed to change directory to /: %s", err)
		}

		file, err := instance.Add(name, instance.SingSubDir)
		if err != nil {
			return err
		}
		file.PPid = os.Getpid()
		file.Image = engine.EngineConfig.GetImage()

		for _, ns := range engine.EngineConfig.OciConfig.Linux.Namespaces {
			if ns.Type == specs.UserNamespace {
				file.UserNs = true
				break
			}
		}

		return file.Update()
	}
	return nil
}
