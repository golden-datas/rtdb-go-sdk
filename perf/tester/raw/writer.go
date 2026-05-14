package raw

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/golden-datas/rtdb-go-sdk"
	"github.com/golden-datas/rtdb-go-sdk/perf/tester"
)

// Writer Raw API 写入测试器
type Writer struct {
	conn   *rtdb_api.RtdbConnect
	points []*rtdb_api.PointInfo
}

// NewWriter 创建 Raw API 写入测试器
func NewWriter(conn *rtdb_api.RtdbConnect, points []*rtdb_api.PointInfo) *Writer {
	return &Writer{
		conn:   conn,
		points: points,
	}
}

// WritePeriodic 周期性写入实时数据（使用 Raw API）
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

	// 获取连接句柄
	handle := w.conn.ConnectHandle

	for {
		select {
		case <-stopTimer.C:
			return nil
		case <-ticker.C:
			start := time.Now()

			// 准备数据
			ids := make([]rtdb_api.PointID, batchSize)
			datetimes := make([]rtdb_api.TimestampType, batchSize)
			subtimes := make([]rtdb_api.SubtimeType, batchSize)
			values := make([]float64, batchSize)
			states := make([]int64, batchSize)
			qualities := make([]rtdb_api.Quality, batchSize)

			timestamp := time.Now()
			rtdbTs, rtdbSubTs := rtdb_api.GoTimeToRtdbTimestamp(timestamp)

			for i := 0; i < batchSize; i++ {
				point := w.points[i]
				ids[i] = point.ID
				datetimes[i] = rtdbTs
				subtimes[i] = rtdbSubTs
				values[i], states[i] = generateRawValue(point.ValueType)
				qualities[i] = 0 // Good quality
			}

			// 使用 Raw API 批量写入
			errs, rte := rtdb_api.RawRtdbsPutSnapshots64Warp(
				handle,
				ids,
				datetimes,
				subtimes,
				values,
				states,
				qualities,
			)

			if !rtdb_api.RteIsOk(rte) {
				m.Add(time.Since(start), 0, rte.GoError())
				continue
			}

			// 检查每个点的错误
			errorCount := 0
			for _, err := range errs {
				if !rtdb_api.RteIsOk(err) {
					errorCount++
				}
			}

			if errorCount > 0 {
				m.Add(time.Since(start), int64(batchSize-errorCount), fmt.Errorf("%d errors", errorCount))
			} else {
				m.Add(time.Since(start), int64(batchSize), nil)
			}
		}
	}
}

// WriteBurst 急速写入实时数据（使用 Raw API）
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

	// 获取连接句柄
	handle := w.conn.ConnectHandle

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

				// 准备数据
				ids := make([]rtdb_api.PointID, len(points))
				datetimes := make([]rtdb_api.TimestampType, len(points))
				subtimes := make([]rtdb_api.SubtimeType, len(points))
				values := make([]float64, len(points))
				states := make([]int64, len(points))
				qualities := make([]rtdb_api.Quality, len(points))

				timestamp := time.Now().Add(time.Duration(batch) * time.Millisecond).Add(time.Duration(workerID) * time.Microsecond)
				rtdbTs, rtdbSubTs := rtdb_api.GoTimeToRtdbTimestamp(timestamp)

				for i, point := range points {
					ids[i] = point.ID
					datetimes[i] = rtdbTs
					subtimes[i] = rtdbSubTs
					values[i], states[i] = generateRawValue(point.ValueType)
					qualities[i] = 0
				}

				// 使用 Raw API 批量写入
				errs, rte := rtdb_api.RawRtdbsPutSnapshots64Warp(
					handle,
					ids,
					datetimes,
					subtimes,
					values,
					states,
					qualities,
				)

				if !rtdb_api.RteIsOk(rte) {
					m.Add(time.Since(start), 0, rte.GoError())
					continue
				}

				// 检查每个点的错误
				errorCount := 0
				for _, err := range errs {
					if !rtdb_api.RteIsOk(err) {
						errorCount++
					}
				}

				if errorCount > 0 {
					m.Add(time.Since(start), int64(len(points)-errorCount), fmt.Errorf("%d errors", errorCount))
				} else {
					m.Add(time.Since(start), int64(len(points)), nil)
				}
			}
		}(w.points[startIdx:endIdx], i)
	}

	wg.Wait()
	return nil
}

// WriteHisPeriodic 周期性写入历史数据（使用 Raw API）
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

	// 获取连接句柄
	handle := w.conn.ConnectHandle
	baseTime := time.Now().Add(-duration)

	for {
		select {
		case <-stopTimer.C:
			return nil
		case <-ticker.C:
			start := time.Now()

			// 准备数据
			ids := make([]rtdb_api.PointID, batchSize)
			datetimes := make([]rtdb_api.TimestampType, batchSize)
			subtimes := make([]rtdb_api.SubtimeType, batchSize)
			values := make([]float64, batchSize)
			states := make([]int64, batchSize)
			qualities := make([]rtdb_api.Quality, batchSize)

			timestamp := baseTime.Add(time.Since(m.StartTime))
			rtdbTs, rtdbSubTs := rtdb_api.GoTimeToRtdbTimestamp(timestamp)

			for i := 0; i < batchSize; i++ {
				point := w.points[i]
				ids[i] = point.ID
				datetimes[i] = rtdbTs
				subtimes[i] = rtdbSubTs
				values[i], states[i] = generateRawValue(point.ValueType)
				qualities[i] = 0
			}

			// 使用 Raw API 批量写入历史数据
			errs, rte := rtdb_api.RawRtdbhPutArchivedValues64Warp(
				handle,
				ids,
				datetimes,
				subtimes,
				values,
				states,
				qualities,
			)

			if !rtdb_api.RteIsOk(rte) {
				m.Add(time.Since(start), 0, rte.GoError())
				continue
			}

			// 检查每个点的错误
			errorCount := 0
			for _, err := range errs {
				if !rtdb_api.RteIsOk(err) {
					errorCount++
				}
			}

			if errorCount > 0 {
				m.Add(time.Since(start), int64(batchSize-errorCount), fmt.Errorf("%d errors", errorCount))
			} else {
				m.Add(time.Since(start), int64(batchSize), nil)
			}
		}
	}
}

// WriteHisBurst 急速写入历史数据（使用 Raw API）
func (w *Writer) WriteHisBurst(batchSize, batchesPerPoint, workers int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	var wg sync.WaitGroup
	pointCount := len(w.points)
	pointsPerWorker := pointCount / workers
	if pointsPerWorker == 0 {
		pointsPerWorker = 1
	}

	// 获取连接句柄
	handle := w.conn.ConnectHandle
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

				// 准备数据
				ids := make([]rtdb_api.PointID, len(points))
				datetimes := make([]rtdb_api.TimestampType, len(points))
				subtimes := make([]rtdb_api.SubtimeType, len(points))
				values := make([]float64, len(points))
				states := make([]int64, len(points))
				qualities := make([]rtdb_api.Quality, len(points))

				timestamp := baseTime.Add(time.Duration(batch) * time.Second).Add(time.Duration(workerID) * time.Millisecond)
				rtdbTs, rtdbSubTs := rtdb_api.GoTimeToRtdbTimestamp(timestamp)

				for i, point := range points {
					ids[i] = point.ID
					datetimes[i] = rtdbTs
					subtimes[i] = rtdbSubTs
					values[i], states[i] = generateRawValue(point.ValueType)
					qualities[i] = 0
				}

				// 使用 Raw API 批量写入历史数据
				errs, rte := rtdb_api.RawRtdbhPutArchivedValues64Warp(
					handle,
					ids,
					datetimes,
					subtimes,
					values,
					states,
					qualities,
				)

				if !rtdb_api.RteIsOk(rte) {
					m.Add(time.Since(start), 0, rte.GoError())
					continue
				}

				// 检查每个点的错误
				errorCount := 0
				for _, err := range errs {
					if !rtdb_api.RteIsOk(err) {
						errorCount++
					}
				}

				if errorCount > 0 {
					m.Add(time.Since(start), int64(len(points)-errorCount), fmt.Errorf("%d errors", errorCount))
				} else {
					m.Add(time.Since(start), int64(len(points)), nil)
				}
			}
		}(w.points[startIdx:endIdx], i)
	}

	wg.Wait()
	return nil
}

// generateRawValue 根据类型生成测试值（返回 value 和 state）
func generateRawValue(valueType rtdb_api.ValueType) (float64, int64) {
	timestamp := time.Now().Unix()

	switch valueType {
	case rtdb_api.ValueTypeInt8:
		return 0, int64(int8(timestamp % 100))
	case rtdb_api.ValueTypeInt16:
		return 0, int64(int16(timestamp % 1000))
	case rtdb_api.ValueTypeInt32:
		return 0, int64(int32(timestamp % 10000))
	case rtdb_api.ValueTypeInt64:
		return 0, timestamp
	case rtdb_api.ValueTypeUint8:
		return 0, int64(uint8(timestamp % 100))
	case rtdb_api.ValueTypeUint16:
		return 0, int64(uint16(timestamp % 1000))
	case rtdb_api.ValueTypeUint32:
		return 0, int64(uint32(timestamp % 10000))
	case rtdb_api.ValueTypeFloat32:
		return float64(float32(math.Sin(float64(timestamp) / 100))), 0
	case rtdb_api.ValueTypeFloat64:
		return math.Sin(float64(timestamp) / 100), 0
	default:
		return float64(timestamp % 10000), 0
	}
}
