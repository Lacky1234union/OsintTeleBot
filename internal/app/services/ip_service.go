package services

import (
	"context"
	"fmt"
	"net"
	"time"

	"OsintTeleBot/internal/app/repositories"
)

// IPService defines the interface for IP-related business logic
type IPService interface {
	FindByIP(ctx context.Context, ip string) (*IPResponse, error)
}

// IPResponse represents the response structure for IP information
type IPResponse struct {
	IP        string    `json:"ip"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	ISP       string    `json:"isp"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ipService implements IPService interface
type ipService struct {
	repo repositories.IPRepository
}

// NewIPService creates a new instance of IPService
func NewIPService(repo repositories.IPRepository) IPService {
	return &ipService{repo: repo}
}

// FindByIP retrieves and processes IP information
func (s *ipService) FindByIP(ctx context.Context, ip string) (*IPResponse, error) {
	// Validate IP format
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid IP address format: %s", ip)
	}

	// Check if IP exists in database
	info, err := s.repo.FindByIP(ctx, ip)
	if err != nil {
		return nil, fmt.Errorf("failed to find IP info: %w", err)
	}

	// If IP not found in database, create new entry
	if info == nil {
		// Here you would typically call an external IP geolocation service
		// For now, we'll create a placeholder entry
		info = &repositories.IPInfo{
			IP:      ip,
			Country: "Unknown",
			City:    "Unknown",
			ISP:     "Unknown",
		}

		if err := s.repo.Create(ctx, info); err != nil {
			return nil, fmt.Errorf("failed to create IP info: %w", err)
		}
	}

	// Convert repository model to response model
	response := &IPResponse{
		IP:        info.IP,
		Country:   info.Country,
		City:      info.City,
		ISP:       info.ISP,
		UpdatedAt: info.UpdatedAt,
	}

	return response, nil
}
