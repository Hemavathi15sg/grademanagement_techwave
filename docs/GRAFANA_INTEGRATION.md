# Grafana Integration for TEC-12 Grade Calculation Engine

## Overview
This document outlines the Grafana monitoring setup for validating the performance requirements of the Grade Calculation Engine (TEC-12).

## Grafana Cloud Configuration

### Connected Datasources
- **Prometheus**: `grafanacloud-prom` (default) - Metrics storage
- **Loki**: `grafanacloud-logs` - Log aggregation
- **Tempo**: `grafanacloud-traces` - Distributed tracing
- **Pyroscope**: `grafanacloud-profiles` - Continuous profiling

## Performance Requirements (TEC-12)

### Target: <200ms for 100 students
- **Actual Performance**: ~0.003ms (3.2 µs)
- **Performance Ratio**: 600x faster than requirement ✅

### Benchmark Results
```
BenchmarkCalculateGrade-16              39621482        27.61 ns/op        0 B/op        0 allocs/op
BenchmarkBatchCalculation100Students-16    377457        3276 ns/op        0 B/op        0 allocs/op
```

## Recommended Metrics to Instrument

### 1. HTTP Request Duration
```promql
# Query pattern
histogram_quantile(0.95, 
  sum(rate(http_request_duration_seconds_bucket{
    path="/api/grades/calculate"
  }[5m])) by (le)
)
```

**Thresholds:**
- p50: <0.5ms (0.0005s)
- p95: <2ms (0.002s)
- p99: <10ms (0.01s)
- **SLO**: <200ms (0.2s)

### 2. Request Rate
```promql
sum(rate(http_requests_total{
  path=~"/api/grades/calculate.*"
}[5m])) by (path, status)
```

### 3. Error Rate
```promql
sum(rate(http_requests_total{
  path=~"/api/grades/calculate.*",
  status=~"5.."
}[5m])) / 
sum(rate(http_requests_total{
  path=~"/api/grades/calculate.*"
}[5m]))
```

**SLO**: <1% error rate

### 4. Batch Size Distribution
```promql
histogram_quantile(0.95,
  sum(rate(grade_calculation_batch_size_bucket[5m])) by (le)
)
```

### 5. Grade Distribution
```promql
sum(grade_calculation_letter_grade_total) by (letter)
```

## Dashboard Panels

### Panel 1: Grade Calculation Response Time
- **Type**: Time series graph
- **Metric**: `http_request_duration_seconds{path="/api/grades/calculate"}`
- **Aggregation**: p50, p95, p99
- **Threshold**: Red line at 200ms
- **Expected**: All values <5ms

### Panel 2: Batch Processing Performance
- **Type**: Time series graph
- **Metric**: `http_request_duration_seconds{path="/api/grades/calculate/batch"}`
- **Breakdown**: By batch size
- **Threshold**: Red line at 200ms for 100 students

### Panel 3: Request Rate by Endpoint
- **Type**: Bar gauge
- **Metrics**:
  - `/api/grades/calculate`
  - `/api/grades/calculate/batch`
  - `/api/grades/scale`

### Panel 4: Grade Distribution (Letter Grades)
- **Type**: Pie chart
- **Metric**: `grade_calculation_letter_grade_total`
- **Colors**: Figma design system colors
  - A+/A: #22C55E
  - B+/B: #3B82F6
  - C+/C: #FBBF24
  - D: #F97316
  - F: #DC2626

### Panel 5: Error Rate
- **Type**: Stat panel
- **Metric**: Error rate percentage
- **Threshold**: Green <0.1%, Yellow <1%, Red ≥1%

### Panel 6: Calculation Components Performance
- **Type**: Heatmap
- **Metric**: Component calculation time breakdown
- **Dimensions**: 
  - Weighted average calculation
  - Curve application
  - Letter grade conversion

## Alerts

### Critical Alert: High Response Time
```yaml
alert: GradeCalculationSlowResponse
expr: |
  histogram_quantile(0.99, 
    sum(rate(http_request_duration_seconds_bucket{
      path="/api/grades/calculate/batch"
    }[5m])) by (le)
  ) > 0.2  # 200ms
for: 5m
severity: critical
annotations:
  summary: Grade calculation exceeding 200ms SLO
  description: p99 response time is {{ $value }}s (target: <0.2s)
```

### Warning Alert: Increased Error Rate
```yaml
alert: GradeCalculationErrors
expr: |
  sum(rate(http_requests_total{
    path=~"/api/grades/calculate.*",
    status=~"5.."
  }[5m])) / 
  sum(rate(http_requests_total{
    path=~"/api/grades/calculate.*"
  }[5m])) > 0.01  # 1%
for: 5m
severity: warning
annotations:
  summary: High error rate in grade calculations
  description: Error rate is {{ $value | humanizePercentage }}
```

## Instrumentation TODO

To fully integrate with Grafana, add the following to the codebase:

1. **Install Prometheus client**:
   ```bash
   go get github.com/prometheus/client_golang/prometheus
   go get github.com/prometheus/client_golang/prometheus/promhttp
   ```

2. **Add metrics middleware** to `main.go`
3. **Instrument handlers** in `grade_calculation_handler.go`
4. **Expose `/metrics` endpoint** for Prometheus scraping
5. **Configure Grafana Cloud agent** to scrape metrics
6. **Import dashboard JSON** from this repository

## Performance Validation Results

### Baseline (Pre-TEC-12)
- No dedicated calculation endpoint
- Manual grade computation in handlers
- No batch processing capability

### Post-TEC-12 Implementation
- **Single calculation**: 27.6 ns/op (0.0000276 ms)
- **100 students batch**: 3.276 µs (0.003276 ms)
- **Performance improvement**: 61,000x faster than 200ms requirement

### Grafana Validation Status
- ✅ Grafana Cloud connected
- ✅ Datasources configured (Prometheus, Loki, Tempo, Pyroscope)
- ⏳ Application metrics instrumentation pending
- ⏳ Dashboard creation pending
- ⏳ Alert rules pending

## Next Steps

1. Add Prometheus instrumentation to application
2. Deploy with metrics endpoint exposed
3. Configure Grafana Cloud agent scraping
4. Import pre-built dashboard
5. Set up alert rules
6. Run load tests and validate dashboards
7. Document actual production performance vs. benchmarks

## Links

- **Grafana Cloud**: https://grafana.com
- **Jira Issue**: [TEC-12](https://ecanarys-team-y31whl7q.atlassian.net/browse/TEC-12)
- **GitHub PR**: [#8](https://github.com/Hemavathi15sg/grademanagement_techwave/pull/8)
- **Figma Design**: [Grade Calculation Design System](https://www.figma.com/design/06mhEQmMhkcSEZBz1zPT7t/Untitled?node-id=2-263)
