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
