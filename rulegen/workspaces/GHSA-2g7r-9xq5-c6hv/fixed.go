package main

	"net/http"
	"net/url"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type Image struct {
		return nil, err
	}

	bodyBytes, err = SanitizeContent(bodyBytes)
	if err != nil {
		return nil, err
	}

	image := &Image{
		Blob:      bodyBytes,
		Mediatype: mediatype,
	}
	return image, nil
}

func SanitizeContent(content []byte) ([]byte, error) {
	bodyString := string(content)

	bm := bluemonday.UGCPolicy()
	return []byte(bm.Sanitize(bodyString)), nil
}
