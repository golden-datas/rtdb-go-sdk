package easy

import (
	"sync"
	"time"

	"github.com/kkbase/rtdb_api"
	"github.com/kkbase/rtdb_api/perf/tester"
)

// Reader Easy API 读取测试器
type Reader struct {
	conn   *rtdb_api.RtdbConnect
	points []*rtdb_api.PointInfo
}

// NewReader 创建 Easy API 读取测试器
func NewReader(conn *rtdb_api.RtdbConnect, points []*rtdb_api.PointInfo) *Reader {
	return &Reader{
		conn:   conn,
		points: points,
	}
}

// ReadLast 读取最新值
func (r *Reader) ReadLast(workers int, samples int, m *tester.Metrics) error {
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
				for _, point := range points {
					if _, err := r.conn.ReadLast(point); err != nil {
						m.Add(time.Since(start), 0, err)
						continue
					}
				}
				m.Add(time.Since(start), int64(len(points)), nil)
			}(r.points[startIdx:endIdx])
		}
		wg.Wait()
		m.Stop()
	}
	return nil
}

// ReadRaw 读取历史原始值
func (r *Reader) ReadRaw(start, end time.Time, workers, pointCount int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	var wg sync.WaitGroup
	actualPointCount := len(r.points)
	if pointCount > actualPointCount {
		pointCount = actualPointCount
	}
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
			for _, point := range points {
				s := time.Now()
				ptvqs, err := r.conn.ReadRange(point, start, end)
				if err != nil {
					m.Add(time.Since(s), 0, err)
					continue
				}
				m.Add(time.Since(s), int64(len(ptvqs.TVQs)), nil)
			}
		}(r.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}

// ReadInterpo 读取插值
func (r *Reader) ReadInterpo(start, end time.Time, count int, workers, pointCount int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	var wg sync.WaitGroup
	actualPointCount := len(r.points)
	if pointCount > actualPointCount {
		pointCount = actualPointCount
	}
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
			for _, point := range points {
				s := time.Now()
				ptvqs, err := r.conn.ReadInterpo(point, int32(count), start, end, "")
				if err != nil {
					m.Add(time.Since(s), 0, err)
					continue
				}
				m.Add(time.Since(s), int64(len(ptvqs.TVQs)), nil)
			}
		}(r.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}

// ReadPlot 读取趋势曲线值
func (r *Reader) ReadPlot(start, end time.Time, interval int, workers, pointCount int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	var wg sync.WaitGroup
	actualPointCount := len(r.points)
	if pointCount > actualPointCount {
		pointCount = actualPointCount
	}
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
			for _, point := range points {
				s := time.Now()
				ptvqs, err := r.conn.ReadPlot(point, int32(interval), start, end)
				if err != nil {
					m.Add(time.Since(s), 0, err)
					continue
				}
				m.Add(time.Since(s), int64(len(ptvqs.TVQs)), nil)
			}
		}(r.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}

// ReadSummary 读取统计值
func (r *Reader) ReadSummary(start, end time.Time, workers, pointCount int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	var wg sync.WaitGroup
	actualPointCount := len(r.points)
	if pointCount > actualPointCount {
		pointCount = actualPointCount
	}
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
			for _, point := range points {
				s := time.Now()
				_, err := r.conn.ReadSummary(point, "", start, end)
				if err != nil {
					m.Add(time.Since(s), 0, err)
					continue
				}
				m.Add(time.Since(s), 1, nil)
			}
		}(r.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}

// ReadSection 读取断面数据
func (r *Reader) ReadSection(timestamp time.Time, workers, pointCount int, m *tester.Metrics) error {
	m.Start()
	defer m.Stop()

	var wg sync.WaitGroup
	actualPointCount := len(r.points)
	if pointCount > actualPointCount {
		pointCount = actualPointCount
	}
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
			s := time.Now()
			_, errs, err := r.conn.ReadSection(points, rtdb_api.RtdbHisModePrevious, timestamp)
			if err != nil {
				m.Add(time.Since(s), 0, err)
				return
			}
			successCount := 0
			for _, e := range errs {
				if e == nil {
					successCount++
				}
			}
			m.Add(time.Since(s), int64(successCount), nil)
		}(r.points[startIdx:endIdx])
	}

	wg.Wait()
	return nil
}
