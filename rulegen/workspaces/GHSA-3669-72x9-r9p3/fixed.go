package main

	"strings"
)

const (
	defaultMaxSize = 16000
)

// NewDecoder returns a new Decoder.
func NewDecoder() *Decoder {
	return &Decoder{cache: newCache(), maxSize: defaultMaxSize}
}

// Decoder decodes values from a map[string][]string to a struct.
	cache             *cache
	zeroEmpty         bool
	ignoreUnknownKeys bool
	maxSize           int
}

// SetAliasTag changes the tag used to locate custom field aliases.
	d.ignoreUnknownKeys = i
}

// MaxSize limits the size of slices for URL nested arrays or object arrays.
// Choose MaxSize carefully; large values may create many zero-value slice elements.
// Example: "items.100000=apple" would create a slice with 100,000 empty strings.
func (d *Decoder) MaxSize(size int) {
	d.maxSize = size
}

// RegisterConverter registers a converter function for a custom type.
func (d *Decoder) RegisterConverter(value interface{}, converterFunc Converter) {
	d.cache.registerConverter(value, converterFunc)
	// Slice of structs. Let's go recursive.
	if len(parts) > 1 {
		idx := parts[0].index
		// a defensive check to avoid creating a large slice based on user input index
		if idx > d.maxSize {
			return fmt.Errorf("%v index %d is larger than the configured maxSize %d", v.Kind(), idx, d.maxSize)
		}
		if v.IsNil() || v.Len() < idx+1 {
			value := reflect.MakeSlice(t, idx+1, idx+1)
			if v.Len() < idx+1 {

type S24e struct {
	*S24
	F2 string `schema:"F2"`
}

func TestUnmarshallToEmbeddedNoData(t *testing.T) {
	s := &S24e{}

	decoder := NewDecoder()
	err := decoder.Decode(s, data)

	expectedErr := `schema: invalid path "F3"`
	if err.Error() != expectedErr {
		t.Fatalf("got %q, want %q", err, expectedErr)
	}
}

type S25ee struct {
	F3 string `schema:"F3"`
}
	F1 string `schema:"F1"`
}

func TestDoubleEmbedded(t *testing.T) {
	data := map[string][]string{
		"F1": {"raw a"},
		"F2": {"raw b"},
		"F3": {"raw c"},
	}

	s := S25{}
	decoder := NewDecoder()

		t.Errorf("decoding should fail with error msg %s got %q", expected, err)
	}
}

func TestDecoder_MaxSize(t *testing.T) {
	t.Parallel()

	type Nested struct {
		Val          int
		NestedValues []struct {
			NVal int
		}
	}
	type NestedSlices struct {
		Values []Nested
	}

	testcases := []struct {
		name            string
		maxSize         int
		decoderInput    func() (dst NestedSlices, src map[string][]string)
		expectedDecoded NestedSlices
		expectedErr     MultiError
	}{
		{
			name:    "no error on decoding under max size",
			maxSize: 10,
			decoderInput: func() (dst NestedSlices, src map[string][]string) {
				return dst, map[string][]string{
					"Values.1.Val":                 {"132"},
					"Values.1.NestedValues.1.NVal": {"1"},
					"Values.1.NestedValues.2.NVal": {"2"},
					"Values.1.NestedValues.3.NVal": {"3"},
				}
			},
			expectedDecoded: NestedSlices{
				Values: []Nested{
					{
						Val:          0,
						NestedValues: nil,
					},
					{
						Val: 132, NestedValues: []struct{ NVal int }{
							{NVal: 0},
							{NVal: 1},
							{NVal: 2},
							{NVal: 3},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:    "error on decoding above max size",
			maxSize: 1,
			decoderInput: func() (dst NestedSlices, src map[string][]string) {
				return dst, map[string][]string{
					"Values.1.Val":                 {"132"},
					"Values.1.NestedValues.1.NVal": {"1"},
					"Values.1.NestedValues.2.NVal": {"2"},
					"Values.1.NestedValues.3.NVal": {"3"},
				}
			},
			expectedErr: MultiError{
				"Values.1.NestedValues.2.NVal": errors.New("slice index 2 is larger than the configured maxSize 1"),
				"Values.1.NestedValues.3.NVal": errors.New("slice index 3 is larger than the configured maxSize 1"),
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dec := NewDecoder()
			dec.MaxSize(tc.maxSize)
			dst, src := tc.decoderInput()
			err := dec.Decode(&dst, src)

			if tc.expectedErr != nil {
				var gotErr MultiError
				if !errors.As(err, &gotErr) {
					t.Errorf("decoder error is not of type %T", gotErr)
				}
				if !reflect.DeepEqual(gotErr, tc.expectedErr) {
					t.Errorf("expected %v, got %v", tc.expectedErr, gotErr)
				}
			} else {
				if !reflect.DeepEqual(dst, tc.expectedDecoded) {
					t.Errorf("expected %v, got %v", tc.expectedDecoded, dst)
				}
			}
		})
	}
}

func TestDecoder_SetMaxSize(t *testing.T) {

	t.Run("default maxsize should be equal to given constant", func(t *testing.T) {
		t.Parallel()
		dec := NewDecoder()
		if !reflect.DeepEqual(dec.maxSize, defaultMaxSize) {
			t.Errorf("unexpected default max size")
		}
	})

	t.Run("configured maxsize should be set properly", func(t *testing.T) {
		t.Parallel()
		configuredMaxSize := 50
		limitedMaxSizeDecoder := NewDecoder()
		limitedMaxSizeDecoder.MaxSize(configuredMaxSize)
		if !reflect.DeepEqual(limitedMaxSizeDecoder.maxSize, configuredMaxSize) {
			t.Errorf("invalid decoder maxsize, expected: %d, got: %d",
				configuredMaxSize, limitedMaxSizeDecoder.maxSize)
		}
	})
}
