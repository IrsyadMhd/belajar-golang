package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetDailyReport() (*models.DailyReport, error) {
	return s.repo.GetDailyReport()
}

func (s *ReportService) GetReportByDateRange(startDate, endDate string) (*models.DailyReport, error) {
	return s.repo.GetReportByDateRange(startDate, endDate)
}
