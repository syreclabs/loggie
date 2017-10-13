package loggie

import (
	"io"
	"sync"
	"time"
)

// These flags define which info will be included to each entry generated by the logger.
// By default, all logging entries are prefixed by date, time, level and logger name.
const (
	Fdate = 1 << iota
	Ftime
	Fmilliseconds
	Futc
	Flevel
	Fname
	Fcolor
	Fdefault = Fdate | Ftime | Fmilliseconds | Flevel | Fname
)

// Colors for date, time and logger name fields.
var (
	Cdate = "\033[0;37m"
	Ctime = "\033[0;37m"
	Cname = "\033[0;36m"
)

const creset = "\033[0m"

// Formatter is the interface implemented by output formatters.
type Formatter interface {
	Format(w io.Writer, lvl int, name string, msg string) error
	Flags() int
	SetFlags(flags int)
}

// Mlevel defines logging level field colorization.
var Mlevel = map[int][2]string{
	Ldebug:   {"DBG", "\033[0;34m"},
	Linfo:    {"INF", "\033[0;32m"},
	Lwarning: {"WRN", "\033[0;33m"},
	Lerror:   {"ERR", "\033[0;31m"},
	Lfatal:   {"FTL", "\033[1;31m"},
	Lpanic:   {"PAN", "\033[1;31m"},
}

type textFormatter struct {
	mu    sync.RWMutex
	flags int
}

// NextTextFormatter returns new default unstructured text formatter.
func NewTextFormatter(flags int) Formatter {
	return &textFormatter{flags: flags}
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
// Based on itoa() implementation in log/log.go.
func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// formatHeader writes log header to buf in following order:
//   * date and/or time (if corresponding flags are provided),
//   * level
//   * name (if it's not blank),
// Based on formatHeader() implementation in log/log.go.
func (f *textFormatter) formatHeader(buf *[]byte, lvl int, name string) {
	if f.flags&(Fdate|Ftime|Fmilliseconds) != 0 {
		now := time.Now()
		if f.flags&Futc != 0 {
			now = now.UTC()
		}
		if f.flags&Fdate != 0 {
			year, month, day := now.Date()
			if f.flags&Fcolor != 0 {
				*buf = append(*buf, Cdate...)
			}
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			if f.flags&Fcolor != 0 {
				*buf = append(*buf, creset...)
			}
			*buf = append(*buf, ' ')
		}
		if f.flags&(Ftime|Fmilliseconds) != 0 {
			hour, min, sec := now.Clock()
			if f.flags&Fcolor != 0 {
				*buf = append(*buf, Cdate...)
			}
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if f.flags&Fmilliseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, now.Nanosecond()/1e6, 3)
			}
			if f.flags&Fcolor != 0 {
				*buf = append(*buf, creset...)
			}
			*buf = append(*buf, ' ')
		}
	}

	if f.flags&Flevel != 0 {
		if f.flags&Fcolor != 0 {
			*buf = append(*buf, Mlevel[lvl][1]...)
		}
		*buf = append(*buf, Mlevel[lvl][0]...)
		if f.flags&Fcolor != 0 {
			*buf = append(*buf, creset...)
		}
		*buf = append(*buf, ' ')
	}

	if f.flags&Fname != 0 && name != "" {
		if f.flags&Fcolor != 0 {
			*buf = append(*buf, Cname...)
		}
		*buf = append(*buf, name...)
		if f.flags&Fcolor != 0 {
			*buf = append(*buf, creset...)
		}
		*buf = append(*buf, ' ')
	}
}

func (f *textFormatter) Format(w io.Writer, lvl int, name string, msg string) error {
	buf := pool.get()
	defer pool.put(buf)
	f.formatHeader(&buf, lvl, name)
	buf = append(buf, msg...)
	if msg == "" || msg[len(msg)-1] != '\n' {
		buf = append(buf, '\n')
	}
	_, err := w.Write(buf)
	return err
}

func (f *textFormatter) Flags() int {
	f.mu.RLock()
	defer f.mu.Unlock()
	return f.flags
}

func (f *textFormatter) SetFlags(flags int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.flags = flags
}

type bufferPool struct {
	sync.Pool
}

var pool = newBufferPool()

func newBufferPool() *bufferPool {
	return &bufferPool{
		Pool: sync.Pool{New: func() interface{} {
			return []byte{}
		}},
	}
}

func (p *bufferPool) get() []byte {
	return p.Pool.Get().([]byte)
}

func (p *bufferPool) put(b []byte) {
	p.Pool.Put(b)
}