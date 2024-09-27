package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	dns := "host=localhost user=postgres password=password dbname=my_ecommerce port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		log.Fatal("Unable to connect to DB")
	}
	return db
}


// In case database causes a problem for foreign key relation between products and user
// ALTER TABLE products DROP CONSTRAINT fk_users_products;
// ALTER TABLE products ADD CONSTRAINT fk_users_products FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE;