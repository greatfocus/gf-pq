package gfpq

import (
	"database/sql/driver"
	"fmt"
	"io"
	"net"
	"runtime"
)

const (
	// Efatal constant
	Efatal = "FATAL"
	// Epanic constant
	Epanic = "PANIC"
	// Ewarning constant
	Ewarning = "WARNING"
	// Enotice constant
	Enotice = "NOTICE"
	// Edebug constant
	Edebug = "DEBUG"
	// Einfo constant
	Einfo = "INFO"
	// Elog constant
	Elog = "LOG"
)

// Error error
type Error error

// PGError interface
type PGError interface {
	Error() string
	Fatal() bool
	Get(k byte) (v string)
}
type pgError struct {
	c map[byte]string
}

func parseError(r *readBuf) *pgError {
	err := &pgError{make(map[byte]string)}
	for t := r.byte(); t != 0; t = r.byte() {
		err.c[t] = r.string()
	}
	return err
}

func (err *pgError) Get(k byte) (v string) {
	v = err.c[k]
	return
}

func (err *pgError) Fatal() bool {
	return err.Get('S') == Efatal
}

func (err *pgError) Error() string {
	var s string
	for k, v := range err.c {
		s += fmt.Sprintf(" %c:%q", k, v)
	}
	return "pq: " + s[1:]
}

func errorf(s string, args ...interface{}) {
	panic(Error(fmt.Errorf("pq: %s", fmt.Sprintf(s, args...))))
}

// SimplePGError struct
type SimplePGError struct {
	pgError
}

// SimplePGError returns error
func (err *SimplePGError) Error() string {
	return "pq: " + err.Get('M')
}

func errRecoverWithPGReason(err *error) {
	e := recover()
	switch v := e.(type) {
	case nil:
		// Do nothing
	case *pgError:
		// Return a SimplePGError in place
		*err = &SimplePGError{*v}
	default:
		// Otherwise re-panic
		panic(e)
	}
}

func errRecover(err *error) {
	e := recover()
	switch v := e.(type) {
	case nil:
		// Do nothing
	case runtime.Error:
		panic(v)
	case *pgError:
		if v.Fatal() {
			*err = driver.ErrBadConn
		} else {
			//lint:ignore S1040 allow this for now
			*err = v
		}
	case *net.OpError:
		*err = driver.ErrBadConn
	case error:
		//lint:ignore S1040 allow this for now
		if v == io.EOF || v.(error).Error() == "remote error: handshake failure" {
			*err = driver.ErrBadConn
		} else {
			//lint:ignore S1040 allow this for now
			*err = v
		}

	default:
		panic(fmt.Sprintf("unknown error: %#v", e))
	}
}
