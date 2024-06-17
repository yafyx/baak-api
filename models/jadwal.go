package models

type Jadwal struct {
	Senin  interface{} `json:"senin"`
	Selasa interface{} `json:"selasa"`
	Rabu   interface{} `json:"rabu"`
	Kamis  interface{} `json:"kamis"`
	Jumat  interface{} `json:"jumat"`
	Sabtu  interface{} `json:"sabtu"`
}

type MataKuliah struct {
	Nama  string `json:"nama"`
	Waktu string `json:"waktu"`
	Ruang string `json:"ruang"`
	Dosen string `json:"dosen"`
}
