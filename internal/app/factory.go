package app

import (
	"biz/internal/audit"
	"biz/internal/platform/config"
	logx "biz/internal/platform/log"
	"biz/internal/policy"
)

func Build(configPath, profileOverride string) (*Runtime, error) {
	cfg, err := config.Load(configPath, profileOverride)
	if err != nil {
		return nil, err
	}
	logger, err := logx.New(cfg.Log, cfg.Profile)
	if err != nil {
		return nil, err
	}
	agentAuthorizer, err := policy.NewAgentAuthorizer(cfg.AgentPolicy)
	if err != nil {
		return nil, err
	}
	var auditor *audit.Writer
	if cfg.Audit.Enabled {
		auditor, err = audit.NewWriter(cfg.Audit)
		if err != nil {
			return nil, err
		}
	}

	return &Runtime{
		Config:          cfg,
		Logger:          logger,
		AgentAuthorizer: agentAuthorizer,
		Auditor:         auditor,
	}, nil
}
