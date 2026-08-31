package main

				}
			}),
		},
		{
			Name: "kube_certificatesigningrequest_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapCSRFunc(func(j *certv1beta1.CertificateSigningRequest) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(j.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_certificatesigningrequest_condition gauge
		# HELP kube_certificatesigningrequest_cert_length Length of the issued cert
		# TYPE kube_certificatesigningrequest_cert_length gauge
		# HELP kube_certificatesigningrequest_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_certificatesigningrequest_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Labels: map[string]string{
						"cert": "test",
					},
					Annotations: map[string]string{
						"cert": "test",
					},
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
				},
				Status: certv1beta1.CertificateSigningRequestStatus{},
				kube_certificatesigningrequest_condition{certificatesigningrequest="certificate-test",condition="denied"} 0
				kube_certificatesigningrequest_labels{certificatesigningrequest="certificate-test",label_cert="test"} 1
				kube_certificatesigningrequest_cert_length{certificatesigningrequest="certificate-test"} 0
				kube_certificatesigningrequest_annotations{certificatesigningrequest="certificate-test",annotation_cert="test"} 1
`,
			MetricNames: []string{"kube_certificatesigningrequest_created", "kube_certificatesigningrequest_condition", "kube_certificatesigningrequest_labels", "kube_certificatesigningrequest_cert_length", "kube_certificatesigningrequest_annotations"},
		},
		{
			Obj: &certv1beta1.CertificateSigningRequest{
					Labels: map[string]string{
						"cert": "test",
					},
					Annotations: map[string]string{
						"cert": "test",
					},
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
				},
				Status: certv1beta1.CertificateSigningRequestStatus{
				kube_certificatesigningrequest_condition{certificatesigningrequest="certificate-test",condition="denied"} 1
				kube_certificatesigningrequest_labels{certificatesigningrequest="certificate-test",label_cert="test"} 1
				kube_certificatesigningrequest_cert_length{certificatesigningrequest="certificate-test"} 0
				kube_certificatesigningrequest_annotations{certificatesigningrequest="certificate-test",annotation_cert="test"} 1
`,
			MetricNames: []string{"kube_certificatesigningrequest_created", "kube_certificatesigningrequest_condition", "kube_certificatesigningrequest_labels", "kube_certificatesigningrequest_cert_length", "kube_certificatesigningrequest_annotations"},
		},
		{
			Obj: &certv1beta1.CertificateSigningRequest{
					Labels: map[string]string{
						"cert": "test",
					},
					Annotations: map[string]string{
						"cert": "test",
					},
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
				},
				Status: certv1beta1.CertificateSigningRequestStatus{
				kube_certificatesigningrequest_condition{certificatesigningrequest="certificate-test",condition="denied"} 0
				kube_certificatesigningrequest_labels{certificatesigningrequest="certificate-test",label_cert="test"} 1
				kube_certificatesigningrequest_cert_length{certificatesigningrequest="certificate-test"} 0
				kube_certificatesigningrequest_annotations{certificatesigningrequest="certificate-test",annotation_cert="test"} 1
`,
			MetricNames: []string{"kube_certificatesigningrequest_created", "kube_certificatesigningrequest_condition", "kube_certificatesigningrequest_labels", "kube_certificatesigningrequest_cert_length", "kube_certificatesigningrequest_annotations"},
		},
		{
			Obj: &certv1beta1.CertificateSigningRequest{
					Labels: map[string]string{
						"cert": "test",
					},
					Annotations: map[string]string{
						"cert": "test",
					},
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
				},
				Status: certv1beta1.CertificateSigningRequestStatus{
				kube_certificatesigningrequest_condition{certificatesigningrequest="certificate-test",condition="denied"} 0
				kube_certificatesigningrequest_labels{certificatesigningrequest="certificate-test",label_cert="test"} 1
				kube_certificatesigningrequest_cert_length{certificatesigningrequest="certificate-test"} 13
				kube_certificatesigningrequest_annotations{certificatesigningrequest="certificate-test",annotation_cert="test"} 1
`,
			MetricNames: []string{"kube_certificatesigningrequest_created", "kube_certificatesigningrequest_condition", "kube_certificatesigningrequest_labels", "kube_certificatesigningrequest_cert_length", "kube_certificatesigningrequest_annotations"},
		},
		{
			Obj: &certv1beta1.CertificateSigningRequest{
					Labels: map[string]string{
						"cert": "test",
					},
					Annotations: map[string]string{
						"cert": "test",
					},
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
				},
				Status: certv1beta1.CertificateSigningRequestStatus{
				kube_certificatesigningrequest_condition{certificatesigningrequest="certificate-test",condition="denied"} 1
				kube_certificatesigningrequest_labels{certificatesigningrequest="certificate-test",label_cert="test"} 1
				kube_certificatesigningrequest_cert_length{certificatesigningrequest="certificate-test"} 0
				kube_certificatesigningrequest_annotations{certificatesigningrequest="certificate-test",annotation_cert="test"} 1
`,
			MetricNames: []string{"kube_certificatesigningrequest_created", "kube_certificatesigningrequest_condition", "kube_certificatesigningrequest_labels", "kube_certificatesigningrequest_cert_length", "kube_certificatesigningrequest_annotations"},
		},
		{
			Obj: &certv1beta1.CertificateSigningRequest{
					Labels: map[string]string{
						"cert": "test",
					},
					Annotations: map[string]string{
						"cert": "test",
					},
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
				},
				Status: certv1beta1.CertificateSigningRequestStatus{
				kube_certificatesigningrequest_condition{certificatesigningrequest="certificate-test",condition="denied"} 2
				kube_certificatesigningrequest_labels{certificatesigningrequest="certificate-test",label_cert="test"} 1
				kube_certificatesigningrequest_cert_length{certificatesigningrequest="certificate-test"} 0
				kube_certificatesigningrequest_annotations{certificatesigningrequest="certificate-test",annotation_cert="test"} 1
`,
			MetricNames: []string{"kube_certificatesigningrequest_created", "kube_certificatesigningrequest_condition", "kube_certificatesigningrequest_labels", "kube_certificatesigningrequest_cert_length", "kube_certificatesigningrequest_annotations"},
		},
	}
	for i, c := range cases {
				}
			}),
		},
		{
			Name: "kube_configmap_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapConfigMapFunc(func(c *v1.ConfigMap) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(c.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

	metav1StartTime := metav1.Unix(int64(startTime), 0)

	const metadata = `
		# HELP kube_configmap_info Information about configmap.
		# TYPE kube_configmap_info gauge
		# HELP kube_configmap_created Unix creation timestamp
		# TYPE kube_configmap_created gauge
		# HELP kube_configmap_metadata_resource_version Resource version representing a specific version of the configmap.
		# TYPE kube_configmap_metadata_resource_version gauge
		# HELP kube_configmap_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_configmap_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Name:            "configmap1",
					Namespace:       "ns1",
					ResourceVersion: "123456",
					Annotations: map[string]string{
						"configmap": "test",
					},
				},
			},
			Want: `
				kube_configmap_info{configmap="configmap1",namespace="ns1"} 1
				kube_configmap_metadata_resource_version{configmap="configmap1",namespace="ns1",resource_version="123456"} 1
				kube_configmap_annotations{configmap="configmap1",namespace="ns1",annotation_configmap="test"} 1
`,
			MetricNames: []string{"kube_configmap_info", "kube_configmap_metadata_resource_version", "kube_configmap_annotations"},
		},
		{
			Obj: &v1.ConfigMap{
					Namespace:         "ns2",
					CreationTimestamp: metav1StartTime,
					ResourceVersion:   "abcdef",
					Annotations: map[string]string{
						"configmap": "test2",
					},
				},
			},
			Want: `
				kube_configmap_info{configmap="configmap2",namespace="ns2"} 1
				kube_configmap_created{configmap="configmap2",namespace="ns2"} 1.501569018e+09
				kube_configmap_metadata_resource_version{configmap="configmap2",namespace="ns2",resource_version="abcdef"} 1
				kube_configmap_annotations{configmap="configmap2",namespace="ns2",annotation_configmap="test2"} 1
				`,
			MetricNames: []string{"kube_configmap_info", "kube_configmap_created", "kube_configmap_metadata_resource_version", "kube_configmap_annotations"},
		},
	}
	for i, c := range cases {
				}
			}),
		},
		{
			Name: "kube_cronjob_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapCronJobFunc(func(j *batchv1beta1.CronJob) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(j.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_cronjob_status_last_schedule_time gauge
		# HELP kube_cronjob_next_schedule_time Next time the cronjob should be scheduled. The time after lastScheduleTime, or after the cron job's creation time if it's never been scheduled. Use this to determine if the job is delayed.
		# TYPE kube_cronjob_next_schedule_time gauge
		# HELP kube_cronjob_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_cronjob_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Labels: map[string]string{
						"app": "example-active-running-1",
					},
					Annotations: map[string]string{
						"app": "example-active-running-1",
					},
				},
				Status: batchv1beta1.CronJobStatus{
					Active:           []v1.ObjectReference{{Name: "FakeJob1"}, {Name: "FakeJob2"}},
				kube_cronjob_spec_suspend{cronjob="ActiveRunningCronJob1",namespace="ns1"} 0
				kube_cronjob_status_active{cronjob="ActiveRunningCronJob1",namespace="ns1"} 2
				kube_cronjob_status_last_schedule_time{cronjob="ActiveRunningCronJob1",namespace="ns1"} 1.520742896e+09
				kube_cronjob_annotations{cronjob="ActiveRunningCronJob1",namespace="ns1",annotation_app="example-active-running-1"} 1
` + fmt.Sprintf("kube_cronjob_next_schedule_time{cronjob=\"ActiveRunningCronJob1\",namespace=\"ns1\"} %ve+09\n",
				float64(ActiveRunningCronJob1NextScheduleTime.Unix())/math.Pow10(9)),
			MetricNames: []string{"kube_cronjob_next_schedule_time", "kube_cronjob_spec_starting_deadline_seconds", "kube_cronjob_status_active", "kube_cronjob_spec_suspend", "kube_cronjob_info", "kube_cronjob_created", "kube_cronjob_labels", "kube_cronjob_status_last_schedule_time", "kube_cronjob_annotations"},
		},
		{
			Obj: &batchv1beta1.CronJob{
					Labels: map[string]string{
						"app": "example-suspended-1",
					},
					Annotations: map[string]string{
						"app": "example-suspended-1",
					},
				},
				Status: batchv1beta1.CronJobStatus{
					Active:           []v1.ObjectReference{},
				kube_cronjob_spec_suspend{cronjob="SuspendedCronJob1",namespace="ns1"} 1
				kube_cronjob_status_active{cronjob="SuspendedCronJob1",namespace="ns1"} 0
				kube_cronjob_status_last_schedule_time{cronjob="SuspendedCronJob1",namespace="ns1"} 1.520762696e+09
				kube_cronjob_annotations{cronjob="SuspendedCronJob1",namespace="ns1",annotation_app="example-suspended-1"} 1
`,
			MetricNames: []string{"kube_cronjob_spec_starting_deadline_seconds", "kube_cronjob_status_active", "kube_cronjob_spec_suspend", "kube_cronjob_info", "kube_cronjob_created", "kube_cronjob_labels", "kube_cronjob_status_last_schedule_time", "kube_cronjob_annotations"},
		},
		{
			Obj: &batchv1beta1.CronJob{
					Labels: map[string]string{
						"app": "example-active-no-last-scheduled-1",
					},
					Annotations: map[string]string{
						"app": "example-active-no-last-scheduled-1",
					},
				},
				Status: batchv1beta1.CronJobStatus{
					Active:           []v1.ObjectReference{},
				kube_cronjob_info{concurrency_policy="Forbid",cronjob="ActiveCronJob1NoLastScheduled",namespace="ns1",schedule="25 * * * *"} 1
				kube_cronjob_created{cronjob="ActiveCronJob1NoLastScheduled",namespace="ns1"} 1.520766296e+09
				kube_cronjob_labels{cronjob="ActiveCronJob1NoLastScheduled",label_app="example-active-no-last-scheduled-1",namespace="ns1"} 1
				kube_cronjob_annotations{cronjob="ActiveCronJob1NoLastScheduled",namespace="ns1",annotation_app="example-active-no-last-scheduled-1"} 1
` +
				fmt.Sprintf("kube_cronjob_next_schedule_time{cronjob=\"ActiveCronJob1NoLastScheduled\",namespace=\"ns1\"} %ve+09\n",
					float64(ActiveCronJob1NoLastScheduledNextScheduleTime.Unix())/math.Pow10(9)),
			// TODO: Do we need to specify metricnames?
			MetricNames: []string{"kube_cronjob_next_schedule_time", "kube_cronjob_spec_starting_deadline_seconds", "kube_cronjob_status_active", "kube_cronjob_spec_suspend", "kube_cronjob_info", "kube_cronjob_created", "kube_cronjob_labels", "kube_cronjob_annotations"},
		},
	}
	for i, c := range cases {
				}
			}),
		},
		{
			Name: "kube_daemonset_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapDaemonSetFunc(func(d *v1.DaemonSet) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(d.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_daemonset_updated_number_scheduled gauge
		# HELP kube_daemonset_labels Kubernetes labels converted to Prometheus labels.
		# TYPE kube_daemonset_labels gauge
		# HELP kube_daemonset_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_daemonset_annotations gauge
`
	cases := []generateMetricsTestCase{
		{
					Labels: map[string]string{
						"app": "example1",
					},
					Annotations: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.DaemonSetStatus{
				kube_daemonset_status_number_unavailable{daemonset="ds1",namespace="ns1"} 0
				kube_daemonset_updated_number_scheduled{daemonset="ds1",namespace="ns1"} 0
				kube_daemonset_labels{daemonset="ds1",label_app="example1",namespace="ns1"} 1
				kube_daemonset_annotations{daemonset="ds1",namespace="ns1",annotation_app="example1"} 1
`,
			MetricNames: []string{
				"kube_daemonset_labels",
				"kube_daemonset_status_number_ready",
				"kube_daemonset_status_number_unavailable",
				"kube_daemonset_updated_number_scheduled",
				"kube_daemonset_annotations",
			},
		},
		{
					Labels: map[string]string{
						"app": "example2",
					},
					Annotations: map[string]string{
						"app": "example2",
					},
					Generation: 14,
				},
				Status: v1.DaemonSetStatus{
				kube_daemonset_status_number_unavailable{daemonset="ds2",namespace="ns2"} 0
				kube_daemonset_updated_number_scheduled{daemonset="ds2",namespace="ns2"} 0
				kube_daemonset_labels{daemonset="ds2",label_app="example2",namespace="ns2"} 1
				kube_daemonset_annotations{daemonset="ds2",namespace="ns2",annotation_app="example2"} 1
`,
			MetricNames: []string{
				"kube_daemonset_created",
				"kube_daemonset_status_number_ready",
				"kube_daemonset_status_number_unavailable",
				"kube_daemonset_updated_number_scheduled",
				"kube_daemonset_annotations",
			},
		},
		{
					Labels: map[string]string{
						"app": "example3",
					},
					Annotations: map[string]string{
						"app": "example3",
					},
					Generation: 15,
				},
				Status: v1.DaemonSetStatus{
				kube_daemonset_status_number_unavailable{daemonset="ds3",namespace="ns3"} 5
				kube_daemonset_updated_number_scheduled{daemonset="ds3",namespace="ns3"} 5
				kube_daemonset_labels{daemonset="ds3",label_app="example3",namespace="ns3"} 1
				kube_daemonset_annotations{daemonset="ds3",namespace="ns3",annotation_app="example3"} 1
`,
			MetricNames: []string{
				"kube_daemonset_created",
				"kube_daemonset_status_number_ready",
				"kube_daemonset_status_number_unavailable",
				"kube_daemonset_updated_number_scheduled",
				"kube_daemonset_annotations",
			},
		},
	}
				}
			}),
		},
		{
			Name: "kube_deployment_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(d.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_deployment_spec_strategy_rollingupdate_max_surge gauge
		# HELP kube_deployment_labels Kubernetes labels converted to Prometheus labels.
		# TYPE kube_deployment_labels gauge
		# HELP kube_deployment_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_deployment_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Labels: map[string]string{
						"app": "example1",
					},
					Annotations: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.DeploymentStatus{
        kube_deployment_status_replicas_unavailable{deployment="depl1",namespace="ns1"} 5
        kube_deployment_status_replicas_updated{deployment="depl1",namespace="ns1"} 2
        kube_deployment_status_replicas{deployment="depl1",namespace="ns1"} 15
		kube_deployment_annotations{annotation_app="example1",deployment="depl1",namespace="ns1"} 1
`,
		},
		{
					Labels: map[string]string{
						"app": "example2",
					},
					Annotations: map[string]string{
						"app": "example2",
					},
					Generation: 14,
				},
				Status: v1.DeploymentStatus{
				},
			},
			Want: `
       	kube_deployment_labels{deployment="depl2",label_app="example2",namespace="ns2"} 1
        kube_deployment_metadata_generation{deployment="depl2",namespace="ns2"} 14
        kube_deployment_spec_paused{deployment="depl2",namespace="ns2"} 1
        kube_deployment_spec_replicas{deployment="depl2",namespace="ns2"} 5
        kube_deployment_status_replicas_unavailable{deployment="depl2",namespace="ns2"} 0
        kube_deployment_status_replicas_updated{deployment="depl2",namespace="ns2"} 1
        kube_deployment_status_replicas{deployment="depl2",namespace="ns2"} 10
        kube_deployment_annotations{annotation_app="example2",deployment="depl2",namespace="ns2"} 1
`,
		},
	}
				}
			}),
		},
		{
			Name: "kube_endpoint_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapEndpointFunc(func(e *v1.Endpoints) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(e.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_endpoint_info gauge
		# HELP kube_endpoint_labels Kubernetes labels converted to Prometheus labels.
		# TYPE kube_endpoint_labels gauge
		# HELP kube_endpoint_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_endpoint_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Labels: map[string]string{
						"app": "foobar",
					},
					Annotations: map[string]string{
						"app": "foobar",
					},
				},
				Subsets: []v1.EndpointSubset{
					{Addresses: []v1.EndpointAddress{
				kube_endpoint_created{endpoint="test-endpoint",namespace="default"} 1.5e+09
				kube_endpoint_info{endpoint="test-endpoint",namespace="default"} 1
				kube_endpoint_labels{endpoint="test-endpoint",label_app="foobar",namespace="default"} 1
				kube_endpoint_annotations{endpoint="test-endpoint",namespace="default",annotation_app="foobar"} 1
			`,
		},
	}
				}
			}),
		},
		{
			Name: "kube_hpa_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapHPAFunc(func(a *autoscaling.HorizontalPodAutoscaler) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(a.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

        # TYPE kube_hpa_status_condition gauge
        # HELP kube_hpa_labels Kubernetes labels converted to Prometheus labels.
        # TYPE kube_hpa_labels gauge
		# HELP kube_hpa_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_hpa_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Labels: map[string]string{
						"app": "foobar",
					},
					Annotations: map[string]string{
						"app": "foobar",
					},
				},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					MaxReplicas: 4,
                kube_hpa_status_condition{condition="unknown",hpa="hpa1",namespace="ns1",status="AbleToScale"} 0
				kube_hpa_status_current_replicas{hpa="hpa1",namespace="ns1"} 2
				kube_hpa_status_desired_replicas{hpa="hpa1",namespace="ns1"} 2
				kube_hpa_annotations{hpa="hpa1",namespace="ns1",annotation_app="foobar"} 1
			`,
			MetricNames: []string{
				"kube_hpa_metadata_generation",
				"kube_hpa_status_desired_replicas",
				"kube_hpa_status_condition",
				"kube_hpa_labels",
				"kube_hpa_annotations",
			},
		},
	}
				}
			}),
		},
		{
			Name: "kube_ingress_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapIngressFunc(func(i *v1beta1.Ingress) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(i.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_ingress_metadata_resource_version gauge
		# HELP kube_ingress_path Ingress host, paths and backend service.
		# TYPE kube_ingress_path gauge
		# HELP kube_ingress_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_ingress_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Name:            "ingress1",
					Namespace:       "ns1",
					ResourceVersion: "000000",
					Annotations: map[string]string{
						"app": "ingress1",
					},
				},
			},
			Want: `
				kube_ingress_info{namespace="ns1",ingress="ingress1"} 1
				kube_ingress_metadata_resource_version{namespace="ns1",resource_version="000000",ingress="ingress1"} 1
				kube_ingress_labels{namespace="ns1",ingress="ingress1"} 1
				kube_ingress_annotations{namespace="ns1",ingress="ingress1",annotation_app="ingress1"} 1
`,
			MetricNames: []string{"kube_ingress_info", "kube_ingress_metadata_resource_version", "kube_ingress_created", "kube_ingress_labels", "kube_ingress_path", "kube_ingress_annotations"},
		},
		{
			Obj: &v1beta1.Ingress{
					Namespace:         "ns2",
					CreationTimestamp: metav1StartTime,
					ResourceVersion:   "123456",
					Annotations: map[string]string{
						"app": "ingress2",
					},
				},
			},
			Want: `
				kube_ingress_created{namespace="ns2",ingress="ingress2"} 1.501569018e+09
				kube_ingress_metadata_resource_version{namespace="ns2",resource_version="123456",ingress="ingress2"} 1
				kube_ingress_labels{namespace="ns2",ingress="ingress2"} 1
				kube_ingress_annotations{namespace="ns2",ingress="ingress2",annotation_app="ingress2"} 1
				`,
			MetricNames: []string{"kube_ingress_info", "kube_ingress_metadata_resource_version", "kube_ingress_created", "kube_ingress_labels", "kube_ingress_path", "kube_ingress_annotations"},
		},
		{
			Obj: &v1beta1.Ingress{
					CreationTimestamp: metav1StartTime,
					Labels:            map[string]string{"test-3": "test-3"},
					ResourceVersion:   "abcdef",
					Annotations: map[string]string{
						"test-3": "test-3",
					},
				},
			},
			Want: `
				kube_ingress_created{namespace="ns3",ingress="ingress3"} 1.501569018e+09
				kube_ingress_metadata_resource_version{namespace="ns3",resource_version="abcdef",ingress="ingress3"} 1
				kube_ingress_labels{label_test_3="test-3",namespace="ns3",ingress="ingress3"} 1
				kube_ingress_annotations{namespace="ns3",ingress="ingress3",annotation_test_3="test-3"} 1
`,
			MetricNames: []string{"kube_ingress_info", "kube_ingress_metadata_resource_version", "kube_ingress_created", "kube_ingress_labels", "kube_ingress_path", "kube_ingress_annotations"},
		},
		{
			Obj: &v1beta1.Ingress{
					CreationTimestamp: metav1StartTime,
					Labels:            map[string]string{"test-4": "test-4"},
					ResourceVersion:   "abcdef",
					Annotations: map[string]string{
						"test-4": "test-4",
					},
				},
				Spec: v1beta1.IngressSpec{
					Rules: []v1beta1.IngressRule{
				kube_ingress_metadata_resource_version{namespace="ns4",resource_version="abcdef",ingress="ingress4"} 1
				kube_ingress_labels{label_test_4="test-4",namespace="ns4",ingress="ingress4"} 1
				kube_ingress_path{namespace="ns4",ingress="ingress4",host="somehost",path="/somepath",service_name="someservice",service_port="1234"} 1
				kube_ingress_annotations{namespace="ns4",ingress="ingress4",annotation_test_4="test-4"} 1
`,
			MetricNames: []string{"kube_ingress_info", "kube_ingress_metadata_resource_version", "kube_ingress_created", "kube_ingress_labels", "kube_ingress_path", "kube_ingress_annotations"},
		},
	}
	for i, c := range cases {
				}
			}),
		},
		{
			Name: "kube_job_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapJobFunc(func(j *v1batch.Job) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(j.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_job_status_start_time gauge
		# HELP kube_job_status_succeeded The number of pods which reached Phase Succeeded.
		# TYPE kube_job_status_succeeded gauge
		# HELP kube_job_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_job_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Labels: map[string]string{
						"app": "example-running-1",
					},
					Annotations: map[string]string{
						"app": "example-running-1",
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind:       "CronJob",
				kube_job_status_failed{job_name="RunningJob1",namespace="ns1"} 0
				kube_job_status_start_time{job_name="RunningJob1",namespace="ns1"} 1.495800007e+09
				kube_job_status_succeeded{job_name="RunningJob1",namespace="ns1"} 0
				kube_job_annotations{job_name="RunningJob1",namespace="ns1",annotation_app="example-running-1"} 1
`,
		},
		{
					Labels: map[string]string{
						"app": "example-successful-1",
					},
					Annotations: map[string]string{
						"app": "example-successful-1",
					},
				},
				Status: v1batch.JobStatus{
					Active:         0,
				kube_job_status_failed{job_name="SuccessfulJob1",namespace="ns1"} 0
				kube_job_status_start_time{job_name="SuccessfulJob1",namespace="ns1"} 1.495800007e+09
				kube_job_status_succeeded{job_name="SuccessfulJob1",namespace="ns1"} 1
				kube_job_annotations{job_name="SuccessfulJob1",namespace="ns1",annotation_app="example-successful-1"} 1
`,
		},
		{
					Labels: map[string]string{
						"app": "example-failed-1",
					},
					Annotations: map[string]string{
						"app": "example-failed-1",
					},
				},
				Status: v1batch.JobStatus{
					Active:         0,
				kube_job_status_failed{job_name="FailedJob1",namespace="ns1"} 1
				kube_job_status_start_time{job_name="FailedJob1",namespace="ns1"} 1.495807207e+09
				kube_job_status_succeeded{job_name="FailedJob1",namespace="ns1"} 0
				kube_job_annotations{job_name="FailedJob1",namespace="ns1",annotation_app="example-failed-1"} 1
`,
		},
		{
					Labels: map[string]string{
						"app": "example-successful-2",
					},
					Annotations: map[string]string{
						"app": "example-successful-2",
					},
				},
				Status: v1batch.JobStatus{
					Active:         0,
				kube_job_owner{job_name="SuccessfulJob2NoActiveDeadlineSeconds",namespace="ns1",owner_is_controller="<none>",owner_kind="<none>",owner_name="<none>"} 1
				kube_job_complete{condition="false",job_name="SuccessfulJob2NoActiveDeadlineSeconds",namespace="ns1"} 0
				kube_job_complete{condition="true",job_name="SuccessfulJob2NoActiveDeadlineSeconds",namespace="ns1"} 1
				kube_job_complete{condition="unknown",job_name="SuccessfulJob2NoActiveDeadlineSeconds",namespace="ns1"} 0
				kube_job_info{job_name="SuccessfulJob2NoActiveDeadlineSeconds",namespace="ns1"} 1
				kube_job_labels{job_name="SuccessfulJob2NoActiveDeadlineSeconds",label_app="example-successful-2",namespace="ns1"} 1
				kube_job_status_failed{job_name="SuccessfulJob2NoActiveDeadlineSeconds",namespace="ns1"} 0
				kube_job_status_start_time{job_name="SuccessfulJob2NoActiveDeadlineSeconds",namespace="ns1"} 1.495800607e+09
				kube_job_status_succeeded{job_name="SuccessfulJob2NoActiveDeadlineSeconds",namespace="ns1"} 1
				kube_job_annotations{job_name="SuccessfulJob2NoActiveDeadlineSeconds",namespace="ns1",annotation_app="example-successful-2"} 1
`,
		},
	}
				}
			}),
		},
		{
			Name: "kube_limitrange_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapLimitRangeFunc(func(r *v1.LimitRange) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(r.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

	# TYPE kube_limitrange_created gauge
	# HELP kube_limitrange Information about limit range.
	# TYPE kube_limitrange gauge
	# HELP kube_limitrange_annotations Kubernetes annotations converted to Prometheus labels.
	# TYPE kube_limitrange_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Name:              "quotaTest",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "testNS",
					Annotations: map[string]string{
						"quota": "test",
					},
				},
				Spec: v1.LimitRangeSpec{
					Limits: []v1.LimitRangeItem{
        kube_limitrange{constraint="max",limitrange="quotaTest",namespace="testNS",resource="memory",type="Pod"} 2.1e+09
        kube_limitrange{constraint="maxLimitRequestRatio",limitrange="quotaTest",namespace="testNS",resource="memory",type="Pod"} 2.1e+09
        kube_limitrange{constraint="min",limitrange="quotaTest",namespace="testNS",resource="memory",type="Pod"} 2.1e+09
		kube_limitrange_annotations{limitrange="quotaTest",namespace="testNS",annotation_quota="test"} 1
		`,
		},
	}
	descNamespaceLabelsHelp          = "Kubernetes labels converted to Prometheus labels."
	descNamespaceLabelsDefaultLabels = []string{"namespace"}

	namespaceMetricFamilies = []metric.FamilyGenerator{
		{
			Name: "kube_namespace_created",
			}),
		},
		{
			Name: "kube_namespace_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapNamespaceFunc(func(n *v1.Namespace) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(n.Annotations)
				return &metric.Family{
			Obj: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nsActiveTest",
					Annotations: map[string]string{
						"test": "testNS",
					},
				},
				Spec: v1.NamespaceSpec{
					Finalizers: []v1.FinalizerName{v1.FinalizerKubernetes},
			},
			Want: `
				kube_namespace_labels{namespace="nsActiveTest"} 1
				kube_namespace_status_phase{namespace="nsActiveTest",phase="Active"} 1
				kube_namespace_status_phase{namespace="nsActiveTest",phase="Terminating"} 0
				kube_namespace_annotations{namespace="nsActiveTest",annotation_test="testNS"} 1
`,
		},
		{
			Obj: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nsTerminateTest",
					Annotations: map[string]string{
						"test": "nsTerminateTest",
					},
				},
				Spec: v1.NamespaceSpec{
					Finalizers: []v1.FinalizerName{v1.FinalizerKubernetes},
			},
			Want: `
				kube_namespace_labels{namespace="nsTerminateTest"} 1
				kube_namespace_status_phase{namespace="nsTerminateTest",phase="Active"} 0
				kube_namespace_status_phase{namespace="nsTerminateTest",phase="Terminating"} 1
				kube_namespace_annotations{namespace="nsTerminateTest",annotation_test="nsTerminateTest"} 1
`,
		},
		{
			Want: `
				kube_namespace_created{namespace="ns1"} 1.5e+09
				kube_namespace_labels{label_app="example1",namespace="ns1"} 1
				kube_namespace_status_phase{namespace="ns1",phase="Active"} 1
				kube_namespace_status_phase{namespace="ns1",phase="Terminating"} 0
				kube_namespace_annotations{annotation_app="example1",namespace="ns1"} 1
`,
		},
		{
				}
			}),
		},
		{
			Name: "kube_node_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapNodeFunc(func(n *v1.Node) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(n.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# HELP kube_node_status_allocatable_memory_bytes The memory resources of a node that are available for scheduling.
		# HELP kube_node_status_condition The condition of a cluster node.
		# TYPE kube_node_status_condition gauge
		# HELP kube_node_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_node_annotations gauge
	`
	cases := []generateMetricsTestCase{
		// Verify populating base metric and that metric for unset fields are skipped.
			Obj: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "127.0.0.1",
					Annotations: map[string]string{
						"kubeadm.alpha.kubernetes.io/ttl": "0",
					},
				},
				Status: v1.NodeStatus{
					NodeInfo: v1.NodeSystemInfo{
				kube_node_info{container_runtime_version="rkt",kernel_version="kernel",kubelet_version="kubelet",kubeproxy_version="kubeproxy",node="127.0.0.1",os_image="osimage",provider_id="provider://i-uniqueid"} 1
				kube_node_labels{node="127.0.0.1"} 1
				kube_node_spec_unschedulable{node="127.0.0.1"} 0
				kube_node_annotations{node="127.0.0.1",annotation_kubeadm_alpha_kubernetes_io_ttl="0"} 1
			`,
		},
		// Verify resource metric.
        kube_node_status_capacity{node="127.0.0.1",resource="nvidia_com_gpu",unit="integer"} 4
        kube_node_status_capacity{node="127.0.0.1",resource="pods",unit="integer"} 1000
        kube_node_status_capacity{node="127.0.0.1",resource="storage",unit="byte"} 3e+09
		kube_node_annotations{node="127.0.0.1"} 1
			`,
		},
		// Verify phase enumerations.
				}
			}),
		},
		{
			Name: "kube_persistentvolume_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapPersistentVolumeFunc(func(p *v1.PersistentVolume) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(p.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

			# TYPE kube_persistentvolume_info gauge
			# HELP kube_persistentvolume_capacity_bytes The size of the Persistentvolume in bytes.
			# TYPE kube_persistentvolume_capacity_bytes gauge
			# HELP kube_persistentvolume_annotations Kubernetes annotations converted to Prometheus labels.
			# TYPE kube_persistentvolume_annotations gauge
	`
	cases := []generateMetricsTestCase{
		// Verify phase enumerations.
			Obj: &v1.PersistentVolume{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pv-available",
					Annotations: map[string]string{
						"pv": "test",
					},
				},
				Status: v1.PersistentVolumeStatus{
					Phase: v1.VolumeAvailable,
					kube_persistentvolume_status_phase{persistentvolume="test-pv-available",phase="Failed"} 0
					kube_persistentvolume_status_phase{persistentvolume="test-pv-available",phase="Pending"} 0
					kube_persistentvolume_status_phase{persistentvolume="test-pv-available",phase="Released"} 0
					kube_persistentvolume_annotations{persistentvolume="test-pv-available",annotation_pv="test"} 1
`,
			MetricNames: []string{
				"kube_persistentvolume_status_phase",
				"kube_persistentvolume_annotations",
			},
		},
		{
			Obj: &v1.PersistentVolume{
			Obj: &v1.PersistentVolume{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pv-pending",
					Annotations: map[string]string{
						"pv-test": "test",
					},
				},
				Status: v1.PersistentVolumeStatus{
					Phase: v1.VolumePending,
					kube_persistentvolume_status_phase{persistentvolume="test-pv-pending",phase="Failed"} 0
					kube_persistentvolume_status_phase{persistentvolume="test-pv-pending",phase="Pending"} 1
					kube_persistentvolume_status_phase{persistentvolume="test-pv-pending",phase="Released"} 0
					kube_persistentvolume_annotations{persistentvolume="test-pv-pending",annotation_pv_test="test"} 1
`,
			MetricNames: []string{
				"kube_persistentvolume_status_phase",
				"kube_persistentvolume_annotations",
			},
		},
		{
				}
			}),
		},
		{
			Name: "kube_persistentvolumeclaim_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapPersistentVolumeClaimFunc(func(p *v1.PersistentVolumeClaim) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(p.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_persistentvolumeclaim_resource_requests_storage_bytes gauge
		# HELP kube_persistentvolumeclaim_access_mode The access mode of the persistent volume.
		# TYPE kube_persistentvolumeclaim_access_mode gauge
		# HELP kube_persistentvolumeclaim_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_persistentvolumeclaim_annotations gauge
	`
	storageClassName := "rbd"
	cases := []generateMetricsTestCase{
					Labels: map[string]string{
						"app": "mysql-server",
					},
					Annotations: map[string]string{
						"app": "mysql-server",
					},
				},
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
				kube_persistentvolumeclaim_resource_requests_storage_bytes{namespace="default",persistentvolumeclaim="mysql-data"} 1.073741824e+09
				kube_persistentvolumeclaim_labels{label_app="mysql-server",namespace="default",persistentvolumeclaim="mysql-data"} 1
				kube_persistentvolumeclaim_access_mode{namespace="default",persistentvolumeclaim="mysql-data",access_mode="ReadWriteOnce"} 1
				kube_persistentvolumeclaim_annotations{namespace="default",persistentvolumeclaim="mysql-data",annotation_app="mysql-server"} 1
`,
			MetricNames: []string{"kube_persistentvolumeclaim_info", "kube_persistentvolumeclaim_status_phase", "kube_persistentvolumeclaim_resource_requests_storage_bytes", "kube_persistentvolumeclaim_labels", "kube_persistentvolumeclaim_access_mode", "kube_persistentvolumeclaim_annotations"},
		},
		{
			Obj: &v1.PersistentVolumeClaim{
	for n := 0; n < b.N; n++ {
		families := f(pod)
		if len(families) != 38 {
			b.Fatalf("expected 39 but got %v", len(families))
		}
	}
}
				}
			}),
		},
		{
			Name: "kube_poddisruptionbudget_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapPodDisruptionBudgetFunc(func(p *v1beta1.PodDisruptionBudget) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(p.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

	# TYPE kube_poddisruptionbudget_status_expected_pods gauge
	# HELP kube_poddisruptionbudget_status_observed_generation Most recent generation observed when updating this PDB status
	# TYPE kube_poddisruptionbudget_status_observed_generation gauge
	# HELP kube_poddisruptionbudget_annotations Kubernetes annotations converted to Prometheus labels.
	# TYPE kube_poddisruptionbudget_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Generation:        21,
					Annotations: map[string]string{
						"pdb1": "pdb1",
					},
				},
				Status: v1beta1.PodDisruptionBudgetStatus{
					CurrentHealthy:        12,
			kube_poddisruptionbudget_status_pod_disruptions_allowed{namespace="ns1",poddisruptionbudget="pdb1"} 2
			kube_poddisruptionbudget_status_expected_pods{namespace="ns1",poddisruptionbudget="pdb1"} 15
			kube_poddisruptionbudget_status_observed_generation{namespace="ns1",poddisruptionbudget="pdb1"} 111
			kube_poddisruptionbudget_annotations{namespace="ns1",poddisruptionbudget="pdb1",annotation_pdb1="pdb1"} 1
			`,
		},
		{
				kube_poddisruptionbudget_status_pod_disruptions_allowed{namespace="ns2",poddisruptionbudget="pdb2"} 0
				kube_poddisruptionbudget_status_expected_pods{namespace="ns2",poddisruptionbudget="pdb2"} 10
				kube_poddisruptionbudget_status_observed_generation{namespace="ns2",poddisruptionbudget="pdb2"} 1111
				kube_poddisruptionbudget_annotations{namespace="ns2",poddisruptionbudget="pdb2"} 1
			`,
		},
	}
				}
			}),
		},
		{
			Name: "kube_replicaset_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapReplicaSetFunc(func(d *v1.ReplicaSet) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(d.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_replicaset_owner gauge
		# HELP kube_replicaset_labels Kubernetes labels converted to Prometheus labels.
		# TYPE kube_replicaset_labels gauge
		# HELP kube_replicaset_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_replicaset_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Labels: map[string]string{
						"app": "example1",
					},
					Annotations: map[string]string{
						"app": "example1",
					},
				},
				Status: v1.ReplicaSetStatus{
					Replicas:             5,
				kube_replicaset_status_ready_replicas{namespace="ns1",replicaset="rs1"} 5
				kube_replicaset_spec_replicas{namespace="ns1",replicaset="rs1"} 5
				kube_replicaset_owner{namespace="ns1",owner_is_controller="true",owner_kind="Deployment",owner_name="dp-name",replicaset="rs1"} 1
				kube_replicaset_annotations{namespace="ns1",replicaset="rs1",annotation_app="example1"} 1
`,
		},
		{
						"app": "example2",
						"env": "ex",
					},
					Annotations: map[string]string{
						"app": "example2",
						"env": "ex",
					},
				},
				Status: v1.ReplicaSetStatus{
					Replicas:             0,
				kube_replicaset_status_ready_replicas{namespace="ns2",replicaset="rs2"} 0
				kube_replicaset_spec_replicas{namespace="ns2",replicaset="rs2"} 0
				kube_replicaset_owner{namespace="ns2",owner_is_controller="<none>",owner_kind="<none>",owner_name="<none>",replicaset="rs2"} 1
				kube_replicaset_annotations{namespace="ns2",replicaset="rs2",annotation_app="example2",annotation_env="ex"} 1
			`,
		},
	}
				}
			}),
		},
		{
			Name: "kube_replicationcontroller_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapReplicationControllerFunc(func(r *v1.ReplicationController) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(r.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_replicationcontroller_status_observed_generation gauge
		# HELP kube_replicationcontroller_spec_replicas Number of desired pods for a ReplicationController.
		# TYPE kube_replicationcontroller_spec_replicas gauge
		# HELP kube_replicationcontroller_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_replicationcontroller_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Generation:        21,
					Annotations: map[string]string{
						"app": "example1",
					},
				},
				Status: v1.ReplicationControllerStatus{
					Replicas:             5,
				kube_replicationcontroller_status_ready_replicas{namespace="ns1",replicationcontroller="rc1"} 5
				kube_replicationcontroller_status_available_replicas{namespace="ns1",replicationcontroller="rc1"} 3
				kube_replicationcontroller_spec_replicas{namespace="ns1",replicationcontroller="rc1"} 5
				kube_replicationcontroller_annotations{namespace="ns1",replicationcontroller="rc1",annotation_app="example1"} 1
`,
		},
		{
					Name:       "rc2",
					Namespace:  "ns2",
					Generation: 14,
					Annotations: map[string]string{
						"app": "rc2",
					},
				},
				Status: v1.ReplicationControllerStatus{
					Replicas:             0,
				kube_replicationcontroller_status_ready_replicas{namespace="ns2",replicationcontroller="rc2"} 0
				kube_replicationcontroller_status_available_replicas{namespace="ns2",replicationcontroller="rc2"} 0
				kube_replicationcontroller_spec_replicas{namespace="ns2",replicationcontroller="rc2"} 0
				kube_replicationcontroller_annotations{namespace="ns2",replicationcontroller="rc2",annotation_app="rc2"} 1
`,
		},
	}
				}
			}),
		},
		{
			Name: "kube_resourcequota_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapResourceQuotaFunc(func(r *v1.ResourceQuota) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(r.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

	# TYPE kube_resourcequota gauge
	# HELP kube_resourcequota_created Unix creation timestamp
	# TYPE kube_resourcequota_created gauge
	# HELP kube_resourcequota_annotations Kubernetes annotations converted to Prometheus labels.
	# TYPE kube_resourcequota_annotations gauge
	`
	cases := []generateMetricsTestCase{
		// Verify populating base metric and that metric for unset fields are skipped.
					Name:              "quotaTest",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "testNS",
					Annotations: map[string]string{
						"rq": "rq",
					},
				},
				Status: v1.ResourceQuotaStatus{},
			},
			Want: `
			kube_resourcequota_created{namespace="testNS",resourcequota="quotaTest"} 1.5e+09
			kube_resourcequota_annotations{namespace="testNS",resourcequota="quotaTest",annotation_rq="rq"} 1
			`,
		},
		// Verify resource metric.
				ObjectMeta: metav1.ObjectMeta{
					Name:      "quotaTest",
					Namespace: "testNS",
					Annotations: map[string]string{
						"RQ": "RQ",
					},
				},
				Spec: v1.ResourceQuotaSpec{
					Hard: v1.ResourceList{
			kube_resourcequota{namespace="testNS",resource="services.nodeports",resourcequota="quotaTest",type="used"} 1
			kube_resourcequota{namespace="testNS",resource="storage",resourcequota="quotaTest",type="hard"} 1e+10
			kube_resourcequota{namespace="testNS",resource="storage",resourcequota="quotaTest",type="used"} 9e+09
			kube_resourcequota_annotations{namespace="testNS",resourcequota="quotaTest",annotation_RQ="RQ"} 1
			`,
		},
	}
				}
			}),
		},
		{
			Name: "kube_secret_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapSecretFunc(func(s *v1.Secret) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(s.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_secret_created gauge
		# HELP kube_secret_metadata_resource_version Resource version representing a specific version of secret.
		# TYPE kube_secret_metadata_resource_version gauge
		# HELP kube_secret_annotations Kubernetes annotations converted to Prometheus labels
		# TYPE kube_secret_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					CreationTimestamp: metav1StartTime,
					Labels:            map[string]string{"test-3": "test-3"},
					ResourceVersion:   "abcdef",
					Annotations: map[string]string{
						"test-3": "test-3",
					},
				},
				Type: v1.SecretTypeDockercfg,
			},
				kube_secret_created{namespace="ns3",secret="secret3"} 1.501569018e+09
				kube_secret_metadata_resource_version{namespace="ns3",resource_version="abcdef",secret="secret3"} 1
				kube_secret_labels{label_test_3="test-3",namespace="ns3",secret="secret3"} 1
				kube_secret_annotations{namespace="ns3",secret="secret3",annotation_test_3="test-3"} 1
`,
			MetricNames: []string{"kube_secret_info", "kube_secret_metadata_resource_version", "kube_secret_created", "kube_secret_labels", "kube_secret_type", "kube_secret_annotations"},
		},
	}
	for i, c := range cases {
				}
			}),
		},
		{
			Name: "kube_service_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapSvcFunc(func(s *v1.Service) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(s.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

		# TYPE kube_service_spec_external_ip gauge
		# HELP kube_service_status_load_balancer_ingress Service load balancer ingress status
		# TYPE kube_service_status_load_balancer_ingress gauge
		# HELP kube_service_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_service_annotations gauge
	`
	cases := []generateMetricsTestCase{
		{
					Labels: map[string]string{
						"app": "example1",
					},
					Annotations: map[string]string{
						"app": "example1",
					},
				},
				Spec: v1.ServiceSpec{
					ClusterIP: "1.2.3.4",
				kube_service_info{cluster_ip="1.2.3.4",external_name="",load_balancer_ip="",namespace="default",service="test-service1"} 1
				kube_service_labels{label_app="example1",namespace="default",service="test-service1"} 1
				kube_service_spec_type{namespace="default",service="test-service1",type="ClusterIP"} 1
				kube_service_annotations{namespace="default",service="test-service1",annotation_app="example1"} 1
`,
			MetricNames: []string{
				"kube_service_created",
				"kube_service_info",
				"kube_service_labels",
				"kube_service_spec_type",
				"kube_service_annotations",
			},
		},
		{
					Labels: map[string]string{
						"app": "example2",
					},
					Annotations: map[string]string{
						"app": "example2",
					},
				},
				Spec: v1.ServiceSpec{
					ClusterIP: "1.2.3.5",
				kube_service_info{cluster_ip="1.2.3.5",external_name="",load_balancer_ip="",namespace="default",service="test-service2"} 1
				kube_service_labels{label_app="example2",namespace="default",service="test-service2"} 1
				kube_service_spec_type{namespace="default",service="test-service2",type="NodePort"} 1
				kube_service_annotations{namespace="default",service="test-service2",annotation_app="example2"} 1
`,
		},
		{
					Labels: map[string]string{
						"app": "example3",
					},
					Annotations: map[string]string{
						"app": "example3",
					},
				},
				Spec: v1.ServiceSpec{
					ClusterIP:      "1.2.3.6",
				kube_service_info{cluster_ip="1.2.3.6",external_name="",load_balancer_ip="1.2.3.7",namespace="default",service="test-service3"} 1
				kube_service_labels{label_app="example3",namespace="default",service="test-service3"} 1
				kube_service_spec_type{namespace="default",service="test-service3",type="LoadBalancer"} 1
				kube_service_annotations{namespace="default",service="test-service3",annotation_app="example3"} 1

`,
		},
		{
					Labels: map[string]string{
						"app": "example4",
					},
					Annotations: map[string]string{
						"app": "example4",
					},
				},
				Spec: v1.ServiceSpec{
					ExternalName: "www.example.com",
				kube_service_info{cluster_ip="",external_name="www.example.com",load_balancer_ip="",namespace="default",service="test-service4"} 1
				kube_service_labels{label_app="example4",namespace="default",service="test-service4"} 1
				kube_service_spec_type{namespace="default",service="test-service4",type="ExternalName"} 1
				kube_service_annotations{namespace="default",service="test-service4",annotation_app="example4"} 1
			`,
		},
		{
					Labels: map[string]string{
						"app": "example5",
					},
					Annotations: map[string]string{
						"app": "example5",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
				kube_service_labels{label_app="example5",namespace="default",service="test-service5"} 1
				kube_service_spec_type{namespace="default",service="test-service5",type="LoadBalancer"} 1
				kube_service_status_load_balancer_ingress{hostname="www.example.com",ip="1.2.3.8",namespace="default",service="test-service5"} 1
				kube_service_annotations{namespace="default",service="test-service5",annotation_app="example5"} 1
			`,
		},
		{
					Labels: map[string]string{
						"app": "example6",
					},
					Annotations: map[string]string{
						"app": "example6",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeClusterIP,
				kube_service_spec_type{namespace="default",service="test-service6",type="ClusterIP"} 1
				kube_service_spec_external_ip{external_ip="1.2.3.9",namespace="default",service="test-service6"} 1
				kube_service_spec_external_ip{external_ip="1.2.3.10",namespace="default",service="test-service6"} 1
				kube_service_annotations{namespace="default",service="test-service6",annotation_app="example6"} 1
			`,
		},
	}
				}
			}),
		},
		{
			Name: "kube_statefulset_annotations",
			Type: metric.Gauge,
			Help: "Kubernetes annotations converted to Prometheus labels.",
			GenerateFunc: wrapStatefulSetFunc(func(s *v1.StatefulSet) *metric.Family {
				annotationKeys, annotationValues := kubeAnnotationsToPrometheusLabels(s.Annotations)
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   annotationKeys,
							LabelValues: annotationValues,
							Value:       1,
						},
					},
				}
			}),
		},
	}
)

 		# TYPE kube_statefulset_metadata_generation gauge
		# HELP kube_statefulset_labels Kubernetes labels converted to Prometheus labels.
		# TYPE kube_statefulset_labels gauge
		# HELP kube_statefulset_annotations Kubernetes annotations converted to Prometheus labels.
		# TYPE kube_statefulset_annotations gauge
 	`
	cases := []generateMetricsTestCase{
		{
					Labels: map[string]string{
						"app": "example1",
					},
					Annotations: map[string]string{
						"app": "example1",
					},
					Generation: 3,
				},
				Spec: v1.StatefulSetSpec{
 				kube_statefulset_replicas{namespace="ns1",statefulset="statefulset1"} 3
 				kube_statefulset_metadata_generation{namespace="ns1",statefulset="statefulset1"} 3
				kube_statefulset_labels{label_app="example1",namespace="ns1",statefulset="statefulset1"} 1
				kube_statefulset_annotations{namespace="ns1",statefulset="statefulset1",annotation_app="example1"} 1
`,
			MetricNames: []string{
				"kube_statefulset_created",
				"kube_statefulset_status_replicas_updated",
				"kube_statefulset_status_update_revision",
				"kube_statefulset_status_current_revision",
				"kube_statefulset_annotations",
			},
		},
		{
					Labels: map[string]string{
						"app": "example2",
					},
					Annotations: map[string]string{
						"app": "example2",
					},
					Generation: 21,
				},
				Spec: v1.StatefulSetSpec{
 				kube_statefulset_metadata_generation{namespace="ns2",statefulset="statefulset2"} 21
				kube_statefulset_labels{label_app="example2",namespace="ns2",statefulset="statefulset2"} 1
				kube_statefulset_status_current_revision{namespace="ns2",revision="cr2",statefulset="statefulset2"} 1
				kube_statefulset_annotations{namespace="ns2",statefulset="statefulset2",annotation_app="example2"} 1
`,
			MetricNames: []string{
				"kube_statefulset_labels",
				"kube_statefulset_status_replicas_updated",
				"kube_statefulset_status_update_revision",
				"kube_statefulset_status_current_revision",
				"kube_statefulset_annotations",
			},
		},
		{
					Labels: map[string]string{
						"app": "example3",
					},
					Annotations: map[string]string{
						"app": "example3",
					},
					Generation: 36,
				},
				Spec: v1.StatefulSetSpec{
 				kube_statefulset_metadata_generation{namespace="ns3",statefulset="statefulset3"} 36
				kube_statefulset_labels{label_app="example3",namespace="ns3",statefulset="statefulset3"} 1
				kube_statefulset_status_current_revision{namespace="ns3",revision="cr3",statefulset="statefulset3"} 1
				kube_statefulset_annotations{namespace="ns3",statefulset="statefulset3",annotation_app="example3"} 1
 			`,
			MetricNames: []string{
				"kube_statefulset_labels",
				"kube_statefulset_status_replicas_updated",
				"kube_statefulset_status_update_revision",
				"kube_statefulset_status_current_revision",
				"kube_statefulset_annotations",
			},
		},
	}
			CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
			Namespace:         "default",
			UID:               types.UID("abc-" + i),
			Annotations: map[string]string{
				"pod": "pod",
			},
		},
		Spec: v1.PodSpec{
			NodeName: "node1",
