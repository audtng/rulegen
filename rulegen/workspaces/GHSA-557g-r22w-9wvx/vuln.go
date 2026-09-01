package main

	"strings"
	"syscall"

	"github.com/sylabs/singularity/internal/pkg/sylog"

	specs "github.com/opencontainers/runtime-spec/specs-go"

	"github.com/sylabs/singularity/internal/pkg/util/user"
	"github.com/sylabs/singularity/pkg/util/fs/proc"
)

const (
	OciSubDir = "oci"
	// SingSubDir represents directory where Singularity instance files are stored
	SingSubDir = "sing"
)

const (
	privPath        = "/var/run/singularity/instances"
	unprivPath      = ".singularity/instances"
	authorizedChars = `^[a-zA-Z0-9._-]+$`
	prognameFormat  = "Singularity instance: %s [%s]"
)

var nsMap = map[specs.LinuxNamespaceType]string{
	specs.PIDNamespace:     "pid",
	specs.UTSNamespace:     "uts",
	specs.IPCNamespace:     "ipc",
	specs.MountNamespace:   "mnt",
	specs.CgroupNamespace:  "cgroup",
	specs.NetworkNamespace: "net",
	specs.UserNamespace:    "user",
}

// File represents an instance file storing instance information
type File struct {
	Path       string `json:"-"`
	Pid        int    `json:"pid"`
	PPid       int    `json:"ppid"`
	Name       string `json:"name"`
	User       string `json:"user"`
	Image      string `json:"image"`
	Privileged bool   `json:"privileged"`
	Config     []byte `json:"config"`
}

// ProcName returns processus name based on instance name
}

// getPath returns the path where searching for instance files
func getPath(privileged bool, username string, subDir string) (string, error) {
	path := ""
	var pw *user.User
	var err error
		}
	}

	if privileged {
		path = filepath.Join(privPath, subDir, pw.Name)
		return path, nil
	}

	containerID, hostID, err := proc.ReadIDMap("/proc/self/uid_map")
	if err != nil {
		return path, err
	} else if containerID == 0 && containerID != hostID {
		if pw, err = user.GetPwUID(hostID); err != nil {
			return path, err
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		return path, err
	}

	path = filepath.Join(pw.Dir, unprivPath, subDir, hostname, pw.Name)
	return path, nil
}

func getDir(privileged bool, name string, subDir string) (string, error) {
	if err := CheckName(name); err != nil {
		return "", err
	}
	path, err := getPath(privileged, "", subDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(path, name), nil
}

// GetDirPrivileged returns directory where instances file will be stored
// if instance is run with privileges
func GetDirPrivileged(name string, subDir string) (string, error) {
	return getDir(true, name, subDir)
}

// GetDirUnprivileged returns directory where instances file will be stored
// if instance is run without privileges
func GetDirUnprivileged(name string, subDir string) (string, error) {
	return getDir(false, name, subDir)
}

// Get returns the instance file corresponding to instance name
func Get(name string, subDir string) (*File, error) {
	if err := CheckName(name); err != nil {

// Add creates an instance file for a named instance in a privileged
// or unprivileged path
func Add(name string, privileged bool, subDir string) (*File, error) {
	if err := CheckName(name); err != nil {
		return nil, err
	}
	if err == nil {
		return nil, fmt.Errorf("instance %s already exists", name)
	}
	i := &File{Name: name, Privileged: privileged}
	i.Path, err = getPath(privileged, "", subDir)
	if err != nil {
		return nil, err
	}
// List returns instance files matching username and/or name pattern
func List(username string, name string, subDir string) ([]*File, error) {
	list := make([]*File, 0)
	privileged := true

	for {
		path, err := getPath(privileged, username, subDir)
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
			}
			if err != nil {
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
		privileged = !privileged
		if privileged {
			break
		}
	}

	return list, nil
}

// PrivilegedPath returns if instance file is stored in privileged path or not
func (i *File) PrivilegedPath() bool {
	return strings.HasPrefix(i.Path, privPath)
}

// Delete deletes instance file
func (i *File) Delete() error {
	path := filepath.Dir(i.Path)

	nspath := filepath.Join(path, "ns")
	if _, err := os.Stat(nspath); err == nil {
		if err := syscall.Unmount(nspath, syscall.MNT_DETACH); err != nil {
			sylog.Errorf("can't umount %s: %s", nspath, err)
		}
	}

	return os.RemoveAll(path)
}

// Update stores instance information in associated instance file
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	if i.PrivilegedPath() {
		pw, err := user.GetPwNam(i.User)
		if err != nil {
			return err
		}
		if err := os.Chmod(path, 0550); err != nil {
			return err
		}
		if err := os.Chown(path, int(pw.UID), 0); err != nil {
			return err
		}
	}
	file, err := os.OpenFile(i.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	b = append(b, '\n')
	if n, err := file.Write(b); err != nil || n != len(b) {
		return fmt.Errorf("failed to write instance file %s: %s", i.Path, err)
	}

	return file.Sync()
}

// MountNamespaces binds /proc/<pid>/ns directory into instance folder
func (i *File) MountNamespaces() error {
	path := filepath.Join(filepath.Dir(i.Path), "ns")

	oldumask := syscall.Umask(0)
	defer syscall.Umask(oldumask)

	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	nspath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}

	src := fmt.Sprintf("/proc/%d/ns", i.Pid)
	if err := syscall.Mount(src, nspath, "", syscall.MS_BIND, ""); err != nil {
		return fmt.Errorf("mounting %s in instance folder failed: %s", src, err)
	}

	return nil
}

// UpdateNamespacesPath updates namespaces path for the provided configuration
func (i *File) UpdateNamespacesPath(configNs []specs.LinuxNamespace) error {
	path := filepath.Join(filepath.Dir(i.Path), "ns")
	nspath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	nsBase := filepath.Join(fmt.Sprintf("/proc/%d/root", i.PPid), nspath)

	procPath := fmt.Sprintf("/proc/%d/cmdline", i.PPid)

	if i.PrivilegedPath() {
		var st syscall.Stat_t

		if err := syscall.Stat(procPath, &st); err != nil {
			return err
		}
		if st.Uid != 0 || st.Gid != 0 {
			return fmt.Errorf("not an instance process")
		}

		uid := os.Geteuid()
		taskPath := fmt.Sprintf("/proc/%d/task", i.PPid)
		if err := syscall.Stat(taskPath, &st); err != nil {
			return err
		}
		if int(st.Uid) != uid {
			return fmt.Errorf("you do not own the instance")
		}
	}

	data, err := ioutil.ReadFile(procPath)
	if err != nil {
		return err
	}

	cmdline := string(data[:len(data)-1])
	procName, err := ProcName(i.Name, i.User)
	if err != nil {
		return err
	}
	if cmdline != procName {
		return fmt.Errorf("no command line match found")
	}

	for i, n := range configNs {
		ns, ok := nsMap[n.Type]
		if !ok {
			configNs[i].Path = ""
			continue
		}
		if n.Path != "" {
			configNs[i].Path = filepath.Join(nsBase, ns)
		}
	}

	return nil
}

// SetLogFile replaces stdout/stderr streams and redirect content
// to log file
func SetLogFile(name string, uid int, subDir string) (*os.File, *os.File, error) {
	path, err := getPath(false, "", subDir)
	if err != nil {
		return nil, nil, err
	}
		return nil, nil, err
	}

	stderr, err := os.OpenFile(stderrPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	stdout, err := os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}
package singularity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sylabs/singularity/internal/pkg/buildcfg"
	"github.com/sylabs/singularity/pkg/util/capabilities"
)

// prepareUserCaps is responsible for checking that user's requested
// capabilities are authorized
func (e *EngineOperations) prepareUserCaps() error {
		return err
	}

	// check if SUID workflow is really used with a privileged instance
	if !file.PrivilegedPath() && starterConfig.GetIsSUID() {
		return fmt.Errorf("try to join unprivileged instance with SUID workflow")
	}

	instanceEngineConfig := singularityConfig.NewConfig()

	// extract configuration from instance file
	instanceConfig := &config.Common{
		EngineConfig: instanceEngineConfig,
	}
		return err
	}

	starterConfig.SetJoinMount(true)

	// set namespaces to join
	if err := file.UpdateNamespacesPath(instanceEngineConfig.OciConfig.Linux.Namespaces); err != nil {
		return err
	}

	if err := starterConfig.SetNsPathFromSpec(instanceEngineConfig.OciConfig.Linux.Namespaces); err != nil {
		return err
	}

	if e.EngineConfig.OciConfig.Process == nil {
		e.EngineConfig.OciConfig.Process = &specs.Process{}
	}
	if e.EngineConfig.OciConfig.Process.Capabilities == nil {
		e.EngineConfig.OciConfig.Process.Capabilities = &specs.LinuxCapabilities{}
	}

	// duplicate instance capabilities
		e.EngineConfig.OciConfig.Process.Capabilities.Ambient = instanceEngineConfig.OciConfig.Process.Capabilities.Ambient
	}

	if os.Getuid() == 0 {
		if err := e.prepareRootCaps(); err != nil {
			return err
		}
		}
	}

	// restore apparmor profile
	param := security.GetParam(e.EngineConfig.GetSecurity(), "apparmor")
	if param != "" {
		sylog.Debugf("Applying Apparmor profile %s", param)
		e.EngineConfig.OciConfig.SetProcessApparmorProfile(instanceEngineConfig.OciConfig.Process.ApparmorProfile)
	}

	// restore selinux context
	param = security.GetParam(e.EngineConfig.GetSecurity(), "selinux")
	if param != "" {
		sylog.Debugf("Applying SELinux context %s", param)
		e.EngineConfig.OciConfig.SetProcessSelinuxLabel(instanceEngineConfig.OciConfig.Process.SelinuxLabel)
	}

	// restore security features
	param = security.GetParam(e.EngineConfig.GetSecurity(), "seccomp")
	if param != "" {
		sylog.Debugf("Applying seccomp rule from %s", param)
			return err
		}
	} else {
		if instanceEngineConfig.OciConfig.Linux != nil {
			if e.EngineConfig.OciConfig.Linux == nil {
				e.EngineConfig.OciConfig.Linux = &specs.Linux{}
			}
			e.EngineConfig.OciConfig.Linux.Seccomp = instanceEngineConfig.OciConfig.Linux.Seccomp
		}
	}

	e.EngineConfig.OciConfig.Process.NoNewPrivileges = instanceEngineConfig.OciConfig.Process.NoNewPrivileges

	return nil
}

	"github.com/sylabs/singularity/internal/pkg/security"

	"github.com/sylabs/singularity/internal/pkg/util/mainthread"
	"github.com/sylabs/singularity/internal/pkg/util/user"

	specs "github.com/opencontainers/runtime-spec/specs-go"

	if engine.EngineConfig.GetInstance() {
		uid := os.Getuid()
		gid := os.Getgid()
		name := engine.CommonConfig.ContainerID
		privileged := true

		if err := os.Chdir("/"); err != nil {
			return fmt.Errorf("failed to change directory to /: %s", err)
		}

		if engine.EngineConfig.OciConfig.Linux != nil {
			for _, ns := range engine.EngineConfig.OciConfig.Linux.Namespaces {
				if ns.Type == specs.UserNamespace {
					privileged = false
					break
				}
			}
		}

		file, err := instance.Add(name, privileged, instance.SingSubDir)
		if err != nil {
			return err
		}
		file.PPid = os.Getpid()
		file.Image = engine.EngineConfig.GetImage()

		if privileged {
			var err error

			mainthread.Execute(func() {
				if err = syscall.Setresuid(0, 0, uid); err != nil {
					err = fmt.Errorf("failed to escalate uid privileges")
					return
				}
				if err = syscall.Setresgid(0, 0, gid); err != nil {
					err = fmt.Errorf("failed to escalate gid privileges")
					return
				}
				if err = file.Update(); err != nil {
					return
				}
				if err = file.MountNamespaces(); err != nil {
					return
				}
				if err = syscall.Setresgid(gid, gid, 0); err != nil {
					err = fmt.Errorf("failed to escalate gid privileges")
					return
				}
				if err = syscall.Setresuid(uid, uid, 0); err != nil {
					err = fmt.Errorf("failed to escalate uid privileges")
					return
				}
			})

			return err
		}

		if err := file.Update(); err != nil {
			return err
		}
		return file.MountNamespaces()
	}
	return nil
}
