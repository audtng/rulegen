package main

							Name:      submariner.NetworkPluginSyncerComponent,
						},
					},
					&rbacv1.ClusterRoleBinding{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: t.Namespace,
							Name:      submariner.NetworkPluginSyncerComponent,
						},
					},
					&corev1.ServiceAccount{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: t.Namespace,
					},
				})

				t.AssertNoResource(&corev1.ServiceAccount{
					ObjectMeta: metav1.ObjectMeta{
						Name: submariner.NetworkPluginSyncerComponent,

	deleteAll := func(objs ...client.Object) error {
		for _, obj := range objs {
			obj.SetName(NetworkPluginSyncerComponent)
			obj.SetNamespace(instance.Namespace)

			err := r.config.ScopedClient.Delete(ctx, obj)
	return deleteAll(
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: instance.Namespace,
				Name:      NetworkPluginSyncerComponent,
			},
		},
		&rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: instance.Namespace,
				Name:      NetworkPluginSyncerComponent,
			},
		},
		&rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: instance.Namespace,
				Name:      NetworkPluginSyncerComponent,
			},
		},
		&rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: instance.Namespace,
				Name:      "ocp-submariner-networkplugin-syncer",
			},
		},
		&rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: instance.Namespace,
				Name:      "ocp-submariner-networkplugin-syncer",
			},
		},
		&corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: instance.Namespace,
				Name:      NetworkPluginSyncerComponent,
			},
		},
	)
    resources:
      # Temporarily needed for network-plugin syncer removal
      - serviceaccounts
    verbs:
      - delete
  - apiGroups:
      # Temporarily needed for network-plugin syncer removal
      - clusterroles
      - clusterrolebindings
    verbs:
      - delete
`
