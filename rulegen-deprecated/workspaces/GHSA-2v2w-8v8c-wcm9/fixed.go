package main

package provisioning

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rancher/norman/types"

	"github.com/sirupsen/logrus"

	"github.com/rancher/shepherd/clients/corral"
	"github.com/rancher/shepherd/clients/rancher"

	apiv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"

	"github.com/rancher/rancher/tests/v2/actions/clusters"
	k3sHardening "github.com/rancher/rancher/tests/v2/actions/hardening/k3s"
	rke1Hardening "github.com/rancher/rancher/tests/v2/actions/hardening/rke1"
	rke2Hardening "github.com/rancher/rancher/tests/v2/actions/hardening/rke2"
	"github.com/rancher/rancher/tests/v2/actions/machinepools"
	"github.com/rancher/rancher/tests/v2/actions/pipeline"
	"github.com/rancher/rancher/tests/v2/actions/provisioninginput"
	nodepools "github.com/rancher/rancher/tests/v2/actions/rke1/nodepools"
	"github.com/rancher/rancher/tests/v2/actions/rke1/nodetemplates"
	"github.com/rancher/rancher/tests/v2/actions/secrets"
	v1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/cloudcredentials"
	"github.com/rancher/shepherd/extensions/cloudcredentials/aws"
	"github.com/rancher/shepherd/extensions/cloudcredentials/azure"
	"github.com/rancher/shepherd/extensions/cloudcredentials/google"
	"github.com/rancher/shepherd/extensions/cloudcredentials/vsphere"
	shepherdclusters "github.com/rancher/shepherd/extensions/clusters"
	"github.com/rancher/shepherd/extensions/clusters/aks"
	"github.com/rancher/shepherd/extensions/clusters/eks"
	"github.com/rancher/shepherd/extensions/clusters/gke"
	"github.com/rancher/shepherd/extensions/defaults"
	"github.com/rancher/shepherd/extensions/defaults/stevetypes"
	"github.com/rancher/shepherd/extensions/etcdsnapshot"
	nodestat "github.com/rancher/shepherd/extensions/nodes"
	"github.com/rancher/shepherd/extensions/tokenregistration"
	"github.com/rancher/shepherd/pkg/environmentflag"
	namegen "github.com/rancher/shepherd/pkg/namegenerator"
	"github.com/rancher/shepherd/pkg/nodes"
	"github.com/rancher/shepherd/pkg/wait"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kwait "k8s.io/apimachinery/pkg/util/wait"

	management "github.com/rancher/shepherd/clients/rancher/generated/management/v3"
)

const (
	active         = "active"
	internalIP     = "alpha.kubernetes.io/provided-node-ip"
	rke1ExternalIP = "rke.cattle.io/external-ip"
	namespace      = "fleet-default"

	rke2k3sAirgapCustomCluster           = "rke2k3sairgapcustomcluster"
	rke2k3sNodeCorralName                = "rke2k3sregisterNode"
	corralPackageAirgapCustomClusterName = "airgapCustomCluster"
	rke1AirgapCustomCluster              = "rke1airgapcustomcluster"
	rke1NodeCorralName                   = "rke1registerNode"
)

var (
	updateConfig = true
)

// CreateProvisioningCluster provisions a non-rke1 cluster, then runs verify checks
func CreateProvisioningCluster(client *rancher.Client, provider Provider, clustersConfig *clusters.ClusterConfig, hostnameTruncation []machinepools.HostnameTruncation) (*v1.SteveAPIObject, error) {
	credentialSpec := cloudcredentials.LoadCloudCredential(string(provider.Name))
	cloudCredential, err := provider.CloudCredFunc(client, credentialSpec)
	if err != nil {
		return nil, err
	}

	if clustersConfig.PSACT == string(provisioninginput.RancherBaseline) {
		err = clusters.CreateRancherBaselinePSACT(client, clustersConfig.PSACT)
		if err != nil {
			return nil, err
		}
	}

	clusterName := namegen.AppendRandomString(provider.Name.String())
	generatedPoolName := fmt.Sprintf("nc-%s-pool1-", clusterName)
	machinePoolConfigs := provider.MachinePoolFunc(generatedPoolName, namespace)

	var machinePoolResponses []v1.SteveAPIObject

	for _, machinePoolConfig := range machinePoolConfigs {
		machinePoolConfigResp, err := client.Steve.
			SteveType(provider.MachineConfigPoolResourceSteveType).
			Create(&machinePoolConfig)
		if err != nil {
			return nil, err
		}
		machinePoolResponses = append(machinePoolResponses, *machinePoolConfigResp)
	}

	if clustersConfig.Registries != nil {
		if clustersConfig.Registries.RKE2Registries != nil {
			if clustersConfig.Registries.RKE2Username != "" && clustersConfig.Registries.RKE2Password != "" {
				steveClient, err := client.Steve.ProxyDownstream("local")
				if err != nil {
					return nil, err
				}

				secretName := fmt.Sprintf("priv-reg-sec-%s", clusterName)
				secretTemplate := secrets.NewSecretTemplate(secretName, namespace, map[string][]byte{
					"password": []byte(clustersConfig.Registries.RKE2Password),
					"username": []byte(clustersConfig.Registries.RKE2Username),
				},
					corev1.SecretTypeBasicAuth,
				)

				registrySecret, err := steveClient.SteveType(secrets.SecretSteveType).Create(secretTemplate)
				if err != nil {
					return nil, err
				}

				for registryName, registry := range clustersConfig.Registries.RKE2Registries.Configs {
					registry.AuthConfigSecretName = registrySecret.Name
					clustersConfig.Registries.RKE2Registries.Configs[registryName] = registry
				}
			}
		}
	}

	var machineConfigs []machinepools.MachinePoolConfig
	var pools []machinepools.Pools
	for _, pool := range clustersConfig.MachinePools {
		machineConfigs = append(machineConfigs, pool.MachinePoolConfig)
		pools = append(pools, pool.Pools)
	}

	machinePools := machinepools.
		CreateAllMachinePools(machineConfigs, pools, machinePoolResponses, provider.Roles, hostnameTruncation)

	if clustersConfig.CloudProvider == provisioninginput.VsphereCloudProviderName.String() {

		vcenterCredentials := map[string]interface{}{
			"datacenters": machinePoolConfigs[0].Object["datacenter"],
			"host":        credentialSpec.VmwareVsphereConfig.Vcenter,
			"password":    vsphere.GetVspherePassword(),
			"username":    credentialSpec.VmwareVsphereConfig.Username,
		}
		clustersConfig.AddOnConfig = &provisioninginput.AddOnConfig{
			ChartValues: &rkev1.GenericMap{
				Data: map[string]interface{}{
					"rancher-vsphere-cpi": map[string]interface{}{
						"vCenter": vcenterCredentials,
					},
					"rancher-vsphere-csi": map[string]interface{}{
						"storageClass": map[string]interface{}{
							"datastoreURL": machinePoolConfigs[0].Object["datastoreUrl"],
						},
						"vCenter": vcenterCredentials,
					},
				},
			},
		}
	}

	cluster := clusters.NewK3SRKE2ClusterConfig(clusterName, namespace, clustersConfig, machinePools, cloudCredential.Namespace+":"+cloudCredential.Name)

	for _, truncatedPool := range hostnameTruncation {
		if truncatedPool.PoolNameLengthLimit > 0 || truncatedPool.ClusterNameLengthLimit > 0 {
			cluster.GenerateName = "t-"
			if truncatedPool.ClusterNameLengthLimit > 0 {
				cluster.Spec.RKEConfig.MachinePoolDefaults.HostnameLengthLimit = truncatedPool.ClusterNameLengthLimit
			}

			break
		}
	}

	_, err = shepherdclusters.CreateK3SRKE2Cluster(client, cluster)
	if err != nil {
		return nil, err
	}

	if client.Flags.GetValue(environmentflag.UpdateClusterName) && updateConfig {
		pipeline.UpdateConfigClusterName(clusterName)
	}

	adminClient, err := rancher.NewClient(client.RancherConfig.AdminToken, client.Session)
	if err != nil {
		return nil, err
	}

	createdCluster, err := adminClient.Steve.
		SteveType(stevetypes.Provisioning).
		ByID(namespace + "/" + clusterName)

	return createdCluster, err
}

// CreateProvisioningCustomCluster provisions a non-rke1 cluster using a 3rd party client for its nodes, then runs verify checks
func CreateProvisioningCustomCluster(client *rancher.Client, externalNodeProvider *ExternalNodeProvider, clustersConfig *clusters.ClusterConfig) (*v1.SteveAPIObject, error) {
	setLogrusFormatter()
	rolesPerNode := []string{}
	quantityPerPool := []int32{}
	rolesPerPool := []string{}
	for _, pool := range clustersConfig.MachinePools {
		var finalRoleCommand string
		if pool.MachinePoolConfig.ControlPlane {
			finalRoleCommand += " --controlplane"
		}

		if pool.MachinePoolConfig.Etcd {
			finalRoleCommand += " --etcd"
		}

		if pool.MachinePoolConfig.Worker {
			finalRoleCommand += " --worker"
		}

		if pool.MachinePoolConfig.Windows {
			finalRoleCommand += " --windows"
		}

		quantityPerPool = append(quantityPerPool, pool.MachinePoolConfig.Quantity)
		rolesPerPool = append(rolesPerPool, finalRoleCommand)
		for i := int32(0); i < pool.MachinePoolConfig.Quantity; i++ {
			rolesPerNode = append(rolesPerNode, finalRoleCommand)
		}
	}

	if clustersConfig.PSACT == string(provisioninginput.RancherBaseline) {
		err := clusters.CreateRancherBaselinePSACT(client, clustersConfig.PSACT)
		if err != nil {
			return nil, err
		}
	}

	nodes, err := externalNodeProvider.NodeCreationFunc(client, rolesPerPool, quantityPerPool)
	if err != nil {
		return nil, err
	}

	clusterName := namegen.AppendRandomString(externalNodeProvider.Name)

	cluster := clusters.NewK3SRKE2ClusterConfig(clusterName, namespace, clustersConfig, nil, "")

	if clustersConfig.Hardened && strings.Contains(clustersConfig.KubernetesVersion, shepherdclusters.RKE2ClusterType.String()) {
		err = rke2Hardening.HardenRKE2Nodes(nodes, rolesPerNode)
		if err != nil {
			return nil, err
		}

		cluster = clusters.HardenRKE2ClusterConfig(clusterName, namespace, clustersConfig, nil, "")
	}

	clusterResp, err := shepherdclusters.CreateK3SRKE2Cluster(client, cluster)
	if err != nil {
		return nil, err
	}

	if client.Flags.GetValue(environmentflag.UpdateClusterName) && updateConfig {
		pipeline.UpdateConfigClusterName(clusterName)
	}

	client, err = client.ReLogin()
	if err != nil {
		return nil, err
	}

	customCluster, err := client.Steve.SteveType(etcdsnapshot.ProvisioningSteveResouceType).ByID(clusterResp.ID)
	if err != nil {
		return nil, err
	}

	clusterStatus := &apiv1.ClusterStatus{}
	err = v1.ConvertToK8sType(customCluster.Status, clusterStatus)
	if err != nil {
		return nil, err
	}

	token, err := tokenregistration.GetRegistrationToken(client, clusterStatus.ClusterName)
	if err != nil {
		return nil, err
	}

	kubeProvisioningClient, err := client.GetKubeAPIProvisioningClient()
	if err != nil {
		return nil, err
	}

	result, err := kubeProvisioningClient.Clusters(namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + clusterName,
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	})
	if err != nil {
		return nil, err
	}

	checkFunc := shepherdclusters.IsProvisioningClusterReady
	var command string
	totalNodesObserved := 0
	for poolIndex, poolRole := range rolesPerPool {
		if strings.Contains(poolRole, "windows") {
			totalNodesObserved += int(quantityPerPool[poolIndex])
			continue
		}
		for nodeIndex := 0; nodeIndex < int(quantityPerPool[poolIndex]); nodeIndex++ {
			node := nodes[totalNodesObserved+nodeIndex]

			logrus.Infof("Execute Registration Command for node %s", node.NodeID)
			logrus.Infof("Linux pool detected, using bash...")

			command = fmt.Sprintf("%s %s", token.InsecureNodeCommand, poolRole)
			if clustersConfig.MachinePools[poolIndex].IsSecure {
				command = fmt.Sprintf("%s %s", token.NodeCommand, poolRole)
			}
			command = createRegistrationCommand(command, node.PublicIPAddress, node.PrivateIPAddress, clustersConfig.MachinePools[poolIndex])
			logrus.Infof("Command: %s", command)

			output, err := node.ExecuteCommand(command)
			if err != nil {
				return nil, err
			}
			logrus.Infof(output)
		}
		totalNodesObserved += int(quantityPerPool[poolIndex])
	}

	err = wait.WatchWait(result, checkFunc)
	if err != nil {
		return nil, err
	}
	totalNodesObserved = 0
	for poolIndex := 0; poolIndex < len(rolesPerPool); poolIndex++ {
		if strings.Contains(rolesPerPool[poolIndex], "windows") {
			for nodeIndex := 0; nodeIndex < int(quantityPerPool[poolIndex]); nodeIndex++ {
				node := nodes[totalNodesObserved+nodeIndex]

				logrus.Infof("Execute Registration Command for node %s", node.NodeID)
				logrus.Infof("Windows pool detected, using powershell.exe...")
				command = fmt.Sprintf("powershell.exe %s ", token.InsecureWindowsNodeCommand)
				if clustersConfig.MachinePools[poolIndex].IsSecure {
					command = fmt.Sprintf("powershell.exe %s ", token.WindowsNodeCommand)
				}
				command = createWindowsRegistrationCommand(command, node.PublicIPAddress, node.PrivateIPAddress, clustersConfig.MachinePools[poolIndex])
				logrus.Infof("Command: %s", command)

				output, err := node.ExecuteCommand(command)
				if err != nil {
					return nil, err
				}
				logrus.Infof(output)
			}
		}
		totalNodesObserved += int(quantityPerPool[poolIndex])
	}

	if clustersConfig.Hardened {
		if strings.Contains(clustersConfig.KubernetesVersion, shepherdclusters.K3SClusterType.String()) {
			err = k3sHardening.HardenK3SNodes(nodes, rolesPerNode, clustersConfig.KubernetesVersion)
			if err != nil {
				return nil, err
			}

			hardenCluster := clusters.HardenK3SClusterConfig(clusterName, namespace, clustersConfig, nil, "")

			_, err := shepherdclusters.UpdateK3SRKE2Cluster(client, clusterResp, hardenCluster)
			if err != nil {
				return nil, err
			}
		} else {
			err = rke2Hardening.PostRKE2HardeningConfig(nodes, rolesPerNode)
			if err != nil {
				return nil, err
			}
		}
	}

	createdCluster, err := client.Steve.
		SteveType(stevetypes.Provisioning).
		ByID(namespace + "/" + clusterName)
	return createdCluster, err
}

// CreateProvisioningRKE1Cluster provisions an rke1 cluster, then runs verify checks
func CreateProvisioningRKE1Cluster(client *rancher.Client, provider RKE1Provider, clustersConfig *clusters.ClusterConfig, nodeTemplate *nodetemplates.NodeTemplate) (*management.Cluster, error) {
	if clustersConfig.PSACT == string(provisioninginput.RancherBaseline) {
		err := clusters.CreateRancherBaselinePSACT(client, clustersConfig.PSACT)
		if err != nil {
			return nil, err
		}
	}

	clusterName := namegen.AppendRandomString(provider.Name.String())
	cluster := clusters.NewRKE1ClusterConfig(clusterName, client, clustersConfig)
	clusterResp, err := shepherdclusters.CreateRKE1Cluster(client, cluster)
	if err != nil {
		return nil, err
	}

	if client.Flags.GetValue(environmentflag.UpdateClusterName) && updateConfig {
		pipeline.UpdateConfigClusterName(clusterName)
	}

	var nodeRoles []nodepools.NodeRoles
	for _, nodes := range clustersConfig.NodePools {
		nodeRoles = append(nodeRoles, nodes.NodeRoles)
	}
	_, err = nodepools.NodePoolSetup(client, nodeRoles, clusterResp.ID, nodeTemplate.ID)
	if err != nil {
		return nil, err
	}

	createdCluster, err := client.Management.Cluster.ByID(clusterResp.ID)
	return createdCluster, err
}

// CreateProvisioningRKE1CustomCluster provisions an rke1 cluster using a 3rd party client for its nodes, then runs verify checks
func CreateProvisioningRKE1CustomCluster(client *rancher.Client, externalNodeProvider *ExternalNodeProvider, clustersConfig *clusters.ClusterConfig) (*management.Cluster, []*nodes.Node, error) {
	setLogrusFormatter()
	quantityPerPool := []int32{}
	rolesPerPool := []string{}
	for _, pool := range clustersConfig.NodePools {
		var finalRoleCommand string
		if pool.NodeRoles.ControlPlane {
			finalRoleCommand += " --controlplane"
		}
		if pool.NodeRoles.Etcd {
			finalRoleCommand += " --etcd"
		}
		if pool.NodeRoles.Worker {
			finalRoleCommand += " --worker"
		}

		quantityPerPool = append(quantityPerPool, int32(pool.NodeRoles.Quantity))
		rolesPerPool = append(rolesPerPool, finalRoleCommand)
	}

	if clustersConfig.PSACT == string(provisioninginput.RancherBaseline) {
		err := clusters.CreateRancherBaselinePSACT(client, clustersConfig.PSACT)
		if err != nil {
			return nil, nil, err
		}
	}

	nodes, err := externalNodeProvider.NodeCreationFunc(client, rolesPerPool, quantityPerPool)
	if err != nil {
		return nil, nil, err
	}

	clusterName := namegen.AppendRandomString(externalNodeProvider.Name)

	cluster := clusters.NewRKE1ClusterConfig(clusterName, client, clustersConfig)

	if clustersConfig.Hardened {
		err = rke1Hardening.HardenRKE1Nodes(nodes, rolesPerPool)
		if err != nil {
			return nil, nil, err
		}

		cluster = clusters.HardenRKE1ClusterConfig(client, clusterName, clustersConfig)
	}

	clusterResp, err := shepherdclusters.CreateRKE1Cluster(client, cluster)
	if err != nil {
		return nil, nil, err
	}

	if client.Flags.GetValue(environmentflag.UpdateClusterName) && updateConfig {
		pipeline.UpdateConfigClusterName(clusterName)
	}

	client, err = client.ReLogin()
	if err != nil {
		return nil, nil, err
	}

	customCluster, err := client.Management.Cluster.ByID(clusterResp.ID)
	if err != nil {
		return nil, nil, err
	}

	token, err := tokenregistration.GetRegistrationToken(client, customCluster.ID)
	if err != nil {
		return nil, nil, err
	}

	adminClient, err := rancher.NewClient(client.RancherConfig.AdminToken, client.Session)
	if err != nil {
		return nil, nil, err
	}

	result, err := adminClient.GetManagementWatchInterface(management.ClusterType, metav1.ListOptions{
		FieldSelector:  "metadata.name=" + customCluster.ID,
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	})
	if err != nil {
		return nil, nil, err
	}

	checkFunc := shepherdclusters.IsHostedProvisioningClusterReady

	var command string
	totalNodesObserved := 0
	for poolIndex, poolRole := range rolesPerPool {
		for nodeIndex := 0; nodeIndex < int(quantityPerPool[poolIndex]); nodeIndex++ {
			node := nodes[totalNodesObserved+nodeIndex]

			logrus.Infof("Execute Registration Command for node %s", node.NodeID)
			logrus.Infof("Linux pool detected, using bash...")

			command = fmt.Sprintf("%s %s", token.NodeCommand, poolRole)
			command = createRKE1RegistrationCommand(command, node.PublicIPAddress, node.PrivateIPAddress, clustersConfig.NodePools[poolIndex])
			logrus.Infof("Command: %s", command)

			if clustersConfig.RKE1CustomClusterDockerInstall != nil && clustersConfig.RKE1CustomClusterDockerInstall.InstallDockerURL != "" {
				_, err := node.ExecuteCommand("curl " + clustersConfig.RKE1CustomClusterDockerInstall.InstallDockerURL + " | sh")
				if err != nil {
					return nil, nil, err
				}

				_, err = node.ExecuteCommand("sudo systemctl start docker")
				if err != nil {
					return nil, nil, err
				}

				_, err = node.ExecuteCommand("sudo chmod 777 /var/run/docker.sock")
				if err != nil {
					return nil, nil, err
				}
			}

			output, err := node.ExecuteCommand(command)
			if err != nil {
				return nil, nil, err
			}
			logrus.Infof(output)
		}
		totalNodesObserved += int(quantityPerPool[poolIndex])
	}

	err = wait.WatchWait(result, checkFunc)
	if err != nil {
		return nil, nil, err
	}

	if clustersConfig.Hardened {
		err = rke1Hardening.PostRKE1HardeningConfig(nodes, rolesPerPool)
		if err != nil {
			return nil, nil, err
		}
	}

	createdCluster, err := client.Management.Cluster.ByID(clusterResp.ID)

	return createdCluster, nodes, err
}

// CreateProvisioningAirgapCustomCluster provisions a non-rke1 cluster using corral to gather its nodes, then runs verify checks
func CreateProvisioningAirgapCustomCluster(client *rancher.Client, clustersConfig *clusters.ClusterConfig, corralPackages *corral.Packages) (*v1.SteveAPIObject, error) {
	setLogrusFormatter()
	quantityPerPool := []int32{}
	rolesPerPool := []string{}
	for _, pool := range clustersConfig.MachinePools {
		var finalRoleCommand string
		if pool.MachinePoolConfig.ControlPlane {
			finalRoleCommand += " --controlplane"
		}

		if pool.MachinePoolConfig.Etcd {
			finalRoleCommand += " --etcd"
		}

		if pool.MachinePoolConfig.Worker {
			finalRoleCommand += " --worker"
		}

		if pool.MachinePoolConfig.Windows {
			finalRoleCommand += " --windows"
		}

		quantityPerPool = append(quantityPerPool, pool.MachinePoolConfig.Quantity)
		rolesPerPool = append(rolesPerPool, finalRoleCommand)

	}

	if clustersConfig.PSACT == string(provisioninginput.RancherBaseline) {
		err := clusters.CreateRancherBaselinePSACT(client, clustersConfig.PSACT)
		if err != nil {
			return nil, err
		}
	}

	clusterName := namegen.AppendRandomString(rke2k3sAirgapCustomCluster)

	cluster := clusters.NewK3SRKE2ClusterConfig(clusterName, namespace, clustersConfig, nil, "")

	clusterResp, err := shepherdclusters.CreateK3SRKE2Cluster(client, cluster)
	if err != nil {
		return nil, err
	}

	client, err = client.ReLogin()
	if err != nil {
		return nil, err
	}

	customCluster, err := client.Steve.SteveType(stevetypes.Provisioning).ByID(clusterResp.ID)
	if err != nil {
		return nil, err
	}

	clusterStatus := &apiv1.ClusterStatus{}
	err = v1.ConvertToK8sType(customCluster.Status, clusterStatus)
	if err != nil {
		return nil, err
	}

	token, err := tokenregistration.GetRegistrationToken(client, clusterStatus.ClusterName)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Register Custom Cluster Through Corral")
	corralsArgs := []corral.Args{}

	for poolIndex, poolRole := range rolesPerPool {

		regCmd := fmt.Sprintf("%s %s", token.InsecureNodeCommand, poolRole)

		// environment variables must be escaped inside original registration command
		regCmd = strings.Replace(regCmd, "\"", "\\\"", -1)

		corralsArgs = append(corralsArgs, corral.Args{
			Name:        namegen.AppendRandomString(rke2k3sNodeCorralName),
			PackageName: corralPackages.CorralPackageImages[corralPackageAirgapCustomClusterName],
			Updates:     map[string]string{"registration_command": regCmd, "node_count": fmt.Sprint(quantityPerPool[poolIndex])},
		})
	}

	_, err = corral.CreateMultipleCorrals(client.Session, corralsArgs, corralPackages.HasDebug, corralPackages.HasCleanup)
	if err != nil {
		return nil, err
	}

	createdCluster, err := client.Steve.SteveType(stevetypes.Provisioning).ByID(namespace + "/" + clusterName)
	return createdCluster, err
}

// CreateProvisioningRKE1AirgapCustomCluster provisions an rke1 cluster using corral to gather its nodes, then runs verify checks
func CreateProvisioningRKE1AirgapCustomCluster(client *rancher.Client, clustersConfig *clusters.ClusterConfig, corralPackages *corral.Packages) (*management.Cluster, error) {
	setLogrusFormatter()
	clusterName := namegen.AppendRandomString(rke1AirgapCustomCluster)
	quantityPerPool := []int32{}
	rolesPerPool := []string{}
	for _, pool := range clustersConfig.NodePools {
		var finalRoleCommand string
		if pool.NodeRoles.ControlPlane {
			finalRoleCommand += " --controlplane"
		}
		if pool.NodeRoles.Etcd {
			finalRoleCommand += " --etcd"
		}
		if pool.NodeRoles.Worker {
			finalRoleCommand += " --worker"
		}

		quantityPerPool = append(quantityPerPool, int32(pool.NodeRoles.Quantity))
		rolesPerPool = append(rolesPerPool, finalRoleCommand)
	}

	if clustersConfig.PSACT == string(provisioninginput.RancherBaseline) {
		err := clusters.CreateRancherBaselinePSACT(client, clustersConfig.PSACT)
		if err != nil {
			return nil, err
		}
	}

	cluster := clusters.NewRKE1ClusterConfig(clusterName, client, clustersConfig)
	clusterResp, err := shepherdclusters.CreateRKE1Cluster(client, cluster)
	if err != nil {
		return nil, err
	}

	client, err = client.ReLogin()
	if err != nil {
		return nil, err
	}

	customCluster, err := client.Management.Cluster.ByID(clusterResp.ID)
	if err != nil {
		return nil, err
	}

	token, err := tokenregistration.GetRegistrationToken(client, customCluster.ID)
	if err != nil {
		return nil, err
	}

	corralsArgs := []corral.Args{}

	logrus.Infof("Register Custom Cluster Through Corral")
	for poolIndex, poolRole := range rolesPerPool {
		// environment variables must be escaped inside original registration command
		escapedCommand := strings.Replace(token.NodeCommand, "\"", "\\\"", -1)

		corralsArgs = append(corralsArgs, corral.Args{
			Name:        namegen.AppendRandomString(rke1NodeCorralName),
			PackageName: corralPackages.CorralPackageImages[corralPackageAirgapCustomClusterName],
			Updates:     map[string]string{"registration_command": fmt.Sprintf("%s %s", escapedCommand, poolRole), "node_count": fmt.Sprint(quantityPerPool[poolIndex])},
		})
	}

	_, err = corral.CreateMultipleCorrals(client.Session, corralsArgs, corralPackages.HasDebug, corralPackages.HasCleanup)
	if err != nil {
		return nil, err
	}

	createdCluster, err := client.Management.Cluster.ByID(clusterResp.ID)
	return createdCluster, err
}

// CreateProvisioningAKSHostedCluster provisions an AKS cluster, then runs verify checks
func CreateProvisioningAKSHostedCluster(client *rancher.Client, aksClusterConfig aks.ClusterConfig) (*management.Cluster, error) {
	cloudCredentialConfig := cloudcredentials.LoadCloudCredential("azure")
	cloudCredential, err := azure.CreateAzureCloudCredentials(client, cloudCredentialConfig)
	if err != nil {
		return nil, err
	}

	clusterName := namegen.AppendRandomString("akshostcluster")
	clusterResp, err := aks.CreateAKSHostedCluster(client, clusterName, cloudCredential.Namespace+":"+cloudCredential.Name, aksClusterConfig, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if client.Flags.GetValue(environmentflag.UpdateClusterName) && updateConfig {
		pipeline.UpdateConfigClusterName(clusterName)
	}

	client, err = client.ReLogin()
	if err != nil {
		return nil, err
	}

	return client.Management.Cluster.ByID(clusterResp.ID)
}

// CreateProvisioningEKSHostedCluster provisions an EKS cluster, then runs verify checks
func CreateProvisioningEKSHostedCluster(client *rancher.Client, eksClusterConfig eks.ClusterConfig) (*management.Cluster, error) {
	cloudCredentialConfig := cloudcredentials.LoadCloudCredential("aws")
	cloudCredential, err := aws.CreateAWSCloudCredentials(client, cloudCredentialConfig)
	if err != nil {
		return nil, err
	}

	clusterName := namegen.AppendRandomString("ekshostcluster")
	clusterResp, err := eks.CreateEKSHostedCluster(client, clusterName, cloudCredential.Namespace+":"+cloudCredential.Name, eksClusterConfig, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if client.Flags.GetValue(environmentflag.UpdateClusterName) && updateConfig {
		pipeline.UpdateConfigClusterName(clusterName)
	}

	client, err = client.ReLogin()
	if err != nil {
		return nil, err
	}

	return client.Management.Cluster.ByID(clusterResp.ID)
}

// CreateProvisioningGKEHostedCluster provisions an GKE cluster, then runs verify checks
func CreateProvisioningGKEHostedCluster(client *rancher.Client, gkeClusterConfig gke.ClusterConfig) (*management.Cluster, error) {
	credentialSpec := cloudcredentials.LoadCloudCredential(provisioninginput.GoogleProviderName.String())
	cloudCredential, err := google.CreateGoogleCloudCredentials(client, credentialSpec)
	if err != nil {
		return nil, err
	}

	clusterName := namegen.AppendRandomString("gkehostcluster")
	clusterResp, err := gke.CreateGKEHostedCluster(client, clusterName, cloudCredential.Namespace+":"+cloudCredential.Name, gkeClusterConfig, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if client.Flags.GetValue(environmentflag.UpdateClusterName) && updateConfig {
		pipeline.UpdateConfigClusterName(clusterName)
	}

	client, err = client.ReLogin()
	if err != nil {
		return nil, err
	}

	return client.Management.Cluster.ByID(clusterResp.ID)
}

func setLogrusFormatter() {
	formatter := &logrus.TextFormatter{}
	formatter.DisableQuote = true
	logrus.SetFormatter(formatter)
}

// createRKE1RegistrationCommand is a helper for rke1 custom clusters to create the registration command with advanced options configured per node
func createRKE1RegistrationCommand(command, publicIP, privateIP string, nodePool provisioninginput.NodePools) string {
	if nodePool.SpecifyCustomPublicIP {
		command += fmt.Sprintf(" --address %s", publicIP)
	}
	if nodePool.SpecifyCustomPrivateIP {
		command += fmt.Sprintf(" --internal-address %s", privateIP)
	}
	if nodePool.CustomNodeNameSuffix != "" {
		command += fmt.Sprintf(" --node-name %s", namegen.AppendRandomString(nodePool.CustomNodeNameSuffix))
	}
	for labelKey, labelValue := range nodePool.NodeLabels {
		command += fmt.Sprintf(" --label %s=%s", labelKey, labelValue)
	}
	for _, taint := range nodePool.NodeTaints {
		command += fmt.Sprintf(" --taints %s=%s:%s", taint.Key, taint.Value, taint.Effect)
	}
	return command
}

// createRegistrationCommand is a helper for rke2/k3s custom clusters to create the registration command with advanced options configured per node
func createRegistrationCommand(command, publicIP, privateIP string, machinePool provisioninginput.MachinePools) string {
	if machinePool.SpecifyCustomPublicIP {
		command += fmt.Sprintf(" --address %s", publicIP)
	}
	if machinePool.SpecifyCustomPrivateIP {
		command += fmt.Sprintf(" --internal-address %s", privateIP)
	}
	if machinePool.CustomNodeNameSuffix != "" {
		command += fmt.Sprintf(" --node-name %s", namegen.AppendRandomString(machinePool.CustomNodeNameSuffix))
	}
	for labelKey, labelValue := range machinePool.NodeLabels {
		command += fmt.Sprintf(" --label %s=%s", labelKey, labelValue)
	}
	for _, taint := range machinePool.NodeTaints {
		command += fmt.Sprintf(" --taints %s=%s:%s", taint.Key, taint.Value, taint.Effect)
	}
	return command
}

// createWindowsRegistrationCommand is a helper for rke2 windows custom clusters to create the registration command with advanced options configured per node
func createWindowsRegistrationCommand(command, publicIP, privateIP string, machinePool provisioninginput.MachinePools) string {
	if machinePool.SpecifyCustomPublicIP {
		command += fmt.Sprintf(" -Address '%s'", publicIP)
	}
	if machinePool.SpecifyCustomPrivateIP {
		command += fmt.Sprintf(" -InternalAddress '%s'", privateIP)
	}
	if machinePool.CustomNodeNameSuffix != "" {
		command += fmt.Sprintf(" -NodeName '%s'", namegen.AppendRandomString(machinePool.CustomNodeNameSuffix))
	}
	// powershell requires only 1 flag per command, so we need to append the custom labels and taints together which is different from linux
	if len(machinePool.NodeLabels) > 0 {
		// there is an existing label for all windows nodes, so we need to insert the custom labels after the existing label
		labelIndex := strings.Index(command, " -Label '") + len(" -Label '")
		customLabels := ""
		for labelKey, labelValue := range machinePool.NodeLabels {
			customLabels += fmt.Sprintf("%s=%s,", labelKey, labelValue)
		}
		command = command[:labelIndex] + customLabels + command[labelIndex:]
	}
	if len(machinePool.NodeTaints) > 0 {
		var customTaints string
		for _, taint := range machinePool.NodeTaints {
			customTaints += fmt.Sprintf("%s=%s:%s,", taint.Key, taint.Value, taint.Effect)
		}
		wrappedTaints := fmt.Sprintf(" -Taint '%s'", customTaints)
		command += wrappedTaints
	}
	return command
}

// AddRKE2K3SCustomClusterNodes is a helper method that will add nodes to the custom RKE2/K3S custom cluster.
func AddRKE2K3SCustomClusterNodes(client *rancher.Client, cluster *v1.SteveAPIObject, nodes []*nodes.Node, rolesPerNode []string) error {
	clusterStatus := &apiv1.ClusterStatus{}
	err := v1.ConvertToK8sType(cluster.Status, clusterStatus)
	if err != nil {
		return err
	}

	token, err := tokenregistration.GetRegistrationToken(client, clusterStatus.ClusterName)
	if err != nil {
		return err
	}

	var command string
	for key, node := range nodes {
		logrus.Infof("Adding node %s to cluster %s", node.NodeID, cluster.Name)
		if strings.Contains(rolesPerNode[key], "windows") {
			command = fmt.Sprintf("powershell.exe %s -Address %s", token.InsecureWindowsNodeCommand, node.PublicIPAddress)
		} else {
			command = fmt.Sprintf("%s %s --address %s", token.InsecureNodeCommand, rolesPerNode[key], node.PublicIPAddress)
		}

		output, err := node.ExecuteCommand(command)
		if err != nil {
			return err
		}

		logrus.Infof(output)
	}

	err = kwait.PollUntilContextTimeout(context.TODO(), 500*time.Millisecond, defaults.ThirtyMinuteTimeout, true, func(ctx context.Context) (done bool, err error) {
		clusterResp, err := client.Steve.SteveType(stevetypes.Provisioning).ByID(cluster.ID)
		if err != nil {
			return false, err
		}

		if clusterResp.ObjectMeta.State.Name == active &&
			nodestat.AllMachineReady(client, cluster.ID, defaults.ThirtyMinuteTimeout) == nil {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteRKE2K3SCustomClusterNodes is a method that will delete nodes from the custom RKE2/K3S custom cluster.
func DeleteRKE2K3SCustomClusterNodes(client *rancher.Client, clusterID string, cluster *v1.SteveAPIObject, nodesToDelete []*nodes.Node) error {
	steveclient, err := client.Steve.ProxyDownstream(clusterID)
	if err != nil {
		return err
	}

	nodesSteveObjList, err := steveclient.SteveType("node").List(nil)
	if err != nil {
		return err
	}

	for _, nodeToDelete := range nodesToDelete {
		for _, node := range nodesSteveObjList.Data {
			snippedIP := strings.Split(node.Annotations[internalIP], ",")[0]

			if snippedIP == nodeToDelete.PrivateIPAddress {
				machine, err := client.Steve.SteveType(machineSteveResourceType).ByID(namespace + "/" + node.Annotations[machineNameAnnotation])
				if err != nil {
					return err
				}

				logrus.Infof("Deleting node %s from cluster %s", nodeToDelete.NodeID, cluster.Name)
				err = client.Steve.SteveType(machineSteveResourceType).Delete(machine)
				if err != nil {
					return err
				}

				err = kwait.PollUntilContextTimeout(context.TODO(), 500*time.Millisecond, defaults.ThirtyMinuteTimeout, true, func(ctx context.Context) (done bool, err error) {
					_, err = client.Steve.SteveType(machineSteveResourceType).ByID(machine.ID)
					if err != nil {
						logrus.Infof("Node has successfully been deleted!")
						return true, nil
					}
					return false, nil
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// AddRKE1CustomClusterNodes is a method that will add nodes to the custom RKE1 custom cluster.
func AddRKE1CustomClusterNodes(client *rancher.Client, cluster *management.Cluster, nodes []*nodes.Node, rolesPerNode []string) error {
	token, err := tokenregistration.GetRegistrationToken(client, cluster.ID)
	if err != nil {
		return err
	}

	var command string
	for key, node := range nodes {
		logrus.Infof("Adding node %s to cluster %s", node.NodeID, cluster.Name)
		command = fmt.Sprintf("%s %s --address %s", token.NodeCommand, rolesPerNode[key], node.PublicIPAddress)

		output, err := node.ExecuteCommand(command)
		if err != nil {
			return err
		}

		logrus.Infof(output)
	}

	err = kwait.PollUntilContextTimeout(context.TODO(), 500*time.Millisecond, defaults.ThirtyMinuteTimeout, true, func(ctx context.Context) (done bool, err error) {
		client, err = client.ReLogin()
		if err != nil {
			return false, err
		}

		clusterResp, err := client.Management.Cluster.ByID(cluster.ID)
		if err != nil {
			return false, err
		}

		if clusterResp.State == active &&
			nodestat.AllManagementNodeReady(client, cluster.ID, defaults.ThirtyMinuteTimeout) == nil {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteRKE1CustomClusterNodes is a helper method that will delete nodes from the custom RKE1 custom cluster.
func DeleteRKE1CustomClusterNodes(client *rancher.Client, cluster *management.Cluster, nodesToDelete []*nodes.Node) error {
	nodes, err := client.Management.Node.ListAll(&types.ListOpts{Filters: map[string]interface{}{
		"clusterId": cluster.ID,
	}})
	if err != nil {
		return err
	}

	for _, nodeToDelete := range nodesToDelete {
		for _, node := range nodes.Data {
			if node.Annotations[rke1ExternalIP] == nodeToDelete.PublicIPAddress {
				machine, err := client.Management.Node.ByID(node.ID)
				if err != nil {
					return err
				}

				logrus.Infof("Deleting node %s from cluster %s", nodeToDelete.NodeID, cluster.Name)
				err = client.Management.Node.Delete(machine)
				if err != nil {
					return err
				}

				err = kwait.PollUntilContextTimeout(context.TODO(), 500*time.Millisecond, defaults.ThirtyMinuteTimeout, true, func(ctx context.Context) (done bool, err error) {
					_, err = client.Management.Node.ByID(machine.ID)
					if err != nil {
						logrus.Infof("Node has successfully been deleted!")
						return true, nil
					}
					return false, nil
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// DisableUpdateConfig is a function that disable cattle config update and clean updateConfig for true to don't affect next tests.
func DisableUpdateConfig(client *rancher.Client) {
	updateConfig = false
	client.Session.RegisterCleanupFunc(func() error {
		updateConfig = true
		return nil
	})
}

// CreateProvisioningRKE1ClusterWithClusterTemplate provisions an rke1 cluster by using the rke1 template and revision ID and other values from the template.
func CreateProvisioningRKE1ClusterWithClusterTemplate(client *rancher.Client, templateID, revisionID string, nodesAndRoles []provisioninginput.NodePools, nodeTemplate *nodetemplates.NodeTemplate, answers *management.Answer) (*management.Cluster, error) {
	clusterName := namegen.AppendRandomString("rke1cluster-template-")

	rke1Cluster := &management.Cluster{
		DockerRootDir:                 "/var/lib/docker",
		Name:                          namegen.AppendRandomString("rketemplate-cluster-"),
		ClusterTemplateID:             templateID,
		ClusterTemplateRevisionID:     revisionID,
		ClusterTemplateAnswers:        answers,
		RancherKubernetesEngineConfig: nil,
	}
	clusterResp, err := shepherdclusters.CreateRKE1Cluster(client, rke1Cluster)
	if err != nil {
		return nil, err
	}

	if client.Flags.GetValue(environmentflag.UpdateClusterName) {
		pipeline.UpdateConfigClusterName(clusterName)
	}

	var nodeRoles []nodepools.NodeRoles
	for _, nodes := range nodesAndRoles {
		nodeRoles = append(nodeRoles, nodes.NodeRoles)
	}
	_, err = nodepools.NodePoolSetup(client, nodeRoles, clusterResp.ID, nodeTemplate.ID)
	if err != nil {
		return nil, err
	}

	createdCluster, err := client.Management.Cluster.ByID(clusterResp.ID)
	return createdCluster, err
}
package clusters

import (
	"fmt"
	"strings"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	apisV1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	provv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/rancher/tests/v2/actions/provisioninginput"
	"github.com/rancher/shepherd/clients/rancher"
	management "github.com/rancher/shepherd/clients/rancher/generated/management/v3"
	steveV1 "github.com/rancher/shepherd/clients/rancher/v1"
	v1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/clusters"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	baseline         = "baseline"
	externalAws      = "external-aws"
	etcdRole         = "etcd-role"
	controlPlaneRole = "control-plane-role"
	workerRole       = "worker-role"

	externalCloudProviderString = "cloud-provider=external"
	kubeletArgKey               = "kubelet-arg"
	kubeletAPIServerArgKey      = "kubeapi-server-arg"
	kubeControllerManagerArgKey = "kube-controller-manager-arg"
	cloudProviderAnnotationName = "cloud-provider-name"
	disableCloudController      = "disable-cloud-controller"
	protectKernelDefaults       = "protect-kernel-defaults"
	localcluster                = "fleet-local/local"
	rancherRestricted           = "rancher-restricted"
	rke1HardenedGID             = 52034
	rke1HardenedUID             = 52034
)

// CreateRancherBaselinePSACT creates custom PSACT called rancher-baseline which sets each PSS to baseline.
func CreateRancherBaselinePSACT(client *rancher.Client, psact string) error {
	_, err := client.Steve.SteveType(clusters.PodSecurityAdmissionSteveResoureType).ByID(psact)
	if err == nil {
		return err
	}

	template := &v3.PodSecurityAdmissionConfigurationTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name: psact,
		},
		Description: "This is a custom baseline Pod Security Admission Configuration Template. " +
			"It defines a minimally restrictive policy which prevents known privilege escalations. " +
			"This policy contains namespace level exemptions for Rancher components.",
		Configuration: v3.PodSecurityAdmissionConfigurationTemplateSpec{
			Defaults: v3.PodSecurityAdmissionConfigurationTemplateDefaults{
				Enforce: baseline,
				Audit:   baseline,
				Warn:    baseline,
			},
			Exemptions: v3.PodSecurityAdmissionConfigurationTemplateExemptions{
				Usernames:      []string{},
				RuntimeClasses: []string{},
				Namespaces: []string{
					"ingress-nginx",
					"kube-system",
					"cattle-system",
					"cattle-epinio-system",
					"cattle-fleet-system",
					"longhorn-system",
					"cattle-neuvector-system",
					"cattle-monitoring-system",
					"rancher-alerting-drivers",
					"cis-operator-system",
					"cattle-csp-adapter-system",
					"cattle-externalip-system",
					"cattle-gatekeeper-system",
					"istio-system",
					"cattle-istio-system",
					"cattle-logging-system",
					"cattle-windows-gmsa-system",
					"cattle-sriov-system",
					"cattle-ui-plugin-system",
					"tigera-operator",
				},
			},
		},
	}

	_, err = client.Steve.SteveType(clusters.PodSecurityAdmissionSteveResoureType).Create(template)
	if err != nil {
		return err
	}

	return nil
}

// NewRKE1ClusterConfig is a constructor for a v3.Cluster object, to be used by the rancher.Client.Provisioning client.
func NewRKE1ClusterConfig(clusterName string, client *rancher.Client, clustersConfig *ClusterConfig) *management.Cluster {
	backupConfigEnabled := true
	criDockerBool := false
	if clustersConfig.CRIDockerd {
		criDockerBool = true
	}
	newConfig := &management.Cluster{
		DockerRootDir: "/var/lib/docker",
		LocalClusterAuthEndpoint: &management.LocalClusterAuthEndpoint{
			Enabled: true,
		},
		Name: clusterName,
		RancherKubernetesEngineConfig: &management.RancherKubernetesEngineConfig{
			DNS: &management.DNSConfig{
				Provider: "coredns",
				Options: map[string]string{
					"stubDomains": "cluster.local",
				},
			},
			EnableCRIDockerd: &criDockerBool,
			Ingress: &management.IngressConfig{
				Provider: "nginx",
			},
			Monitoring: &management.MonitoringConfig{
				Provider: "metrics-server",
			},
			Network: &management.NetworkConfig{
				Plugin:  clustersConfig.CNI,
				MTU:     0,
				Options: map[string]string{},
			},
			Services: &management.RKEConfigServices{
				Etcd: &management.ETCDService{
					BackupConfig: &management.BackupConfig{
						Enabled:       &backupConfigEnabled,
						IntervalHours: 12,
						Retention:     6,
						SafeTimestamp: true,
						Timeout:       120,
					},
				},
			},
			Version: clustersConfig.KubernetesVersion,
		},
	}
	newConfig.ClusterAgentDeploymentCustomization = clustersConfig.ClusterAgent
	newConfig.FleetAgentDeploymentCustomization = clustersConfig.FleetAgent
	newConfig.AgentEnvVars = clustersConfig.AgentEnvVarsRKE1

	if clustersConfig.Registries != nil {
		if clustersConfig.Registries.RKE1Registries != nil {
			newConfig.RancherKubernetesEngineConfig.PrivateRegistries = clustersConfig.Registries.RKE1Registries
			for _, registry := range clustersConfig.Registries.RKE1Registries {
				if registry.ECRCredentialPlugin != nil {
					awsAccessKeyID := fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", registry.ECRCredentialPlugin.AwsAccessKeyID)
					awsSecretAccessKey := fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", registry.ECRCredentialPlugin.AwsSecretAccessKey)
					extraEnv := []string{awsAccessKeyID, awsSecretAccessKey}
					newConfig.RancherKubernetesEngineConfig.Services = &management.RKEConfigServices{
						Kubelet: &management.KubeletService{
							ExtraEnv: extraEnv,
						},
					}
					break
				}
			}
		}
	}

	if clustersConfig.CloudProvider != "" {
		newConfig.RancherKubernetesEngineConfig.CloudProvider = &management.CloudProvider{
			Name: clustersConfig.CloudProvider,
		}
		if clustersConfig.CloudProvider == externalAws {
			trueBoolean := true
			newConfig.RancherKubernetesEngineConfig.CloudProvider.UseInstanceMetadataHostname = &trueBoolean
		}
	}

	if clustersConfig.ETCDRKE1 != nil {
		newConfig.RancherKubernetesEngineConfig.Services.Etcd = clustersConfig.ETCDRKE1
	}

	if clustersConfig.PSACT != "" {
		newConfig.DefaultPodSecurityAdmissionConfigurationTemplateName = clustersConfig.PSACT
	}

	if clustersConfig.AgentEnvVars != nil {
		newConfig.AgentEnvVars = clustersConfig.AgentEnvVarsRKE1
	}

	return newConfig
}

// UpdateRKE1ClusterConfig is a constructor for a v3.Cluster object, to be used by the rancher.Client.Provisioning client.
func UpdateRKE1ClusterConfig(clusterName string, client *rancher.Client, clustersConfig *ClusterConfig) *management.Cluster {
	newConfig := &management.Cluster{
		Name: clusterName,
		RancherKubernetesEngineConfig: &management.RancherKubernetesEngineConfig{
			Network: &management.NetworkConfig{
				Plugin: clustersConfig.CNI,
			},
			Version: clustersConfig.KubernetesVersion,
		},
	}

	newConfig.ClusterAgentDeploymentCustomization = clustersConfig.ClusterAgent
	newConfig.FleetAgentDeploymentCustomization = clustersConfig.FleetAgent
	newConfig.AgentEnvVars = clustersConfig.AgentEnvVarsRKE1

	if clustersConfig.Registries != nil {
		if clustersConfig.Registries.RKE1Registries != nil {
			newConfig.RancherKubernetesEngineConfig.PrivateRegistries = clustersConfig.Registries.RKE1Registries
			for _, registry := range clustersConfig.Registries.RKE1Registries {
				if registry.ECRCredentialPlugin != nil {
					awsAccessKeyID := fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", registry.ECRCredentialPlugin.AwsAccessKeyID)
					awsSecretAccessKey := fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", registry.ECRCredentialPlugin.AwsSecretAccessKey)
					extraEnv := []string{awsAccessKeyID, awsSecretAccessKey}
					newConfig.RancherKubernetesEngineConfig.Services = &management.RKEConfigServices{
						Kubelet: &management.KubeletService{
							ExtraEnv: extraEnv,
						},
					}

					break
				}
			}
		}
	}

	if clustersConfig.CloudProvider != "" {
		newConfig.RancherKubernetesEngineConfig.CloudProvider = &management.CloudProvider{
			Name: clustersConfig.CloudProvider,
		}
		if clustersConfig.CloudProvider == externalAws {
			trueBoolean := true
			newConfig.RancherKubernetesEngineConfig.CloudProvider.UseInstanceMetadataHostname = &trueBoolean
		}
	}

	if clustersConfig.ETCDRKE1 != nil {
		newConfig.RancherKubernetesEngineConfig.Services.Etcd = clustersConfig.ETCDRKE1
	}

	if clustersConfig.AgentEnvVars != nil {
		newConfig.AgentEnvVars = clustersConfig.AgentEnvVarsRKE1
	}

	if clustersConfig.PSACT != "" {
		newConfig.DefaultPodSecurityAdmissionConfigurationTemplateName = clustersConfig.PSACT
	}

	return newConfig
}

// NewK3SRKE2ClusterConfig is a constructor for a apisV1.Cluster object, to be used by the rancher.Client.Provisioning client.
func NewK3SRKE2ClusterConfig(clusterName, namespace string, clustersConfig *ClusterConfig, machinePools []apisV1.RKEMachinePool, cloudCredentialSecretName string) *apisV1.Cluster {
	typeMeta := metav1.TypeMeta{
		Kind:       "Cluster",
		APIVersion: "provisioning.cattle.io/v1",
	}

	//metav1.ObjectMeta
	objectMeta := metav1.ObjectMeta{
		Name:      clusterName,
		Namespace: namespace,
	}

	etcd := &rkev1.ETCD{
		SnapshotRetention:    5,
		SnapshotScheduleCron: "0 */5 * * *",
	}
	if clustersConfig.ETCD != nil {
		etcd = clustersConfig.ETCD
	}

	chartValuesMap := rkev1.GenericMap{
		Data: map[string]interface{}{},
	}
	chartAdditionalManifest := ""
	if clustersConfig.AddOnConfig != nil {
		if clustersConfig.AddOnConfig.ChartValues != nil {
			chartValuesMap = *clustersConfig.AddOnConfig.ChartValues
		}
		chartAdditionalManifest = clustersConfig.AddOnConfig.AdditionalManifest
	}

	machineGlobalConfigMap := rkev1.GenericMap{
		Data: map[string]interface{}{
			"cni":                 clustersConfig.CNI,
			"disable-kube-proxy":  false,
			"etcd-expose-metrics": false,
			"profile":             nil,
		},
	}
	machineSelectorConfigs := []rkev1.RKESystemConfig{}
	if clustersConfig.Advanced != nil {
		if clustersConfig.Advanced.MachineGlobalConfig != nil {
			for k, v := range clustersConfig.Advanced.MachineGlobalConfig.Data {
				machineGlobalConfigMap.Data[k] = v
			}
		}

		if clustersConfig.Advanced.MachineSelectors != nil {
			machineSelectorConfigs = *clustersConfig.Advanced.MachineSelectors
		}
	}

	localClusterAuthEndpoint := rkev1.LocalClusterAuthEndpoint{
		CACerts: "",
		Enabled: false,
		FQDN:    "",
	}
	if clustersConfig.Networking != nil {
		if clustersConfig.Networking.LocalClusterAuthEndpoint != nil {
			localClusterAuthEndpoint = *clustersConfig.Networking.LocalClusterAuthEndpoint
		}
	}

	upgradeStrategy := rkev1.ClusterUpgradeStrategy{
		ControlPlaneConcurrency:  "10%",
		ControlPlaneDrainOptions: rkev1.DrainOptions{},
		WorkerConcurrency:        "10%",
		WorkerDrainOptions:       rkev1.DrainOptions{},
	}
	if clustersConfig.UpgradeStrategy != nil {
		upgradeStrategy = *clustersConfig.UpgradeStrategy
	}

	clusterAgentDeploymentCustomization := &apisV1.AgentDeploymentCustomization{}
	if clustersConfig.ClusterAgent != nil {
		clusterAgentOverrides := ResourceConfigHelper(clustersConfig.ClusterAgent.OverrideResourceRequirements)
		clusterAgentDeploymentCustomization.OverrideResourceRequirements = clusterAgentOverrides
		v1ClusterTolerations := []corev1.Toleration{}
		for _, t := range clustersConfig.ClusterAgent.AppendTolerations {
			v1ClusterTolerations = append(v1ClusterTolerations, corev1.Toleration{
				Key:      t.Key,
				Operator: corev1.TolerationOperator(t.Operator),
				Value:    t.Value,
				Effect:   corev1.TaintEffect(t.Effect),
			})
		}
		clusterAgentDeploymentCustomization.AppendTolerations = v1ClusterTolerations
		clusterAgentDeploymentCustomization.OverrideAffinity = AgentAffinityConfigHelper(clustersConfig.ClusterAgent.OverrideAffinity)
	}

	fleetAgentDeploymentCustomization := &apisV1.AgentDeploymentCustomization{}
	if clustersConfig.FleetAgent != nil {
		fleetAgentOverrides := ResourceConfigHelper(clustersConfig.FleetAgent.OverrideResourceRequirements)
		fleetAgentDeploymentCustomization.OverrideResourceRequirements = fleetAgentOverrides
		v1FleetTolerations := []corev1.Toleration{}
		for _, t := range clustersConfig.FleetAgent.AppendTolerations {
			v1FleetTolerations = append(v1FleetTolerations, corev1.Toleration{
				Key:      t.Key,
				Operator: corev1.TolerationOperator(t.Operator),
				Value:    t.Value,
				Effect:   corev1.TaintEffect(t.Effect),
			})
		}
		fleetAgentDeploymentCustomization.AppendTolerations = v1FleetTolerations
		fleetAgentDeploymentCustomization.OverrideAffinity = AgentAffinityConfigHelper(clustersConfig.FleetAgent.OverrideAffinity)
	}
	var registries *rkev1.Registry
	if clustersConfig.Registries != nil {
		registries = clustersConfig.Registries.RKE2Registries
	}

	if clustersConfig.CloudProvider == provisioninginput.AWSProviderName.String() {
		machineSelectorConfigs = append(machineSelectorConfigs, OutOfTreeSystemConfig(clustersConfig.CloudProvider)...)
	} else if strings.Contains(clustersConfig.CloudProvider, "-in-tree") {
		machineSelectorConfigs = append(machineSelectorConfigs, InTreeSystemConfig(strings.Split(clustersConfig.CloudProvider, "-in-tree")[0])...)
	}

	if clustersConfig.CloudProvider == provisioninginput.VsphereCloudProviderName.String() {
		machineSelectorConfigs = append(machineSelectorConfigs,
			RKESystemConfigTemplate(map[string]interface{}{
				cloudProviderAnnotationName: provisioninginput.VsphereCloudProviderName.String(),
				protectKernelDefaults:       false,
			},
				nil),
		)
	}

	rkeSpecCommon := rkev1.RKEClusterSpecCommon{
		UpgradeStrategy:       upgradeStrategy,
		ChartValues:           chartValuesMap,
		MachineGlobalConfig:   machineGlobalConfigMap,
		MachineSelectorConfig: machineSelectorConfigs,
		AdditionalManifest:    chartAdditionalManifest,
		Registries:            registries,
		ETCD:                  etcd,
	}
	rkeConfig := &apisV1.RKEConfig{
		RKEClusterSpecCommon: rkeSpecCommon,
		MachinePools:         machinePools,
	}

	agentEnvVars := []rkev1.EnvVar{}

	spec := apisV1.ClusterSpec{
		CloudCredentialSecretName:           cloudCredentialSecretName,
		KubernetesVersion:                   clustersConfig.KubernetesVersion,
		LocalClusterAuthEndpoint:            localClusterAuthEndpoint,
		RKEConfig:                           rkeConfig,
		ClusterAgentDeploymentCustomization: clusterAgentDeploymentCustomization,
		FleetAgentDeploymentCustomization:   fleetAgentDeploymentCustomization,
		AgentEnvVars:                        agentEnvVars,
	}

	if clustersConfig.AgentEnvVars != nil {
		spec.AgentEnvVars = clustersConfig.AgentEnvVars
	}

	if clustersConfig.PSACT != "" {
		spec.DefaultPodSecurityAdmissionConfigurationTemplateName = clustersConfig.PSACT
	}

	v1Cluster := &apisV1.Cluster{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
		Spec:       spec,
	}

	return v1Cluster
}

// UpdateK3SRKE2ClusterConfig is a constructor for a apisV1.Cluster object, to be used by the rancher.Client.Provisioning client.
func UpdateK3SRKE2ClusterConfig(cluster *v1.SteveAPIObject, clustersConfig *ClusterConfig) *v1.SteveAPIObject {
	clusterSpec := &provv1.ClusterSpec{}
	err := steveV1.ConvertToK8sType(cluster.Spec, clusterSpec)
	if err != nil {
		return nil
	}

	if clustersConfig.ETCD != nil {
		clusterSpec.RKEConfig.ETCD = clustersConfig.ETCD
	}

	if clustersConfig.KubernetesVersion != "" {
		clusterSpec.KubernetesVersion = clustersConfig.KubernetesVersion
	}

	if clustersConfig.CNI != "" {
		clusterSpec.RKEConfig.MachineGlobalConfig.Data["cni"] = clustersConfig.CNI
	}

	if clustersConfig.AddOnConfig != nil {
		if clustersConfig.AddOnConfig.ChartValues != nil {
			clusterSpec.RKEConfig.ChartValues = *clustersConfig.AddOnConfig.ChartValues
		}
		clusterSpec.RKEConfig.AdditionalManifest = clustersConfig.AddOnConfig.AdditionalManifest
	}

	if clustersConfig.Advanced != nil {
		if clustersConfig.Advanced.MachineGlobalConfig != nil {
			for k, v := range clustersConfig.Advanced.MachineGlobalConfig.Data {
				clusterSpec.RKEConfig.MachineGlobalConfig.Data[k] = v
			}
		}

		if clustersConfig.Advanced.MachineSelectors != nil {
			clusterSpec.RKEConfig.MachineSelectorConfig = *clustersConfig.Advanced.MachineSelectors
		}
	}

	if clustersConfig.Networking != nil {
		if clustersConfig.Networking.LocalClusterAuthEndpoint != nil {
			clusterSpec.LocalClusterAuthEndpoint = *clustersConfig.Networking.LocalClusterAuthEndpoint
		}
	}

	if clustersConfig.UpgradeStrategy != nil {
		clusterSpec.RKEConfig.UpgradeStrategy = *clustersConfig.UpgradeStrategy
	}

	if clustersConfig.ClusterAgent != nil {
		clusterAgentOverrides := ResourceConfigHelper(clustersConfig.ClusterAgent.OverrideResourceRequirements)
		clusterSpec.ClusterAgentDeploymentCustomization.OverrideResourceRequirements = clusterAgentOverrides
		v1ClusterTolerations := []corev1.Toleration{}
		for _, t := range clustersConfig.ClusterAgent.AppendTolerations {
			v1ClusterTolerations = append(v1ClusterTolerations, corev1.Toleration{
				Key:      t.Key,
				Operator: corev1.TolerationOperator(t.Operator),
				Value:    t.Value,
				Effect:   corev1.TaintEffect(t.Effect),
			})
		}
		clusterSpec.ClusterAgentDeploymentCustomization.AppendTolerations = v1ClusterTolerations
		clusterSpec.ClusterAgentDeploymentCustomization.OverrideAffinity = AgentAffinityConfigHelper(clustersConfig.ClusterAgent.OverrideAffinity)
	}

	if clustersConfig.FleetAgent != nil {
		fleetAgentOverrides := ResourceConfigHelper(clustersConfig.FleetAgent.OverrideResourceRequirements)
		clusterSpec.ClusterAgentDeploymentCustomization.OverrideResourceRequirements = fleetAgentOverrides
		v1FleetTolerations := []corev1.Toleration{}
		for _, t := range clustersConfig.FleetAgent.AppendTolerations {
			v1FleetTolerations = append(v1FleetTolerations, corev1.Toleration{
				Key:      t.Key,
				Operator: corev1.TolerationOperator(t.Operator),
				Value:    t.Value,
				Effect:   corev1.TaintEffect(t.Effect),
			})
		}
		clusterSpec.FleetAgentDeploymentCustomization.AppendTolerations = v1FleetTolerations
		clusterSpec.FleetAgentDeploymentCustomization.OverrideAffinity = AgentAffinityConfigHelper(clustersConfig.FleetAgent.OverrideAffinity)
	}

	if clustersConfig.Registries != nil {
		clusterSpec.RKEConfig.Registries = clustersConfig.Registries.RKE2Registries
	}

	if clustersConfig.CloudProvider == provisioninginput.AWSProviderName.String() {
		clusterSpec.RKEConfig.MachineSelectorConfig = append(clusterSpec.RKEConfig.MachineSelectorConfig, OutOfTreeSystemConfig(clustersConfig.CloudProvider)...)
	} else if strings.Contains(clustersConfig.CloudProvider, "-in-tree") {
		clusterSpec.RKEConfig.MachineSelectorConfig = append(clusterSpec.RKEConfig.MachineSelectorConfig, InTreeSystemConfig(strings.Split(clustersConfig.CloudProvider, "-in-tree")[0])...)
	}

	if clustersConfig.CloudProvider == provisioninginput.VsphereCloudProviderName.String() {
		clusterSpec.RKEConfig.MachineSelectorConfig = append(clusterSpec.RKEConfig.MachineSelectorConfig,
			RKESystemConfigTemplate(map[string]interface{}{
				cloudProviderAnnotationName: provisioninginput.VsphereCloudProviderName.String(),
				protectKernelDefaults:       false,
			},
				nil),
		)
	}

	if clustersConfig.AgentEnvVars != nil {
		clusterSpec.AgentEnvVars = clustersConfig.AgentEnvVars
	}

	if clustersConfig.PSACT != "" {
		clusterSpec.DefaultPodSecurityAdmissionConfigurationTemplateName = clustersConfig.PSACT
	}

	cluster.Spec = clusterSpec

	return cluster
}

// ResourceConfigHelper is a "helper" function that is used to convert the management.ResourceRequirements struct
// to a corev1.ResourceRequirements struct.
func ResourceConfigHelper(advancedClusterResourceRequirements *management.ResourceRequirements) *corev1.ResourceRequirements {
	agentOverrides := corev1.ResourceRequirements{}
	agentOverrides.Limits = corev1.ResourceList{}
	agentOverrides.Requests = corev1.ResourceList{}
	if advancedClusterResourceRequirements.Limits[string(corev1.ResourceCPU)] != "" {
		agentOverrides.Limits[corev1.ResourceCPU] = resource.MustParse(advancedClusterResourceRequirements.Limits[string(corev1.ResourceCPU)])
	}
	if advancedClusterResourceRequirements.Limits[string(corev1.ResourceMemory)] != "" {
		agentOverrides.Limits[corev1.ResourceMemory] = resource.MustParse(advancedClusterResourceRequirements.Limits[string(corev1.ResourceMemory)])
	}
	if advancedClusterResourceRequirements.Requests[string(corev1.ResourceCPU)] != "" {
		agentOverrides.Requests[corev1.ResourceCPU] = resource.MustParse(advancedClusterResourceRequirements.Requests[string(corev1.ResourceCPU)])
	}
	if advancedClusterResourceRequirements.Requests[string(corev1.ResourceMemory)] != "" {
		agentOverrides.Requests[corev1.ResourceMemory] = resource.MustParse(advancedClusterResourceRequirements.Requests[string(corev1.ResourceMemory)])
	}
	return &agentOverrides
}

// AgentAffinityConfigHelper is a "helper" function that converts a management.Affinity struct and returns a corev1.Affinity struct.
func AgentAffinityConfigHelper(advancedClusterAffinity *management.Affinity) *corev1.Affinity {
	agentAffinity := &corev1.Affinity{}
	if advancedClusterAffinity.NodeAffinity != nil {
		agentAffinity.NodeAffinity = &corev1.NodeAffinity{}
		if advancedClusterAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
			agentAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &corev1.NodeSelector{}
			agentAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = []corev1.NodeSelectorTerm{}
			for _, term := range advancedClusterAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms {
				agentMatchExpressions := []corev1.NodeSelectorRequirement{}
				if term.MatchExpressions != nil {
					for _, match := range term.MatchExpressions {
						newMatchExpression := corev1.NodeSelectorRequirement{}
						newMatchExpression.Key = match.Key
						newMatchExpression.Operator = corev1.NodeSelectorOperator(match.Operator)
						newMatchExpression.Values = match.Values
						agentMatchExpressions = append(agentMatchExpressions, newMatchExpression)
					}
				}
				agentMatchFields := []corev1.NodeSelectorRequirement{}
				if term.MatchFields != nil {
					for _, match := range term.MatchFields {
						newMatchExpression := corev1.NodeSelectorRequirement{}
						newMatchExpression.Key = match.Key
						newMatchExpression.Operator = corev1.NodeSelectorOperator(match.Operator)
						newMatchExpression.Values = match.Values
						agentMatchFields = append(agentMatchFields, newMatchExpression)
					}
				}
				agentAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = append(agentAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms, corev1.NodeSelectorTerm{
					MatchExpressions: agentMatchExpressions,
					MatchFields:      agentMatchFields,
				})
			}
		}
		if advancedClusterAffinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
			agentAffinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []corev1.PreferredSchedulingTerm{}
			for _, preferred := range advancedClusterAffinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
				termPreferences := corev1.NodeSelectorTerm{}
				if preferred.Preference.MatchExpressions != nil {
					termPreferences.MatchExpressions = []corev1.NodeSelectorRequirement{}
					for _, match := range preferred.Preference.MatchExpressions {
						newMatchExpression := corev1.NodeSelectorRequirement{}
						newMatchExpression.Key = match.Key
						newMatchExpression.Operator = corev1.NodeSelectorOperator(match.Operator)
						newMatchExpression.Values = match.Values
						termPreferences.MatchExpressions = append(termPreferences.MatchExpressions, newMatchExpression)
					}
				}
				if preferred.Preference.MatchFields != nil {
					termPreferences.MatchFields = []corev1.NodeSelectorRequirement{}
					for _, match := range preferred.Preference.MatchFields {
						newMatchExpression := corev1.NodeSelectorRequirement{}
						newMatchExpression.Key = match.Key
						newMatchExpression.Operator = corev1.NodeSelectorOperator(match.Operator)
						newMatchExpression.Values = match.Values
						termPreferences.MatchFields = append(termPreferences.MatchFields, newMatchExpression)
					}
				}
				agentAffinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(agentAffinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution, corev1.PreferredSchedulingTerm{
					Weight:     int32(preferred.Weight),
					Preference: termPreferences,
				})
			}
		}
	}
	if advancedClusterAffinity.PodAffinity != nil {
		agentAffinity.PodAffinity = &corev1.PodAffinity{}
		if advancedClusterAffinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
			agentAffinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []corev1.PodAffinityTerm{}
			for _, term := range advancedClusterAffinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution {
				matchExpressions := []metav1.LabelSelectorRequirement{}
				if term.LabelSelector != nil {
					for _, expression := range term.LabelSelector.MatchExpressions {
						newExpression := metav1.LabelSelectorRequirement{}
						newExpression.Key = expression.Key
						newExpression.Operator = metav1.LabelSelectorOperator(expression.Operator)
						newExpression.Values = expression.Values
						matchExpressions = append(matchExpressions, newExpression)
					}
				}
				matchNamespaces := metav1.LabelSelector{}
				if term.NamespaceSelector != nil {
					if term.NamespaceSelector.MatchLabels != nil {
						matchNamespaces.MatchLabels = term.NamespaceSelector.MatchLabels
					}
					for _, expression := range term.NamespaceSelector.MatchExpressions {
						newExpression := metav1.LabelSelectorRequirement{}
						newExpression.Key = expression.Key
						newExpression.Operator = metav1.LabelSelectorOperator(expression.Operator)
						newExpression.Values = expression.Values
						matchNamespaces.MatchExpressions = append(matchNamespaces.MatchExpressions, newExpression)
					}
				}
				newAffinityTerms := corev1.PodAffinityTerm{
					TopologyKey: term.TopologyKey,
				}
				if len(term.Namespaces) > 0 {
					newAffinityTerms.Namespaces = term.Namespaces
				}
				if term.LabelSelector != nil {
					newAffinityTerms.LabelSelector = &metav1.LabelSelector{}
					if term.LabelSelector.MatchLabels != nil {
						newAffinityTerms.LabelSelector.MatchLabels = term.LabelSelector.MatchLabels
					}
					if len(matchExpressions) > 0 {
						newAffinityTerms.LabelSelector.MatchExpressions = matchExpressions
					}
				}
				if matchNamespaces.MatchLabels != nil || len(matchNamespaces.MatchExpressions) > 0 {
					newAffinityTerms.NamespaceSelector = &matchNamespaces
				}
				agentAffinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(agentAffinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution, newAffinityTerms)
			}
		}
		if advancedClusterAffinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
			agentAffinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []corev1.WeightedPodAffinityTerm{}
			for _, preferred := range advancedClusterAffinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
				matchExpressions := []metav1.LabelSelectorRequirement{}
				if preferred.PodAffinityTerm.LabelSelector != nil {
					for _, expression := range preferred.PodAffinityTerm.LabelSelector.MatchExpressions {
						newExpression := metav1.LabelSelectorRequirement{}
						newExpression.Key = expression.Key
						newExpression.Operator = metav1.LabelSelectorOperator(expression.Operator)
						newExpression.Values = expression.Values
						matchExpressions = append(matchExpressions, newExpression)
					}
				}
				matchNamespaces := metav1.LabelSelector{}
				if preferred.PodAffinityTerm.NamespaceSelector != nil {
					if preferred.PodAffinityTerm.NamespaceSelector.MatchLabels == nil {
						matchNamespaces.MatchLabels = preferred.PodAffinityTerm.NamespaceSelector.MatchLabels
					}
					for _, expression := range preferred.PodAffinityTerm.NamespaceSelector.MatchExpressions {
						newExpression := metav1.LabelSelectorRequirement{}
						newExpression.Key = expression.Key
						newExpression.Operator = metav1.LabelSelectorOperator(expression.Operator)
						newExpression.Values = expression.Values
						matchNamespaces.MatchExpressions = append(matchNamespaces.MatchExpressions, newExpression)
					}
				}
				newAffinityTerms := corev1.WeightedPodAffinityTerm{
					Weight: int32(preferred.Weight),
					PodAffinityTerm: corev1.PodAffinityTerm{
						TopologyKey: preferred.PodAffinityTerm.TopologyKey,
					},
				}
				// add in optional variables if they exist
				if preferred.PodAffinityTerm.Namespaces != nil {
					newAffinityTerms.PodAffinityTerm.Namespaces = preferred.PodAffinityTerm.Namespaces
				}
				if matchNamespaces.MatchLabels != nil || matchNamespaces.MatchExpressions != nil {
					newAffinityTerms.PodAffinityTerm.NamespaceSelector = &matchNamespaces
				}
				if preferred.PodAffinityTerm.LabelSelector != nil {
					newAffinityTerms.PodAffinityTerm.LabelSelector = &metav1.LabelSelector{}
					if preferred.PodAffinityTerm.LabelSelector.MatchLabels != nil {
						newAffinityTerms.PodAffinityTerm.LabelSelector.MatchLabels = preferred.PodAffinityTerm.LabelSelector.MatchLabels
					}
					if preferred.PodAffinityTerm.LabelSelector.MatchExpressions != nil {
						newAffinityTerms.PodAffinityTerm.LabelSelector.MatchExpressions = matchExpressions
					}
				}
				agentAffinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(agentAffinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution, newAffinityTerms)
			}
		}
	}
	if advancedClusterAffinity.PodAntiAffinity != nil {
		agentAffinity.PodAntiAffinity = &corev1.PodAntiAffinity{}
		if advancedClusterAffinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
			agentAffinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []corev1.PodAffinityTerm{}
			for _, term := range advancedClusterAffinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution {
				matchExpressions := []metav1.LabelSelectorRequirement{}
				if term.LabelSelector != nil {
					for _, expression := range term.LabelSelector.MatchExpressions {
						newExpression := metav1.LabelSelectorRequirement{}
						newExpression.Key = expression.Key
						newExpression.Operator = metav1.LabelSelectorOperator(expression.Operator)
						newExpression.Values = expression.Values
						matchExpressions = append(matchExpressions, newExpression)
					}
				}
				matchNamespaces := metav1.LabelSelector{}
				if term.NamespaceSelector != nil {
					if term.NamespaceSelector.MatchLabels != nil {
						matchNamespaces.MatchLabels = term.NamespaceSelector.MatchLabels
					}
					for _, expression := range term.NamespaceSelector.MatchExpressions {
						newExpression := metav1.LabelSelectorRequirement{}
						newExpression.Key = expression.Key
						newExpression.Operator = metav1.LabelSelectorOperator(expression.Operator)
						newExpression.Values = expression.Values
						matchNamespaces.MatchExpressions = append(matchNamespaces.MatchExpressions, newExpression)
					}
				}
				newAffinityTerms := corev1.PodAffinityTerm{
					TopologyKey: term.TopologyKey,
				}
				if len(term.Namespaces) > 0 {
					newAffinityTerms.Namespaces = term.Namespaces
				}
				if term.LabelSelector != nil {
					newAffinityTerms.LabelSelector = &metav1.LabelSelector{}
					if term.LabelSelector.MatchLabels != nil {
						newAffinityTerms.LabelSelector.MatchLabels = term.LabelSelector.MatchLabels
					}
					if len(matchExpressions) > 0 {
						newAffinityTerms.LabelSelector.MatchExpressions = matchExpressions
					}
				}
				if matchNamespaces.MatchLabels != nil || len(matchNamespaces.MatchExpressions) > 0 {
					newAffinityTerms.NamespaceSelector = &matchNamespaces
				}
				agentAffinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(agentAffinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution, newAffinityTerms)
			}
		}
		if advancedClusterAffinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
			agentAffinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []corev1.WeightedPodAffinityTerm{}
			for _, preferred := range advancedClusterAffinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
				matchExpressions := []metav1.LabelSelectorRequirement{}
				if preferred.PodAffinityTerm.LabelSelector != nil {
					for _, expression := range preferred.PodAffinityTerm.LabelSelector.MatchExpressions {
						newExpression := metav1.LabelSelectorRequirement{}
						newExpression.Key = expression.Key
						newExpression.Operator = metav1.LabelSelectorOperator(expression.Operator)
						newExpression.Values = expression.Values
						matchExpressions = append(matchExpressions, newExpression)
					}
				}
				matchNamespaces := metav1.LabelSelector{}
				if preferred.PodAffinityTerm.NamespaceSelector != nil {
					if preferred.PodAffinityTerm.NamespaceSelector.MatchLabels == nil {
						matchNamespaces.MatchLabels = preferred.PodAffinityTerm.NamespaceSelector.MatchLabels
					}
					for _, expression := range preferred.PodAffinityTerm.NamespaceSelector.MatchExpressions {
						newExpression := metav1.LabelSelectorRequirement{}
						newExpression.Key = expression.Key
						newExpression.Operator = metav1.LabelSelectorOperator(expression.Operator)
						newExpression.Values = expression.Values
						matchNamespaces.MatchExpressions = append(matchNamespaces.MatchExpressions, newExpression)
					}
				}
				newAffinityTerms := corev1.WeightedPodAffinityTerm{
					Weight: int32(preferred.Weight),
					PodAffinityTerm: corev1.PodAffinityTerm{
						TopologyKey: preferred.PodAffinityTerm.TopologyKey,
					},
				}
				// add in optional variables if they exist
				if preferred.PodAffinityTerm.Namespaces != nil {
					newAffinityTerms.PodAffinityTerm.Namespaces = preferred.PodAffinityTerm.Namespaces
				}
				if matchNamespaces.MatchLabels != nil || matchNamespaces.MatchExpressions != nil {
					newAffinityTerms.PodAffinityTerm.NamespaceSelector = &matchNamespaces
				}
				if preferred.PodAffinityTerm.LabelSelector != nil {
					newAffinityTerms.PodAffinityTerm.LabelSelector = &metav1.LabelSelector{}
					if preferred.PodAffinityTerm.LabelSelector.MatchLabels != nil {
						newAffinityTerms.PodAffinityTerm.LabelSelector.MatchLabels = preferred.PodAffinityTerm.LabelSelector.MatchLabels
					}
					if preferred.PodAffinityTerm.LabelSelector.MatchExpressions != nil {
						newAffinityTerms.PodAffinityTerm.LabelSelector.MatchExpressions = matchExpressions
					}
				}
				agentAffinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(agentAffinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution, newAffinityTerms)
			}
		}
	}
	return agentAffinity
}

// HardenK3SClusterConfig is a function that modifies the cluster configuration to be hardened according to the CIS benchmark.
func HardenK3SClusterConfig(clusterName, namespace string, clustersConfig *ClusterConfig, machinePools []apisV1.RKEMachinePool, cloudCredentialSecretName string) *apisV1.Cluster {
	v1Cluster := NewK3SRKE2ClusterConfig(clusterName, namespace, clustersConfig, machinePools, cloudCredentialSecretName)
	v1Cluster.Spec.DefaultPodSecurityAdmissionConfigurationTemplateName = rancherRestricted

	v1Cluster.Spec.RKEConfig.MachineGlobalConfig.Data["kube-apiserver-arg"] = []string{
		"audit-policy-file=/var/lib/rancher/k3s/server/audit.yaml",
		"audit-log-path=/var/lib/rancher/k3s/server/logs/audit.log",
		"audit-log-maxage=30",
		"audit-log-maxbackup=10",
		"audit-log-maxsize=100",
		"request-timeout=300s",
		"service-account-lookup=true",
	}

	v1Cluster.Spec.RKEConfig.MachineSelectorConfig = []rkev1.RKESystemConfig{
		{
			Config: rkev1.GenericMap{
				Data: map[string]interface{}{
					"kubelet-arg": []string{
						"make-iptables-util-chains=true",
					},
					protectKernelDefaults: true,
				},
			},
		},
	}

	return v1Cluster
}

// HardenRKE1ClusterConfig is a function that modifies the cluster configuration to be hardened according to the CIS benchmark.
func HardenRKE1ClusterConfig(client *rancher.Client, clusterName string, clustersConfig *ClusterConfig) *management.Cluster {
	cluster := NewRKE1ClusterConfig(clusterName, client, clustersConfig)

	cluster.DefaultPodSecurityAdmissionConfigurationTemplateName = rancherRestricted
	cluster.RancherKubernetesEngineConfig.Services.Etcd.GID = rke1HardenedGID
	cluster.RancherKubernetesEngineConfig.Services.Etcd.UID = rke1HardenedUID

	return cluster
}

// HardenRKE2ClusterConfig is a function that modifies the cluster configuration to be hardened according to the CIS benchmark.
func HardenRKE2ClusterConfig(clusterName, namespace string, clustersConfig *ClusterConfig, machinePools []apisV1.RKEMachinePool, cloudCredentialSecretName string) *apisV1.Cluster {
	v1Cluster := NewK3SRKE2ClusterConfig(clusterName, namespace, clustersConfig, machinePools, cloudCredentialSecretName)
	v1Cluster.Spec.DefaultPodSecurityAdmissionConfigurationTemplateName = rancherRestricted

	v1Cluster.Spec.RKEConfig.MachineSelectorConfig = []rkev1.RKESystemConfig{
		{
			Config: rkev1.GenericMap{
				Data: map[string]interface{}{
					"profile":             "cis",
					protectKernelDefaults: true,
				},
			},
		},
	}

	return v1Cluster
}

// CheckServiceAccountTokenSecret verifies if a serviceAccountTokenSecret exists or not in the cluster.
func CheckServiceAccountTokenSecret(client *rancher.Client, clusterName string) (success bool, err error) {
	clusterID, err := clusters.GetClusterIDByName(client, clusterName)
	if err != nil {
		return false, err
	}

	cluster, err := client.Management.Cluster.ByID(clusterID)
	if err != nil {
		return false, err
	}

	if cluster.ServiceAccountTokenSecret == "" {
		logrus.Warn("warning: serviceAccountTokenSecret does not exist in this cluster!")
		return false, nil
	}

	logrus.Infof("serviceAccountTokenSecret in this cluster is: %s", cluster.ServiceAccountTokenSecret)
	return true, nil
}

// OutOfTreeSystemConfig constructs the proper rkeSystemConfig slice for enabling the aws cloud provider
// out-of-tree services
func OutOfTreeSystemConfig(providerName string) (rkeConfig []rkev1.RKESystemConfig) {
	roles := []string{etcdRole, controlPlaneRole, workerRole}

	for _, role := range roles {
		selector := &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"rke.cattle.io/" + role: "true",
			},
		}
		configData := map[string]interface{}{}

		configData[kubeletArgKey] = []string{externalCloudProviderString}

		if role == controlPlaneRole {
			configData[kubeletAPIServerArgKey] = []string{externalCloudProviderString}
			configData[kubeControllerManagerArgKey] = []string{externalCloudProviderString}
		}

		if role == workerRole || role == controlPlaneRole {
			configData[disableCloudController] = true
		}

		rkeConfig = append(rkeConfig, RKESystemConfigTemplate(configData, selector))
	}

	configData := map[string]interface{}{
		cloudProviderAnnotationName: providerName,
		protectKernelDefaults:       false,
	}

	rkeConfig = append(rkeConfig, RKESystemConfigTemplate(configData, nil))
	return
}

// InTreeSystemConfig constructs the proper rkeSystemConfig slice for enabling cloud provider
// in-tree services.
// Vsphere deprecated 1.21+
// AWS deprecated 1.27+
// Azure deprecated 1.28+
func InTreeSystemConfig(providerName string) (rkeConfig []rkev1.RKESystemConfig) {
	configData := map[string]interface{}{
		cloudProviderAnnotationName: providerName,
		protectKernelDefaults:       false,
	}
	rkeConfig = append(rkeConfig, RKESystemConfigTemplate(configData, nil))
	return
}

// RKESYstemConfigTemplate constructs an RKESystemConfig object given config data and a selector
func RKESystemConfigTemplate(config map[string]interface{}, selector *metav1.LabelSelector) rkev1.RKESystemConfig {
	return rkev1.RKESystemConfig{
		Config: rkev1.GenericMap{
			Data: config,
		},
		MachineLabelSelector: selector,
	}
}
package etcdsnapshot

import (
	"errors"
	"fmt"
	"strings"
	"time"

	apisV1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/rancher/tests/v2/actions/scalinginput"
	"github.com/rancher/rancher/tests/v2/actions/services"
	deploy "github.com/rancher/rancher/tests/v2/actions/workloads/deployment"
	"github.com/rancher/shepherd/clients/rancher"
	management "github.com/rancher/shepherd/clients/rancher/generated/management/v3"
	steveV1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/clusters"
	"github.com/rancher/shepherd/extensions/clusters/kubernetesversions"
	extdefault "github.com/rancher/shepherd/extensions/defaults"
	"github.com/rancher/shepherd/extensions/defaults/stevetypes"
	shepherdsnapshot "github.com/rancher/shepherd/extensions/etcdsnapshot"
	extensionsingress "github.com/rancher/shepherd/extensions/ingresses"
	nodestat "github.com/rancher/shepherd/extensions/nodes"
	"github.com/rancher/shepherd/extensions/workloads"
	"github.com/rancher/shepherd/extensions/workloads/pods"
	namegen "github.com/rancher/shepherd/pkg/namegenerator"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	InitialIngress  = "ingress-before-restore"
	InitialWorkload = "wload-before-restore"

	all                          = "all"
	nginxImage                   = "nginx"
	containerName                = "nginx"
	defaultNamespace             = "default"
	DeploymentSteveType          = "apps.deployment"
	isCattleLabeled              = true
	ingressSteveType             = "networking.k8s.io.ingress"
	ingressPath                  = "/index.html"
	K3S                          = "k3s"
	kubernetesVersion            = "kubernetesVersion"
	namespace                    = "fleet-default"
	port                         = "port"
	postWorkload                 = "wload-after-backup"
	ProvisioningSteveResouceType = "provisioning.cattle.io.cluster"
	RKE1                         = "rke1"
	RKE2                         = "rke2"
	serviceAppendName            = "service-"
	serviceType                  = "service"
	windowsContainerImage        = "mcr.microsoft.com/windows/servercore/iis"
	windowsContainerName         = "iis"
)

// CreateAndValidateSnapshotRestore is an e2e helper that determines the engine type of the cluster, then takes a snapshot, and finally restores the cluster to the original snapshot
func CreateAndValidateSnapshotRestore(client *rancher.Client, clusterName string, etcdRestore *Config, containerImage string) error {

	clusterID, err := clusters.GetClusterIDByName(client, clusterName)
	if err != nil {
		return err
	}

	steveclient, err := client.Steve.ProxyDownstream(clusterID)
	if err != nil {
		return err
	}

	var isRKE1 = false

	clusterObject, _, _ := clusters.GetProvisioningClusterByName(client, clusterName, namespace)
	if clusterObject == nil {
		_, err := client.Management.Cluster.ByID(clusterID)
		if err != nil {
			return err
		}

		isRKE1 = true
	}

	podTemplate, deploymentTemplate, deploymentResp, serviceResp, ingressResp, err := createAndVerifyResources(steveclient, containerImage)
	if err != nil {
		return err
	}

	if isRKE1 {
		cluster, snapshotName, postDeploymentResp, postServiceResp, err := CreateAndValidateSnapshotRKE1(client, podTemplate, deploymentTemplate, clusterName, clusterID, etcdRestore, isRKE1)
		if err != nil {
			return err
		}

		err = RestoreAndValidateSnapshotRKE1(client, snapshotName, etcdRestore, cluster, clusterID)
		if err != nil {
			return err
		}

		_, err = steveclient.SteveType(DeploymentSteveType).ByID(postDeploymentResp.ID)
		if err == nil {
			return errors.New("expecting cluster restore to remove resource")
		}

		_, err = steveclient.SteveType(serviceType).ByID(postServiceResp.ID)
		if err == nil {
			return errors.New("expecting cluster restore to remove resource")
		}

	} else {
		cluster, snapshotName, postDeploymentResp, postServiceResp, err := CreateAndValidateSnapshotV2Prov(client, podTemplate, deploymentTemplate, clusterName, clusterID, etcdRestore, isRKE1)
		if err != nil {
			return err
		}

		err = RestoreAndValidateSnapshotV2Prov(client, snapshotName, etcdRestore, cluster, clusterID)
		if err != nil {
			return err
		}

		_, err = steveclient.SteveType(DeploymentSteveType).ByID(postDeploymentResp.ID)
		if err == nil {
			return errors.New("expecting cluster restore to remove resource")
		}

		_, err = steveclient.SteveType(serviceType).ByID(postServiceResp.ID)
		if err == nil {
			return errors.New("expecting cluster restore to remove resource")
		}
	}

	logrus.Infof("Deleting created workloads...")
	err = steveclient.SteveType(DeploymentSteveType).Delete(deploymentResp)
	if err != nil {
		return err
	}

	err = steveclient.SteveType(serviceType).Delete(serviceResp)
	if err != nil {
		return err
	}

	err = steveclient.SteveType(ingressSteveType).Delete(ingressResp)
	if err != nil {
		return err
	}
	return err
}

// CreateAndValidateSnapshotRKE1 is a helper that takes a snapshot of a given rke1 cluster and validates is resources after the snapshot
func CreateAndValidateSnapshotRKE1(client *rancher.Client, podTemplate *corev1.PodTemplateSpec, deployment *v1.Deployment, clusterName, clusterID string,
	etcdRestore *Config, isRKE1 bool) (*management.Cluster, string, *steveV1.SteveAPIObject, *steveV1.SteveAPIObject, error) {

	createdSnapshots, err := shepherdsnapshot.CreateRKE1Snapshot(client, clusterName)
	if err != nil {
		return nil, "", nil, nil, err
	}

	cluster, err := client.Management.Cluster.ByID(clusterID)
	if err != nil {
		return nil, "", nil, nil, err
	}

	if etcdRestore.ReplaceRoles != nil && cluster.RancherKubernetesEngineConfig.Services.Etcd.BackupConfig.S3BackupConfig != nil {
		err = scalinginput.ReplaceRKE1Nodes(client, clusterName, etcdRestore.ReplaceRoles.Etcd, etcdRestore.ReplaceRoles.ControlPlane, etcdRestore.ReplaceRoles.Worker)
		if err != nil {
			return nil, "", nil, nil, err
		}
	}

	snapshotToRestore := createdSnapshots[0].ID
	createdSnapshotIDs := []string{}
	isSnapshotS3 := false

	// prioritize s3 snapshots over local.
	for _, snapshot := range createdSnapshots {
		if snapshot.BackupConfig.S3BackupConfig != nil {
			snapshotToRestore = snapshot.ID
			isSnapshotS3 = true
		}
		createdSnapshotIDs = append(createdSnapshotIDs, snapshot.ID)
	}

	if cluster.RancherKubernetesEngineConfig.Services.Etcd.BackupConfig.S3BackupConfig != nil && !isSnapshotS3 {
		return nil, "", nil, nil, fmt.Errorf("s3 is enabled for the cluster, but selected snapshot is not from s3")
	}

	podErrors := pods.StatusPods(client, clusterID)
	if len(podErrors) != 0 {
		return nil, "", nil, nil, errors.New("cluster's pods not in good health post snapshot")
	}

	postDeploymentResp, postServiceResp, err := createPostBackupWorkloads(client, clusterID, *podTemplate, deployment)
	if err != nil {
		return nil, "", nil, nil, err
	}

	err = VerifySnapshots(client, clusterName, createdSnapshotIDs, isRKE1)
	if err != nil {
		return nil, "", nil, nil, err
	}

	if etcdRestore.SnapshotRestore == kubernetesVersion || etcdRestore.SnapshotRestore == all {
		clusterID, err := clusters.GetClusterIDByName(client, clusterName)
		if err != nil {
			return nil, "", nil, nil, err
		}

		clusterResp, err := client.Management.Cluster.ByID(clusterID)
		if err != nil {
			return nil, "", nil, nil, err
		}

		if etcdRestore.UpgradeKubernetesVersion == "" {
			defaultVersion, err := kubernetesversions.Default(client, clusters.RKE1ClusterType.String(), nil)
			etcdRestore.UpgradeKubernetesVersion = defaultVersion[0]
			if err != nil {
				return nil, "", nil, nil, err
			}
		}

		clusterResp.RancherKubernetesEngineConfig.Version = etcdRestore.UpgradeKubernetesVersion

		if etcdRestore.SnapshotRestore == all && etcdRestore.ControlPlaneUnavailableValue != "" && etcdRestore.WorkerUnavailableValue != "" {
			clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableControlplane = etcdRestore.ControlPlaneUnavailableValue
			clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableWorker = etcdRestore.WorkerUnavailableValue
		}

		_, err = client.Management.Cluster.Update(clusterResp, &clusterResp)
		if err != nil {
			return nil, "", nil, nil, err
		}

		err = clusters.WaitClusterToBeUpgraded(client, clusterID)
		if err != nil {
			return nil, "", nil, nil, err
		}

		logrus.Infof("Cluster version is upgraded to: %s", clusterResp.RancherKubernetesEngineConfig.Version)

		nodestat.AllManagementNodeReady(client, clusterResp.ID, extdefault.ThirtyMinuteTimeout)

		// getting a false positive when restoring rke1. fixing by re-checking the upgrade
		err = clusters.WaitClusterToBeUpgraded(client, clusterID)
		if err != nil {
			return nil, "", nil, nil, err
		}
		nodestat.AllManagementNodeReady(client, clusterResp.ID, extdefault.ThirtyMinuteTimeout)

		podErrors := pods.StatusPods(client, clusterID)
		if len(podErrors) != 0 {
			return nil, "", nil, nil, errors.New("cluster's pods not in good health post upgrade")
		}

		if etcdRestore.UpgradeKubernetesVersion != clusterResp.RancherKubernetesEngineConfig.Version {
			return nil, "", nil, nil, fmt.Errorf("K8s Version after upgrade %s does not match expected version %s", clusterResp.RancherKubernetesEngineConfig.Version, etcdRestore.UpgradeKubernetesVersion)
		}

		if etcdRestore.SnapshotRestore == all && etcdRestore.ControlPlaneUnavailableValue != "" && etcdRestore.WorkerUnavailableValue != "" {
			logrus.Infof("Control plane unavailable value is set to: %s", clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableControlplane)
			logrus.Infof("Worker unavailable value is set to: %s", clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableWorker)

			if etcdRestore.ControlPlaneUnavailableValue != clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableControlplane {
				return nil, "", nil, nil, fmt.Errorf("cpUnavailable after upgrade %s does not match expected version %s", clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableControlplane, etcdRestore.ControlPlaneUnavailableValue)
			}

			if etcdRestore.WorkerUnavailableValue != clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableWorker {
				return nil, "", nil, nil, fmt.Errorf("cpUnavailable after upgrade %s does not match expected version %s", clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableWorker, etcdRestore.WorkerUnavailableValue)
			}
		}
	}

	return cluster, snapshotToRestore, postDeploymentResp, postServiceResp, nil
}

// RestoreAndValidateSnapshotRKE1 restores a given snapshot for an rke1 cluster and validates its resources after the restore against the original cluster object
func RestoreAndValidateSnapshotRKE1(client *rancher.Client, snapshotName string, etcdRestore *Config, oldCluster *management.Cluster, clusterID string) error {
	// Give the option to restore the same snapshot multiple times. By default, it is set to 1.
	for i := 0; i < etcdRestore.RecurringRestores; i++ {
		snapshotRKE1Restore := &management.RestoreFromEtcdBackupInput{
			EtcdBackupID:     snapshotName,
			RestoreRkeConfig: etcdRestore.SnapshotRestore,
		}

		err := shepherdsnapshot.RestoreRKE1Snapshot(client, oldCluster.Name, snapshotRKE1Restore)
		if err != nil {
			return err
		}

		nodestat.AllManagementNodeReady(client, oldCluster.ID, extdefault.ThirtyMinuteTimeout)

		clusterResp, err := client.Management.Cluster.ByID(clusterID)
		if err != nil {
			return err
		}

		if oldCluster.RancherKubernetesEngineConfig.Version != clusterResp.RancherKubernetesEngineConfig.Version {
			return fmt.Errorf("K8s version after restore %s does not match expected version %s", clusterResp.RancherKubernetesEngineConfig.Version, oldCluster.RancherKubernetesEngineConfig.Version)
		}

		logrus.Infof("Cluster version is restored to: %s", clusterResp.RancherKubernetesEngineConfig.Version)

		client, err = client.ReLogin()
		if err != nil {
			return err
		}

		podErrors := pods.StatusPods(client, clusterID)
		if len(podErrors) != 0 {
			return errors.New("cluster's pods not in good health post restore")
		}

		if etcdRestore.SnapshotRestore == all && etcdRestore.ControlPlaneUnavailableValue != "" && etcdRestore.WorkerUnavailableValue != "" {
			logrus.Infof("Control plane unavailable value is restored to: %s", clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableControlplane)
			logrus.Infof("Worker unavailable value is restored to: %s", clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableWorker)

			if oldCluster.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableControlplane != clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableControlplane {
				return fmt.Errorf("cpUnavailable after restore %s does not match expected version %s", clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableControlplane, oldCluster.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableControlplane)
			}
			if oldCluster.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableWorker != clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableWorker {
				return fmt.Errorf("workerUnavailable after restore %s does not match expected version %s", clusterResp.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableWorker, oldCluster.RancherKubernetesEngineConfig.UpgradeStrategy.MaxUnavailableWorker)
			}
		}
	}
	return nil
}

// CreateAndValidateSnapshotV2Prov is a helper that takes a snapshot of a given v2prov cluster and validates is resources after the snapshot
func CreateAndValidateSnapshotV2Prov(client *rancher.Client, podTemplate *corev1.PodTemplateSpec, deployment *v1.Deployment, clusterName, clusterID string,
	etcdRestore *Config, isRKE1 bool) (*apisV1.Cluster, string, *steveV1.SteveAPIObject, *steveV1.SteveAPIObject, error) {

	createdSnapshots, err := shepherdsnapshot.CreateRKE2K3SSnapshot(client, clusterName)
	if err != nil {
		return nil, "", nil, nil, err
	}

	snapshotToRestore := createdSnapshots[0].ID
	createdSnapshotIDs := []string{}
	// prioritize s3 snapshots over local.
	for _, snapshot := range createdSnapshots {
		if strings.Contains(snapshot.ID, "-s3") {
			snapshotToRestore = snapshot.ID
		}
		createdSnapshotIDs = append(createdSnapshotIDs, snapshot.ID)
	}

	cluster, _, err := clusters.GetProvisioningClusterByName(client, clusterName, namespace)
	if err != nil {
		return nil, "", nil, nil, err
	}

	if cluster.Spec.RKEConfig.ETCD.S3 != nil && !strings.Contains(snapshotToRestore, "-s3") {
		return nil, "", nil, nil, fmt.Errorf("s3 is enabled for the cluster, but selected snapshot is not from s3")
	}

	if etcdRestore.ReplaceRoles != nil && cluster.Spec.RKEConfig.ETCD.S3 != nil {
		err = scalinginput.ReplaceNodes(client, clusterName, etcdRestore.ReplaceRoles.Etcd, etcdRestore.ReplaceRoles.ControlPlane, etcdRestore.ReplaceRoles.Worker)
		if err != nil {
			return nil, "", nil, nil, err
		}
	}

	podErrors := pods.StatusPods(client, clusterID)
	if len(podErrors) != 0 {
		return nil, "", nil, nil, errors.New("cluster's pods not in good health post snapshot")
	}

	postDeploymentResp, postServiceResp, err := createPostBackupWorkloads(client, clusterID, *podTemplate, deployment)
	if err != nil {
		return nil, "", nil, nil, err
	}

	err = VerifySnapshots(client, clusterName, createdSnapshotIDs, isRKE1)
	if err != nil {
		return nil, "", nil, nil, err
	}

	if etcdRestore.SnapshotRestore == kubernetesVersion || etcdRestore.SnapshotRestore == all {
		clusterObject, clusterResponse, err := clusters.GetProvisioningClusterByName(client, clusterName, namespace)
		if err != nil {
			return nil, "", nil, nil, err
		}

		initialKubernetesVersion := clusterObject.Spec.KubernetesVersion

		if etcdRestore.UpgradeKubernetesVersion == "" {
			if strings.Contains(initialKubernetesVersion, RKE2) {
				defaultVersion, err := kubernetesversions.Default(client, clusters.RKE2ClusterType.String(), nil)
				etcdRestore.UpgradeKubernetesVersion = defaultVersion[0]
				if err != nil {
					return nil, "", nil, nil, err
				}
			} else if strings.Contains(initialKubernetesVersion, K3S) {
				defaultVersion, err := kubernetesversions.Default(client, clusters.K3SClusterType.String(), nil)
				etcdRestore.UpgradeKubernetesVersion = defaultVersion[0]
				if err != nil {
					return nil, "", nil, nil, err
				}
			}
		}

		clusterObject.Spec.KubernetesVersion = etcdRestore.UpgradeKubernetesVersion

		if etcdRestore.SnapshotRestore == all && etcdRestore.ControlPlaneConcurrencyValue != "" && etcdRestore.WorkerConcurrencyValue != "" {
			clusterObject.Spec.RKEConfig.UpgradeStrategy.ControlPlaneConcurrency = etcdRestore.ControlPlaneConcurrencyValue
			clusterObject.Spec.RKEConfig.UpgradeStrategy.WorkerConcurrency = etcdRestore.WorkerConcurrencyValue
		}

		_, err = client.Steve.SteveType(ProvisioningSteveResouceType).Update(clusterResponse, clusterObject)
		if err != nil {
			return nil, "", nil, nil, err
		}

		err = clusters.WaitClusterToBeUpgraded(client, clusterID)
		if err != nil {
			return nil, "", nil, nil, err
		}

		logrus.Infof("Cluster version is upgraded to: %s", clusterObject.Spec.KubernetesVersion)

		podErrors := pods.StatusPods(client, clusterID)
		if len(podErrors) != 0 {
			return nil, "", nil, nil, errors.New("cluster's pods not in good health post upgrade")
		}

		if etcdRestore.UpgradeKubernetesVersion != clusterObject.Spec.KubernetesVersion {
			return nil, "", nil, nil, fmt.Errorf("K8s Version after upgrade %s does not match expected version %s", clusterObject.Spec.KubernetesVersion, etcdRestore.UpgradeKubernetesVersion)
		}

		if etcdRestore.SnapshotRestore == all && etcdRestore.ControlPlaneConcurrencyValue != "" && etcdRestore.WorkerConcurrencyValue != "" {
			logrus.Infof("Control plane concurrency value is set to: %s", clusterObject.Spec.RKEConfig.UpgradeStrategy.ControlPlaneConcurrency)
			logrus.Infof("Worker concurrency value is set to: %s", clusterObject.Spec.RKEConfig.UpgradeStrategy.WorkerConcurrency)

			if etcdRestore.ControlPlaneConcurrencyValue != clusterObject.Spec.RKEConfig.UpgradeStrategy.ControlPlaneConcurrency {
				return nil, "", nil, nil, fmt.Errorf("controlPlaneConcurrency after upgrade %s does not match expected version %s", clusterObject.Spec.RKEConfig.UpgradeStrategy.ControlPlaneConcurrency, etcdRestore.ControlPlaneConcurrencyValue)
			}

			if etcdRestore.WorkerConcurrencyValue != clusterObject.Spec.RKEConfig.UpgradeStrategy.WorkerConcurrency {
				return nil, "", nil, nil, fmt.Errorf("wokerConcurrency after upgrade %s does not match expected version %s", clusterObject.Spec.RKEConfig.UpgradeStrategy.WorkerConcurrency, etcdRestore.WorkerUnavailableValue)
			}
		}
		// sometimes we get a false positive on the cluster's state where it briefly goes 'active'. This is a way to mitigate that.
		clusterSteveObject, err := client.Steve.SteveType(ProvisioningSteveResouceType).ByID(clusterID)
		if err != nil {
			return nil, "", nil, nil, err
		}

		if clusterSteveObject.State == nil {
			err = clusters.WaitClusterUntilUpgrade(client, clusterID)
			if err != nil {
				return nil, "", nil, nil, err
			}

			podErrors := pods.StatusPods(client, clusterID)
			if len(podErrors) != 0 {
				return nil, "", nil, nil, errors.New("cluster's pods not in good health post upgrade")
			}
		}
	}

	return cluster, snapshotToRestore, postDeploymentResp, postServiceResp, err
}

// RestoreAndValidateSnapshotV2Prov restores a given snapshot for a v2prov cluster and validates its resources
// after the restore against the original cluster object
func RestoreAndValidateSnapshotV2Prov(client *rancher.Client, snapshotID string, etcdRestore *Config, cluster *apisV1.Cluster, clusterID string) error {
	clusterObject, _, err := clusters.GetProvisioningClusterByName(client, cluster.Name, namespace)
	if err != nil {
		return err
	}

	// Give the option to restore the same snapshot multiple times. By default, it is set to 1.
	for i := 0; i < etcdRestore.RecurringRestores; i++ {
		generation := int(1)
		if clusterObject.Spec.RKEConfig.ETCDSnapshotRestore != nil {
			generation = clusterObject.Spec.RKEConfig.ETCDSnapshotRestore.Generation + 1
		}

		splitSnapshot := strings.Split(snapshotID, "/")
		snapshotID = splitSnapshot[0]

		if len(splitSnapshot) > 1 {
			snapshotID = splitSnapshot[1]
		}

		snapshotRKE2K3SRestore := &rkev1.ETCDSnapshotRestore{
			Name:             snapshotID,
			Generation:       generation,
			RestoreRKEConfig: etcdRestore.SnapshotRestore,
		}

		err := shepherdsnapshot.RestoreRKE2K3SSnapshot(client, snapshotRKE2K3SRestore, clusterObject.Name)
		if err != nil {
			return err
		}

		clusterObject, _, err = clusters.GetProvisioningClusterByName(client, cluster.Name, namespace)
		if err != nil {
			return err
		}

		err = clusters.WaitClusterToBeUpgraded(client, clusterID)
		if err != nil {
			return err
		}

		podErrors := pods.StatusPods(client, clusterID)
		if len(podErrors) != 0 {
			return errors.New("cluster's pods not in good health post restore")
		}

		if cluster.Spec.KubernetesVersion != clusterObject.Spec.KubernetesVersion {
			return fmt.Errorf("K8s Version after upgrade %s does not match expected version %s after restore", clusterObject.Spec.KubernetesVersion, etcdRestore.UpgradeKubernetesVersion)
		}

		if etcdRestore.SnapshotRestore == all && etcdRestore.ControlPlaneConcurrencyValue != "" && etcdRestore.WorkerConcurrencyValue != "" {
			logrus.Infof("Control plane concurrency value is restored to: %s", clusterObject.Spec.RKEConfig.UpgradeStrategy.ControlPlaneConcurrency)
			logrus.Infof("Worker concurrency value is restored to: %s", clusterObject.Spec.RKEConfig.UpgradeStrategy.WorkerConcurrency)

			if cluster.Spec.RKEConfig.UpgradeStrategy.ControlPlaneConcurrency != clusterObject.Spec.RKEConfig.UpgradeStrategy.ControlPlaneConcurrency {
				return fmt.Errorf("controlPlaneConcurrency after restore %s does not match expected version %s", clusterObject.Spec.RKEConfig.UpgradeStrategy.ControlPlaneConcurrency, cluster.Spec.RKEConfig.UpgradeStrategy.ControlPlaneConcurrency)
			}

			if cluster.Spec.RKEConfig.UpgradeStrategy.WorkerConcurrency != clusterObject.Spec.RKEConfig.UpgradeStrategy.WorkerConcurrency {
				return fmt.Errorf("wokerConcurrency after restore %s does not match expected version %s", clusterObject.Spec.RKEConfig.UpgradeStrategy.WorkerConcurrency, cluster.Spec.RKEConfig.UpgradeStrategy.ControlPlaneConcurrency)
			}
		}
	}

	return nil
}

// This function waits for retentionlimit+1 automatic snapshots to be taken before verifying that the retention limit is respected
func CreateSnapshotsUntilRetentionLimit(client *rancher.Client, clusterName string, retentionLimit int, timeBetweenSnapshots int) error {
	v1ClusterID, err := clusters.GetV1ProvisioningClusterByName(client, clusterName)
	if v1ClusterID == "" {
		v3ClusterID, err := clusters.GetClusterIDByName(client, clusterName)
		if err != nil {
			return err
		}
		v1ClusterID = "fleet-default/" + v3ClusterID
	}
	if err != nil {
		return err
	}

	fleetCluster, err := client.Steve.SteveType(stevetypes.FleetCluster).ByID(v1ClusterID)
	if err != nil {
		return err
	}

	provider := fleetCluster.ObjectMeta.Labels["provider.cattle.io"]
	if provider == "rke" {
		sleepNum := (retentionLimit + 1) * timeBetweenSnapshots
		logrus.Infof("Waiting %v hours for %v automatic snapshots to be taken", sleepNum, (retentionLimit + 1))
		time.Sleep(time.Duration(sleepNum)*time.Hour + time.Minute*5)

		err := RKE1RetentionLimitCheck(client, clusterName)
		if err != nil {
			return err
		}

	} else {
		sleepNum := (retentionLimit + 1) * timeBetweenSnapshots
		logrus.Infof("Waiting %v minutes for %v automatic snapshots to be taken", sleepNum, (retentionLimit + 1))
		time.Sleep(time.Duration(sleepNum)*time.Minute + time.Minute*5)

		err := RKE2K3SRetentionLimitCheck(client, clusterName)
		if err != nil {
			return err
		}
	}

	return nil
}

func createPostBackupWorkloads(client *rancher.Client, clusterID string, podTemplate corev1.PodTemplateSpec, deployment *v1.Deployment) (*steveV1.SteveAPIObject, *steveV1.SteveAPIObject, error) {
	workloadNamePostBackup := namegen.AppendRandomString(postWorkload)

	postBackupDeployment := workloads.NewDeploymentTemplate(workloadNamePostBackup, defaultNamespace, podTemplate, isCattleLabeled, nil)
	postBackupService := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAppendName + workloadNamePostBackup,
			Namespace: defaultNamespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: port,
					Port: 80,
				},
			},
			Selector: deployment.Spec.Template.Labels,
		},
	}

	steveclient, err := client.Steve.ProxyDownstream(clusterID)
	if err != nil {
		return nil, nil, err
	}

	postDeploymentResp, err := createDeployment(steveclient, workloadNamePostBackup, postBackupDeployment)
	if err != nil {
		return nil, nil, err
	}

	err = deploy.VerifyDeployment(steveclient, postDeploymentResp)
	if err != nil {
		return nil, nil, err
	}

	if workloadNamePostBackup != postDeploymentResp.ObjectMeta.Name {
		return nil, nil, fmt.Errorf("PostBackup deployment name %s does not match created deployment %s ", workloadNamePostBackup, postDeploymentResp.ObjectMeta.Name)
	}

	postServiceResp, err := services.CreateService(steveclient, postBackupService)
	if err != nil {
		return nil, nil, err
	}

	err = services.VerifyService(steveclient, postServiceResp)
	if err != nil {
		return nil, nil, err
	}

	if serviceAppendName+workloadNamePostBackup != postServiceResp.ObjectMeta.Name {
		return nil, nil, fmt.Errorf("PostBackup service name %s does not match created deployment %s ", serviceAppendName+workloadNamePostBackup, postServiceResp.ObjectMeta.Name)
	}

	return postDeploymentResp, postServiceResp, nil
}

func createAndVerifyResources(steveclient *steveV1.Client, containerImage string) (*corev1.PodTemplateSpec, *v1.Deployment, *steveV1.SteveAPIObject, *steveV1.SteveAPIObject, *steveV1.SteveAPIObject, error) {
	var containerTemplate corev1.Container
	initialIngressName := namegen.AppendRandomString(InitialIngress)
	initialWorkloadName := namegen.AppendRandomString(InitialWorkload)

	containerTemplate = workloads.NewContainer(containerName, containerImage, corev1.PullAlways, []corev1.VolumeMount{}, []corev1.EnvFromSource{}, nil, nil, nil)

	podTemplate := workloads.NewPodTemplate([]corev1.Container{containerTemplate}, []corev1.Volume{}, []corev1.LocalObjectReference{}, nil, map[string]string{})
	deploymentTemplate := workloads.NewDeploymentTemplate(initialWorkloadName, defaultNamespace, podTemplate, isCattleLabeled, nil)

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAppendName + initialWorkloadName,
			Namespace: defaultNamespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: port,
					Port: 80,
				},
			},
			Selector: deploymentTemplate.Spec.Template.Labels,
		},
	}

	deploymentResp, err := createDeployment(steveclient, initialWorkloadName, deploymentTemplate)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = deploy.VerifyDeployment(steveclient, deploymentResp)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if initialWorkloadName != deploymentResp.ObjectMeta.Name {
		return nil, nil, nil, nil, nil, errors.New("deployment name doesn't match spec")
	}

	serviceResp, err := services.CreateService(steveclient, service)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = services.VerifyService(steveclient, serviceResp)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if serviceAppendName+initialWorkloadName != serviceResp.ObjectMeta.Name {
		return nil, nil, nil, nil, nil, errors.New("service name doesn't match spec")
	}

	path := extensionsingress.NewIngressPathTemplate(networking.PathTypeExact, ingressPath, serviceAppendName+initialWorkloadName, 80)
	ingressTemplate := extensionsingress.NewIngressTemplate(initialIngressName, defaultNamespace, "", []networking.HTTPIngressPath{path})

	ingressResp, err := extensionsingress.CreateIngress(steveclient, initialIngressName, ingressTemplate)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = extensionsingress.WaitIngress(steveclient, ingressResp, initialIngressName)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if initialIngressName != ingressResp.ObjectMeta.Name {
		return nil, nil, nil, nil, nil, errors.New("ingress name doesn't match spec")
	}

	return &podTemplate, deploymentTemplate, deploymentResp, serviceResp, ingressResp, nil
}

func createDeployment(steveclient *steveV1.Client, wlName string, deployment *v1.Deployment) (*steveV1.SteveAPIObject, error) {
	logrus.Infof("Creating deployment: %s", wlName)
	deploymentResp, err := steveclient.SteveType(DeploymentSteveType).Create(deployment)
	if err != nil {
		return nil, err
	}

	return deploymentResp, err
}
