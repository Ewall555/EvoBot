package client

const (
	KeyWecomCursorPrefix     = "wecom:cursor:"
	KeyWecomCursorLockPrefix = "wecom:cursor_lock:"

	WecomCallbackMsgTypeEvent      = "event"
	WecomCallbackEventKFMsgOrEvent = "kf_msg_or_event"

	WecomMsgTypeText = "text"

	WecomMsgTypeEvent                     = "event"
	WecomMsgTypeEnterSessionEvent         = "enter_session"          // 用户进入会话事件
	WecomMsgTypeMsgSendFailEvent          = "msg_send_fail"          // 消息发送失败事件
	WecomMsgTypeServicerStatusChangeEvent = "servicer_status_change" // 接待人员接待状态变更事件
	WecomMsgTypeSessionStatusChangeEvent  = "session_status_change"  // 会话状态变更事件
	WecomMsgTypeUserRecallMsgEvent        = "user_recall_msg"        // 用户撤回消息事件
	WecomMsgTypeServicerRecallMsgEvent    = "servicer_recall_msg"    // 接待人员撤回消息事件
)
const (
	WecomEventChangeTypeJoinSession     uint32 = iota + 1 // 从接待池接入会话
	WecomEventChangeTypeTransferSession                   // 转接会话
	WecomEventChangeTypeEndSession                        // 结束会话
	WecomEventChangeTypeRejoinSession                     // 重新接入已结束/已转接会话
)

const (
	SessionStatusNew               int = iota // 0 未处理
	SessionStatusHandled                      // 1 由智能助手接待
	SessionStatusWaiting                      // 2 待接入池排队中
	SessionStatusInProgress                   // 3 由人工接待
	SessionStatusEndedOrNotStarted            // 4 已结束/未开始
)

const (
	MessageTypeCustomer            uint32 = iota + 3 // 3 微信客户发送的消息
	MessageTypeSystemEvent                           // 4 系统推送的事件消息
	MessageTypeReceptionistMessage                   // 5 接待人员在企业微信客户端发送的消息
)
