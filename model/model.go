package model

import (
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

type Customer struct {
	Username  string    `gorm:"primary_key;not null;unique" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Nama      string    `gorm:"not null" json:"nama"`
	Umur      uint      `json:"umur"`
	Email     string    `gorm:"not null;unique" json:"email"`
	Telp      string    `gorm:"not null;unique" json:"telp"`
	Alamat    string    `gorm:"not null" json:"alamat"`
	CreatedAt time.Time `json:"created_at"`
}

type Admin struct {
	Username string `gorm:"primary_key;not null;unique" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Nama     string `gorm:"not null" json:"nama"`
	Email    string `gorm:"not null;unique" json:"email"`
	Telp     string `gorm:"not null;unique" json:"telp"`
}

type Produk struct {
	gorm.Model
	Nama      string `gorm:"not null;unique" json:"nama"`
	Gambar    string `gorm:"not null" json:"gambar"`
	Stok      uint   `json:"stok"`
	Deskripsi string `json:"deskripsi"`
	Harga     uint   `gorm:"not null" json:"harga"`
}

type Keranjang struct {
	gorm.Model
	Username   string `gorm:"not null" json:"username"`
	TotalHarga uint   `json:"total_harga"`
	Status     string `gorm:"not null" json:"status"`
}

type Produk_Keranjang struct {
	IdKeranjang  uint `gorm:"not null" json:"id_keranjang"`
	IdProduk     uint `gorm:"not null" json:"id_produk"`
	JumlahProduk uint `gorm:"not null" json:"jumlah_produk"`
	TotalHarga   uint `gorm:"not null" json:"total_harga"`
}

type Produksi struct {
	gorm.Model
	Date          string `gorm:"not null" json:"date"`
	AdminUsername string `gorm:"not null" json:"admin_username"`
	NamaBarang    string `gorm:"not null" json:"nama_barang"`
	TotalBiaya    uint   `gorm:"not null" json:"total_biaya"`
	Gambar        string `gorm:"not null" json:"gambar"`
}

type Pemesanan struct {
	gorm.Model
	IdKeranjang      uint   `gorm:"not null" json:"id_keranjang"`
	CustomerUsername string `gorm:"not null" json:"customer_username"`
	JumlahBarang     uint   `gorm:"not null" json:"jumlah_barang"`
	TotalHarga       uint   `gorm:"not null" json:"total_harga"`
	Status           string `gorm:"not null" json:"status"`
	Alamat           string `gorm:"not null" json:"alamat"`
	Kurir            string `gorm:"not null" json:"kurir"`
	BuktiPembayaran  string `json:"bukti_pembayaran"`
	DiValidasiOleh   string `json:"di_validasi_oleh"`
}

type Admin_Pemesanan struct {
	IdPemesanan         uint      `gorm:"not null" json:"id_pemesanan"`
	UsernameAdmin       string    `gorm:"not null" json:"username_admin"`
	UpdateStatusOrderTo string    `gorm:"not null" json:"update_status_order_to"`
	TanggalValidasi     time.Time `json:"tanggal_validasi"`
}

type Feedback_Pemesanan struct {
	IdFeedback  uint      `gorm:"not null" json:"id_feedback"`
	IdPemesanan uint      `gorm:"not null" json:"id_pemesanan"`
	Username    string    `gorm:"not null" json:"username"`
	Tanggal     time.Time `json:"tanggal"`
}

type Feedback struct {
	gorm.Model
	IdProduk    uint   `gorm:"not null" json:"id_produk"`
	IsiFeedback string `json:"isi_feedback"`
	Rating      uint   `json:"rating"`
}

type Laporan_Keuangan struct {
	Tanggal          string `gorm:"primary_key;not null" json:"tanggal"`
	TotalPemasukan   uint   `json:"total_pemasukan"`
	TotalPengeluaran uint   `json:"total_pengeluaran"`
}
