package model

import (
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

type Customer struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Nama      string    `json:"nama"`
	Umur      uint      `json:"umur"`
	Email     string    `json:"email"`
	Telp      string    `json:"telp"`
	Alamat    string    `json:"alamat"`
	CreatedAt time.Time `json:"created_at"`
}

type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Telp     string `json:"telp"`
}

type Produk struct {
	gorm.Model
	Nama      string `json:"nama"`
	Gambar    string `json:"gambar"`
	Stok      uint   `json:"stok"`
	Deskripsi string `json:"deskripsi"`
	Harga     uint   `json:"harga"`
}

type Keranjang struct {
	gorm.Model
	Username   string `json:"username"`
	TotalHarga uint   `json:"total_harga"`
	Status     bool   `json:"status"`
}

type Produk_Keranjang struct {
	IdKeranjang  uint `json:"id_keranjang"`
	IdProduk     uint `json:"id_produk"`
	JumlahProduk uint `json:"jumlah_produk"`
	TotalHarga   uint `json:"total_harga"`
}

type Produksi struct {
	Id            uint      `json:"id_produk"`
	AdminUsername string    `json:"admin_username"`
	Tanggal       time.Time `json:"tanggal"`
	JumlahBarang  uint      `json:"jumlah_barang"`
	TotalBiaya    uint      `json:"total_biaya"`
}

type Produk_Produksi struct {
	IdProduk     uint `json:"id_produk"`
	IdProduksi   uint `json:"id_produksi"`
	JumlahProduk uint `json:"jumlah_produk"`
	TotalBiaya   uint `json:"total_biaya"`
}

type Pemesanan struct {
	Id               uint      `json:"id"`
	CustomerUsername string    `json:"customer_username"`
	Tanggal          time.Time `json:"tanggal"`
	JumlahBarang     uint      `json:"jumlah_barang"`
	TotalHarga       uint      `json:"total_harga"`
	Status           string    `json:"status"`
	Alamat           string    `json:"alamat"`
	Resi             string    `json:"resi"`
}

type Produk_Pemesanan struct {
	IdProduk     uint `json:"id_produk"`
	IdPemesanan  uint `json:"id_pemesanan"`
	JumlahProduk uint `json:"jumlah_produk"`
	TotalHarga   uint `json:"total_harga"`
}

type Feedback_Pemesanan struct {
	IdFeedback  uint      `json:"id_feedback"`
	IdPemesanan uint      `json:"id_pemesanan"`
	Username    string    `json:"username"`
	Tanggal     time.Time `json:"tanggal"`
}

type Feedback struct {
	gorm.Model
	IdProduk    uint   `json:"id_produk"`
	IsiFeedback string `json:"isi_feedback"`
	Rating      uint   `json:"rating"`
}
