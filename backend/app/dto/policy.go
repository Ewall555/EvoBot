package dto

type Policy struct {
	PolicyName string     `json:"policyname"`
	Repeat     int        `json:"repeat"`
	Week       string     `json:"week"`
	MaxCount   int        `json:"max_count"`
	WorkTimes  []WorkTime `json:"work_times"`
}

type WorkTime struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
