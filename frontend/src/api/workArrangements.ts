import { apiRequest } from './client';
import type { WorkArrangement, FilterParams, ImportResult } from '../types/workArrangement';

// Work Arrangements CRUD
export function fetchAll() {
  return apiRequest<WorkArrangement[]>('/work-arrangements');
}

export function fetchFiltered(filters: FilterParams) {
  const params: Record<string, string> = {};
  if (filters.date_from) params.date_from = filters.date_from;
  if (filters.date_to) params.date_to = filters.date_to;
  if (filters.customer) params.customer = filters.customer;
  if (filters.project) params.project = filters.project;
  if (filters.work_type) params.work_type = filters.work_type;
  if (filters.progress) params.progress = filters.progress;
  return apiRequest<WorkArrangement[]>('/work-arrangements', { params });
}

export function fetchById(id: number) {
  return apiRequest<WorkArrangement>(`/work-arrangements/${id}`);
}

export function createOne(record: WorkArrangement) {
  return apiRequest<WorkArrangement>('/work-arrangements', { method: 'POST', body: record });
}

export function updateOne(record: WorkArrangement) {
  return apiRequest<WorkArrangement>(`/work-arrangements/${record.id}`, { method: 'PUT', body: record });
}

export function deleteOne(id: number) {
  return apiRequest<void>(`/work-arrangements/${id}`, { method: 'DELETE' });
}

export function bulkCreate(records: WorkArrangement[]) {
  return apiRequest<ImportResult>('/work-arrangements/bulk', { method: 'POST', body: records });
}

// Reference data
export function fetchCustomers() {
  return apiRequest<string[]>('/reference/customers');
}

export function fetchProjects() {
  return apiRequest<string[]>('/reference/projects');
}

// Import
export function importParse(file: File) {
  const formData = new FormData();
  formData.append('file', file);
  return apiRequest<ImportResult & { Records?: WorkArrangement[] }>('/import/parse', {
    method: 'POST',
    body: formData,
  });
}

export function importConfirm(records: WorkArrangement[]) {
  return apiRequest<ImportResult>('/import/confirm', { method: 'POST', body: records });
}

// Auto-start
export function getAutoStartStatus() {
  return apiRequest<{ enabled: boolean }>('/settings/autostart');
}

export function setAutoStart(enabled: boolean) {
  return apiRequest<{ enabled: boolean }>('/settings/autostart', { method: 'PUT', body: { enabled } });
}
