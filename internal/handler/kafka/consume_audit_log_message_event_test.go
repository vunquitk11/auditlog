package kafka

import (
	"context"
	"errors"
	"testing"

	"github.com/audit-log-service/internal/controller/auditlog"
	"github.com/audit-log-service/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func Test_ConsumeAuditLogMessageEvent(t *testing.T) {
	type mockCtrl struct {
		in  model.AuditLog
		err error
	}

	type arg struct {
		givenInput *model.AuditMessage
		mockCtrl   mockCtrl
		expErr     error
	}

	tcs := map[string]arg{
		"error when handle audit log message": {
			givenInput: &model.AuditMessage{
				ID:          "fail-id",
				UserID:      "user_fail",
				ServiceName: "svc_fail",
				Action:      "fail_action",
				Details: map[string]interface{}{
					"amount":   10000,
					"currency": "USD",
					"status":   "fail",
				},
			},
			mockCtrl: mockCtrl{
				in: model.AuditLog{
					ID:          "fail-id",
					ServiceName: "svc_fail",
					Username:    "user_fail",
					Action:      "fail_action",
					BeforeState: datatypes.JSON("{}"),
					AfterState:  datatypes.JSON(`{"amount":10000,"currency":"USD","status":"fail"}`),
					Metadata:    datatypes.JSON(`{"ip_address":"","user_agent":""}`),
				},
				err: errors.New("some error"),
			},
			expErr: errors.New("some error"),
		},
		"success when handle audit log message": {
			givenInput: &model.AuditMessage{
				ID:          "ok-id",
				UserID:      "user_ok",
				ServiceName: "svc_ok",
				Action:      "ok_action",
				Details: map[string]interface{}{
					"amount":         150000,
					"currency":       "VND",
					"merchant_id":    "mch_123",
					"order_id":       "order_456",
					"payment_method": "credit_card",
					"status":         "success",
				},
			},
			mockCtrl: mockCtrl{
				in: model.AuditLog{
					ID:          "ok-id",
					ServiceName: "svc_ok",
					Username:    "user_ok",
					Action:      "ok_action",
					BeforeState: datatypes.JSON("{}"),
					AfterState:  datatypes.JSON(`{"amount":150000,"currency":"VND","merchant_id":"mch_123","order_id":"order_456","payment_method":"credit_card","status":"success"}`),
					Metadata:    datatypes.JSON(`{"ip_address":"","user_agent":""}`),
				},
				err: nil,
			},
			expErr: nil,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Mock
			mockAuditCtrl := auditlog.NewMockController(t)
			mockAuditCtrl.On("HandleAuditLogMessage", mock.Anything, tc.mockCtrl.in).
				Return(tc.mockCtrl.err)

			// When:
			instance := New(mockAuditCtrl)
			err := instance.ConsumeAuditLogMessageEvent(context.Background(), tc.givenInput)

			// Then:
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
