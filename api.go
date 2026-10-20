package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	store *Store
}

func NewHandler() http.Handler {
	server := &Server{
		store: NewStore(),
	}

	mux := http.NewServeMux()

	// B1: Sumber Daya
	mux.HandleFunc("POST /sumber-daya", server.handlePostSumberDaya)
	mux.HandleFunc("GET /sumber-daya", server.handleGetSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", server.handleGetSumberDayaDetail)

	// B2 & B3: Peminjaman
	mux.HandleFunc("POST /peminjaman", server.handlePostPeminjaman)
	mux.HandleFunc("GET /peminjaman", server.handleGetPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", server.handleGetPeminjamanDetail)

	return mux
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// --- Handlers B1: Sumber Daya ---

func (s *Server) handlePostSumberDaya(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nama  string `json:"nama"`
		Jenis string `json:"jenis"`
		Stok  int    `json:"stok"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Body JSON tidak valid")
		return
	}

	nama := strings.TrimSpace(req.Nama)
	jenis := strings.TrimSpace(req.Jenis)

	if nama == "" || utf8RuneCount(nama) > 50 {
		writeError(w, http.StatusBadRequest, "Nama wajib diisi dan maksimal 50 karakter")
		return
	}

	if jenis != "ruang" && jenis != "alat" {
		writeError(w, http.StatusBadRequest, "Jenis harus 'ruang' atau 'alat'")
		return
	}

	if req.Stok < 1 {
		writeError(w, http.StatusBadRequest, "Stok minimal 1")
		return
	}

	sd := s.store.CreateSumberDaya(nama, jenis, req.Stok)
	writeJSON(w, http.StatusCreated, sd)
}

func (s *Server) handleGetSumberDaya(w http.ResponseWriter, r *http.Request) {
	list := s.store.GetAllSumberDaya()
	writeJSON(w, http.StatusOK, list)
}

func (s *Server) handleGetSumberDayaDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	sd, found := s.store.GetSumberDayaByID(id)
	if !found {
		writeError(w, http.StatusNotFound, "Sumber daya tidak ditemukan")
		return
	}

	writeJSON(w, http.StatusOK, sd)
}

// --- Handlers B2 & B3: Peminjaman ---

func (s *Server) handlePostPeminjaman(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SumberDayaID int    `json:"sumber_daya_id"`
		Peminjam     string `json:"peminjam"`
		Tanggal      string `json:"tanggal"`
		SlotMulai    int    `json:"slot_mulai"`
		SlotSelesai  int    `json:"slot_selesai"`
		Jumlah       int    `json:"jumlah"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Body JSON tidak valid")
		return
	}

	peminjam := strings.TrimSpace(req.Peminjam)
	tanggal := strings.TrimSpace(req.Tanggal)

	if peminjam == "" || utf8RuneCount(peminjam) > 50 {
		writeError(w, http.StatusBadRequest, "Peminjam wajib diisi dan maksimal 50 karakter")
		return
	}

	if !isValidDate(tanggal) {
		writeError(w, http.StatusBadRequest, "Format tanggal tidak valid atau tanggal tidak ada")
		return
	}

	if req.SlotMulai < 1 || req.SlotSelesai > 8 || req.SlotMulai > req.SlotSelesai {
		writeError(w, http.StatusBadRequest, "Slot harus antara 1 sampai 8 dan slot_mulai <= slot_selesai")
		return
	}

	sd, found := s.store.GetSumberDayaByID(req.SumberDayaID)
	if !found {
		writeError(w, http.StatusBadRequest, "Sumber daya tidak ditemukan")
		return
	}

	if req.Jumlah < 1 || req.Jumlah > sd.Stok {
		writeError(w, http.StatusBadRequest, "Jumlah melebihi stok atau kurang dari 1")
		return
	}

	pem := s.store.CreatePeminjaman(req.SumberDayaID, peminjam, tanggal, req.SlotMulai, req.SlotSelesai, req.Jumlah)
	writeJSON(w, http.StatusCreated, pem)
}

func (s *Server) handleGetPeminjaman(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filterSDStr := query.Get("sumber_daya_id")
	filterTanggal := query.Get("tanggal")
	filterStatus := query.Get("status")

	var filterSDID *int
	if filterSDStr != "" {
		id, err := strconv.Atoi(filterSDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "sumber_daya_id filter tidak valid")
			return
		}
		filterSDID = &id
	}

	list := s.store.GetPeminjamanFiltered(filterSDID, filterTanggal, filterStatus)
	writeJSON(w, http.StatusOK, list)
}

func (s *Server) handleGetPeminjamanDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	pem, found := s.store.GetPeminjamanByID(id)
	if !found {
		writeError(w, http.StatusNotFound, "Peminjaman tidak ditemukan")
		return
	}

	writeJSON(w, http.StatusOK, pem)
}