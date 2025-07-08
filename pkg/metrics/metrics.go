package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// AuditLogsTotal counts the total number of processed audit logs
	AuditLogsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "audit_logs_total",
			Help: "Total number of audit logs processed",
		},
		[]string{"service_name", "action"},
	)

	// AuditLogsFailedTotal counts the total number of audit logs that failed to process
	AuditLogsFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "audit_logs_failed_total",
			Help: "Total number of audit logs failed to process",
		},
		[]string{"service_name", "action"},
	)
)

// RegisterMetrics registers custom metrics with Prometheus
func RegisterMetrics() {
	prometheus.MustRegister(AuditLogsTotal)
	prometheus.MustRegister(AuditLogsFailedTotal)
}
