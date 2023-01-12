package process

import (
	"fmt"
	"strings"
	"time"
)

type Process struct {
	Total     int
	Current   int
	maxWidth  int
	mode      string
	pad       string
	endpoint  []string
	frequency time.Duration
	str       chan string
	done      chan struct{}
}

func NewProcess(total int, options ...Option) *Process {
	opt := applyOptions(options...)
	p := &Process{
		Total:     total,
		Current:   0,
		mode:      opt.mode,
		pad:       opt.pad,
		maxWidth:  opt.maxWidth,
		frequency: opt.frequency,
		endpoint:  opt.endpoint,
		str:       make(chan string, 1),
		done:      make(chan struct{}, 1),
	}

	return p
}

type options struct {
	maxWidth  int
	frequency time.Duration
	mode      string
	endpoint  []string
	pad       string
}

type Option func(r *options)

func WithMaxWidth(width int) Option {
	return func(r *options) {
		r.maxWidth = width
	}
}

func WithFrequency(frequency time.Duration) Option {
	return func(r *options) {
		r.frequency = frequency
	}
}

func WithMode(mode string) Option {
	return func(r *options) {
		r.mode = mode
	}
}

func WithPad(pad string) Option {
	return func(r *options) {
		r.pad = pad
	}
}

func (p *Process) Go() {
	str := ""
	go func() {
		fill := ""
		count, percent := 0, 0
		numWidth := len(fmt.Sprint(p.Total))
		middleStr := ">"
		length := len(p.endpoint)
		index := 0
		for p.Current <= p.Total {

			if p.mode != ">" {
				index += 1
				if index >= length {
					index = 0
				}

				middleStr = p.endpoint[index]
			}

			currentProcess := p.strPad(fmt.Sprint(p.Current), numWidth, " ", "LEFT")
			fill = strings.Repeat("-", count) + middleStr + strings.Repeat(p.pad, p.maxWidth-count)
			fill = p.strPad(fill, p.maxWidth, p.pad, "RIGHT")
			str = fmt.Sprintf("%s/%d [%s] %d%s", currentProcess, p.Total, fill, percent, "%")

			time.Sleep(p.frequency)
			p.str <- str
			percent = int(float32(p.Current) / float32(p.Total) * 100)
			count = int(float32(p.Current) / float32(p.Total) * float32(p.maxWidth))
		}

		fill = strings.Repeat("-", count)
		str = fmt.Sprintf("%d/%d [%s] %d%s", p.Current, p.Total, fill, percent, "%")
		p.str <- str
		p.done <- struct{}{}
	}()
}

func (p *Process) Finish() {
	close(p.str)
	close(p.done)
}

func (p *Process) Print() {
	go func() {
		for {
			select {
			case <-p.str:
				fmt.Printf("\r%s", <-p.str)
			case <-p.done:
				return
			}
		}
	}()
}

func (p *Process) strPad(input string, padLength int, padString string, padType string) string {

	output := ""
	inputLen := len(input)

	if inputLen >= padLength {
		return input
	}

	padStringLen := len(padString)
	needFillLen := padLength - inputLen

	if diffLen := padStringLen - needFillLen; diffLen > 0 {
		padString = padString[diffLen:]
	}

	for i := 1; i <= needFillLen; i += padStringLen {
		output += padString
	}

	switch padType {
	case "LEFT":
		return output + input
	default:
		return input + output
	}
}

func applyOptions(opts ...Option) *options {
	opt := &options{
		maxWidth:  50,
		mode:      ">",
		frequency: time.Second / 10,
		endpoint:  []string{"/", "-", "\\", "|"},
		pad:       " ",
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}
