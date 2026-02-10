package utils

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGenerateID(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		_ = dbMock.Close()
	}()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dbMock,
	}), &gorm.Config{})
	assert.NoError(t, err)

	tests := []struct {
		name        string
		tableName   string
		prefix      string
		mockSetup   func()
		expected    string
		expectError bool
	}{
		{
			name:      "Generate first ID",
			tableName: "cases",
			prefix:    "CS",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("SELECT id FROM \"cases\" WHERE id LIKE \\$1").
					WithArgs("CS-%").
					WillReturnRows(rows)
			},
			expected: "CS-00001",
		},
		{
			name:      "Generate next ID",
			tableName: "cases",
			prefix:    "CS",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("CS-00001").
					AddRow("CS-00005").
					AddRow("CS-00002")
				mock.ExpectQuery("SELECT id FROM \"cases\" WHERE id LIKE \\$1").
					WithArgs("CS-%").
					WillReturnRows(rows)
			},
			expected: "CS-00006",
		},
		{
			name:      "Ignore invalid IDs",
			tableName: "cases",
			prefix:    "CS",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("CS-INVALID").
					AddRow("CS-00010")
				mock.ExpectQuery("SELECT id FROM \"cases\" WHERE id LIKE \\$1").
					WithArgs("CS-%").
					WillReturnRows(rows)
			},
			expected: "CS-00011",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := GenerateID(gormDB, tt.tableName, tt.prefix)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
