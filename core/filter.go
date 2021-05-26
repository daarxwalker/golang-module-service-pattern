package core

type Filter struct {
	FilterType string      `json:"filterType"`
	Value      interface{} `json:"value"`
}

type Filters = map[string]Filter

var FilterTypes = struct {
	Select string
	Date   string
}{
	Select: "select",
	Date:   "date",
}
