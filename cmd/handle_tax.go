package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	tax "github.com/huyffs/shelly"
)

func (s *server) handleTax() http.HandlerFunc {
	type response struct {
		Type responseType `json:"type"`
		Data tax.Data     `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		fn := fmt.Sprintf("%s%s.json", s.taxDataDir, r.URL.Path[4:])
		f, err := os.Open(fn)
		if err != nil {
			s.NotFound(w, r, err)
			return
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			s.InternalServerError(w, err)
			return
		}
		var td tax.Data
		err = json.Unmarshal(b, &td)
		if err != nil {
			s.InternalServerError(w, err)
			return
		}

		res := response{
			Type: taxData,
			Data: td,
		}
		j, err := json.Marshal(res)
		if err != nil {
			s.InternalServerError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}
}
