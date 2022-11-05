package model

import "time"

type Laporan_Keuangann struct {
	Id               uint      `json:"id"`
	Tanggal          time.Time `json:"tanggal"`
	TotalPemasukan   uint      `json:"total_pemasukan"`
	TotalPengeluaran uint      `json:"total_pengeluaran"`
}

type Daftar_Pemesanan struct {
	IdPemesanan     uint               `json:"id_pemesanan"`
	IdCustomer      uint               `json:"id_customer"`
	Status          string             `json:"status"`
	DaftarPemesanan []Produk_Pemesanan `json:"daftar_pemesanan"`
}

type Keranjang_View struct {
	Id         int                     `json:"id"`
	Username   string                  `json:"username"`
	TotalHarga int                     `json:"total_harga"`
	Produk     []Produk_Keranjang_View `json:"produk"`
}

type Produk_Keranjang_View struct {
	IdKeranjang  int         `json:"id_keranjang"`
	IdProduk     int         `json:"id_produk"`
	JumlahProduk int         `json:"jumlah_produk"`
	TotalHarga   int         `json:"total_harga"`
	Produk       Produk_View `json:"produk"`
}

type Produk_View struct {
	Id        int    `json:"id"`
	Nama      string `json:"nama"`
	Gambar    string `json:"gambar"`
	Deskripsi string `json:"deskripsi"`
	Harga     int    `json:"harga"`
}

type Detail_Produk_View struct {
	Id             int             `json:"id"`
	Nama           string          `json:"nama"`
	Gambar         string          `json:"gambar"`
	Deskripsi      string          `json:"deskripsi"`
	Harga          int             `json:"harga"`
	Stok           int             `json:"stok"`
	DaftarFeedback []Feedback_View `json:"daftar_feedback"`
}

type Feedback_View struct {
	Username string   `json:"username"`
	Feedback Feedback `json:"feedback"`
}

type History_View struct {
	Pemesanan []Pemesanan `json:"pemesanan"`
}

type Detail_History_View struct {
	IdPemesanan int                     `json:"id_pemesanan"`
	Tanggal     string                  `json:"tanggal"`
	Status      string                  `json:"status"`
	Alamat      string                  `json:"alamat"`
	Resi        string                  `json:"resi"`
	Produk      []Produk_Pemesanan_View `json:"produk"`
}

type Produk_Pemesanan_View struct {
	Produk       Produk_View `json:"produk"`
	JumlahProduk int         `json:"jumlah_produk"`
	TotalHarga   int         `json:"total_harga"`
}

type Laporan_Bulanan_View struct {
	Bulan   string              `json:"bulan"`
	Tahun   string              `json:"tahun"`
	Laporan []Laporan_Keuangann `json:"laporan"`
}

type Laporan_Tahunan_View struct {
	Tahun   string              `json:"tahun"`
	Laporan []Laporan_Keuangann `json:"laporan"`
}
