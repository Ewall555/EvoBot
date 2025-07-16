export interface IPolicy {
    uuid: string;
    policyname?: string; // 策略名称
    max_count?: number; // 最大接待数量
    repeat?: number; // 重复策略
    week?: string; // 周工作日标记，例如 "1111100" 表示周一到周五工作
    work_times?: WorkTime[]; // 工作时间数组
}

export interface WorkTime {
    start_time?: string; // 工作开始时间，例如 "09:00:00"
    end_time?: string; // 工作结束时间，例如 "12:00:00"
}