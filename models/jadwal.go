package models

type Jadwal struct {
	Search string `json:"search"`
	Jadwal []Hari `json:"jadwal"`
}

type Hari struct {
	Hari string `json:"hari"`
	Data []Data `json:"data"`
}

type Data struct {
	Kelas      string `json:"kelas"`
	MataKuliah string `json:"mata_kuliah"`
	Waktu      string `json:"waktu"`
	Ruang      string `json:"ruang"`
	Dosen      string `json:"dosen"`
}
