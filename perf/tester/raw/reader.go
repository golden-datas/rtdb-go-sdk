package raw

import (
	"sync"
	"time"

	"github.com/golden-datas/rtdb-go-sdk"
	"github.com/golden-datas/rtdb-go-sdk/perf/tester"
)

// Reader Raw API 读取测试器
type Reader struct {
	conn   *rtdb_api.RtdbConnect
	points []*rtdb_api.PointInfo
}

// NewReader 创建 Raw API 读取测试器
func NewReader(conn *rtdb_api.RtdbConnect, points []*rtdb_api.PointInfo) *Reader {
	return &Reader{
		conn:   conn,
		points: points,
	}
}

// ReadLast 读取最新值（使用 Raw API）
func (r *Reader) ReadLast(workers int, samples int, m *tester.Metrics) error {
	// 获取连接句柄
	handle := r.conn.ConnectHandle

	for s := 0; s < samples; s++ {
		m.Start()
		var wg sync.WaitGroup
		pointCount := len(r.points)
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
				start := time.Now()

				// 准备点ID
				ids := make([]rtdb_api.PointID, len(points))
				for i, point := range points {
					ids[i] = point.ID
				}

				// 使用 Raw API 批量读取
				_, _, _, _, _, errs, rte := rtdb_api.RawRtdbsGetSnapshots64Warp(handle, ids)

				if !rtdb_api.RteIsOk(rte) {
					m.Add(time.Since(start), 0, rte.GoError())
					return
				}

				successCount := 0
				for _, err := range errs {
					if rtdb_api.RteIsOk(err) {
						successCount++
					}
				}
				m.Add(time.Since(start), int64(successCount), nil)
			}(r.points[startIdx:endIdx])
		}
		wg.Wait()
		m.Stop()
	}
	return nil
}

// ReadRaw 读取历史原始值（使用 Raw API）
func (r *Reader) ReadRaw(start, end time.Time, workers, pointCount int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	// 获取连接句柄
	handle := r.conn.ConnectHandle

	var wg sync.WaitGroup
	actualPointCount := len(r.points)
	if pointCount > actualPointCount {
		pointCount = actualPointCount
	}
	pointsPerWorker := pointCount / workers
	if pointsPerWorker == 0 {
		pointsPerWorker = 1
	}

	startTs, startSubTs := rtdb_api.GoTimeToRtdbTimestamp(start)
	endTs, endSubTs := rtdb_api.GoTimeToRtdbTimestamp(end)

	for i := 0; i < workers && i*pointsPerWorker < pointCount; i++ {
		wg.Add(1)
		startIdx := i * pointsPerWorker
		endIdx := startIdx + pointsPerWorker
		if i == workers-1 || endIdx > pointCount {
			endIdx = pointCount
		}

		go func(points []*rtdb_api.PointInfo) {
			defer wg.Done()
			for _, point := range points {
				s := time.Now()

				// 使用 Raw API 读取历史数据
				_, _, _, _, _, rte := rtdb_api.RawRtdbhGetArchivedValues64Warp(
					handle,
					point.ID,
					10000, // maxCount
					startTs,
					startSubTs,
					endTs,
					endSubTs,
				)

				if !rtdb_api.RteIsOk(rte) {
					m.Add(time.Since(s), 0, rte.GoError())
					continue
				}
				m.Add(time.Since(s), 1, nil)
			}
		}(r.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}

// ReadInterpo 读取插值（使用 Raw API）
func (r *Reader) ReadInterpo(start, end time.Time, count int, workers, pointCount int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	// 获取连接句柄
	handle := r.conn.ConnectHandle

	var wg sync.WaitGroup
	actualPointCount := len(r.points)
	if pointCount > actualPointCount {
		pointCount = actualPointCount
	}
	pointsPerWorker := pointCount / workers
	if pointsPerWorker == 0 {
		pointsPerWorker = 1
	}

	startTs, startSubTs := rtdb_api.GoTimeToRtdbTimestamp(start)
	endTs, endSubTs := rtdb_api.GoTimeToRtdbTimestamp(end)

	for i := 0; i < workers && i*pointsPerWorker < pointCount; i++ {
		wg.Add(1)
		startIdx := i * pointsPerWorker
		endIdx := startIdx + pointsPerWorker
		if i == workers-1 || endIdx > pointCount {
			endIdx = pointCount
		}

		go func(points []*rtdb_api.PointInfo) {
			defer wg.Done()
			for _, point := range points {
				s := time.Now()

				// 使用 Raw API 读取插值
				_, _, _, _, _, rte := rtdb_api.RawRtdbhGetInterpoValues64Warp(
					handle,
					point.ID,
					int32(count),
					startTs,
					startSubTs,
					endTs,
					endSubTs,
				)

				if !rtdb_api.RteIsOk(rte) {
					m.Add(time.Since(s), 0, rte.GoError())
					continue
				}
				m.Add(time.Since(s), int64(count), nil)
			}
		}(r.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}

// ReadPlot 读取趋势曲线值（使用 Raw API）
func (r *Reader) ReadPlot(start, end time.Time, interval int, workers, pointCount int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	// 获取连接句柄
	handle := r.conn.ConnectHandle

	var wg sync.WaitGroup
	actualPointCount := len(r.points)
	if pointCount > actualPointCount {
		pointCount = actualPointCount
	}
	pointsPerWorker := pointCount / workers
	if pointsPerWorker == 0 {
		pointsPerWorker = 1
	}

	startTs, startSubTs := rtdb_api.GoTimeToRtdbTimestamp(start)
	endTs, endSubTs := rtdb_api.GoTimeToRtdbTimestamp(end)

	for i := 0; i < workers && i*pointsPerWorker < pointCount; i++ {
		wg.Add(1)
		startIdx := i * pointsPerWorker
		endIdx := startIdx + pointsPerWorker
		if i == workers-1 || endIdx > pointCount {
			endIdx = pointCount
		}

		go func(points []*rtdb_api.PointInfo) {
			defer wg.Done()
			for _, point := range points {
				s := time.Now()

				// 使用 Raw API 读取趋势曲线
				_, _, _, _, _, rte := rtdb_api.RawRtdbhGetPlotValues64Warp(
					handle,
					point.ID,
					int32(interval),
					startTs,
					startSubTs,
					endTs,
					endSubTs,
				)

				if !rtdb_api.RteIsOk(rte) {
					m.Add(time.Since(s), 0, rte.GoError())
					continue
				}
				m.Add(time.Since(s), int64(interval), nil)
			}
		}(r.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}

// ReadSummary 读取统计值（使用 Raw API）
func (r *Reader) ReadSummary(start, end time.Time, workers, pointCount int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	// 获取连接句柄
	handle := r.conn.ConnectHandle

	var wg sync.WaitGroup
	actualPointCount := len(r.points)
	if pointCount > actualPointCount {
		pointCount = actualPointCount
	}
	pointsPerWorker := pointCount / workers
	if pointsPerWorker == 0 {
		pointsPerWorker = 1
	}

	startTs, startSubTs := rtdb_api.GoTimeToRtdbTimestamp(start)
	endTs, endSubTs := rtdb_api.GoTimeToRtdbTimestamp(end)

	for i := 0; i < workers && i*pointsPerWorker < pointCount; i++ {
		wg.Add(1)
		startIdx := i * pointsPerWorker
		endIdx := startIdx + pointsPerWorker
		if i == workers-1 || endIdx > pointCount {
			endIdx = pointCount
		}

		go func(points []*rtdb_api.PointInfo) {
			defer wg.Done()
			for _, point := range points {
				s := time.Now()

				// 使用 Raw API 读取统计值
				_, rte := rtdb_api.RawRtdbhSummaryDataWarp(
					handle,
					point.ID,
					startTs,
					startSubTs,
					endTs,
					endSubTs,
				)

				if !rtdb_api.RteIsOk(rte) {
					m.Add(time.Since(s), 0, rte.GoError())
					continue
				}
				m.Add(time.Since(s), 1, nil)
			}
		}(r.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}

// ReadSection 读取断面数据（使用 Raw API）
func (r *Reader) ReadSection(timestamp time.Time, workers, pointCount int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	// 获取连接句柄
	handle := r.conn.ConnectHandle

	var wg sync.WaitGroup
	actualPointCount := len(r.points)
	if pointCount > actualPointCount {
		pointCount = actualPointCount
	}
	pointsPerWorker := pointCount / workers
	if pointsPerWorker == 0 {
		pointsPerWorker = 1
	}

	ts, subTs := rtdb_api.GoTimeToRtdbTimestamp(timestamp)

	for i := 0; i < workers && i*pointsPerWorker < pointCount; i++ {
		wg.Add(1)
		startIdx := i * pointsPerWorker
		endIdx := startIdx + pointsPerWorker
		if i == workers-1 || endIdx > pointCount {
			endIdx = pointCount
		}

		go func(points []*rtdb_api.PointInfo) {
			defer wg.Done()
			s := time.Now()

			// 准备点ID
			ids := make([]rtdb_api.PointID, len(points))
			for i, point := range points {
				ids[i] = point.ID
			}

			// 使用 Raw API 读取断面数据
			_, _, _, _, _, errs, rte := rtdb_api.RawRtdbhGetCrossSectionValues64Warp(
				handle,
				ids,
				rtdb_api.RtdbHisModePrevious,
				ts,
				subTs,
			)

			if !rtdb_api.RteIsOk(rte) {
				m.Add(time.Since(s), 0, rte.GoError())
				return
			}

			successCount := 0
			for _, err := range errs {
				if rtdb_api.RteIsOk(err) {
					successCount++
				}
			}
			m.Add(time.Since(s), int64(successCount), nil)
		}(r.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}
