package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// IPInfo represents the IP information stored in the database
type IPInfo struct {
	ID        int64
	IP        string
	Country   string
	City      string
	ISP       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IPRepository defines the interface for IP-related database operations
type IPRepository interface {
	FindByIP(ctx context.Context, ip string) (*IPInfo, error)
	Create(ctx context.Context, info *IPInfo) error
	Update(ctx context.Context, info *IPInfo) error
}

// ipRepository implements IPRepository interface
type ipRepository struct {
	db *sql.DB
}

// NewIPRepository creates a new instance of IPRepository
func NewIPRepository(db *sql.DB) IPRepository {
	return &ipRepository{db: db}
}

// FindByIP retrieves IP information from the database
func (r *ipRepository) FindByIP(ctx context.Context, ip string) (*IPInfo, error) {
	query := `
		SELECT id, ip, country, city, isp, created_at, updated_at
		FROM ip_info
		WHERE ip = $1
	`

	var info IPInfo
	err := r.db.QueryRowContext(ctx, query, ip).Scan(
		&info.ID,
		&info.IP,
		&info.Country,
		&info.City,
		&info.ISP,
		&info.CreatedAt,
		&info.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find IP info: %w", err)
	}

	return &info, nil
}

// Create stores new IP information in the database
func (r *ipRepository) Create(ctx context.Context, info *IPInfo) error {
	query := `
		INSERT INTO ip_info (ip, country, city, isp, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	info.CreatedAt = now
	info.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		info.IP,
		info.Country,
		info.City,
		info.ISP,
		info.CreatedAt,
		info.UpdatedAt,
	).Scan(&info.ID)

	if err != nil {
		return fmt.Errorf("failed to create IP info: %w", err)
	}

	return nil
}

// Update modifies existing IP information in the database
func (r *ipRepository) Update(ctx context.Context, info *IPInfo) error {
	query := `
		UPDATE ip_info
		SET country = $1, city = $2, isp = $3, updated_at = $4
		WHERE ip = $5
	`

	info.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		info.Country,
		info.City,
		info.ISP,
		info.UpdatedAt,
		info.IP,
	)

	if err != nil {
		return fmt.Errorf("failed to update IP info: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no IP info found for IP: %s", info.IP)
	}

	return nil
}
