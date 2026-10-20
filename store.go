package main

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
	PaketID      int    `json:"paket_id"`
}

type Store struct {
	mu         sync.Mutex
	sumberDaya []SumberDaya
	peminjaman []Peminjaman
	nextSDID   int
	nextPemID  int
}

func NewStore() *Store {
	return &Store{
		sumberDaya: make([]SumberDaya, 0),
		peminjaman: make([]Peminjaman, 0),
		nextSDID:   1,
		nextPemID:  1,
	}
}

// --- Sumber Daya Ops ---

func (s *Store) CreateSumberDaya(nama, jenis string, stok int) SumberDaya {
	s.mu.Lock()
	defer s.mu.Unlock()

	sd := SumberDaya{
		ID:    s.nextSDID,
		Nama:  nama,
		Jenis: jenis,
		Stok:  stok,
	}
	s.nextSDID++
	s.sumberDaya = append(s.sumberDaya, sd)
	return sd
}

func (s *Store) GetAllSumberDaya() []SumberDaya {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sumberDaya == nil {
		return []SumberDaya{}
	}
	result := make([]SumberDaya, len(s.sumberDaya))
	copy(result, s.sumberDaya)
	return result
}

func (s *Store) GetSumberDayaByID(id int) (SumberDaya, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sd := range s.sumberDaya {
		if sd.ID == id {
			return sd, true
		}
	}
	return SumberDaya{}, false
}

// --- Peminjaman Ops ---

func (s *Store) CreatePeminjaman(sdID int, peminjam, tanggal string, slotMulai, slotSelesai, jumlah int) Peminjaman {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := Peminjaman{
		ID:           s.nextPemID,
		SumberDayaID: sdID,
		Peminjam:     peminjam,
		Tanggal:      tanggal,
		SlotMulai:    slotMulai,
		SlotSelesai:  slotSelesai,
		Jumlah:       jumlah,
		Status:       "aktif",
		PaketID:      0,
	}
	s.nextPemID++
	s.peminjaman = append(s.peminjaman, p)
	return p
}

func (s *Store) GetPeminjamanFiltered(sdIDFilter *int, tanggalFilter, statusFilter string) []Peminjaman {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []Peminjaman{}
	for _, p := range s.peminjaman {
		if sdIDFilter != nil && p.SumberDayaID != *sdIDFilter {
			continue
		}
		if tanggalFilter != "" && p.Tanggal != tanggalFilter {
			continue
		}
		if statusFilter != "" && p.Status != statusFilter {
			continue
		}
		result = append(result, p)
	}
	return result
}

func (s *Store) GetPeminjamanByID(id int) (Peminjaman, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.peminjaman {
		if p.ID == id {
			return p, true
		}
	}
	return Peminjaman{}, false
}

// --- Helper Utilities ---

func utf8RuneCount(s string) int {
	return len([]rune(s))
}

func isValidDate(dateStr string) bool {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	year, month, day := t.Date()
	parsedBack := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	return parsedBack == dateStr
}