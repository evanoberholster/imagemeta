package transforms

import "sync"

// msgDCT1D is a Message for a DCT1D transformation
type msgDCT1D struct {
	wg    *sync.WaitGroup
	input *[]float64
	i     int
}

// WorkerPool is a DCT WorkerPool
type WorkerPool struct {
	wgPool     sync.Pool
	dCT1DColCh chan msgDCT1D
	dCT1DRowCh chan msgDCT1D
	doneCh     chan struct{}
}

// StartWorkerPool starts a DCT WorkerPool
func StartWorkerPool(workers int) *WorkerPool {
	wp := &WorkerPool{
		dCT1DColCh: make(chan msgDCT1D, workers/2),
		dCT1DRowCh: make(chan msgDCT1D, workers/2),
		doneCh:     make(chan struct{}),
		wgPool: sync.Pool{
			New: func() interface{} { return new(sync.WaitGroup) },
		},
	}

	for i := 0; i < workers; i++ {
		go wp.WorkerFunc()
	}
	return wp
}

func (wp *WorkerPool) sendDCT1DCol(wg *sync.WaitGroup, pixels *[]float64, i int) {
	wp.dCT1DColCh <- msgDCT1D{wg: wg, input: pixels, i: i}
}

func (wp *WorkerPool) sendDCT1DRow(wg *sync.WaitGroup, pixels *[]float64, i int) {
	wp.dCT1DRowCh <- msgDCT1D{wg: wg, input: pixels, i: i}
}

// Close closes a workerpool
func (wp *WorkerPool) Close() {
	close(wp.dCT1DColCh)
	close(wp.dCT1DRowCh)
	close(wp.doneCh)
}

func (wp *WorkerPool) WorkerFunc() {
	for {
		temp := make([]float64, 64)
		row := make([]float64, 64)
		select {
		case task := <-wp.dCT1DRowCh:
			if task.input == nil {
				continue
			}
			for j := 0; j < pHashSize; j++ {
				row[j] = (*task.input)[task.i+(j*pHashSize)]
			}
			forwardTransform(row, temp, len(row))
			for j := 0; j < len(row); j++ {
				(*task.input)[task.i+(j*pHashSize)] = row[j]
			}
			task.wg.Done()
		case task := <-wp.dCT1DColCh:
			if task.input == nil {
				continue
			}
			forwardTransform((*task.input)[task.i*pHashSize:(task.i*pHashSize)+pHashSize], temp, len(temp))
			task.wg.Done()
		case <-wp.doneCh:
			return
		}
	}
}
