package constant

const (
	TokenKey = "Authorization"

	KeyUserTokenPrefix = "evobot:user:token:"
	KeyPostTimeZSet    = "evobot:post:time"

	KeyStaffWeightPrefix   = "evobot:staff:weight:"
	KeyWecomKFStaffPrefix  = "wecom:kf:staff:"
	KeyWecomBotStaffPrefix = "wecom:bot:staff:"
)

const (
	KFStatusRobotToHuman int = iota + 1 // 1 机器人可转人工
	KFStatusOnlyRobot                   // 2 仅机器人
	KFStatusOnlyHuman                   // 3 仅人工
)

const (
	KFReceiveRuleRoundRobin int = iota + 1 // 1 轮流接待
	KFReceiveRuleIdle                      // 2 空闲接待
)

const (
	KHStatusUnprocessed     int = iota + 1 // 未处理
	KHStatusVerification                   // 验证
	KHStatusUserInfoConfirm                // 用户信息确认
)
