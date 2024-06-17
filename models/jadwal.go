package models

type Jadwal struct {
	Jadwal []Hari `json:"jadwal"`
}

type Hari struct {
	Hari       string       `json:"hari"`
	Kelas      string       `json:"kelas"`
	MataKuliah []MataKuliah `json:"matakuliah"`
}

type MataKuliah struct {
	Nama  string `json:"nama"`
	Waktu string `json:"waktu"`
	Ruang string `json:"ruang"`
	Dosen string `json:"dosen"`
}
