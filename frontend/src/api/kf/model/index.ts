import type { IStaff } from '@/api/staff/model';

export interface IKF {
    uuid?: string;
    kfname?: string;
    kfid?: string;
    kfplatform?: string;
    botid?: string;
    botplatform?: string;
    status?: number;
    receive_priority?: number;
    receive_rule?: number;
    chat_timeout?: number;
    bot_timeout?: number;
    bot_timeout_msg?: string;
    bot_welcome_msg?: string;
    staff_welcome_msg?: string;
    unmanned_msg?: string;
    chatend_msg?: string;
    transfer_keywords?: string;
    staff_list?: string[];
    staffs?: Array<IStaff>;
}
