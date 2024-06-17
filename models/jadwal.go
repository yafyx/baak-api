// models/models.go
package models

type Jadwal struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	Jadwal []Hari `json:"jadwal"`
}

type Hari struct {
	Hari       string       `json:"hari"`
	Kelas      string       `json:"kelas"`
	MataKuliah []MataKuliah `json:"mata_kuliah"`
}

type MataKuliah struct {
	Nama  string `json:"nama"`
	Waktu string `json:"waktu"`
	Ruang string `json:"ruang"`
	Dosen string `json:"dosen"`
}
