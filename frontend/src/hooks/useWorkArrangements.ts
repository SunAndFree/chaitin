import { useState, useCallback } from 'react';
import { message } from 'antd';
import * as api from '../api/workArrangements';
import type { WorkArrangement, FilterParams } from '../types/workArrangement';

export function useWorkArrangements() {
  const [data, setData] = useState<WorkArrangement[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchAll = useCallback(async () => {
    setLoading(true);
    try {
      const result = await api.fetchAll();
      setData(result || []);
    } catch (err: any) {
      message.error(`获取数据失败: ${err?.message || err}`);
      setData([]);
    } finally {
      setLoading(false);
    }
  }, []);

  const filter = useCallback(async (filters: FilterParams) => {
    setLoading(true);
    try {
      const hasFilters = Object.values(filters).some((v) => v !== '' && v !== undefined);
      const result = hasFilters ? await api.fetchFiltered(filters) : await api.fetchAll();
      setData(result || []);
    } catch (err: any) {
      message.error(`筛选失败: ${err?.message || err}`);
    } finally {
      setLoading(false);
    }
  }, []);

  const create = useCallback(async (record: WorkArrangement): Promise<WorkArrangement | null> => {
    try {
      const created = await api.createOne(record);
      if (created) {
        setData((prev) => [created, ...prev]);
      }
      return created;
    } catch (err: any) {
      message.error(`创建失败: ${err?.message || err}`);
      return null;
    }
  }, []);

  const update = useCallback(async (record: WorkArrangement): Promise<WorkArrangement | null> => {
    try {
      const updated = await api.updateOne(record);
      if (updated) {
        setData((prev) => prev.map((item) => (item.id === updated.id ? updated : item)));
      }
      return updated;
    } catch (err: any) {
      message.error(`更新失败: ${err?.message || err}`);
      return null;
    }
  }, []);

  const remove = useCallback(async (id: number) => {
    try {
      await api.deleteOne(id);
      setData((prev) => prev.filter((item) => item.id !== id));
      message.success('删除成功');
    } catch (err: any) {
      message.error(`删除失败: ${err?.message || err}`);
    }
  }, []);

  const copyText = useCallback(async (record: WorkArrangement) => {
    try {
      const text =
        record.partner === '是'
          ? `【${record.location}】【${record.work_type}】【生态】${record.project}-${record.content}`
          : `【${record.location}】【${record.work_type}】${record.project}-${record.content}`;
      await navigator.clipboard.writeText(text);
      message.success('已复制到剪贴板');
    } catch {
      message.error('复制失败，请重试');
    }
  }, []);

  return { data, loading, fetchAll, filter, create, update, remove, copyText };
}
