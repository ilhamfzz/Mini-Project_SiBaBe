CREATE TABLE Customer(
    username VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    nama VARCHAR(50) NOT NULL,
    umur INT NOT NULL,
    email VARCHAR(50) NOT NULL,
    telp VARCHAR(50) NOT NULL,
    alamat VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (username)
);

CREATE TABLE Admin(
    username VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    nama VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    telp VARCHAR(50) NOT NULL,
    PRIMARY KEY (username)
);

CREATE TABLE Produk(
    id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    nama VARCHAR(50) NOT NULL,
    gambar VARCHAR(50) NOT NULL,
    stok INT NOT NULL,
    deskripsi VARCHAR(50) NOT NULL,
    harga INT NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE Keranjang(
    id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    username VARCHAR(50) NOT NULL,
    total_harga INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (username) REFERENCES Customer(username)
);

CREATE TABLE Produk_Keranjang(
    id_keranjang INT NOT NULL,
    id_produk INT NOT NULL,
    jumlah_produk INT NOT NULL,
    total_harga INT NOT NULL,
    PRIMARY KEY (id_keranjang, id_produk),
    FOREIGN KEY (id_keranjang) REFERENCES Keranjang(id),
    FOREIGN KEY (id_produk) REFERENCES Produk(id)
);

CREATE TABLE Produksi(
    id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    admin_username VARCHAR(50) NOT NULL,
    total_biaya INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (admin_username) REFERENCES Admin(username)
);

CREATE TABLE Pemesanan(
    id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    id_keranjang INT NOT NULL,
    customer_username VARCHAR(50) NOT NULL,
    jumlah_barang INT NOT NULL,
    total_harga INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    alamat VARCHAR(50) NOT NULL,
    kurir VARCHAR(50) NOT NULL,
    bukti_pembayaran VARCHAR(50) NOT NULL,
    di_validasi_oleh VARCHAR(50) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (id_keranjang) REFERENCES Keranjang(id),
    FOREIGN KEY (customer_username) REFERENCES Customer(username)
);

CREATE TABLE Admin_Pemesanan(
    id_pemesanan INT NOT NULL,
    username_admin VARCHAR(50) NOT NULL,
    update_status_order_to VARCHAR(50) NOT NULL,
    tanggal_validasi TIMESTAMP NOT NULL,
    PRIMARY KEY (id_pemesanan, username_admin),
    FOREIGN KEY (id_pemesanan) REFERENCES Pemesanan(id),
    FOREIGN KEY (username_admin) REFERENCES Admin(username)
);

CREATE TABLE Feedback(
    id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    id_produk INT NOT NULL,
    isi_feedback VARCHAR(50) NOT NULL,
    rating INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (id_produk) REFERENCES Produk(id)
);

CREATE TABLE Feedback_Pemesanan(
    id_feedback INT NOT NULL,
    id_pemesanan INT NOT NULL,
    username VARCHAR(50) NOT NULL,
    tanggal TIMESTAMP NOT NULL,
    PRIMARY KEY (id_feedback, id_pemesanan),
    FOREIGN KEY (id_feedback) REFERENCES Feedback(id),
    FOREIGN KEY (id_pemesanan) REFERENCES Pemesanan(id)
);

CREATE TABLE Laporan_Keuangan(
    tanggal VARCHAR(50) NOT NULL,
    total_pemasukan INT NOT NULL,
    total_pengeluaran INT NOT NULL,
    PRIMARY KEY (tanggal),
    FOREIGN KEY (tanggal) REFERENCES Produksi(created_at),
    FOREIGN KEY (tanggal) REFERENCES Admin_Pemesanan(tanggal_validasi)
);

-- Initial Assets Data
INSERT INTO admins VALUES ('admin', 'admin', 'admin', 'admintest@gmail.com', '081234567890');
-- Product Data from Backend Inserting data from Postman


-- Insert Dummy data for testing chart monthly report

INSERT INTO laporan_keuangans (tanggal, total_pemasukan, total_pengeluaran)
    VALUES
        ('2022-12-15', 345000, 72000),
