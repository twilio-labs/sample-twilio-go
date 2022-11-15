package metric

import "github.com/prometheus/client_golang/prometheus"

type InviteMetrics struct {
	Sent     prometheus.Counter
	Accepted prometheus.Counter
	Declined prometheus.Counter
}

func NewInviteMetrics(sent, accepted, declined prometheus.Counter) *InviteMetrics {
	return &InviteMetrics{sent, accepted, declined}
}
