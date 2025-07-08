package auditlog

import (
	"context"
	"errors"
	"testing"

	"github.com/audit-log-service/internal/model"
	"github.com/audit-log-service/internal/repository"
	auditlogRepo "github.com/audit-log-service/internal/repository/auditlog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func Test_GetAuditLogs(t *testing.T) {
	type mockRepo struct {
		in  model.FilterCriteria
		out []model.AuditLog
		err error
	}

	type arg struct {
		criteria model.FilterCriteria
		mockRepo mockRepo
		expOut   []model.AuditLog
		expErr   error
	}

	tcs := map[string]arg{
		"error": {
			criteria: model.FilterCriteria{
				ServiceName: "nonexistent",
			},
			mockRepo: mockRepo{
				in: model.FilterCriteria{
					ServiceName: "nonexistent",
				},
				err: errors.New("some error"),
			},
			expErr: errors.New("some error"),
		},
		"success": {
			criteria: model.FilterCriteria{
				ServiceName: "svc1",
				Username:    "user1",
			},
			mockRepo: mockRepo{
				in: model.FilterCriteria{
					ServiceName: "svc1",
					Username:    "user1",
				},
				out: []model.AuditLog{
					{
						ID:           "log1",
						ServiceName:  "svc1",
						Username:     "user1",
						Action:       "create",
						ResourceType: "type1",
						ResourceID:   "rid1",
						BeforeState:  datatypes.JSON(`{"status":"inactive"}`),
						AfterState:   datatypes.JSON(`{"status":"active"}`),
						Metadata:     datatypes.JSON(`{"ip":"127.0.0.1"}`),
					},
				},
			},
			expOut: []model.AuditLog{
				{
					ID:           "log1",
					ServiceName:  "svc1",
					Username:     "user1",
					Action:       "create",
					ResourceType: "type1",
					ResourceID:   "rid1",
					BeforeState:  datatypes.JSON(`{"status":"inactive"}`),
					AfterState:   datatypes.JSON(`{"status":"active"}`),
					Metadata:     datatypes.JSON(`{"ip":"127.0.0.1"}`),
				},
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Mock
			mockRegistry := repository.NewMockRegistry(t)
			mockAuditRepo := auditlogRepo.NewMockRepository(t)
			mockRegistry.On("GetAuditLog").Return(mockAuditRepo)

			mockAuditRepo.On("GetAuditLogs", mock.Anything, tc.mockRepo.in).
				Return(tc.mockRepo.out, tc.mockRepo.err)

			// When:
			instance := New(mockRegistry)
			logs, err := instance.GetAuditLogs(context.Background(), tc.criteria)

			// Then:
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, logs, tc.expOut)
			}
		})
	}
}
