package logic

import "EvoBot/backend/app/repo"

var (
	commonRepo = repo.NewCommonRepo()

	authRepo   = repo.NewISettingRepo()
	wecomRepo  = repo.NewIWecomRepo()
	llmappRepo = repo.NewILLMAppRepo()
	kHRepo     = repo.NewIKHRepo()
	kFRepo     = repo.NewIKFRepo()
	policyRepo = repo.NewIPolicyRepo()
	staffRepo  = repo.NewIStaffRepo()
)
