package timgr

import "github.com/2997215859/gomdsdk/timescale"

type Option func(mgr *TiMgr)

func WithTiCallback(cb TiCallback) Option {
	return func(mgr *TiMgr) {
		mgr.tiCallback = cb
	}
}

func WithTiSeqCallback(cb TiCallback) Option {
	return func(mgr *TiMgr) {
		mgr.tiSeqCallback = cb
	}
}

func WithTimescale(scale timescale.TimeScale) Option {
	return func(mgr *TiMgr) {
		mgr.timescale = &scale
	}
}
