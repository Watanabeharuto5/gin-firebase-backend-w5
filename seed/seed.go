package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/Watanabeharuto5/gin-firebase-backend-w5/config"
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	// Init database
	config.InitDatabase()

	// Seed data K-Pop Album
	products := []models.Product{
		{
			Name:        "TREASURE - 2ND FULL ALBUM [REBOOT] PHOTOBOOK VER. (Set)",
			Price:       650000,
			Category:    "Album",
			Stock:       8,
			Description: "Album full TREASURE versi photobook set dengan berbagai collectible item dan konsep comeback terbaru.",
			ImageURL:    "https://i.ibb.co.com/zWvcRg9S/reboot.png",
		},
		{
			Name:        "BTS - 'ARIRANG' (Set) + 'ARIRANG' (Living Legend Ver.) + 'ARIRANG' (Weverse Albums ver.) Set",
			Price:       900000,
			Category:    "Album",
			Stock:       12,
			Description: "Set album BTS dengan beberapa versi berbeda, termasuk edisi Weverse dan Living Legend.",
			ImageURL:    "https://i.ibb.co.com/Y4sBGRJm/arirang2.png",
		},
		{
			Name:        "NCT DREAM - [DREAM( )SCAPE] (QR Ver.)",
			Price:       140000,
			Category:    "Album",
			Stock:       10,
			Description: "Album versi QR dengan cover, QR card, 8 image card, sticker, dan photocard random (1 dari 7 versi). Detail tambahan akan diperbarui.",
			ImageURL:    "https://i.ibb.co.com/x8qWVf1t/nctdream.jpg",
		},
		{
			Name:        "NCT 127 The 1st Album [TASTE] (Full Spread Ver.) Random",
			Price:       200000,
			Category:    "Album",
			Stock:       7,
			Description: "Termasuk photobook 96 halaman, CD, poster lipat, sticker, photofilm random, dan photocard random (1 dari 3 versi).",
			ImageURL:    "https://i.ibb.co.com/GQLdVtWv/nct127.png",
		},
		{
			Name:        "WayV - Winter Special Album [白色定格 (Eternal White)] (Ornament Ver.) Random",
			Price:       350000,
			Category:    "Album",
			Stock:       6,
			Description: "Album spesial dengan ornament set, QR card, image card, slide film random, sticker, dan photocard random.",
			ImageURL:    "https://i.ibb.co.com/1xFV5Kn/wayv.png",
		},
		{
			Name:        "SEVENTEEN - 에스쿱스X민규 1st Mini Album ‘HYPE VIBES’ (Set)",
			Price:       400000,
			Category:    "Album",
			Stock:       9,
			Description: "Mini album unit SEVENTEEN (S.Coups & Mingyu) dengan konsep powerful dan collectible set.",
			ImageURL:    "https://i.ibb.co.com/CpGw1VZD/seventeen.png",
		},
		{
			Name:        "BLACKPINK - BLACKPINK 3rd MINI ALBUM [DEADLINE] BLACK Ver.",
			Price:       250000,
			Category:    "Album",
			Stock:       5,
			Description: "Termasuk photobook 72 halaman, CD, selfie photocard, sticker set, dan poster lipat grup.",
			ImageURL:    "https://i.ibb.co.com/BV42zmdr/blackpink.png",
		},
		{
			Name:        "BABYMONSTER - BABYMONSTER 2nd MINI ALBUM [WE GO UP] GO Ver.",
			Price:       220000,
			Category:    "Album",
			Stock:       12,
			Description: "Mini album dengan photobook 72 halaman, CD, poster, 6 photocard, bookmark, dan logo sticker.",
			ImageURL:    "https://i.ibb.co.com/BHG9YvT3/baemoon.png",
		},
		{
			Name:        "TXT - 7TH YEAR: A Moment of Stillness in the Thorn (Set)",
			Price:       700000,
			Category:    "Album",
			Stock:       12,
			Description: "Set album TXT dengan photobook, poster, lyric book, postcard random, sticker, photocard, dan CD.",
			ImageURL:    "https://i.ibb.co.com/wFJDNSs2/txt.png",
		},
		{
			Name:        "CORTIS - [STUDIO CHOOM GIFT] The 2nd EP [GREENGREEN] (Set)",
			Price:       650000,
			Category:    "Album",
			Stock:       6,
			Description: "EP kedua CORTIS edisi set dengan benefit Studio Choom gift, berisi berbagai collectible item dan konsep eksklusif.",
			ImageURL:    "https://i.ibb.co.com/v6wJP6ZZ/coer.png",
},
	}

	// Insert ke database
	for _, p := range products {
		if err := config.DB.Create(&p).Error; err != nil {
			log.Printf("Gagal insert produk: %v", err)
		}
	}

	log.Printf("Seed berhasil: %d produk ditambahkan", len(products))
}