package core

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/go-hclog"
)

type Hood struct {
	Name      string       `yaml:"name"`
	IsEnabled bool         `yaml:"enabled"`
	Path      string       `yaml:"path"`
	Interval  int          `yaml:"interval"`
	Reporting bool         `yaml:"reporting"`
	Rules     []*HoodRule  `yaml:"rules"`
	Hq        HeadQuarter  `yaml:"-"`
	Inspector Inspector    `yaml:"-"`
	l         hclog.Logger `yaml:"-"`
}

type HoodRule struct {
	Name        string         `yaml:"name"`
	Descriotion string         `yaml:"description"`
	Regex       string         `yaml:"regex"`
	Reporting   bool           `yaml:"reporting"`
	Tokens      []string       `yaml:"tokens"`
	Regexp      *regexp.Regexp `yaml:"-"`
}

func (h *Hood) HoodInit(hq HeadQuarter, l hclog.Logger) error {
	// Check if file exist and a valid file
	fi, err := os.Stat(h.Path)
	if err != nil {
		return err
	} else if fi.IsDir() {
		return fmt.Errorf("%s is a directory", h.Path)
	}

	h.Hq = hq
	h.l = l

	for _, r := range h.Rules {
		r.Regexp = regexp.MustCompile(r.Regex)
	}

	h.Inspector = NewInvestigator(hq.ReportCrime(), hq.ReportError(), h.Path, h.Rules, h.Interval, l)
	h.l.Debug("Hood initialised", "Name", h.Name)
	return err
}

func (h *Hood) Watch() {
	go func() {
		h.l.Info("Start watching", "Hood", h.Name)
		h.Hq.Listen()
		h.Inspector.Investigate()
	}()
}
