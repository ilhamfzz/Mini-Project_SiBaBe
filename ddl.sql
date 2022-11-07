CREATE TABLE Customer(
	c_username VARCHAR(16) PRIMARY KEY,
	c_password VARCHAR(24) NOT NULL,
	c_nama VARCHAR(64) NOT NULL,
	c_umur INTEGER NOT NULL,
	c_email VARCHAR(64) NOT NULL,
	c_no_telp VARCHAR(16) NOT NULL,
	c_alamat VARCHAR(64) NOT NULL,
    c_created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Pemesanan(
	pmsn_id CHAR(6) PRIMARY KEY,
    k_id CHAR(6) NOT NULL,
	c_username VARCHAR(16) NOT NULL,
	pmsn_tanggal DATE NOT NULL,
	pmsn_jumlah_barang INTEGER NOT NULL,
	pmsn_total_harga INTEGER NOT NULL,
	pmsn_status VARCHAR(16) NOT NULL,
	pmsn_alamat VARCHAR(512) NOT NULL,
	pmsn_kurir VARCHAR(32) NOT NULL,
    pmsn_bukti_pembayaran VARCHAR(64),
    pmsn_divalidasi_oleh VARCHAR(16),
	FOREIGN KEY (c_username) REFERENCES Customer(c_username),
    FOREIGN KEY (k_id) REFERENCES Keranjang(k_id)
);

CREATE TABLE Feedback(
	f_id CHAR(6) PRIMARY KEY,
    p_id CHAR(6) NOT NULL,
    f_isi DATE NOT NULL,
    f_rating INTEGER NOT NULL,
    FOREIGN KEY (p_id) REFERENCES Produk(p_id)
);

CREATE TABLE Feedback_Pemesanan(
	f_id CHAR(6) PRIMARY KEY,
    pmsn_id CHAR(6) NOT NULL,
    fp_username VARCHAR(16) NOT NULL,
    fp_tanggal DATE NOT NULL,
    FOREIGN KEY (f_id) REFERENCES Feedback(f_id),
    FOREIGN KEY (pmsn_id) REFERENCES Pemesanan(pmsn_id)
);

CREATE TABLE Produk(
    p_id CHAR(6) PRIMARY KEY,
    p_nama VARCHAR(128) NOT NULL,
    p_gambar VARCHAR(256) NOT NULL,
    p_stok INTEGER NOT NULL,
    p_deskripsi VARCHAR(512) NOT NULL,
    p_harga INTEGER NOT NULL
);

CREATE TABLE Keranjang(
    k_id CHAR(6) PRIMARY KEY,
    c_username VARCHAR(16) NOT NULL,
    k_total_harga INTEGER NOT NULL,
    k_status VARCHAR(16) NOT NULL,
    FOREIGN KEY (c_username) REFERENCES Customer(c_username)
);

CREATE TABLE Produk_Keranjang(
    k_id CHAR(6) NOT NULL,
    p_id CHAR(6) NOT NULL,
    pk_jumlah_produk INTEGER NOT NULL,
    pk_total_harga INTEGER NOT NULL,
    FOREIGN KEY (p_id) REFERENCES Produk(p_id),
    FOREIGN KEY (k_id) REFERENCES Keranjang(k_id)
);

CREATE TABLE Admin(
    adm_username VARCHAR(16) PRIMARY KEY,
    adm_password VARCHAR(24) NOT NULL,
    adm_nama VARCHAR(128) NOT NULL,
    adm_email VARCHAR(64) NOT NULL,
    adm_no_telp VARCHAR(16) NOT NULL
);

CREATE TABLE Admin_Pemesanan(
    adm_username VARCHAR(16) NOT NULL,
    pmsn_id CHAR(6) NOT NULL,
    ap_tanggal DATE NOT NULL,
    FOREIGN KEY (adm_username) REFERENCES Admin(adm_username),
    FOREIGN KEY (pmsn_id) REFERENCES Pemesanan(pmsn_id)
);

CREATE TABLE Produksi(
    pr_id CHAR(6) PRIMARY KEY,
    adm_username VARCHAR(16) NOT NULL,
    pr_tanggal DATE NOT NULL,
    pr_jumlah_barang INTEGER NOT NULL,
    pr_total_biaya INTEGER NOT NULL,
    FOREIGN KEY (adm_username) REFERENCES Admin(adm_username)
);

CREATE TABLE Produk_Produksi(
    p_id CHAR(6) NOT NULL,
    pr_id CHAR(6) NOT NULL,
    ppr_jumlah INTEGER NOT NULL,
    ppr_biaya INTEGER NOT NULL,
    FOREIGN KEY (p_id) REFERENCES Produk(p_id),
    FOREIGN KEY (pr_id) REFERENCES Produksi(pr_id)
);