package server

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Bauka07/AP2/internal/model"
	"github.com/Bauka07/AP2/internal/store"
	"github.com/Bauka07/AP2/internal/utils"
)

type Server struct{
	mux *http.ServeMux
	St *store.Store[string, string]
	startedAt time.Time
	ReqCounter atomic.Int64
	WorkerStop chan struct{}
	WorkerDone chan struct{}
}

func NewServer(store *store.Store[string, string]) *Server {
	return &Server{
		mux: http.NewServeMux(),
		St: store,
		startedAt: time.Now(),
		WorkerStop: make(chan struct{}),
		WorkerDone: make(chan struct{}),
	}
}

func (s *Server) Router() http.Handler {
	s.mux.HandleFunc("/data", s.handleData)
	s.mux.HandleFunc("/data/", s.handleDataByKey)
	s.mux.HandleFunc("/stats", s.handleStats)
	return s.mux
}


//Handle data
func (s *Server) handleData(w http.ResponseWriter, r *http.Request) {
	s.ReqCounter.Add(1)
	switch r.Method {
	//get all
	case http.MethodGet:
		utils.WriteJSON(w, 200, s.St.Snapshot())
		return
	//create
	case http.MethodPost:
		var req model.Data
		if err := utils.ReadJSON(r, &req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if req.Key == "" || req.Value == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "key and value must be non-empty"})
			return
		}
		s.St.Set(req.Key, req.Value)
		utils.WriteJSON(w, 201, req)
		return
	default:
		w.Header().Set("Allow", "GET, POST")
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
}

//Handle data By Key
func (s *Server) handleDataByKey(w http.ResponseWriter, r *http.Request) {
	s.ReqCounter.Add(1)
	key, err := utils.GetKey(r.URL.Path);
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	switch r.Method {
	//get by key
	case http.MethodGet:
		v, ok := s.St.Get(key)
		if !ok {
			utils.WriteJSON(w, 404, map[string]string{"message": "Key does not exist"})
			return
		}
		utils.WriteJSON(w, http.StatusOK, model.Data{Key: key, Value: v})
		return
	//delete
	case http.MethodDelete:
		ok := s.St.Delete(key)
		if !ok {
			utils.WriteJSON(w, 404, map[string]string{"message": "Key does not exist"})
			return
		}
		utils.WriteJSON(w, 200, map[string]string{"message": "Successfully deleted"})
		return
	default:
		w.Header().Set("Allow", "GET, DELETE")
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
}

//get stats

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	s.ReqCounter.Add(1)

	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	statusRes := model.Stats {
		Requests: s.ReqCounter.Load(),
		Keys: s.St.Len(),
		UptimeSeconds: int64(time.Since(s.startedAt).Seconds()),
	}
	utils.WriteJSON(w, 200, statusRes)
}
