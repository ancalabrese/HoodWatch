package core

import (
	"github.com/hashicorp/go-hclog"
)

type Hq struct {
	eventBus chan *Crime
	errorBus chan error
	l        hclog.Logger
}

type HeadQuarter interface {
	Listen()
	ReportCrime() chan *Crime
	ReportError() chan error
}

var hq = (*Hq)(nil)

func NewHq(l hclog.Logger) *Hq {
	if hq != nil {
		return hq
	}

	l.Debug("New Hood HQ requested")
	return &Hq{
		eventBus: make(chan *Crime),
		errorBus: make(chan error),
		l:        l,
	}
}

func (hq *Hq) Listen() {
	// on each event archieve it and publish it
	go func() {
		for {
			crime := <-hq.eventBus
			hq.onNewCrime(crime)
		}
	}()

	go func() {
		for {
			err := <-hq.errorBus
			hq.onError(err)
		}
	}()
}

func (hq *Hq) ReportCrime() chan *Crime {
	return hq.eventBus
}

func (hq *Hq) ReportError() chan error {
	return hq.errorBus
}

func (hq *Hq) onNewCrime(c *Crime) {
	hq.l.Info("New crime reported", "crime", c)
	// TODO on crime notify press
}

func (hq *Hq) onError(err error) {
	hq.l.Error("Error received", "error", err)
}
