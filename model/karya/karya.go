package karya

type Karya struct {
	KaryaId     int    `json:"id"`
	Judul       string `json:"judul"`
	PelukisId   int    `json:"pelukis_id"`
	TahunDibuat int    `json:"tahun_dibuat"`
	Media       string `json:"media"`
}
