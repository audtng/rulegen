package main

							Name:      submariner.NetworkPluginSyncerComponent,
						},
					},
					&rbacv1.ClusterRole{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: t.Namespace,
							Name:      "ocp-submariner-networkplugin-syncer",
						},
					},
					&rbacv1.ClusterRoleBinding{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: t.Namespace,
							Name:      submariner.NetworkPluginSyncerComponent,
						},
					},
					&rbacv1.ClusterRoleBinding{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: t.Namespace,
							Name:      "ocp-submariner-networkplugin-syncer",
						},
					},
					&corev1.ServiceAccount{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: t.Namespace,
					},
				})

				t.AssertNoResource(&rbacv1.ClusterRole{
					ObjectMeta: metav1.ObjectMeta{
						Name: "ocp-submariner-networkplugin-syncer",
					},
				})

				t.AssertNoResource(&rbacv1.ClusterRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name: "ocp-submariner-networkplugin-syncer",
					},
				})

				t.AssertNoResource(&corev1.ServiceAccount{
					ObjectMeta: metav1.ObjectMeta{
						Name: submariner.NetworkPluginSyncerComponent,

	deleteAll := func(objs ...client.Object) error {
		for _, obj := range objs {
			obj.SetNamespace(instance.Namespace)

			err := r.config.ScopedClient.Delete(ctx, obj)
	return deleteAll(
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: NetworkPluginSyncerComponent,
			},
		},
		&rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: NetworkPluginSyncerComponent,
			},
		},
		&rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: NetworkPluginSyncerComponent,
			},
		},
		&rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ocp-submariner-networkplugin-syncer",
			},
		},
		&rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ocp-submariner-networkplugin-syncer",
			},
		},
		&corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: NetworkPluginSyncerComponent,
			},
		},
	)
    resources:
      # Temporarily needed for network-plugin syncer removal
      - serviceaccounts
    resourceNames:
      - submariner-networkplugin-syncer
    verbs:
      - delete
  - apiGroups:
      # Temporarily needed for network-plugin syncer removal
      - clusterroles
      - clusterrolebindings
    resourceNames:
      - ocp-submariner-networkplugin-syncer
      - submariner-networkplugin-syncer
    verbs:
      - delete
`
