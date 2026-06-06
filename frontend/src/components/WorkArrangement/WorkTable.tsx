import React from 'react';
import { Table, Tag, Button, Space, Popconfirm, Tooltip, Typography } from 'antd';
import { EditOutlined, CopyOutlined, DeleteOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { WorkArrangement } from '../../types/workArrangement';

const { Text } = Typography;

interface WorkTableProps {
  data: WorkArrangement[];
  loading: boolean;
  onEdit: (record: WorkArrangement) => void;
  onCopy: (record: WorkArrangement) => void;
  onDelete: (id: number) => void;
}

const workTypeColor: Record<string, string> = {
  '售前': 'blue',
  '售后': 'green',
  '测试': 'orange',
};

const progressColor: Record<string, string> = {
  '未开始': 'default',
  '进行中': 'processing',
  '已完成': 'success',
  '已暂停': 'warning',
  '已取消': 'error',
};

const WorkTable: React.FC<WorkTableProps> = ({ data, loading, onEdit, onCopy, onDelete }) => {
  const columns: ColumnsType<WorkArrangement> = [
    {
      title: '日期', dataIndex: 'date', key: 'date', width: 110,
      render: (date: string) => <Text style={{ fontFamily: 'monospace' }}>{date}</Text>,
    },
    { title: '客户名称', dataIndex: 'customer', key: 'customer', width: 140, ellipsis: true },
    { title: '项目名称', dataIndex: 'project', key: 'project', width: 140, ellipsis: true },
    {
      title: 'ID', dataIndex: 'project_id', key: 'project_id', width: 80, align: 'center',
      render: (id: number) => id > 0 ? <Text style={{ fontFamily: 'monospace' }}>{id}</Text> : <Text type="secondary">-</Text>,
    },
    {
      title: '工作类型', dataIndex: 'work_type', key: 'work_type', width: 90,
      render: (t: string) => <Tag color={workTypeColor[t] || 'default'}>{t}</Tag>,
    },
    {
      title: '地点', dataIndex: 'location', key: 'location', width: 70,
      render: (loc: string) => <Tag color={loc === '远程' ? 'cyan' : 'purple'}>{loc}</Tag>,
    },
    {
      title: '伙伴', dataIndex: 'partner', key: 'partner', width: 70,
      render: (p: string) => <Tag color={p === '是' ? 'green' : 'default'}>{p}</Tag>,
    },
    { title: '工作内容', dataIndex: 'content', key: 'content', width: 200, ellipsis: true },
    { title: '耗时(h)', dataIndex: 'duration', key: 'duration', width: 80, align: 'center', render: (d: number) => (d > 0 ? d : '-') },
    {
      title: '进度', dataIndex: 'progress', key: 'progress', width: 90,
      render: (p: string) => <Tag color={progressColor[p] || 'default'}>{p}</Tag>,
    },
    { title: '备注', dataIndex: 'notes', key: 'notes', width: 150, ellipsis: true, render: (n: string) => n || '-' },
    {
      title: '操作', key: 'actions', width: 130, fixed: 'right',
      render: (_: any, record: WorkArrangement) => (
        <Space size="small">
          <Tooltip title="编辑">
            <Button type="link" size="small" icon={<EditOutlined />} onClick={() => onEdit(record)} />
          </Tooltip>
          <Tooltip title="复制">
            <Button type="link" size="small" icon={<CopyOutlined />} onClick={() => onCopy(record)} />
          </Tooltip>
          <Popconfirm title="确认删除？" description="删除后无法恢复" onConfirm={() => onDelete(record.id)} okText="确认" cancelText="取消">
            <Tooltip title="删除">
              <Button type="link" size="small" danger icon={<DeleteOutlined />} />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Table
      columns={columns}
      dataSource={data}
      rowKey="id"
      loading={loading}
      scroll={{ x: 1500 }}
      size="middle"
      pagination={{
        defaultPageSize: 50,
        showSizeChanger: true,
        pageSizeOptions: ['20', '50', '100', '200'],
        showTotal: (total) => `共 ${total} 条记录`,
        position: ['bottomRight'],
      }}
      locale={{
        emptyText: (
          <div style={{ padding: 40 }}>
            <Text type="secondary" style={{ fontSize: 14 }}>暂无数据，点击"新增"开始记录 ✨</Text>
          </div>
        ),
      }}
    />
  );
};

export default WorkTable;
