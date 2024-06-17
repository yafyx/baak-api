package models

type Jadwal struct {
	Senin  interface{} `json:"Senin"`
	Selasa interface{} `json:"Selasa"`
	Rabu   interface{} `json:"Rabu"`
	Kamis  interface{} `json:"Kamis"`
	Jumat  interface{} `json:"Jumat"`
	Sabtu  interface{} `json:"Sabtu"`
}

type MataKuliah struct {
	Nama  string `json:"nama"`
	Waktu string `json:"waktu"`
	Ruang string `json:"ruang"`
	Dosen string `json:"dosen"`
}
