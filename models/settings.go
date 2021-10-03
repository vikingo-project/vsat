package models

type Setting struct {
	Required bool `json:"required"`
}

type TextView struct {
	Setting
	Required bool   `json:"required"`
	Value    string `json:"value"`
	Label    string `json:"label"`
}

type InputStringView struct {
	Required bool   `json:"required"`
	Value    string `json:"value"`
	Default  string `json:"default"`
	Label    string `json:"label"`
}

type InputNumberView struct {
	Required bool   `json:"required"`
	Value    int64  `json:"value"`
	Default  int64  `json:"default"`
	Label    string `json:"label"`
}
