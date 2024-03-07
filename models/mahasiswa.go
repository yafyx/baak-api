package models

type Mahasiswa struct {
	NPM       string `json:"npm"`
	Nama      string `json:"nama"`
	KelasLama string `json:"kelas_lama"`
	KelasBaru string `json:"kelas_baru"`
}
