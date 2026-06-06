import { useState, useEffect, useCallback } from 'react';
import {
  ConfigProvider, theme, App as AntApp, Button, Space, Dropdown,
  Switch, Card, Typography, message,
} from 'antd';
import {
  PlusOutlined, UploadOutlined, DownloadOutlined,
  FileExcelOutlined, FileTextOutlined,
} from '@ant-design/icons';
import zhCN from 'antd/locale/zh_CN';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import AppLayout from './components/Layout/AppLayout';
import FilterBar from './components/WorkArrangement/FilterBar';
import WorkTable from './components/WorkArrangement/WorkTable';
import WorkForm from './components/WorkArrangement/WorkForm';
import ImportModal from './components/WorkArrangement/ImportModal';
import { useWorkArrangements } from './hooks/useWorkArrangements';
import * as api from './api/workArrangements';
import { downloadFile } from './api/client';
import type { WorkArrangement, FilterParams } from './types/workArrangement';
import './style.css';

dayjs.locale('zh-cn');

const API_BASE = '/api';

function App() {
  const [menuKey, setMenuKey] = useState('work-arrangements');
  const [filters, setFilters] = useState<FilterParams>({});
  const [formVisible, setFormVisible] = useState(false);
  const [editingRecord, setEditingRecord] = useState<WorkArrangement | null>(null);
  const [importModalVisible, setImportModalVisible] = useState(false);
  const [autoStart, setAutoStartState] = useState(false);
  const [autoStartLoading, setAutoStartLoading] = useState(false);

  const { data, loading, fetchAll, filter, create, update, remove, copyText } =
    useWorkArrangements();

  useEffect(() => {
    fetchAll();
    loadAutoStartStatus();
  }, []);

  const loadAutoStartStatus = async () => {
    try {
      const { enabled } = await api.getAutoStartStatus();
      setAutoStartState(enabled);
    } catch { /* ignore */ }
  };

  const handleAutoStartToggle = async (enabled: boolean) => {
    setAutoStartLoading(true);
    try {
      await api.setAutoStart(enabled);
      setAutoStartState(enabled);
      message.success(enabled ? '已开启开机自启' : '已关闭开机自启');
    } catch (err: any) {
      message.error(`操作失败: ${err?.message || err}`);
    } finally {
      setAutoStartLoading(false);
    }
  };

  const handleSearch = useCallback(() => { filter(filters); }, [filters, filter]);
  const handleReset = useCallback(() => { setFilters({}); fetchAll(); }, [fetchAll]);

  const handleAdd = useCallback(() => {
    setEditingRecord(null);
    setFormVisible(true);
  }, []);

  const handleEdit = useCallback((record: WorkArrangement) => {
    setEditingRecord(record);
    setFormVisible(true);
  }, []);

  const handleFormSubmit = useCallback(
    async (values: WorkArrangement) => {
      if (editingRecord?.id) {
        await update({ ...values, id: editingRecord.id });
      } else {
        await create(values);
      }
    },
    [editingRecord, create, update]
  );

  const handleDelete = useCallback((id: number) => { remove(id); }, [remove]);
  const handleCopy = useCallback((record: WorkArrangement) => { copyText(record); }, [copyText]);

  // Export
  const buildExportUrl = useCallback((format: string) => {
    const params = new URLSearchParams();
    if (filters.date_from) params.set('date_from', filters.date_from);
    if (filters.date_to) params.set('date_to', filters.date_to);
    if (filters.customer) params.set('customer', filters.customer);
    if (filters.project) params.set('project', filters.project);
    if (filters.work_type) params.set('work_type', filters.work_type);
    if (filters.progress) params.set('progress', filters.progress);
    const qs = params.toString();
    return `${API_BASE}/export/${format}${qs ? '?' + qs : ''}`;
  }, [filters]);

  const handleExportExcel = useCallback(() => {
    downloadFile(buildExportUrl('excel'), '工作安排.xlsx');
  }, [buildExportUrl]);
  const handleExportJSON = useCallback(() => {
    downloadFile(buildExportUrl('json'), '工作安排.json');
  }, [buildExportUrl]);

  const handleImportSuccess = useCallback(() => { fetchAll(); }, [fetchAll]);

  const exportItems = {
    items: [
      { key: 'excel', label: '导出 Excel (.xlsx)', icon: <FileExcelOutlined />, onClick: handleExportExcel },
      { key: 'json', label: '导出 JSON (.json)', icon: <FileTextOutlined />, onClick: handleExportJSON },
    ],
  };

  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        algorithm: theme.defaultAlgorithm,
        token: {
          colorPrimary: '#1677ff',
          borderRadius: 6,
        },
      }}
    >
      <AntApp>
        <AppLayout activeKey={menuKey} onMenuClick={setMenuKey}>
          {menuKey === 'settings' ? (
            <div style={{ padding: 24 }}>
              <Typography.Title level={4} style={{ marginBottom: 24 }}>系统设置</Typography.Title>
              <Card style={{ maxWidth: 500 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <div>
                    <Typography.Text strong>开机自启动</Typography.Text>
                    <br />
                    <Typography.Text type="secondary">开启后，系统启动时将自动运行本应用</Typography.Text>
                  </div>
                  <Switch checked={autoStart} loading={autoStartLoading} onChange={handleAutoStartToggle} />
                </div>
              </Card>
            </div>
          ) : (
            <div>
              <FilterBar filters={filters} onChange={setFilters} onSearch={handleSearch} onReset={handleReset} />

              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
                <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>新增</Button>
                <Space>
                  <Button icon={<UploadOutlined />} onClick={() => setImportModalVisible(true)}>导入数据</Button>
                  <Dropdown menu={exportItems}>
                    <Button icon={<DownloadOutlined />}>导出数据</Button>
                  </Dropdown>
                </Space>
              </div>

              <WorkTable
                data={data}
                loading={loading}
                onEdit={handleEdit}
                onCopy={handleCopy}
                onDelete={handleDelete}
              />

              <WorkForm
                visible={formVisible}
                editingRecord={editingRecord}
                onCancel={() => setFormVisible(false)}
                onSubmit={handleFormSubmit}
              />

              <ImportModal
                visible={importModalVisible}
                onCancel={() => setImportModalVisible(false)}
                onSuccess={handleImportSuccess}
              />
            </div>
          )}
        </AppLayout>
      </AntApp>
    </ConfigProvider>
  );
}

export default App;
