package timgr

import (
	"sync"

	"github.com/2997215859/gomdsdk/timescale"
)

type TiCallback func(ti int)

type TiMgr struct {
	latestTimestampMtx sync.Mutex
	latestTimestamp    int64

	latestTiMtx sync.Mutex
	latestTi    int

	tiCallback TiCallback

	stopChan      chan struct{}
	tiSeqChannel  chan int
	tiSeqCallback TiCallback

	timescale *timescale.TimeScale
}

func NewTiMgr(opts ...Option) *TiMgr {
	mgr := &TiMgr{
		latestTimestamp: -1,
		latestTi:        0,
		timescale:       timescale.DefaultTimeScale,
		stopChan:        make(chan struct{}),
	}

	for _, o := range opts {
		o(mgr)
	}

	mgr.Start()

	return mgr
}

func (mgr *TiMgr) SetTimestamp(timestamp int64) {
	mgr.latestTimestampMtx.Lock()
	defer mgr.latestTimestampMtx.Unlock()

	mgr.latestTimestamp = timestamp
}

func (mgr *TiMgr) GetTimestamp() int64 {
	mgr.latestTimestampMtx.Lock()
	defer mgr.latestTimestampMtx.Unlock()

	return mgr.latestTimestamp
}

func (mgr *TiMgr) UpdateLatestTi(ti int) bool {
	mgr.latestTiMtx.Lock()
	defer mgr.latestTiMtx.Unlock()

	if ti > mgr.latestTi {
		mgr.latestTi = ti
		return true
	}
	return false
}

func (mgr *TiMgr) Run() {
	if mgr.tiSeqCallback == nil {
		return
	}

	mgr.tiSeqChannel = make(chan int, 1024)
	go func() {
		for {
			select {
			case ti := <-mgr.tiSeqChannel:
				mgr.tiSeqCallback(ti)
			case <-mgr.stopChan:
				return
			}
		}
	}()
}

func (mgr *TiMgr) Update(updateTime string) {
	ti := mgr.timescale.GetTi(updateTime)

	if ti < 0 {
		return
	}

	if mgr.UpdateLatestTi(ti) { // 如果更新成功，则调用 1 分钟回调
		if mgr.tiCallback != nil {
			go mgr.tiCallback(ti)
		}
		if mgr.tiSeqCallback != nil {
			mgr.tiSeqChannel <- ti
		}
	}
}

func (mgr *TiMgr) Start() {
	go mgr.Run()
}

func (mgr *TiMgr) Stop() {
	close(mgr.stopChan)
}
