package models

type Jadwal struct {
	Senin  []MataKuliah `json:"senin"`
	Selasa []MataKuliah `json:"selasa"`
	Rabu   []MataKuliah `json:"rabu"`
	Kamis  []MataKuliah `json:"kamis"`
	Jumat  []MataKuliah `json:"jumat"`
	Sabtu  []MataKuliah `json:"sabtu"`
}

type MataKuliah struct {
	Nama  string `json:"nama"`
	Waktu string `json:"waktu"`
	Jam   string `json:"jam"`
	Ruang string `json:"ruang"`
	Dosen string `json:"dosen"`
}

type Kegiatan struct {
	Kegiatan string `json:"kegiatan"`
	Tanggal  string `json:"tanggal"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

type Mahasiswa struct {
	NPM       string `json:"npm"`
	Nama      string `json:"nama"`
	KelasLama string `json:"kelas_lama"`
	KelasBaru string `json:"kelas_baru"`
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type UTS struct {
	Nama  string `json:"nama"`
	Waktu string `json:"waktu"`
	Ruang string `json:"ruang"`
	Dosen string `json:"dosen"`
}
