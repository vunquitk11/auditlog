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

func Test_HandleAuditLogMessage(t *testing.T) {
	type mockRepo struct {
		in  model.AuditLog
		err error
	}

	type arg struct {
		givenIn  model.AuditLog
		mockRepo mockRepo
		expErr   error
	}

	tcs := map[string]arg{
		"error": {
			givenIn: model.AuditLog{
				ServiceName:  "svc1",
				Username:     "user1",
				Action:       "create",
				ResourceType: "type1",
				ResourceID:   "rid1",
				BeforeState:  datatypes.JSON(`{"status":"inactive"}`),
				AfterState:   datatypes.JSON(`{"status":"active"}`),
				Metadata:     datatypes.JSON(`{"ip":"127.0.0.1"}`),
			},
			mockRepo: mockRepo{
				in: model.AuditLog{
					ServiceName:  "svc1",
					Username:     "user1",
					Action:       "create",
					ResourceType: "type1",
					ResourceID:   "rid1",
					BeforeState:  datatypes.JSON(`{"status":"inactive"}`),
					AfterState:   datatypes.JSON(`{"status":"active"}`),
					Metadata:     datatypes.JSON(`{"ip":"127.0.0.1"}`),
				},
				err: errors.New("some error"),
			},
			expErr: errors.New("some error"),
		},
		"success": {
			givenIn: model.AuditLog{
				ServiceName:  "svc1",
				Username:     "user1",
				Action:       "create",
				ResourceType: "type1",
				ResourceID:   "rid1",
				BeforeState:  datatypes.JSON(`{"status":"inactive"}`),
				AfterState:   datatypes.JSON(`{"status":"active"}`),
				Metadata:     datatypes.JSON(`{"ip":"127.0.0.1"}`),
			},
			mockRepo: mockRepo{
				in: model.AuditLog{
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

			mockAuditRepo.On("InsertAuditLog", mock.Anything, tc.mockRepo.in).
				Return(model.AuditLog{}, tc.mockRepo.err)

			// When:
			instance := New(mockRegistry)
			err := instance.HandleAuditLogMessage(context.Background(), tc.givenIn)

			// Then:
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
