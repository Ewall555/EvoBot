package logic

import (
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	"context"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
)

func getHighestWeightStaff(staffIDs []string) (string, int) {
	maxWeightvalue := -1
	selectedstaffid := ""
	for _, staffid := range staffIDs {
		staffweightkey := constant.KeyStaffWeightPrefix + staffid
		staffweightvalue, err := global.RDS.Get(context.Background(), staffweightkey).Int()
		if err != nil {
			global.ZAPLOG.Error("获取权重失败", zap.String("staffweightkey", staffweightkey), zap.Error(err))
			continue
		}
		if staffweightvalue > maxWeightvalue {
			maxWeightvalue = staffweightvalue
			selectedstaffid = staffid
		}
	}
	return selectedstaffid, maxWeightvalue
}

func isStaffWorkByStaffID(staffid string) (bool, error) {
	staffinfo, err := staffRepo.Get(staffRepo.WithStaffID(staffid))
	if err != nil {
		return false, err
	}
	for _, policy := range staffinfo.Policies {
		if isWithinTime(policy.Repeat, policy.Week) {
			for _, worktime := range policy.WorkTimes {
				if isTimeInRange(worktime.StartTime, worktime.EndTime) {
					ctx := context.Background()
					staffweightkey := constant.KeyStaffWeightPrefix + staffid
					_, err := global.RDS.Get(ctx, staffweightkey).Result()
					if err != nil {
						global.ZAPLOG.Error("redis get error", zap.Error(err))
						global.ZAPLOG.Info("初始化权重缓存")
						err = global.RDS.Set(ctx, staffweightkey, policy.MaxCount, 0).Err()
						if err != nil {
							global.ZAPLOG.Error("redis set error", zap.Error(err))
							return false, err
						}
					}
					return true, nil
				}
			}
			global.ZAPLOG.Info("该接待人员的工作时间导致目前无法接待", zap.String("StaffName", staffinfo.StaffName))
			return false, nil
		} else {
			global.ZAPLOG.Info("该接待人员的策略导致目前无法接待", zap.String("StaffName", staffinfo.StaffName))
			return false, nil
		}
	}
	return false, nil
}

func isTimeInRange(startTimeStr, endTimeStr string) bool {
	layout := "15:04:05"
	startTime, err1 := time.Parse(layout, startTimeStr)
	endTime, err2 := time.Parse(layout, endTimeStr)
	if err1 != nil || err2 != nil {
		global.ZAPLOG.Error("Error parsing time:", zap.Error(err1))
		global.ZAPLOG.Error("Error parsing time:", zap.Error(err2))
		return false
	}
	now := time.Now()
	currentTime, _ := time.Parse(layout, now.Format(layout))
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}
	return false
}

func isWithinTime(repeat int, week string) bool {
	currentTime := time.Now()
	currentWeekday := int(currentTime.Weekday()) // 0: Sunday, 1: Monday, ..., 6: Saturday
	switch repeat {
	case 1:
		if len(week) != 7 {
			return false
		}
		return week[currentWeekday] == '1'
	case 2:
		return true
	case 3:
		return currentWeekday >= 1 && currentWeekday <= 5
	case 4:
		// 法定工作日有效，跳过法定节假日（此处需要法定节假日的额外数据支持，简单实现为周一到周五）
		// 简单假设法定工作日是周一到周五，实际情况可能要处理节假日
		return currentWeekday >= 1 && currentWeekday <= 5
	case 5:
		// 法定节假日有效，跳过法定工作日（此处需要法定节假日的额外数据支持）
		// 简单假设法定节假日是周六和周日，实际情况可能要处理节假日
		return currentWeekday == 0 || currentWeekday == 6
	default:
		return false
	}
}

func MarkdownToText(markdown string) string {
	var buffer strings.Builder
	lines := strings.Split(markdown, "\n")
	inList := false
	for _, line := range lines {
		if titleMatch := regexp.MustCompile(`^(#{1,6})\s*(.*)$`).FindStringSubmatch(line); titleMatch != nil {
			buffer.WriteString(titleMatch[2] + "\n")
			continue
		}
		if listMatch := regexp.MustCompile(`^(\s*)-\s+(.*)$`).FindStringSubmatch(line); listMatch != nil {
			if !inList {
				buffer.WriteString("\n")
				inList = true
			}
			buffer.WriteString(strings.Repeat(" ", len(listMatch[1])) + "- " + listMatch[2] + "\n")
			continue
		}
		if inList && strings.TrimSpace(line) == "" {
			inList = false
			continue
		}
		buffer.WriteString(line + "\n")
	}
	text := buffer.String()
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`__(.*?)__`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`\*(.*?)\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`\[.*?\]\((.*?)\)`).ReplaceAllString(text, "$1")
	text = strings.ReplaceAll(text, "\n\n", "\n")
	reNewline := regexp.MustCompile(`\n+`)
	text = reNewline.ReplaceAllString(text, "\n")
	return text
}
