package main


import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
}

// podLocalEphemeralStorageUsage aggregates pod local ephemeral storage usage and inode consumption for the specified stats to measure.
func podLocalEphemeralStorageUsage(podStats statsapi.PodStats, pod *v1.Pod, statsToMeasure []fsStatsType, etcHostsPath string) (v1.ResourceList, error) {
	disk := resource.Quantity{Format: resource.BinarySI}
	inodes := resource.Quantity{Format: resource.DecimalSI}

		disk.Add(podLocalVolumeUsageList[v1.ResourceEphemeralStorage])
		inodes.Add(podLocalVolumeUsageList[resourceInodes])
	}
	if len(etcHostsPath) > 0 {
		if stat, err := os.Stat(etcHostsPath); err == nil {
			disk.Add(*resource.NewQuantity(int64(stat.Size()), resource.BinarySI))
			inodes.Add(*resource.NewQuantity(int64(1), resource.DecimalSI))
		}
	}
	return v1.ResourceList{
		v1.ResourceEphemeralStorage: disk,
		resourceInodes:              inodes,

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/tools/record"
	thresholdNotifiers []ThresholdNotifier
	// thresholdsLastUpdated is the last time the thresholdNotifiers were updated.
	thresholdsLastUpdated time.Time
	// etcHostsPath is a function that will get the etc-hosts file's path for a pod given its UID
	etcHostsPath func(podUID types.UID) string
}

// ensure it implements the required interface
	recorder record.EventRecorder,
	nodeRef *v1.ObjectReference,
	clock clock.Clock,
	etcHostsPath func(types.UID) string,
) (Manager, lifecycle.PodAdmitHandler) {
	manager := &managerImpl{
		clock:                        clock,
		thresholdsFirstObservedAt:    thresholdsObservedAt{},
		dedicatedImageFs:             nil,
		thresholdNotifiers:           []ThresholdNotifier{},
		etcHostsPath:                 etcHostsPath,
	}
	return manager, manager
}
	} else {
		fsStatsSet = []fsStatsType{fsStatsRoot, fsStatsLogs, fsStatsLocalVolumeSource}
	}
	podEphemeralUsage, err := podLocalEphemeralStorageUsage(podStats, pod, fsStatsSet, m.etcHostsPath(pod.UID))
	if err != nil {
		klog.Errorf("eviction manager: error getting pod disk usage %v", err)
		return false
	}
}

// getEtcHostsPath returns the full host-side path to a pod's generated /etc/hosts file
func getEtcHostsPath(podDir string) string {
	return path.Join(podDir, "etc-hosts")
}

// makeHostsMount makes the mountpoint for the hosts file that the containers
// in a pod are injected with. podIPs is provided instead of podIP as podIPs
// are present even if dual-stack feature flag is not enabled.
func makeHostsMount(podDir string, podIPs []string, hostName, hostDomainName string, hostAliases []v1.HostAlias, useHostNetwork bool) (*kubecontainer.Mount, error) {
	hostsFilePath := getEtcHostsPath(podDir)
	if err := ensureHostsFile(hostsFilePath, podIPs, hostName, hostDomainName, hostAliases, useHostNetwork); err != nil {
		return nil, err
	}
