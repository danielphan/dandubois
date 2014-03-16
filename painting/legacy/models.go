package legacy

type Category struct {
	Pk     int    `json:"pk"`
	Model  string `json:"model"`
	Fields struct {
		Enabled   int    `json:"enabled"`
		Name      string `json:"name"`
		ShortForm string `json:"short_form"`
	} `json:"fields"`
}

type Medium struct {
	Pk     int    `json:"pk"`
	Model  string `json:"model"`
	Fields struct {
		Enabled   int    `json:"enabled"`
		Name      string `json:"name"`
		ShortForm string `json:"short_form"`
	} `json:"fields"`
}

type Painting struct {
	Pk     int    `json:"pk"`
	Model  string `json:"model"`
	Fields struct {
		Status      string `json:"status"`
		Medium      int    `json:"medium"`
		Description string `json:"description"`
		Title       string `json:"title"`
		Price       int    `json:"price"`
		Enabled     int    `json:"enabled"`
		Number      int    `json:"number"`
		Height      int    `json:"height"`
		Width       int    `json:"width"`
		Featured    int    `json:"featured"`
		Year        int    `json:"year"`
		Framed      int    `json:"framed"`
		Image       string `json:"image"`
		Categories  []int  `json:"categories"`
	} `json:"fields"`
}
