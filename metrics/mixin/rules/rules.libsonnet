{
  prometheusRules+:: {
    groups+: [
      {
        name: 'ODF_standardized_metrics.rules',
        rules: [
          {
            record: 'odf_system_health_status',
            expr: |||
              ceph_health_status
            ||| % $._config,
            labels: {
              subsystem_vendor: 'Red Hat',
              subsystem_type : 'OCS',
            },
          },
          {
            record: 'odf_system_raw_capacity_total_bytes',
            expr: |||
              ceph_cluster_total_bytes
            ||| % $._config,
            labels: {
              subsystem_vendor: 'Red Hat',
              subsystem_type : 'OCS',
            },
          },
          {
            record: 'odf_system_raw_capacity_used_bytes',
            expr: |||
              ceph_cluster_total_used_raw_bytes
            ||| % $._config,
            labels: {
              subsystem_vendor: 'Red Hat',
              subsystem_type : 'OCS',
            },
          },
          {
            record: 'odf_system_iops_total_bytes',
            expr: |||
              sum by (namespace, subsystem_name, job, service) (rate(ceph_pool_wr[1m]) + rate(ceph_pool_rd[1m]))
            ||| % $._config,
            labels: {
              subsystem_vendor: 'Red Hat',
              subsystem_type : 'OCS',
            },
          },
          {
            record: 'odf_system_throughput_total_bytes',
            expr: |||
              sum by (namespace, subsystem_name, job, service) (rate(ceph_pool_wr_bytes[1m]) + rate(ceph_pool_rd_bytes[1m]))
            ||| % $._config,
            labels: {
              subsystem_vendor: 'Red Hat',
              subsystem_type : 'OCS',
            },
          },
          {
            record: 'odf_system_latency_seconds',
            expr: |||
              avg by (namespace, subsystem_name, job, service)
              (
                topk by (ceph_daemon) (1, label_replace(label_replace(ceph_disk_occupation{job="rook-ceph-mgr"}, "instance", "$1", "exported_instance", "(.*)"), "device", "$1", "device", "/dev/(.*)")) 
                * on(instance, device) group_right(ceph_daemon, subsystem_name) topk by (instance,device) 
                (1,
                  (
                    (  
                        rate(node_disk_read_time_seconds_total{device=~"sd.*", }[1m]) / (clamp_min(rate(node_disk_reads_completed_total[1m]), 1))
                    ) +
                    (
                        rate(node_disk_write_time_seconds_total{device=~"sd.*" }[1m]) / (clamp_min(rate(node_disk_reads_completed_total[1m]), 1))
                    )
                  )
                )
              )
            ||| % $._config,
            labels: {
              subsystem_vendor: 'Red Hat',
              subsystem_type : 'OCS',
            },
          },
        ],
      },
    ],
  },
}
