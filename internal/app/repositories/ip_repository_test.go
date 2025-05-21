package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Connect to test database
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/appdb?sslmode=disable")
	require.NoError(t, err)

	// Create test table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ip_info (
			id BIGSERIAL PRIMARY KEY,
			ip VARCHAR(45) NOT NULL UNIQUE,
			country VARCHAR(100),
			city VARCHAR(100),
			isp VARCHAR(100),
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	require.NoError(t, err)

	// Load test data
	_, err = db.Exec(`
		TRUNCATE TABLE ip_info;
		INSERT INTO ip_info (ip, country, city, isp, created_at, updated_at)
		VALUES 
			('8.8.8.8', 'United States', 'Mountain View', 'Google', NOW(), NOW()),
			('1.1.1.1', 'Australia', 'Sydney', 'Cloudflare', NOW(), NOW()),
			('208.67.222.222', 'United States', 'San Francisco', 'OpenDNS', NOW(), NOW())
	`)
	require.NoError(t, err)

	return db
}

func TestIPRepository_FindByIP(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewIPRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		ip      string
		want    *IPInfo
		wantErr bool
	}{
		{
			name: "Find existing IP",
			ip:   "8.8.8.8",
			want: &IPInfo{
				IP:      "8.8.8.8",
				Country: "United States",
				City:    "Mountain View",
				ISP:     "Google",
			},
			wantErr: false,
		},
		{
			name:    "Find non-existing IP",
			ip:      "192.168.1.1",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Invalid IP format",
			ip:      "invalid-ip",
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByIP(ctx, tt.ip)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tt.want == nil {
				assert.Nil(t, got)
				return
			}
			assert.Equal(t, tt.want.IP, got.IP)
			assert.Equal(t, tt.want.Country, got.Country)
			assert.Equal(t, tt.want.City, got.City)
			assert.Equal(t, tt.want.ISP, got.ISP)
		})
	}
}

func TestIPRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewIPRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		info    *IPInfo
		wantErr bool
	}{
		{
			name: "Create new IP info",
			info: &IPInfo{
				IP:      "9.9.9.9",
				Country: "United States",
				City:    "Reston",
				ISP:     "Quad9",
			},
			wantErr: false,
		},
		{
			name: "Create duplicate IP",
			info: &IPInfo{
				IP:      "8.8.8.8",
				Country: "United States",
				City:    "Mountain View",
				ISP:     "Google",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.info)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotZero(t, tt.info.ID)
			assert.NotZero(t, tt.info.CreatedAt)
			assert.NotZero(t, tt.info.UpdatedAt)
		})
	}
}

func TestIPRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewIPRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		info    *IPInfo
		wantErr bool
	}{
		{
			name: "Update existing IP",
			info: &IPInfo{
				IP:      "8.8.8.8",
				Country: "United States",
				City:    "New York",
				ISP:     "Google",
			},
			wantErr: false,
		},
		{
			name: "Update non-existing IP",
			info: &IPInfo{
				IP:      "192.168.1.1",
				Country: "Unknown",
				City:    "Unknown",
				ISP:     "Unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.info)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify update
			got, err := repo.FindByIP(ctx, tt.info.IP)
			assert.NoError(t, err)
			assert.Equal(t, tt.info.Country, got.Country)
			assert.Equal(t, tt.info.City, got.City)
			assert.Equal(t, tt.info.ISP, got.ISP)
		})
	}
}
