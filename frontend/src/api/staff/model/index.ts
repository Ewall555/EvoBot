import type { IPolicy } from '@/api/policy/model';

export interface IStaff {
    uuid: string;
    staffname?: string;
    staffid?: string;
    number?: string;
    email?: string;
    role?: string;
    policy_list?: string[];
    policies?: Array<IPolicy>;
}
