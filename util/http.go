package util

import (
	"encoding/json"
	"net/http"
)

func ParseQuery[T any, R *T](r *http.Request) (R, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	m := map[string]string{}
	for k, v := range r.Form {
		m[k] = v[0]
	}

	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	filter := new(T)
	if err = json.Unmarshal(data, filter); err != nil {
		return nil, err
	}

	return filter, nil
}
