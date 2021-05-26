package core

type ControllerPayload = interface{}

type Controller = interface{}

type BaseController = interface {
	ProvideService
}

type formData = map[string]interface{}

type body struct {
	Param Param `json:"param,omitempty"`
}

type OrderParam struct {
	Key       string `json:"key,omitempty"`
	Direction string `json:"direction,omitempty"`
}

type Param struct {
	ID        int          `json:"id,omitempty"`
	IDs       []int        `json:"ids,omitempty"`
	Form      bool         `json:"form,omitempty"`
	Fulltext  string       `json:"fulltext,omitempty"`
	Order     []OrderParam `json:"order,omitempty"`
	Offset    int          `json:"offset,omitempty"`
	ZeroLimit bool         `json:"zeroLimit,omitempty"`
	Limit     int          `json:"limit,omitempty"`
	All       bool         `json:"all,omitempty"`
	Filter    Filters      `json:"filter,omitempty"`
}
