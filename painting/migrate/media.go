package migrate

import "github.com/danielphan/dandubois/painting/legacy"
import "net/http"
import "io"
import "encoding/json"
import "strings"

const mediaURL = baseURL + "media.json"

type media map[int][]string

func fetchMedia(c *http.Client) (media, error) {
	res, err := c.Get(mediaURL)
	if err != nil {
		return nil, err
	}
	return newMedia(res.Body)
}

func newMedia(r io.Reader) (media, error) {
	var ms []*legacy.Medium
	err := json.NewDecoder(r).Decode(&ms)
	if err != nil {
		return nil, err
	}

	m := media{}
	for _, old := range ms {
		namesSet := map[string]bool{}
		for _, name := range strings.Split(old.Fields.Name, " ") {
			if name == "" || name == " " ||
				name == "on" || name == "and" {
				continue
			}
			name = strings.Title(name)
			namesSet[name] = true
		}

		var names []string
		for name, _ := range namesSet {
			names = append(names, name)
		}

		m[old.Pk] = names
	}
	return m, nil
}

func (m media) convert(id int) []string {
	return m[id]
}
