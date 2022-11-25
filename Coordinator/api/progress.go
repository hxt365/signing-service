package api

import (
	"coordinator/utils"
	"net/http"
	"strconv"
)

func (s *Server) handleProgress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			s.handleGetNextRange()(w, r)
		case "POST":
			s.handleSetCommittedOffset()(w, r)
		default:
			respondHTTPErr(w, r, http.StatusMethodNotAllowed)
		}
	}
}

func (s *Server) handleGetNextRange() http.HandlerFunc {
	type response struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		keyParam := r.URL.Query().Get("key_id")
		keyID, err := strconv.Atoi(keyParam)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, "malformed request format")
			return
		}

		start, end, err := s.pu.GetNextRange(keyID)
		if err != nil {
			if e, ok := err.(utils.ValidationErr); ok {
				respondHTTPErr(w, r, http.StatusBadRequest, e)
				return
			}
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		respond(w, r, http.StatusOK, response{
			Start: start,
			End:   end,
		})
	}
}

func (s *Server) handleSetCommittedOffset() http.HandlerFunc {
	type request struct {
		KeyID  int `json:"key_id"`
		Offset int `json:"offset"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeBody(r, &req); err != nil {
			respondErr(w, r, http.StatusBadRequest, "malformed request format")
			return
		}

		if err := s.pu.SetCommittedOffset(req.KeyID, req.Offset); err != nil {
			if e, ok := err.(utils.ValidationErr); ok {
				respondHTTPErr(w, r, http.StatusBadRequest, e)
				return
			}
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		respond(w, r, http.StatusOK, req)
	}
}
