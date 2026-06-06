import React, { useEffect } from 'react';
import { Modal, Form, Input, DatePicker, Select, InputNumber, message } from 'antd';
import dayjs from 'dayjs';
import type { WorkArrangement } from '../../types/workArrangement';
import { WORK_TYPES, LOCATIONS, PROGRESSES } from '../../types/workArrangement';

const { TextArea } = Input;

interface WorkFormProps {
  visible: boolean;
  editingRecord?: WorkArrangement | null;
  onCancel: () => void;
  onSubmit: (values: WorkArrangement) => Promise<void>;
}

const WorkForm: React.FC<WorkFormProps> = ({ visible, editingRecord, onCancel, onSubmit }) => {
  const [form] = Form.useForm();
  const [loading, setLoading] = React.useState(false);
  const isEdit = !!editingRecord;

  useEffect(() => {
    if (visible) {
      if (editingRecord) {
        form.setFieldsValue({
          ...editingRecord,
          date: dayjs(editingRecord.date),
        });
      } else {
        form.setFieldsValue({
          date: dayjs(),
          work_type: '测试',
          location: '远程',
          partner: '否',
          duration: 0,
          progress: '未开始',
        });
      }
    }
  }, [visible, editingRecord, form]);

  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      const record: WorkArrangement = {
        id: 0,
        project_id: values.project_id || 0,
        date: values.date.format('YYYY-MM-DD'),
        customer: values.customer || '',
        project: values.project || '',
        work_type: values.work_type,
        location: values.location,
        partner: values.partner,
        content: values.content || '',
        duration: values.duration || 0,
        progress: values.progress || '未开始',
        notes: values.notes || '',
        created_at: editingRecord?.created_at || '',
        updated_at: editingRecord?.updated_at || '',
      };

      await onSubmit(record);
      message.success(isEdit ? '修改成功' : '创建成功');
      form.resetFields();
      onCancel();
    } catch (err: any) {
      if (err?.errorFields) return;
      message.error(`操作失败: ${err?.message || err}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title={isEdit ? '编辑工作安排' : '新增工作安排'}
      open={visible}
      onOk={handleOk}
      onCancel={() => { form.resetFields(); onCancel(); }}
      confirmLoading={loading}
      okText={isEdit ? '保存' : '创建'}
      cancelText="取消"
      width={680}
      destroyOnClose
    >
      <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0 16px' }}>
          <Form.Item name="project_id" label="ID">
            <InputNumber
              style={{ width: '100%' }}
              min={0}
              placeholder="项目ID/工单ID"
            />
          </Form.Item>

          <Form.Item name="date" label="日期" rules={[{ required: true, message: '请选择日期' }]}>
            <DatePicker style={{ width: '100%' }} format="YYYY-MM-DD" />
          </Form.Item>

          <Form.Item name="duration" label="工作耗时（小时）">
            <InputNumber style={{ width: '100%' }} min={0} max={24} step={0.5} />
          </Form.Item>

          <Form.Item name="customer" label="客户名称" rules={[{ required: true, message: '请输入客户名称' }]}>
            <Input placeholder="请输入客户名称" />
          </Form.Item>

          <Form.Item name="project" label="项目名称" rules={[{ required: true, message: '请输入项目名称' }]}>
            <Input placeholder="请输入项目名称" />
          </Form.Item>

          <Form.Item name="work_type" label="工作类型" rules={[{ required: true, message: '请选择工作类型' }]}>
            <Select options={WORK_TYPES.map((t) => ({ label: t, value: t }))} />
          </Form.Item>

          <Form.Item name="location" label="工作地点" rules={[{ required: true, message: '请选择工作地点' }]}>
            <Select options={LOCATIONS.map((l) => ({ label: l, value: l }))} />
          </Form.Item>

          <Form.Item name="partner" label="伙伴" rules={[{ required: true, message: '请选择' }]}>
            <Select options={[{ label: '是', value: '是' }, { label: '否', value: '否' }]} />
          </Form.Item>

          <Form.Item name="progress" label="工作进度" rules={[{ required: true, message: '请选择工作进度' }]}>
            <Select options={PROGRESSES.map((p) => ({ label: p, value: p }))} />
          </Form.Item>
        </div>

        <Form.Item name="content" label="工作内容">
          <TextArea rows={3} placeholder="请输入工作内容" />
        </Form.Item>

        <Form.Item name="notes" label="备注">
          <TextArea rows={2} placeholder="可选备注信息" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default WorkForm;
