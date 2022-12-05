package model

type Daftar_Pemesanan struct {
	IdPemesanan      uint          `json:"id_pemesanan"`
	IdKeranjang      uint          `json:"id_keranjang"`
	Status           string        `json:"status"`
	Daftar_Pemesanan []Produk_View `json:"daftar_pemesanan"`
}

type Checkout_Binding struct {
	Kurir  string `json:"kurir"`
	Alamat string `json:"alamat"`
}

type Payment_Binding struct {
	BuktiPembayaran string `json:"bukti_pembayaran"`
}

type Checkout struct {
	Alamat      string         `json:"alamat"`
	Kurir       string         `json:"kurir"`
	OngkosKirim uint           `json:"ongkos_kirim"`
	Total_Harga uint           `json:"total_harga"`
	Keranjang   Keranjang_View `json:"keranjang"`
}

type Keranjang_View struct {
	Id           uint                    `json:"id"`
	Username     string                  `json:"username"`
	JumlahBarang uint                    `json:"jumlah_barang"`
	TotalHarga   uint                    `json:"total_harga"`
	Produk       []Produk_Keranjang_View `json:"produk"`
}

type Produk_Keranjang_View struct {
	IdKeranjang  uint        `json:"id_keranjang"`
	IdProduk     uint        `json:"id_produk"`
	JumlahProduk uint        `json:"jumlah_produk"`
	TotalHarga   uint        `json:"total_harga"`
	Produk       Produk_View `json:"produk"`
}

type Produk_View struct {
	Id        uint   `json:"id"`
	Nama      string `json:"nama"`
	Gambar    string `json:"gambar"`
	Deskripsi string `json:"deskripsi"`
	Harga     uint   `json:"harga"`
}

type Detail_Produk_View struct {
	Id             uint            `json:"id"`
	Nama           string          `json:"nama"`
	Gambar         string          `json:"gambar"`
	Deskripsi      string          `json:"deskripsi"`
	Harga          uint            `json:"harga"`
	Stok           uint            `json:"stok"`
	DaftarFeedback []Feedback_View `json:"daftar_feedback"`
}

type Feedback_View struct {
	Username string   `json:"username"`
	IdProduk uint     `json:"id_produk"`
	Feedback Feedback `json:"feedback"`
}

type History_View struct {
	Pemesanan []Pemesanan `json:"pemesanan"`
}

type Detail_History_View struct {
	IdPemesanan uint                    `json:"id_pemesanan"`
	IdKeranjang uint                    `json:"id_keranjang"`
	Status      string                  `json:"status"`
	Alamat      string                  `json:"alamat"`
	Kurir       string                  `json:"kurir"`
	Produk      []Produk_Pemesanan_View `json:"produk"`
}

type Produk_Pemesanan_View struct {
	JumlahProduk uint        `json:"jumlah_produk"`
	TotalHarga   uint        `json:"total_harga"`
	Produk       Produk_View `json:"produk"`
}

type Laporan_Bulanan_View struct {
	Bulan   string              `json:"bulan"`
	Tahun   int                 `json:"tahun"`
	Laporan []Laporan_Keuangann `json:"laporan"`
}

type Update_Order_Status_Binding struct {
	Status string `json:"status"`
}

type Produksi_Binding struct {
	NamaProduk string `json:"nama_produk"`
	TotalBiaya uint   `json:"total_biaya"`
}

type Produk_View_Integrated struct {
	Id          uint          `json:"id"`
	Name        string        `json:"name"`
	Price       uint          `json:"price"`
	Stock       uint          `json:"stock"`
	Image       string        `json:"image"`
	Description string        `json:"description"`
	Reviews     []Review_View `json:"reviews"`
}

type Review_View struct {
	Username string `json:"username"`
	Feedback string `json:"feedback"`
	Rating   uint   `json:"rating"`
}
