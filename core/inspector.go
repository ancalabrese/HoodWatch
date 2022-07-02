package core

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
)

type Inspector interface {
	Investigate()
	OnError(err error)
	ReportCrime(c *Crime)
}

type Investigator struct {
	eventBus chan *Crime
	errBus   chan error
	fPath    string
	rules    []*HoodRule
	lastPos  int64
	freq     int
	l        hclog.Logger
}

func NewInvestigator(event chan *Crime, e chan error, fPath string, rules []*HoodRule, freq int, l hclog.Logger) *Investigator {
	return &Investigator{
		eventBus: event,
		errBus:   e,
		rules:    rules,
		fPath:    fPath,
		lastPos:  0,
		freq:     freq,
		l:        l,
	}
}

func (i *Investigator) Investigate() {
	for {
		fd, err := os.Open(i.fPath)
		if err != nil {
			i.OnError(err)
			time.Sleep(time.Duration(i.freq) * time.Second)
			continue
		}

		_, err = fd.Seek(i.lastPos, io.SeekStart)
		if err != nil {
			i.OnError(err)
			time.Sleep(time.Duration(i.freq) * time.Second)
			continue
		}

		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			ln := scanner.Text()
			// Looping through the rules for this 'hood for a match
			// TODO: there must be a better way
			for _, r := range i.rules {
				if match := r.Regexp.MatchString(ln); match {
					matches := r.Regexp.FindAllStringSubmatch(ln, -1)[0]
					payload := make(map[string]string)
					payload["event"] = matches[0]
					for index := 1; index < len(matches); index++ {
						payload[r.Tokens[index-1]] = matches[index]
					}
					i.ReportCrime(NewCrime(r.Name, payload))
				}
			}
		}

		if scanner.Err() != nil {
			i.OnError(err)
			time.Sleep(time.Duration(i.freq) * time.Second)
			continue
		}

		// Update last position to current last line so we'll only check
		// new lines next time
		i.lastPos, _ = fd.Seek(0, io.SeekCurrent)
		fd.Close()
		time.Sleep(time.Duration(i.freq) * time.Second)
	}
}

func (i *Investigator) OnError(e error) {
	i.errBus <- e
}

func (i *Investigator) ReportCrime(c *Crime) {
	i.eventBus <- c
}
