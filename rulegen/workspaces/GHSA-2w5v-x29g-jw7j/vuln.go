package main


	defer metrics.MeasureSince([]string{"nomad", "volume", "register"}, time.Now())

	if !allowVolume(aclObj, args.RequestNamespace()) || !aclObj.AllowPluginRead() {
		return structs.ErrPermissionDenied
	}

		return fmt.Errorf("missing volume definition")
	}

	// This is the only namespace we ACL checked, force all the volumes to use it.
	// We also validate that the plugin exists for each plugin, and validate the
	// capabilities when the plugin has a controller.
	for _, vol := range args.Volumes {

		snap, err := v.srv.State().Snapshot()
		if err != nil {
			return err
		}
		if vol.Namespace == "" {
			vol.Namespace = args.RequestNamespace()
		}
		if err = vol.Validate(); err != nil {
			return err
		}

		ws := memdb.NewWatchSet()
		existingVol, err := snap.CSIVolumeByID(ws, vol.Namespace, vol.ID)
		if err != nil {
			return err
		}
		return err
	}

	if !allowVolume(aclObj, args.RequestNamespace()) || !aclObj.AllowPluginRead() {
		return structs.ErrPermissionDenied
	}

	}
	validatedVols := []validated{}

	// This is the only namespace we ACL checked, force all the volumes to use it.
	// We also validate that the plugin exists for each plugin, and validate the
	// capabilities when the plugin has a controller.
	for _, vol := range args.Volumes {
		if vol.Namespace == "" {
			vol.Namespace = args.RequestNamespace()
		}
		if err = vol.Validate(); err != nil {
			return err
		}
		}

		// if the volume already exists, we'll update it instead
		snap, err := v.srv.State().Snapshot()
		if err != nil {
			return err
		}
		// current will be nil if it does not exist.
		current, err := snap.CSIVolumeByID(nil, vol.Namespace, vol.ID)
		if err != nil {

func TestCSIVolumeEndpoint_Register(t *testing.T) {
	ci.Parallel(t)
	srv, shutdown := TestServer(t, func(c *Config) {
		c.NumSchedulers = 0 // Prevent automatic dequeue
	})
	defer shutdown()

	id0 := uuid.Generate()

	// Create the register request
	ns := mock.Namespace()
	store.UpsertNamespaces(900, []*structs.Namespace{ns})

	// Create the node and plugin
	node := mock.Node()
			NodeInfo: &structs.CSINodeInfo{},
		},
	}
	require.NoError(t, store.UpsertNode(structs.MsgTypeTestSetup, 1000, node))

	// Create the volume
	vols := []*structs.CSIVolume{{
		ID:             id0,
		Namespace:      ns.Name,
		PluginID:       "minnie",
		AccessMode:     structs.CSIVolumeAccessModeSingleNodeReader, // legacy field ignored
		AttachmentMode: structs.CSIVolumeAttachmentModeBlockDevice,  // legacy field ignored
	}}

	// Create the register request
	req1 := &structs.CSIVolumeRegisterRequest{
		Volumes: vols,
		WriteRequest: structs.WriteRequest{
			Region:    "global",
			Namespace: "",
		},
	}
	resp1 := &structs.CSIVolumeRegisterResponse{}
	err := msgpackrpc.CallWithCodec(codec, "CSIVolume.Register", req1, resp1)
	require.NoError(t, err)
	require.NotEqual(t, uint64(0), resp1.Index)

	// Get the volume back out
	req2 := &structs.CSIVolumeGetRequest{
		ID: id0,
		QueryOptions: structs.QueryOptions{
			Region:    "global",
			Namespace: ns.Name,
		},
	}
	resp2 := &structs.CSIVolumeGetResponse{}
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Get", req2, resp2)
	require.NoError(t, err)
	require.Equal(t, resp1.Index, resp2.Index)
	require.Equal(t, vols[0].ID, resp2.Volume.ID)
	require.Equal(t, "csi.CSISecrets(map[mysecret:[REDACTED]])",
		resp2.Volume.Secrets.String())
	require.Equal(t, "csi.CSIOptions(FSType: ext4, MountFlags: [REDACTED])",
		resp2.Volume.MountOptions.String())

	// Registration does not update
	req1.Volumes[0].PluginID = "adam"
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Register", req1, resp1)
	require.Error(t, err, "exists")

	// Deregistration works
	req3 := &structs.CSIVolumeDeregisterRequest{
		VolumeIDs: []string{id0},
		WriteRequest: structs.WriteRequest{
			Region:    "global",
			Namespace: ns.Name,
		},
	}
	resp3 := &structs.CSIVolumeDeregisterResponse{}
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Deregister", req3, resp3)
	require.NoError(t, err)

	// Volume is missing
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Get", req2, resp2)
	require.NoError(t, err)
	require.Nil(t, resp2.Volume)
}

// TestCSIVolumeEndpoint_Claim exercises the VolumeClaim RPC, verifying that claims
func TestCSIVolumeEndpoint_Create(t *testing.T) {
	ci.Parallel(t)
	var err error
	srv, shutdown := TestServer(t, func(c *Config) {
		c.NumSchedulers = 0 // Prevent automatic dequeue
	})
	defer shutdown()

	req0 := &structs.NodeRegisterRequest{
		Node:         node,
		WriteRequest: structs.WriteRequest{Region: "global"},
	}
	var resp0 structs.NodeUpdateResponse
	err = client.RPC("Node.Register", req0, &resp0)
	require.NoError(t, err)

	testutil.WaitForResult(func() (bool, error) {
		nodes := srv.connectedNodes()
		t.Fatalf("should have a client")
	})

	ns := structs.DefaultNamespace

	state := srv.fsm.State()
	codec := rpcClient(t, srv)
	index := uint64(1000)

		}
	}).Node
	index++
	require.NoError(t, state.UpsertNode(structs.MsgTypeTestSetup, index, node))

	// Create the volume
	volID := uuid.Generate()
	vols := []*structs.CSIVolume{{
		ID:             volID,
		Name:           "vol",
		Namespace:      "", // overriden by WriteRequest
		PluginID:       "minnie",
		AccessMode:     structs.CSIVolumeAccessModeSingleNodeReader, // legacy field ignored
		AttachmentMode: structs.CSIVolumeAttachmentModeBlockDevice,  // legacy field ignored
		},
	}}

	// Create the create request
	req1 := &structs.CSIVolumeCreateRequest{
		Volumes: vols,
		WriteRequest: structs.WriteRequest{
			Region:    "global",
			Namespace: ns,
		},
	}
	resp1 := &structs.CSIVolumeCreateResponse{}
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Create", req1, resp1)
	require.NoError(t, err)

	// Get the volume back out
	req2 := &structs.CSIVolumeGetRequest{
		ID: volID,
		QueryOptions: structs.QueryOptions{
			Region: "global",
		},
	}
	resp2 := &structs.CSIVolumeGetResponse{}
	err = msgpackrpc.CallWithCodec(codec, "CSIVolume.Get", req2, resp2)
	require.NoError(t, err)
	require.Equal(t, resp1.Index, resp2.Index)

	vol := resp2.Volume
	require.NotNil(t, vol)
	require.Equal(t, volID, vol.ID)

	// these fields are set from the args
	require.Equal(t, "csi.CSISecrets(map[mysecret:[REDACTED]])",
		vol.Secrets.String())
	require.Equal(t, "csi.CSIOptions(FSType: ext4, MountFlags: [REDACTED])",
		vol.MountOptions.String())
	require.Equal(t, ns, vol.Namespace)
	require.Len(t, vol.RequestedCapabilities, 1)

	// these fields are set from the plugin and should have been written to raft
	require.Equal(t, "vol-12345", vol.ExternalID)
	require.Equal(t, int64(42), vol.Capacity)
	require.Equal(t, "bar", vol.Context["plugincontext"])
	require.Equal(t, "", vol.Context["mycontext"])
	require.Equal(t, map[string]string{"rack": "R1"}, vol.Topologies[0].Segments)
}

func TestCSIVolumeEndpoint_Delete(t *testing.T) {
