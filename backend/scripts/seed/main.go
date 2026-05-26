package main

import (
	"log"
	"matcha-comic-backend/config"
	"matcha-comic-backend/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load(".env")
	db := config.ConnectDB()

	// Auto Migrate before seeding
	db.AutoMigrate(&models.User{}, &models.Story{}, &models.Chapter{}, &models.Genre{}, &models.Category{})

	// Create admin user
	db.Unscoped().Where("username = ?", "superadmin").Delete(&models.User{})
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	admin := models.User{
		Username: "superadmin",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Println("Error creating admin user:", err)
	} else {
		log.Println("Admin user created/reset: superadmin / password")
	}

	// Create some mock genres
	genres := []models.Genre{
		{Name: "Action", Slug: "action"},
		{Name: "Romance", Slug: "romance"},
		{Name: "Fantasy", Slug: "fantasy"},
	}

	for _, g := range genres {
		db.FirstOrCreate(&g, models.Genre{Slug: g.Slug})
	}

	// Seed website navigation Categories (Danh mục)
	log.Println("Seeding website categories...")
	
	// 1. Parent Categories
	comicCat := models.Category{Name: "Truyện Tranh", Slug: "truyen-tranh", Path: "/truyen-tranh", Icon: "🎨", Order: 1}
	db.FirstOrCreate(&comicCat, models.Category{Slug: "truyen-tranh"})

	novelCat := models.Category{Name: "Truyện Chữ", Slug: "truyen-chu", Path: "/truyen-chu", Icon: "📖", Order: 2}
	db.FirstOrCreate(&novelCat, models.Category{Slug: "truyen-chu"})

	genreCat := models.Category{Name: "Thể Loại", Slug: "the-loai", Path: "/the-loai", Icon: "🏷️", Order: 3}
	db.FirstOrCreate(&genreCat, models.Category{Slug: "the-loai"})

	// 2. Child Categories under Truyện Tranh
	comicChildren := []models.Category{
		{Name: "Tất cả truyện", Slug: "comic-all", Path: "/truyen-tranh", Icon: "📚", ParentID: &comicCat.ID, Order: 1},
		{Name: "Mới cập nhật", Slug: "comic-latest", Path: "/truyen-tranh?sort=latest", Icon: "🆕", ParentID: &comicCat.ID, Order: 2},
		{Name: "Top thịnh hành", Slug: "comic-hot", Path: "/truyen-tranh?sort=views", Icon: "🔥", ParentID: &comicCat.ID, Order: 3},
		{Name: "Đã hoàn thành", Slug: "comic-completed", Path: "/truyen-tranh?status=completed", Icon: "✨", ParentID: &comicCat.ID, Order: 4},
	}
	for _, child := range comicChildren {
		db.FirstOrCreate(&child, models.Category{Slug: child.Slug})
	}

	// 3. Child Categories under Truyện Chữ
	novelChildren := []models.Category{
		{Name: "Tất cả tiểu thuyết", Slug: "novel-all", Path: "/truyen-chu", Icon: "📖", ParentID: &novelCat.ID, Order: 1},
		{Name: "Mới cập nhật", Slug: "novel-latest", Path: "/truyen-chu?sort=latest", Icon: "🆕", ParentID: &novelCat.ID, Order: 2},
		{Name: "Top thịnh hành", Slug: "novel-hot", Path: "/truyen-chu?sort=views", Icon: "🔥", ParentID: &novelCat.ID, Order: 3},
		{Name: "Đã hoàn thành", Slug: "novel-completed", Path: "/truyen-chu?status=completed", Icon: "✨", ParentID: &novelCat.ID, Order: 4},
	}
	for _, child := range novelChildren {
		db.FirstOrCreate(&child, models.Category{Slug: child.Slug})
	}

	log.Println("Seeding completed!")
}
