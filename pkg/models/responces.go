package models

import "time"

type Responce struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
}

type Data struct {
	Status []Status
}

type Status struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Code   string `json:"code"`
}

type UpdateListResponse struct {
	Success   bool        `json:"success"`
	RrUpdates []RrUpdates `json:"rr_updates"`
	Total     int         `json:"total"`
}

type RrUpdates struct {
	Name          string    `json:"name"`
	Hw            []string  `json:"hw"`
	Sw            []string  `json:"sw"`
	Latest        bool      `json:"latest"`
	Link          string    `json:"link"`
	Size          int       `json:"size"`
	Date          time.Time `json:"date"`
	Sha1          string    `json:"sha1"`
	Sha512        string    `json:"sha512"`
	UpdateVersion string    `json:"update_version,omitempty"`
}

// type SortByDate []RrUpdates

// func (inst SortByDate) Len() int { return len(inst) }
// func (inst SortByDate) Less(i, j int) bool {
// 	return inst[i].Date.After(inst[j].Date)
// }

// func (inst SortByDate) Swap(i, j int) {
// 	inst[i], inst[j] = inst[j], inst[i]
// }
