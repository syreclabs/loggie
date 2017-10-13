Simple unstructured leveled logging.

[![GoDoc](https://godoc.org/syreclabs.com/go/loggie?status.svg)](https://godoc.org/syreclabs.com/go/loggie)
[![Build Status](https://travis-ci.org/dmgk/loggie.svg?branch=master)](https://travis-ci.org/dmgk/loggie)

### Installation

    go get -u syreclabs.com/go/loggie

### Usage

    import "syreclabs.com/go/loggie"

    var log = loggie.New("mylogger")

    func main() {
        log.Debug("debug message")
        log.Infof("info: %s", time.Now())
        log.Printf(loggie.Lwarning, "warning: %d", 42)
    }

Output:

    2018/06/06 08:59:29.639 DBG mylogger debug message
    2018/06/06 08:59:29.639 INF mylogger info: 2018-06-06 08:59:29.639226176 -0500 -05 m=+0.000262528
    2018/06/06 08:59:29.639 WRN mylogger warning: 42
