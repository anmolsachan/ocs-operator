{
  prometheusRules+:: {
    groups+: [
      {
        name: 'ODF_standardized_metrics.rules',
        rules: [
          {
            record: 'odf_system_health',
            expr: |||
              ceph_health_status
            ||| % $._config,
            labels: {
              subsystem_type : 'OCS',
            },
          },
        ],
      },
    ],
  },
}
