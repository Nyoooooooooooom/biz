package app

import (
	"biz/internal/audit"
	"biz/internal/invoice"
	"biz/internal/platform/config"
	"biz/internal/policy"
	"biz/internal/records"
	"biz/internal/tax"
	"go.uber.org/zap"
)

type Runtime struct {
	Config          config.Config
	Logger          *zap.Logger
	Tax             tax.Fragment
	Invoice         invoice.Fragment
	Records         records.Fragment
	AgentAuthorizer *policy.AgentAuthorizer
	Auditor         *audit.Writer
}
