package main

	"strings"
)

// NewDecoder returns a new Decoder.
func NewDecoder() *Decoder {
	return &Decoder{cache: newCache()}
}

// Decoder decodes values from a map[string][]string to a struct.
	cache             *cache
	zeroEmpty         bool
	ignoreUnknownKeys bool
}

// SetAliasTag changes the tag used to locate custom field aliases.
	d.ignoreUnknownKeys = i
}

// RegisterConverter registers a converter function for a custom type.
func (d *Decoder) RegisterConverter(value interface{}, converterFunc Converter) {
	d.cache.registerConverter(value, converterFunc)
	// Slice of structs. Let's go recursive.
	if len(parts) > 1 {
		idx := parts[0].index
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
	err := decoder.Decode(s, data);
	
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

func TestDoubleEmbedded(t *testing.T){
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
