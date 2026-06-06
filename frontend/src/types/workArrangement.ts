export interface WorkArrangement {
  id: number;
  project_id: number;
  date: string;
  customer: string;
  project: string;
  work_type: string;
  location: string;
  partner: string;
  content: string;
  duration: number;
  progress: string;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface FilterParams {
  date_from?: string;
  date_to?: string;
  customer?: string;
  project?: string;
  work_type?: string;
  progress?: string;
}

export interface ImportResult {
  created: number;
  skipped: number;
  errors: string[];
}

export const WORK_TYPES = ['测试', '交付', '售后'] as const;
export const LOCATIONS = ['远程', '现场'] as const;
export const PARTNERS = ['是', '否'] as const;
export const PROGRESSES = ['未开始', '进行中', '已完成', '已暂停', '已取消'] as const;
