package transforms

import "sync"

type msg_DCT1DCol struct {
	wg    *sync.WaitGroup
	input []float64
}

type msg_DCT1DRow struct {
	wg    *sync.WaitGroup
	input []float64
	i     int
}

type WorkerPool struct {
	dCT1DCol chan msg_DCT1DCol
	dCT1DRow chan msg_DCT1DRow
	doneCh   chan struct{}
}

func StartWorkerPool(workers int) *WorkerPool {
	pool := &WorkerPool{
		dCT1DCol: make(chan msg_DCT1DCol, workers/2),
		dCT1DRow: make(chan msg_DCT1DRow, workers/2),
		doneCh:   make(chan struct{}),
	}
	for i := 0; i < workers; i++ {
		go pool.workerFunc()
	}
	return pool
}

func (wp *WorkerPool) CloseWorkerPool() {
	close(wp.dCT1DCol)
	close(wp.dCT1DRow)
	close(wp.doneCh)
}

func (wp *WorkerPool) workerFunc() {
	for {
		temp := make([]float64, 64)
		row := make([]float64, 64)
		select {
		case task := <-wp.dCT1DRow:
			for j := 0; j < pHashSize; j++ {
				row[j] = task.input[task.i+(j*pHashSize)]
			}
			forwardTransform(row, temp, len(row))
			for j := 0; j < len(row); j++ {
				task.input[task.i+(j*pHashSize)] = row[j]
			}
			task.wg.Done()
		case task := <-wp.dCT1DCol:
			forwardTransform(task.input, temp, len(task.input))
			task.wg.Done()
		case <-wp.doneCh:
			return
		}
	}
}
