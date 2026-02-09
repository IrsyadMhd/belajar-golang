package repositories

import (
	"database/sql"
	"kasir-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetDailyReport() (*models.DailyReport, error) {
	today := time.Now().Format("2006-01-02")

	// Get total revenue and total transactions for today
	var totalRevenue sql.NullInt64
	var totalTransaksi int
	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) = $1
	`, today).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	// Get best selling product for today
	var bestProduct *models.BestSellingProduct
	var productName sql.NullString
	var qtyTerjual sql.NullInt64
	err = repo.db.QueryRow(`
		SELECT p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) = $1
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, today).Scan(&productName, &qtyTerjual)
	
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if productName.Valid && qtyTerjual.Valid {
		bestProduct = &models.BestSellingProduct{
			Nama:       productName.String,
			QtyTerjual: int(qtyTerjual.Int64),
		}
	}

	return &models.DailyReport{
		TotalRevenue:   int(totalRevenue.Int64),
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: bestProduct,
	}, nil
}

func (repo *ReportRepository) GetReportByDateRange(startDate, endDate string) (*models.DailyReport, error) {
	// Get total revenue and total transactions for date range
	var totalRevenue sql.NullInt64
	var totalTransaksi int
	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
	`, startDate, endDate).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	// Get best selling product for date range
	var bestProduct *models.BestSellingProduct
	var productName sql.NullString
	var qtyTerjual sql.NullInt64
	err = repo.db.QueryRow(`
		SELECT p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) >= $1 AND DATE(t.created_at) <= $2
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&productName, &qtyTerjual)
	
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if productName.Valid && qtyTerjual.Valid {
		bestProduct = &models.BestSellingProduct{
			Nama:       productName.String,
			QtyTerjual: int(qtyTerjual.Int64),
		}
	}

	return &models.DailyReport{
		TotalRevenue:   int(totalRevenue.Int64),
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: bestProduct,
	}, nil
}
