package main

import (
	"Mini-Project_SiBaBe/config"
	"Mini-Project_SiBaBe/model"
	"Mini-Project_SiBaBe/route"
	"Mini-Project_SiBaBe/service"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		db          = config.InitDatabase()
		PORT        = os.Getenv("PORT")
		customerSvc = service.NewCustomerService(db)
		adminSvc    = service.NewAdminService(db)
	)
	model.DB.AutoMigrate(
		&model.Customer{}, &model.Admin{}, &model.Produk{},
		&model.Keranjang{}, &model.Produk_Keranjang{},
		&model.Produksi{}, &model.Pemesanan{},
		&model.Admin_Pemesanan{}, &model.Feedback_Pemesanan{},
		&model.Feedback{}, &model.Laporan_Keuangan{},
	)
	app := route.New(customerSvc, adminSvc)

	app.Start(":" + PORT)
}
