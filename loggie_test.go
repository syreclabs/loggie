package loggie

import (
	"bytes"
	"regexp"
	"testing"
)

// Ensure that defaultLogger implements Logger
var _ Logger = &defaultLogger{}

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	f := NewTextFormatter(Fdefault)
	l := NewLogger(&buf, "TestPrint", f)

	t.Run("debug", func(t *testing.T) {
		buf.Reset()
		l.Print(Ldebug, "print debug")
		rx := regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} DBG TestPrint print debug` + "\n" + `\z`)
		b := buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Printf(Ldebug, "printf: %s", "debug")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} DBG TestPrint printf: debug` + "\n" + `\z`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Debug("debug debug")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} DBG TestPrint debug debug` + "\n" + `\z`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Debugf("debugf: %s", "debug")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} DBG TestPrint debugf: debug` + "\n" + `\z`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}
	})

	t.Run("info", func(t *testing.T) {
		buf.Reset()
		l.Print(Linfo, "print info")
		rx := regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} INF TestPrint print info` + "\n" + `\z`)
		b := buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Printf(Linfo, "printf: %s", "info")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} INF TestPrint printf: info` + "\n" + `\z`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Info("info info")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} INF TestPrint info info` + "\n" + `\z`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Infof("infof: %s", "info")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} INF TestPrint infof: info` + "\n" + `\z`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}
	})

	t.Run("warning", func(t *testing.T) {
		buf.Reset()
		l.Print(Lwarning, "print warning")
		rx := regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} WRN TestPrint print warning` + "\n" + `\z`)
		b := buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Printf(Lwarning, "Printf: %s", "printf warning")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} WRN TestPrint Printf: printf warning` + "\n" + `\z`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Warning("warning warning")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} WRN TestPrint warning warning` + "\n" + `\z`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Warningf("Printf: %s", "warning")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} WRN TestPrint Printf: warning` + "\n" + `\z`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}
	})

	t.Run("error", func(t *testing.T) {
		buf.Reset()
		l.Print(Lerror, "print error")
		rx := regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} ERR TestPrint print error` + "\n" + `.+`)
		b := buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Printf(Lerror, "Printf: %s", "printf error")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} ERR TestPrint Printf: printf error` + "\n" + `.+`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Error("error error")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} ERR TestPrint error error` + "\n" + `.+`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}

		buf.Reset()
		l.Errorf("Printf: %s", "error")
		rx = regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} ERR TestPrint Printf: error` + "\n" + `.+`)
		b = buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}
	})
}

func TestLevel(t *testing.T) {
	var buf bytes.Buffer
	f := NewTextFormatter(Fdefault)
	l := NewLogger(&buf, "TestPrint", f)

	t.Run("default", func(t *testing.T) {
		buf.Reset()
		l.Print(Ldebug, "debug")
		if len(buf.Bytes()) == 0 {
			t.Error("expected output not to be empty")
		}

		buf.Reset()
		l.Print(Lerror, "error")
		if len(buf.Bytes()) == 0 {
			t.Error("expected output not to be empty")
		}
	})

	t.Run("warning", func(t *testing.T) {
		l.SetLevel(Lwarning)

		buf.Reset()
		l.Print(Ldebug, "debug")
		if len(buf.Bytes()) != 0 {
			t.Error("expected output to be empty")
		}

		buf.Reset()
		l.Print(Lwarning, "warning")
		if len(buf.Bytes()) == 0 {
			t.Error("expected output not to be empty")
		}
	})

	t.Run("error", func(t *testing.T) {
		l.SetLevel(Lerror)

		buf.Reset()
		l.Print(Ldebug, "info")
		if len(buf.Bytes()) != 0 {
			t.Error("expected output to be empty")
		}

		buf.Reset()
		l.Print(Lwarning, "warning")
		if len(buf.Bytes()) != 0 {
			t.Error("expected output to be empty")
		}
	})
}

func TestPanic(t *testing.T) {
	var buf bytes.Buffer

	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic")
		}
		rx := regexp.MustCompile(`\A\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} PAN TestPanic panic !!` + "\n" + `\z`)
		b := buf.Bytes()
		if !rx.Match(buf.Bytes()) {
			t.Errorf("expected %q to match %#q", b, rx)
		}
	}()

	f := NewTextFormatter(Fdefault)
	l := NewLogger(&buf, "TestPanic", f)

	l.Panic("panic !!")
}

func TestNone(t *testing.T) {
	var buf bytes.Buffer

	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic")
		}
		if len(buf.Bytes()) != 0 {
			t.Errorf("expected output to be empty")
		}
	}()

	f := NewTextFormatter(Fdefault)
	l := NewLogger(&buf, "TestLevelNone", f)
	l.SetLevel(Lnone)

	l.Panic("panic")
}
