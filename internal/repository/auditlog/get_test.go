package auditlog

import (
	"context"
	"testing"
	"time"

	"github.com/audit-log-service/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupTestDBForGet(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Dialector{
		DSN:        "file::memory:?cache=shared",
		DriverName: "sqlite",
	}, &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&model.AuditLog{})
	require.NoError(t, err)
	return db
}

func TestRepository_GetAuditLogs(t *testing.T) {
	db := setupTestDBForGet(t)
	repo := New(db)
	ctx := context.Background()

	// Insert sample logs
	logs := []model.AuditLog{
		{
			ID:           uuid.New().String(),
			ServiceName:  "svc1",
			Username:     "user1",
			Action:       "create",
			ResourceType: "type1",
			ResourceID:   "rid1",
			Timestamp:    time.Now(),
			BeforeState:  datatypes.JSON(`{"status":"inactive"}`),
			AfterState:   datatypes.JSON(`{"status":"active"}`),
			Metadata:     datatypes.JSON(`{"ip":"127.0.0.1"}`),
		},
		{
			ID:           uuid.New().String(),
			ServiceName:  "svc2",
			Username:     "user2",
			Action:       "update",
			ResourceType: "type2",
			ResourceID:   "rid2",
			Timestamp:    time.Now(),
			BeforeState:  datatypes.JSON(`{"status":"active"}`),
			AfterState:   datatypes.JSON(`{"status":"inactive"}`),
			Metadata:     datatypes.JSON(`{"ip":"127.0.0.2"}`),
		},
	}
	require.NoError(t, db.Create(&logs).Error)

	tcs := map[string]struct {
		criteria   model.FilterCriteria
		expCount   int
		expService string
	}{
		"all": {
			criteria: model.FilterCriteria{},
			expCount: 2,
		},
		"filter by ServiceName": {
			criteria:   model.FilterCriteria{ServiceName: "svc1"},
			expCount:   1,
			expService: "svc1",
		},
		"filter by Username": {
			criteria:   model.FilterCriteria{Username: "user2"},
			expCount:   1,
			expService: "svc2",
		},
		"filter by Action": {
			criteria:   model.FilterCriteria{Action: "create"},
			expCount:   1,
			expService: "svc1",
		},
		"filter by ResourceType": {
			criteria:   model.FilterCriteria{ResourceType: "type2"},
			expCount:   1,
			expService: "svc2",
		},
		"filter by ResourceID": {
			criteria:   model.FilterCriteria{ResourceID: "rid1"},
			expCount:   1,
			expService: "svc1",
		},
		"no match": {
			criteria: model.FilterCriteria{ServiceName: "notfound"},
			expCount: 0,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			result, err := repo.GetAuditLogs(ctx, tc.criteria)
			require.NoError(t, err)
			require.Equal(t, tc.expCount, len(result))
			if tc.expCount == 1 {
				require.Equal(t, tc.expService, result[0].ServiceName)
			}
		})
	}
}
