package transforms

import "sync"

// msgDCT1DCol is a Message for DCT1D column transformation
type msgDCT1DCol struct {
	wg    *sync.WaitGroup
	input []float64
}

// msgDCT1DRow is a Message for DCT1D row transformation
type msgDCT1DRow struct {
	wg    *sync.WaitGroup
	input []float64
	i     int
}

// WorkerPool is a DCT WorkerPool
type WorkerPool struct {
	dCT1DColCh chan msgDCT1DCol
	dCT1DRowCh chan msgDCT1DRow
	doneCh     chan struct{}
}

// StartWorkerPool starts a DCT WorkerPool
func StartWorkerPool(workers int) *WorkerPool {
	pool := &WorkerPool{
		dCT1DColCh: make(chan msgDCT1DCol, workers/2),
		dCT1DRowCh: make(chan msgDCT1DRow, workers/2),
		doneCh:     make(chan struct{}),
	}
	for i := 0; i < workers; i++ {
		go pool.workerFunc()
	}
	return pool
}

func (wp *WorkerPool) DCT1DCol(wg *sync.WaitGroup, pixels []float64) {
	wp.dCT1DRowCh <- msgDCT1DRow{wg: wg, input: pixels}
}

func (wp *WorkerPool) DCT1DRow(wg *sync.WaitGroup, pixels []float64, i int) {
	wp.dCT1DRowCh <- msgDCT1DRow{wg: wg, input: pixels, i: i}
}

// Close closes a workerpool
func (wp *WorkerPool) Close() {
	close(wp.dCT1DColCh)
	close(wp.dCT1DRowCh)
	close(wp.doneCh)
}

func (wp *WorkerPool) workerFunc() {
	for {
		temp := make([]float64, 64)
		row := make([]float64, 64)
		select {
		case task := <-wp.dCT1DRowCh:
			for j := 0; j < pHashSize; j++ {
				row[j] = task.input[task.i+(j*pHashSize)]
			}
			forwardTransform(row, temp, len(row))
			for j := 0; j < len(row); j++ {
				task.input[task.i+(j*pHashSize)] = row[j]
			}
			task.wg.Done()
		case task := <-wp.dCT1DColCh:
			forwardTransform(task.input, temp, len(task.input))
			task.wg.Done()
		case <-wp.doneCh:
			return
		}
	}
}
