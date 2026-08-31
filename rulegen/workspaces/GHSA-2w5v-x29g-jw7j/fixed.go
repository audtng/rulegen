package main


	defer metrics.MeasureSince([]string{"nomad", "volume", "register"}, time.Now())

	// permission for the volume namespaces will be checked below
	if !aclObj.AllowPluginRead() {
		return structs.ErrPermissionDenied
	}

		return fmt.Errorf("missing volume definition")
	}

	snap, err := v.srv.State().Snapshot()
	if err != nil {
		return err
	}

	// Validate ACLs, that the plugin exists for each volume, and validate the
	// capabilities when the plugin has a controller.
	for _, vol := range args.Volumes {
		if vol.Namespace == "" {
			vol.Namespace = args.RequestNamespace()
		}
		if !allowVolume(aclObj, vol.Namespace) {
			return structs.ErrPermissionDenied
		}
		if err = vol.Validate(); err != nil {
			return err
		}

		existingVol, err := snap.CSIVolumeByID(nil, vol.Namespace, vol.ID)
		if err != nil {
			return err
		}
		return err
	}

	// permission for the volume namespaces will be checked below
	if !aclObj.AllowPluginRead() {
		return structs.ErrPermissionDenied
	}

	}
	validatedVols := []validated{}

	snap, err := v.srv.State().Snapshot()
	if err != nil {
		return err
	}

	// Validate ACLs, that the plugin exists for each volume, and validate the
	// capabilities when the plugin has a controller.
	for _, vol := range args.Volumes {
		if vol.Namespace == "" {
			vol.Namespace = args.RequestNamespace()
		}
		if !allowVolume(aclObj, vol.Namespace) {
			return structs.ErrPermissionDenied
		}
		if err = vol.Validate(); err != nil {
			return err
		}
		}

		// if the volume already exists, we'll update it instead
		// current will be nil if it does not exist.
		current, err := snap.CSIVolumeByID(nil, vol.Namespace, vol.ID)
		if err != nil {

func TestCSIVolumeEndpoint_Register(t *testing.T) {
	ci.Parallel(t)
	srv, _, shutdown := TestACLServer(t, func(c *Config) {
		c.NumSchedulers = 0 // Prevent automatic dequeue
	})
	defer shutdown()

	id0 := uuid.Generate()

	ns := "prod"
	otherNS := "other"
	index := uint64(1000)
	must.NoError(t, store.UpsertNamespaces(index, []*structs.Namespace{{Name: ns}, {Name: otherNS}}))

	// Create the node and plugin
	node := mock.Node()
			NodeInfo: &structs.CSINodeInfo{},
		},
	}
	must.NoError(t, store.UpsertNode(structs.MsgTypeTestSetup, index, node))

	index++
	validToken := mock.CreatePolicyAndToken(t, store, index, "csi-access-ns",
		`namespace "prod" { capabilities = ["csi-write-volume", "csi-read-volume"] }
         namespace "default" { capabilities = ["csi-write-volume"] }
         plugin { policy = "read" }
         node { policy = "read" }`).SecretID

	index++
	invalidToken := mock.CreatePolicyAndToken(t, store, index, "csi-access-other",
		`namespace "other" { capabilities = ["csi-write-volume"] }
         plugin { policy = "read" }
         node { policy = "read" }`).SecretID

	vols := []*structs.CSIVolume{{
		ID:             id0,
		Namespace:      ns,
		PluginID:       "minnie",
		AccessMode:     structs.CSIVolumeAccessModeSingleNodeReader, // legacy field ignored
		AttachmentMode: structs.CSIVolumeAttachmentModeBlockDevice,  // legacy field ignored
	}}

	// Create the register request
	// The token has access to the request namespace but not the volume namespace
	req1 := &structs.CSIVolumeRegisterRequest{
		Volumes: vols,
		WriteRequest: structs.WriteRequest{
			Region:    "global",
			AuthToken: invalidToken,
			Namespace: otherNS,
		},
	}
	resp1 := &structs.CSIVolumeRegisterResponse{}
	err := msgpackrpc.CallWithCodec(codec, "CSIVolume.Register", req1, resp1)
	must.EqError(t, err, "Permission denied")

	// Switch to a token that has access to the volume's namespace, but switch
	// the request namespace to one that will be overwritten by the vol spec
	req1.AuthToken = validToken
	req1.Namespace = structs.DefaultNamespace
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Register", req1, resp1)
	must.NoError(t, err)
	must.NotEq(t, uint64(0), resp1.Index)

	// Get the volume back out
	req2 := &structs.CSIVolumeGetRequest{
		ID: id0,
		QueryOptions: structs.QueryOptions{
			Region:    "global",
			Namespace: ns,
			AuthToken: validToken,
		},
	}
	resp2 := &structs.CSIVolumeGetResponse{}
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Get", req2, resp2)
	must.NoError(t, err)
	must.Eq(t, resp1.Index, resp2.Index)
	must.Eq(t, vols[0].ID, resp2.Volume.ID)
	must.Eq(t, "csi.CSISecrets(map[mysecret:[REDACTED]])",
		resp2.Volume.Secrets.String())
	must.Eq(t, "csi.CSIOptions(FSType: ext4, MountFlags: [REDACTED])",
		resp2.Volume.MountOptions.String())

	// Registration does not update
	req1.Volumes[0].PluginID = "adam"
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Register", req1, resp1)
	must.ErrorContains(t, err, "no CSI plugin named")

	// Deregistration works
	req3 := &structs.CSIVolumeDeregisterRequest{
		VolumeIDs: []string{id0},
		WriteRequest: structs.WriteRequest{
			Region:    "global",
			Namespace: ns,
			AuthToken: validToken,
		},
	}
	resp3 := &structs.CSIVolumeDeregisterResponse{}
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Deregister", req3, resp3)
	must.NoError(t, err)

	// Volume is missing
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Get", req2, resp2)
	must.NoError(t, err)
	must.Nil(t, resp2.Volume)
}

// TestCSIVolumeEndpoint_Claim exercises the VolumeClaim RPC, verifying that claims
func TestCSIVolumeEndpoint_Create(t *testing.T) {
	ci.Parallel(t)
	var err error
	srv, rootToken, shutdown := TestACLServer(t, func(c *Config) {
		c.NumSchedulers = 0 // Prevent automatic dequeue
	})
	defer shutdown()

	req0 := &structs.NodeRegisterRequest{
		Node:         node,
		WriteRequest: structs.WriteRequest{Region: "global", AuthToken: rootToken.SecretID},
	}
	var resp0 structs.NodeUpdateResponse
	err = client.RPC("Node.Register", req0, &resp0)
	must.NoError(t, err)

	testutil.WaitForResult(func() (bool, error) {
		nodes := srv.connectedNodes()
		t.Fatalf("should have a client")
	})

	ns := "prod"
	otherNS := "other"

	store := srv.fsm.State()
	codec := rpcClient(t, srv)
	index := uint64(1000)

		}
	}).Node
	index++
	must.NoError(t, store.UpsertNode(structs.MsgTypeTestSetup, index, node))

	index++
	must.NoError(t, store.UpsertNamespaces(index, []*structs.Namespace{{Name: ns}, {Name: otherNS}}))

	// Create the volume
	volID := uuid.Generate()
	vols := []*structs.CSIVolume{{
		ID:             volID,
		Name:           "vol",
		Namespace:      ns,
		PluginID:       "minnie",
		AccessMode:     structs.CSIVolumeAccessModeSingleNodeReader, // legacy field ignored
		AttachmentMode: structs.CSIVolumeAttachmentModeBlockDevice,  // legacy field ignored
		},
	}}

	index++
	validToken := mock.CreatePolicyAndToken(t, store, index, "csi-access-ns",
		`namespace "prod" { capabilities = ["csi-write-volume", "csi-read-volume"] }
         namespace "default" { capabilities = ["csi-write-volume"] }
         plugin { policy = "read" }
         node { policy = "read" }`).SecretID

	index++
	invalidToken := mock.CreatePolicyAndToken(t, store, index, "csi-access-other",
		`namespace "other" { capabilities = ["csi-write-volume"] }
         plugin { policy = "read" }
         node { policy = "read" }`).SecretID

	// Create the create request
	// The token has access to the request namespace but not the volume namespace
	req1 := &structs.CSIVolumeCreateRequest{
		Volumes: vols,
		WriteRequest: structs.WriteRequest{
			Region:    "global",
			AuthToken: invalidToken,
			Namespace: otherNS,
		},
	}
	resp1 := &structs.CSIVolumeCreateResponse{}
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Create", req1, resp1)
	must.EqError(t, err, "Permission denied")

	// Switch to a token that has access to the volume's namespace, but switch
	// the request namespace to one that will be overwritten by the vol spec
	req1.AuthToken = validToken
	req1.Namespace = structs.DefaultNamespace
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Create", req1, resp1)
	must.NoError(t, err)
	must.NotEq(t, uint64(0), resp1.Index)

	// Get the volume back out
	req2 := &structs.CSIVolumeGetRequest{
		ID: volID,
		QueryOptions: structs.QueryOptions{
			Region:    "global",
			Namespace: ns,
			AuthToken: validToken,
		},
	}
	resp2 := &structs.CSIVolumeGetResponse{}
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Get", req2, resp2)
	must.NoError(t, err)
	must.Eq(t, resp1.Index, resp2.Index)

	vol := resp2.Volume
	must.NotNil(t, vol)
	must.Eq(t, volID, vol.ID)

	// these fields are set from the args
	must.Eq(t, "csi.CSISecrets(map[mysecret:[REDACTED]])",
		vol.Secrets.String())
	must.Eq(t, "csi.CSIOptions(FSType: ext4, MountFlags: [REDACTED])",
		vol.MountOptions.String())
	must.Eq(t, ns, vol.Namespace)
	must.Len(t, 1, vol.RequestedCapabilities)

	// these fields are set from the plugin and should have been written to raft
	must.Eq(t, "vol-12345", vol.ExternalID)
	must.Eq(t, int64(42), vol.Capacity)
	must.Eq(t, "bar", vol.Context["plugincontext"])
	must.Eq(t, "", vol.Context["mycontext"])
	must.Eq(t, map[string]string{"rack": "R1"}, vol.Topologies[0].Segments)
}

func TestCSIVolumeEndpoint_Delete(t *testing.T) {
