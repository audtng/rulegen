package main

	v1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/log"

	"kubevirt.io/kubevirt/pkg/safepath"
	"kubevirt.io/kubevirt/pkg/util"
	"kubevirt.io/kubevirt/pkg/virt-controller/services"
	cmdclient "kubevirt.io/kubevirt/pkg/virt-handler/cmd-client"
	}
	defer sock.Close()

	_, err = safepath.NewPathNoFollow(socket)
	if err != nil {
		return -1, err
	}
	ufile, err := sock.(*net.UnixConn).File()
	if err != nil {
		return -1, err
