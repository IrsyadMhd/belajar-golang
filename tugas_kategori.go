package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Produk makanan dan minuman"},
	{ID: 2, Name: "Elektronik", Description: "Barang-barang elektronik"},
	{ID: 3, Name: "Pakaian", Description: "Baju, celana, dan aksesoris"},
}

func main() {
	// GET /categories -> Ambil semua kategori
	// POST /categories -> Tambah kategori
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategories(w, r)
		} else if r.Method == "POST" {
			createCategory(w, r)
		}
	})

	// GET /categories/{id} -> Ambil detail satu kategori
	// PUT /categories/{id} -> Update kategori
	// DELETE /categories/{id} -> Hapus kategori
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Category API Running",
		})
	})

	fmt.Println("Starting Category API server di localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

// GET /categories -> Ambil semua kategori
func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// POST /categories -> Tambah kategori
func createCategory(w http.ResponseWriter, r *http.Request) {
	var categoryBaru Category
	err := json.NewDecoder(r.Body).Decode(&categoryBaru)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Generate ID baru
	categoryBaru.ID = len(categories) + 1
	categories = append(categories, categoryBaru)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(categoryBaru)
}

// GET /categories/{id} -> Ambil detail satu kategori
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	// URL: /categories/123 -> ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Cari category dengan ID tersebut
	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	// Kalau tidak found
	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// PUT /categories/{id} -> Update kategori
func updateCategory(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Baca data update dari request body
	var updatedCategory Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Loop categories, cari ID, update datanya
	for i := range categories {
		if categories[i].ID == id {
			updatedCategory.ID = id // Pastikan ID tidak berubah
			categories[i] = updatedCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCategory)
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// DELETE /categories/{id} -> Hapus kategori
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Loop categories, cari ID, hapus dari slice
	for i, c := range categories {
		if c.ID == id {
			// Hapus elemen dari slice
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Category berhasil dihapus",
			})
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}
