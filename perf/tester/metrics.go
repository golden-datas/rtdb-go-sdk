package tester

import (
	"fmt"
	"sort"
	"time"
)

// Metrics 性能指标收集
type Metrics struct {
	StartTime   time.Time
	EndTime     time.Time
	TotalCount  int64
	TotalBytes  int64
	Latencies   []time.Duration
	Errors      int64
	mutex       chan struct{}
}

// NewMetrics 创建新的指标收集器
func NewMetrics() *Metrics {
	return &Metrics{
		Latencies: make([]time.Duration, 0, 10000),
		mutex:     make(chan struct{}, 1),
	}
}

// Start 开始计时
func (m *Metrics) Start() {
	m.StartTime = time.Now()
}

// Stop 停止计时
func (m *Metrics) Stop() {
	m.EndTime = time.Now()
}

// Add 添加一次操作记录
func (m *Metrics) Add(latency time.Duration, count int64, err error) {
	m.mutex <- struct{}{}
	defer func() { <-m.mutex }()

	if err != nil {
		m.Errors++
	} else {
		m.TotalCount += count
		m.Latencies = append(m.Latencies, latency)
	}
}

// TotalDuration 获取总耗时
func (m *Metrics) TotalDuration() time.Duration {
	if m.EndTime.IsZero() {
		return time.Since(m.StartTime)
	}
	return m.EndTime.Sub(m.StartTime)
}

// Throughput 获取吞吐量（万条/秒）
func (m *Metrics) Throughput() float64 {
	duration := m.TotalDuration().Seconds()
	if duration == 0 {
		return 0
	}
	return float64(m.TotalCount) / 10000 / duration
}

// AvgLatency 获取平均延迟
func (m *Metrics) AvgLatency() time.Duration {
	if len(m.Latencies) == 0 {
		return 0
	}
	var total time.Duration
	for _, l := range m.Latencies {
		total += l
	}
	return total / time.Duration(len(m.Latencies))
}

// P50 获取P50延迟
func (m *Metrics) P50() time.Duration {
	return m.percentile(0.5)
}

// P95 获取P95延迟
func (m *Metrics) P95() time.Duration {
	return m.percentile(0.95)
}

// P99 获取P99延迟
func (m *Metrics) P99() time.Duration {
	return m.percentile(0.99)
}

// MaxLatency 获取最大延迟
func (m *Metrics) MaxLatency() time.Duration {
	if len(m.Latencies) == 0 {
		return 0
	}
	max := m.Latencies[0]
	for _, l := range m.Latencies {
		if l > max {
			max = l
		}
	}
	return max
}

// percentile 计算百分位数
func (m *Metrics) percentile(p float64) time.Duration {
	if len(m.Latencies) == 0 {
		return 0
	}
	sorted := make([]time.Duration, len(m.Latencies))
	copy(sorted, m.Latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	idx := int(float64(len(sorted)-1) * p)
	return sorted[idx]
}

// Print 打印指标报告
func (m *Metrics) Print() {
	fmt.Println("\n========== 性能测试报告 ==========")
	fmt.Printf("总耗时: %v\n", m.TotalDuration())
	fmt.Printf("总记录数: %d\n", m.TotalCount)
	fmt.Printf("错误数: %d\n", m.Errors)
	fmt.Printf("吞吐量: %.2f 万条/秒\n", m.Throughput())
	fmt.Printf("平均延迟: %v\n", m.AvgLatency())
	fmt.Printf("P50延迟: %v\n", m.P50())
	fmt.Printf("P95延迟: %v\n", m.P95())
	fmt.Printf("P99延迟: %v\n", m.P99())
	fmt.Printf("最大延迟: %v\n", m.MaxLatency())
	fmt.Println("==================================")
}
