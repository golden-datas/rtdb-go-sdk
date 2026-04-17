package easy

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/kkbase/rtdb_api"
	"github.com/kkbase/rtdb_api/perf/tester"
)

// Writer Easy API 写入测试器
type Writer struct {
	conn   *rtdb_api.RtdbConnect
	points []*rtdb_api.PointInfo
}

// NewWriter 创建 Easy API 写入测试器
func NewWriter(conn *rtdb_api.RtdbConnect, points []*rtdb_api.PointInfo) *Writer {
	return &Writer{
		conn:   conn,
		points: points,
	}
}

// WritePeriodic 周期性写入实时数据
// batchSize: 每批写入的测点数
// interval: 写入间隔
// duration: 总持续时间
func (w *Writer) WritePeriodic(batchSize int, interval time.Duration, duration time.Duration, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	stopTimer := time.NewTimer(duration)
	defer stopTimer.Stop()

	pointCount := len(w.points)
	if batchSize > pointCount {
		batchSize = pointCount
	}

	for {
		select {
		case <-stopTimer.C:
			return nil
		case <-ticker.C:
			start := time.Now()
			// 随机选择一批点写入
			for i := 0; i < batchSize; i++ {
				point := w.points[i]
				value := generateValue(point.ValueType)
				tvq := point.NewNowTVQ(value, rtdb_api.Quality(0))
				if err := w.conn.WriteValue(point, false, tvq); err != nil {
					m.Add(time.Since(start), 0, err)
					continue
				}
			}
			m.Add(time.Since(start), int64(batchSize), nil)
		}
	}
}

// WriteBurst 急速写入实时数据
// batchSize: 每批写入的测点数
// batchesPerPoint: 每个测点写入的批次数
// workers: 并发工作线程数
func (w *Writer) WriteBurst(batchSize, batchesPerPoint, workers int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	var wg sync.WaitGroup
	pointCount := len(w.points)
	pointsPerWorker := pointCount / workers
	if pointsPerWorker == 0 {
		pointsPerWorker = 1
	}

	for i := 0; i < workers && i*pointsPerWorker < pointCount; i++ {
		wg.Add(1)
		startIdx := i * pointsPerWorker
		endIdx := startIdx + pointsPerWorker
		if i == workers-1 || endIdx > pointCount {
			endIdx = pointCount
		}

		go func(points []*rtdb_api.PointInfo) {
			defer wg.Done()
			for batch := 0; batch < batchesPerPoint; batch++ {
				start := time.Now()
				for _, point := range points {
					value := generateValue(point.ValueType)
					tvq := point.NewNowTVQ(value, rtdb_api.Quality(0))
					if err := w.conn.WriteValue(point, false, tvq); err != nil {
						m.Add(time.Since(start), 0, err)
						continue
					}
				}
				m.Add(time.Since(start), int64(len(points)), nil)
			}
		}(w.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}

// WriteHisPeriodic 周期性写入历史数据
func (w *Writer) WriteHisPeriodic(batchSize int, interval time.Duration, duration time.Duration, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	stopTimer := time.NewTimer(duration)
	defer stopTimer.Stop()

	pointCount := len(w.points)
	if batchSize > pointCount {
		batchSize = pointCount
	}

	baseTime := time.Now().Add(-duration) // 从历史时间点开始写入

	for {
		select {
		case <-stopTimer.C:
			return nil
		case <-ticker.C:
			start := time.Now()
			var ptvqs []rtdb_api.PTVQ

			for i := 0; i < batchSize; i++ {
				point := w.points[i]
				value := generateValue(point.ValueType)
				timestamp := baseTime.Add(time.Since(m.StartTime))
				ptvqs = append(ptvqs, point.NewPTVQ(timestamp, value, rtdb_api.Quality(0)))
			}

			if _, err := w.conn.WriteSection(false, ptvqs); err != nil {
				m.Add(time.Since(start), 0, err)
				continue
			}
			m.Add(time.Since(start), int64(batchSize), nil)
		}
	}
}

// WriteHisBurst 急速写入历史数据
func (w *Writer) WriteHisBurst(batchSize, batchesPerPoint, workers int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	var wg sync.WaitGroup
	pointCount := len(w.points)
	pointsPerWorker := pointCount / workers
	if pointsPerWorker == 0 {
		pointsPerWorker = 1
	}

	baseTime := time.Now().Add(-time.Duration(batchesPerPoint) * time.Second)

	for i := 0; i < workers && i*pointsPerWorker < pointCount; i++ {
		wg.Add(1)
		startIdx := i * pointsPerWorker
		endIdx := startIdx + pointsPerWorker
		if i == workers-1 || endIdx > pointCount {
			endIdx = pointCount
		}

		go func(points []*rtdb_api.PointInfo, workerID int) {
			defer wg.Done()
			for batch := 0; batch < batchesPerPoint; batch++ {
				start := time.Now()
				var ptvqs []rtdb_api.PTVQ

				timestamp := baseTime.Add(time.Duration(batch) * time.Second).Add(time.Duration(workerID) * time.Millisecond)

				for _, point := range points {
					value := generateValue(point.ValueType)
					ptvqs = append(ptvqs, point.NewPTVQ(timestamp, value, rtdb_api.Quality(0)))
				}

				if _, err := w.conn.WriteSection(false, ptvqs); err != nil {
					m.Add(time.Since(start), 0, err)
					continue
				}
				m.Add(time.Since(start), int64(len(points)), nil)
			}
		}(w.points[startIdx:endIdx], i)
	}

	wg.Wait()
	return nil
}

// generateValue 根据类型生成测试值
func generateValue(valueType rtdb_api.ValueType) interface{} {
	switch valueType {
	case rtdb_api.ValueTypeInt8:
		return int8(time.Now().Unix() % 100)
	case rtdb_api.ValueTypeInt16:
		return int16(time.Now().Unix() % 1000)
	case rtdb_api.ValueTypeInt32:
		return int32(time.Now().Unix() % 10000)
	case rtdb_api.ValueTypeInt64:
		return time.Now().Unix()
	case rtdb_api.ValueTypeUint8:
		return uint8(time.Now().Unix() % 100)
	case rtdb_api.ValueTypeUint16:
		return uint16(time.Now().Unix() % 1000)
	case rtdb_api.ValueTypeUint32:
		return uint32(time.Now().Unix() % 10000)
	case rtdb_api.ValueTypeFloat32:
		return float32(math.Sin(float64(time.Now().Unix()) / 100))
	case rtdb_api.ValueTypeFloat64:
		return math.Sin(float64(time.Now().Unix()) / 100)
	case rtdb_api.ValueTypeString:
		return fmt.Sprintf("test_%d", time.Now().Unix())
	default:
		return float64(time.Now().Unix() % 10000)
	}
}
