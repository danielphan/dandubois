package migrate

import "github.com/danielphan/dandubois/painting/legacy"
import "net/http"
import "io"
import "encoding/json"
import "strings"

const categoriesURL = baseURL + "categories.json"

type categories map[int]string

func fetchCategories(c *http.Client) (categories, error) {
	res, err := c.Get(categoriesURL)
	if err != nil {
		return nil, err
	}
	return newCategories(res.Body)
}

func newCategories(r io.Reader) (categories, error) {
	var cs []*legacy.Category
	err := json.NewDecoder(r).Decode(&cs)
	if err != nil {
		return nil, err
	}

	c := categories{}
	for _, old := range cs {
		name := old.Fields.Name
		name = strings.Title(name)
		if name == "Urban" {
			name = "City"
		}
		c[old.Pk] = name
	}
	return c, nil
}

func (c categories) convert(ids []int) []string {
	namesSet := map[string]bool{}
	for _, id := range ids {
		namesSet[c[id]] = true
	}

	var names []string
	for name, _ := range namesSet {
		names = append(names, name)
	}
	return names
}
