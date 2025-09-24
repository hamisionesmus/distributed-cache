package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all Prometheus metrics
type Metrics struct {
	// Cache metrics
	cacheHits         prometheus.Counter
	cacheMisses       prometheus.Counter
	cacheEvictions    prometheus.Counter
	cacheKeysTotal    prometheus.Gauge
	cacheMemoryUsage  prometheus.Gauge

	// Request metrics
	requestsTotal     *prometheus.CounterVec
	requestDuration   *prometheus.HistogramVec
	activeConnections prometheus.Gauge

	// Cluster metrics
	clusterNodes      prometheus.Gauge
	clusterReplicas   prometheus.Gauge
	clusterLeader     prometheus.Gauge

	// System metrics
	goRoutines        prometheus.Gauge
	memoryAllocated   prometheus.Gauge
	gcPauseTime       prometheus.Gauge

	// Custom metrics
	operationsTotal   *prometheus.CounterVec
	errorsTotal       *prometheus.CounterVec

	registry         *prometheus.Registry
	mu               sync.RWMutex
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	m := &Metrics{
		registry: prometheus.NewRegistry(),
	}

	m.initCacheMetrics()
	m.initRequestMetrics()
	m.initClusterMetrics()
	m.initSystemMetrics()
	m.initCustomMetrics()

	return m
}

// initCacheMetrics initializes cache-related metrics
func (m *Metrics) initCacheMetrics() {
	m.cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "Total number of cache hits",
	})
	m.cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "Total number of cache misses",
	})
	m.cacheEvictions = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_evictions_total",
		Help: "Total number of cache evictions",
	})
	m.cacheKeysTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cache_keys_total",
		Help: "Total number of keys in cache",
	})
	m.cacheMemoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cache_memory_usage_bytes",
		Help: "Current memory usage of cache",
	})

	m.registry.MustRegister(
		m.cacheHits,
		m.cacheMisses,
		m.cacheEvictions,
		m.cacheKeysTotal,
		m.cacheMemoryUsage,
	)
}

// initRequestMetrics initializes request-related metrics
func (m *Metrics) initRequestMetrics() {
	m.requestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "endpoint", "status"})

	m.requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "endpoint"})

	m.activeConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "active_connections",
		Help: "Number of active connections",
	})

	m.registry.MustRegister(
		m.requestsTotal,
		m.requestDuration,
		m.activeConnections,
	)
}

// initClusterMetrics initializes cluster-related metrics
func (m *Metrics) initClusterMetrics() {
	m.clusterNodes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cluster_nodes",
		Help: "Number of nodes in cluster",
	})
	m.clusterReplicas = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cluster_replicas",
		Help: "Number of replicas in cluster",
	})
	m.clusterLeader = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cluster_leader",
		Help: "Whether this node is the cluster leader (1=yes, 0=no)",
	})

	m.registry.MustRegister(
		m.clusterNodes,
		m.clusterReplicas,
		m.clusterLeader,
	)
}

// initSystemMetrics initializes system-related metrics
func (m *Metrics) initSystemMetrics() {
	m.goRoutines = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_goroutines",
		Help: "Number of goroutines",
	})
	m.memoryAllocated = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_memory_allocated_bytes",
		Help: "Allocated memory in bytes",
	})
	m.gcPauseTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_gc_pause_time_seconds",
		Help: "GC pause time in seconds",
	})

	m.registry.MustRegister(
		m.goRoutines,
		m.memoryAllocated,
		m.gcPauseTime,
	)
}

// initCustomMetrics initializes custom application metrics
func (m *Metrics) initCustomMetrics() {
	m.operationsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_operations_total",
		Help: "Total number of cache operations",
	}, []string{"operation", "result"})

	m.errorsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_errors_total",
		Help: "Total number of errors",
	}, []string{"type", "operation"})

	m.registry.MustRegister(
		m.operationsTotal,
		m.errorsTotal,
	)
}

// RecordCacheHit records a cache hit
func (m *Metrics) RecordCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheHits.Inc()
}

// RecordCacheMiss records a cache miss
func (m *Metrics) RecordCacheMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheMisses.Inc()
}

// RecordCacheEviction records a cache eviction
func (m *Metrics) RecordCacheEviction() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheEvictions.Inc()
}

// SetCacheKeys sets the total number of keys in cache
func (m *Metrics) SetCacheKeys(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheKeysTotal.Set(float64(count))
}

// SetCacheMemoryUsage sets the current memory usage
func (m *Metrics) SetCacheMemoryUsage(bytes int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheMemoryUsage.Set(float64(bytes))
}

// RecordRequest records an HTTP request
func (m *Metrics) RecordRequest(method, endpoint string, statusCode int, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	status := strconv.Itoa(statusCode)
	m.requestsTotal.WithLabelValues(method, endpoint, status).Inc()
	m.requestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// SetActiveConnections sets the number of active connections
func (m *Metrics) SetActiveConnections(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeConnections.Set(float64(count))
}

// SetClusterNodes sets the number of cluster nodes
func (m *Metrics) SetClusterNodes(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clusterNodes.Set(float64(count))
}

// SetClusterReplicas sets the number of cluster replicas
func (m *Metrics) SetClusterReplicas(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clusterReplicas.Set(float64(count))
}

// SetClusterLeader sets whether this node is the cluster leader
func (m *Metrics) SetClusterLeader(isLeader bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if isLeader {
		m.clusterLeader.Set(1)
	} else {
		m.clusterLeader.Set(0)
	}
}

// UpdateSystemMetrics updates Go runtime metrics
func (m *Metrics) UpdateSystemMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Note: In a real implementation, you'd use runtime metrics
	// For now, we'll use placeholder values
	m.goRoutines.Set(42) // runtime.NumGoroutine()
	m.memoryAllocated.Set(1024 * 1024 * 50) // runtime memory stats
	m.gcPauseTime.Set(0.001) // GC pause time
}

// RecordOperation records a cache operation
func (m *Metrics) RecordOperation(operation, result string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.operationsTotal.WithLabelValues(operation, result).Inc()
}

// RecordError records an error
func (m *Metrics) RecordError(errorType, operation string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorsTotal.WithLabelValues(errorType, operation).Inc()
}

// StartMetricsServer starts the metrics HTTP server
func (m *Metrics) StartMetricsServer(port int) error {
	http.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/health", m.healthHandler)
	http.HandleFunc("/status", m.statusHandler)

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}

// healthHandler handles health check requests
func (m *Metrics) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}

// statusHandler handles status requests
func (m *Metrics) statusHandler(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"cache": map[string]interface{}{
			"hits":    m.cacheHits.Desc().String(),
			"misses":  m.cacheMisses.Desc().String(),
			"keys":    m.cacheKeysTotal.Desc().String(),
			"memory":  m.cacheMemoryUsage.Desc().String(),
		},
		"system": map[string]interface{}{
			"goroutines": m.goRoutines.Desc().String(),
			"memory":     m.memoryAllocated.Desc().String(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In a real implementation, you'd marshal the status map
	w.Write([]byte(`{"status": "ok"}`))
}

// GetMetricsSummary returns a summary of current metrics
func (m *Metrics) GetMetricsSummary() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Gather metrics from the registry
	metricsFamilies, err := m.registry.Gather()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	summary := make(map[string]interface{})

	for _, mf := range metricsFamilies {
		name := mf.GetName()
		metric := mf.GetMetric()[0] // Get first metric

		switch mf.GetType() {
		case prometheus.MetricType_COUNTER:
			summary[name] = metric.GetCounter().GetValue()
		case prometheus.MetricType_GAUGE:
			summary[name] = metric.GetGauge().GetValue()
		case prometheus.MetricType_HISTOGRAM:
			hist := metric.GetHistogram()
			summary[name] = map[string]interface{}{
				"count": hist.GetSampleCount(),
				"sum":   hist.GetSampleSum(),
			}
		}
	}

	return summary
}

// Reset resets all metrics (useful for testing)
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Reset counters
	m.cacheHits.Reset()
	m.cacheMisses.Reset()
	m.cacheEvictions.Reset()

	// Reset gauges
	m.cacheKeysTotal.Set(0)
	m.cacheMemoryUsage.Set(0)
	m.activeConnections.Set(0)
	m.clusterNodes.Set(0)
	m.clusterReplicas.Set(0)
	m.clusterLeader.Set(0)
	m.goRoutines.Set(0)
	m.memoryAllocated.Set(0)
	m.gcPauseTime.Set(0)
}