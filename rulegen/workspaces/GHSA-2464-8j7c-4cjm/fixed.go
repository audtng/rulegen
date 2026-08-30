package main

		}

		// Convert it by parsing
		d, err := time.ParseDuration(data.(string))

		return d, wrapTimeParseDurationError(err)
	}
}

		}

		// Convert it by parsing
		u, err := url.Parse(data.(string))

		return u, wrapUrlError(err)
	}
}

		// Convert it by parsing
		ip := net.ParseIP(data.(string))
		if ip == nil {
			return net.IP{}, fmt.Errorf("failed parsing ip")
		}

		return ip, nil

		// Convert it by parsing
		_, net, err := net.ParseCIDR(data.(string))
		return net, wrapNetParseError(err)
	}
}

		}

		// Convert it by parsing
		ti, err := time.Parse(layout, data.(string))

		return ti, wrapTimeParseError(err)
	}
}

		}

		// Convert it by parsing
		addr, err := netip.ParseAddr(data.(string))

		return addr, wrapNetIPParseAddrError(err)
	}
}

		}

		// Convert it by parsing
		addrPort, err := netip.ParseAddrPort(data.(string))

		return addrPort, wrapNetIPParseAddrPortError(err)
	}
}

		}

		// Convert it by parsing
		prefix, err := netip.ParsePrefix(data.(string))

		return prefix, wrapNetIPParsePrefixError(err)
	}
}


		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 8)
		return int8(i64), wrapStrconvNumError(err)
	}
}


		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 8)
		return uint8(u64), wrapStrconvNumError(err)
	}
}


		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 16)
		return int16(i64), wrapStrconvNumError(err)
	}
}


		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 16)
		return uint16(u64), wrapStrconvNumError(err)
	}
}


		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 32)
		return int32(i64), wrapStrconvNumError(err)
	}
}


		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 32)
		return uint32(u64), wrapStrconvNumError(err)
	}
}

		}

		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 64)
		return int64(i64), wrapStrconvNumError(err)
	}
}

		}

		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 64)
		return uint64(u64), wrapStrconvNumError(err)
	}
}


		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 0)
		return int(i64), wrapStrconvNumError(err)
	}
}


		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 0)
		return uint(u64), wrapStrconvNumError(err)
	}
}


		// Convert it by parsing
		f64, err := strconv.ParseFloat(data.(string), 32)
		return float32(f64), wrapStrconvNumError(err)
	}
}

		}

		// Convert it by parsing
		f64, err := strconv.ParseFloat(data.(string), 64)
		return f64, wrapStrconvNumError(err)
	}
}

		}

		// Convert it by parsing
		b, err := strconv.ParseBool(data.(string))
		return b, wrapStrconvNumError(err)
	}
}


		// Convert it by parsing
		c128, err := strconv.ParseComplex(data.(string), 64)
		return complex64(c128), wrapStrconvNumError(err)
	}
}

		}

		// Convert it by parsing
		c128, err := strconv.ParseComplex(data.(string), 128)
		return c128, wrapStrconvNumError(err)
	}
}

func TestErrorLeakageDecodeHook(t *testing.T) {
	u := unmarshaler{}
	_ = u

	cases := []struct {
		value         interface{}
		{[]uint8{0x00}, "string", WeaklyTypedHook, true},
		{uint(0), "string", WeaklyTypedHook, true},
		{struct{}{}, struct{}{}, RecursiveStructToMapHookFunc(), true},
		{"testing", netip.Addr{}, StringToNetIPAddrHookFunc(), false},
		// {"testing", u, TextUnmarshallerHookFunc(), false},
		// case 15
		{"testing:testing", netip.AddrPort{}, StringToNetIPAddrPortHookFunc(), false},
		{"testing", netip.Prefix{}, StringToNetIPPrefixHookFunc(), false},
		{"testing", int8(0), StringToInt8HookFunc(), false},
package mapstructure

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Error interface is implemented by all errors emitted by mapstructure.
}

func (*UnconvertibleTypeError) mapstructure() {}

func wrapStrconvNumError(err error) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(*strconv.NumError); ok {
		return &strconvNumError{Err: err}
	}

	return err
}

type strconvNumError struct {
	Err *strconv.NumError
}

func (e *strconvNumError) Error() string {
	return "strconv." + e.Err.Func + ": " + e.Err.Err.Error()
}

func (e *strconvNumError) Unwrap() error { return e.Err }

func wrapUrlError(err error) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(*url.Error); ok {
		return &urlError{Err: err}
	}

	return err
}

type urlError struct {
	Err *url.Error
}

func (e *urlError) Error() string {
	return fmt.Sprintf("%s", e.Err.Err)
}

func (e *urlError) Unwrap() error { return e.Err }

func wrapNetParseError(err error) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(*net.ParseError); ok {
		return &netParseError{Err: err}
	}

	return err
}

type netParseError struct {
	Err *net.ParseError
}

func (e *netParseError) Error() string {
	return "invalid " + e.Err.Type
}

func (e *netParseError) Unwrap() error { return e.Err }

func wrapTimeParseError(err error) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(*time.ParseError); ok {
		return &timeParseError{Err: err}
	}

	return err
}

type timeParseError struct {
	Err *time.ParseError
}

func (e *timeParseError) Error() string {
	if e.Err.Message == "" {
		return fmt.Sprintf("parsing time as %q: cannot parse as %q", e.Err.Layout, e.Err.LayoutElem)
	}

	return "parsing time " + e.Err.Message
}

func (e *timeParseError) Unwrap() error { return e.Err }

func wrapNetIPParseAddrError(err error) error {
	if err == nil {
		return nil
	}

	if errMsg := err.Error(); strings.HasPrefix(errMsg, "ParseAddr") {
		errPieces := strings.Split(errMsg, ": ")

		return fmt.Errorf("ParseAddr: %s", errPieces[len(errPieces)-1])
	}

	return err
}

func wrapNetIPParseAddrPortError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	if strings.HasPrefix(errMsg, "invalid port ") {
		return errors.New("invalid port")
	} else if strings.HasPrefix(errMsg, "invalid ip:port ") {
		return errors.New("invalid ip:port")
	}

	return err
}

func wrapNetIPParsePrefixError(err error) error {
	if err == nil {
		return nil
	}

	if errMsg := err.Error(); strings.HasPrefix(errMsg, "netip.ParsePrefix") {
		errPieces := strings.Split(errMsg, ": ")

		return fmt.Errorf("netip.ParsePrefix: %s", errPieces[len(errPieces)-1])
	}

	return err
}

func wrapTimeParseDurationError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	if strings.HasPrefix(errMsg, "time: unknown unit ") {
		return errors.New("time: unknown unit")
	} else if strings.HasPrefix(errMsg, "time: ") {
		idx := strings.LastIndex(errMsg, " ")

		return errors.New(errMsg[:idx])
	}

	return err
}
			return newDecodeError(name, &ParseError{
				Expected: val,
				Value:    data,
				Err:      wrapStrconvNumError(err),
			})
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
			return newDecodeError(name, &ParseError{
				Expected: val,
				Value:    data,
				Err:      wrapStrconvNumError(err),
			})
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
			return newDecodeError(name, &ParseError{
				Expected: val,
				Value:    data,
				Err:      wrapStrconvNumError(err),
			})
		}
		val.SetUint(i)
			return newDecodeError(name, &ParseError{
				Expected: val,
				Value:    data,
				Err:      wrapStrconvNumError(err),
			})
		}
	default:
			return newDecodeError(name, &ParseError{
				Expected: val,
				Value:    data,
				Err:      wrapStrconvNumError(err),
			})
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
