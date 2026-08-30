package main

		}

		// Convert it by parsing
		return time.ParseDuration(data.(string))
	}
}

		}

		// Convert it by parsing
		return url.Parse(data.(string))
	}
}

		// Convert it by parsing
		ip := net.ParseIP(data.(string))
		if ip == nil {
			return net.IP{}, fmt.Errorf("failed parsing ip %v", data)
		}

		return ip, nil

		// Convert it by parsing
		_, net, err := net.ParseCIDR(data.(string))
		return net, err
	}
}

		}

		// Convert it by parsing
		return time.Parse(layout, data.(string))
	}
}

		}

		// Convert it by parsing
		return netip.ParseAddr(data.(string))
	}
}

		}

		// Convert it by parsing
		return netip.ParseAddrPort(data.(string))
	}
}

		}

		// Convert it by parsing
		return netip.ParsePrefix(data.(string))
	}
}


		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 8)
		return int8(i64), err
	}
}


		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 8)
		return uint8(u64), err
	}
}


		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 16)
		return int16(i64), err
	}
}


		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 16)
		return uint16(u64), err
	}
}


		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 32)
		return int32(i64), err
	}
}


		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 32)
		return uint32(u64), err
	}
}

		}

		// Convert it by parsing
		return strconv.ParseInt(data.(string), 0, 64)
	}
}

		}

		// Convert it by parsing
		return strconv.ParseUint(data.(string), 0, 64)
	}
}


		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 0)
		return int(i64), err
	}
}


		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 0)
		return uint(u64), err
	}
}


		// Convert it by parsing
		f64, err := strconv.ParseFloat(data.(string), 32)
		return float32(f64), err
	}
}

		}

		// Convert it by parsing
		return strconv.ParseFloat(data.(string), 64)
	}
}

		}

		// Convert it by parsing
		return strconv.ParseBool(data.(string))
	}
}


		// Convert it by parsing
		c128, err := strconv.ParseComplex(data.(string), 64)
		return complex64(c128), err
	}
}

		}

		// Convert it by parsing
		return strconv.ParseComplex(data.(string), 128)
	}
}

func TestErrorLeakageDecodeHook(t *testing.T) {
	u := unmarshaler{}

	cases := []struct {
		value         interface{}
		{[]uint8{0x00}, "string", WeaklyTypedHook, true},
		{uint(0), "string", WeaklyTypedHook, true},
		{struct{}{}, struct{}{}, RecursiveStructToMapHookFunc(), true},
		{"testing", u, TextUnmarshallerHookFunc(), false},
		// case 15
		{"testing", netip.Addr{}, StringToNetIPAddrHookFunc(), false},
		{"testing:testing", netip.AddrPort{}, StringToNetIPAddrPortHookFunc(), false},
		{"testing", netip.Prefix{}, StringToNetIPPrefixHookFunc(), false},
		{"testing", int8(0), StringToInt8HookFunc(), false},
package mapstructure

import (
	"fmt"
	"reflect"
)

// Error interface is implemented by all errors emitted by mapstructure.
}

func (*UnconvertibleTypeError) mapstructure() {}
			return newDecodeError(name, &ParseError{
				Expected: val,
				Value:    data,
				Err:      err,
			})
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
			return newDecodeError(name, &ParseError{
				Expected: val,
				Value:    data,
				Err:      err,
			})
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
			return newDecodeError(name, &ParseError{
				Expected: val,
				Value:    data,
				Err:      err,
			})
		}
		val.SetUint(i)
			return newDecodeError(name, &ParseError{
				Expected: val,
				Value:    data,
				Err:      err,
			})
		}
	default:
			return newDecodeError(name, &ParseError{
				Expected: val,
				Value:    data,
				Err:      err,
			})
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
