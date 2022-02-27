package utils

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits) // 10位worker最大值
	numberMask  int64 = -1 ^ (-1 << numberBits) // 12位number序号最大值, 4095
	timeShift   uint8 = workerBits + numberBits // 左移10 + 12 == 22位,才是时间戳位的开始
	workerShift uint8 = numberBits
	startTime   int64 = 1525705533000 // 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker ID excess of quantity")
	}
	// 生成一个新节点
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}
func (w *Worker) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number = (w.number + 1) & numberMask
		if w.number == 0 {
			// number数量大于numberMax = 4095时, 就要循环的去等待下一毫秒, 同时将number的值更新为0
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 新的1毫秒使用新的雪花ID
		w.number = 0
	}
	w.timestamp = now
	ID := int64((now-startTime)<<timeShift | (w.workerId << workerShift) | (w.number))
	return ID
}
