import React, { useState, useRef } from 'react';
import { Modal, Table, Tag, Button, Typography, message, Alert, Upload } from 'antd';
import { InboxOutlined } from '@ant-design/icons';
import type { RcFile } from 'antd/es/upload';
import { importParse, importConfirm } from '../../api/workArrangements';
import type { WorkArrangement, ImportResult } from '../../types/workArrangement';

const { Text } = Typography;

interface ImportModalProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
}

const ImportModal: React.FC<ImportModalProps> = ({ visible, onCancel, onSuccess }) => {
  const [previewData, setPreviewData] = useState<WorkArrangement[]>([]);
  const [parseResult, setParseResult] = useState<{ created: number; skipped: number; errors: string[] } | null>(null);
  const [loading, setLoading] = useState(false);
  const [fileSelected, setFileSelected] = useState(false);

  const handleFileSelect = async (file: RcFile): Promise<false> => {
    setLoading(true);
    try {
      const result: any = await importParse(file);
      setPreviewData(result.Records || []);
      setParseResult({
        created: result.Created ?? result.created ?? 0,
        skipped: result.Skipped ?? result.skipped ?? 0,
        errors: result.Errors ?? result.errors ?? [],
      });
      setFileSelected(true);
    } catch (err: any) {
      message.error(`解析失败: ${err?.message || err}`);
    } finally {
      setLoading(false);
    }
    return false; // Prevent default upload
  };

  const handleConfirmImport = async () => {
    if (!previewData.length) {
      message.warning('没有可导入的数据');
      return;
    }
    setLoading(true);
    try {
      const result = await importConfirm(previewData);
      if (result.errors && result.errors.length > 0) {
        message.warning(`导入完成: ${result.created} 条成功, ${result.skipped} 条失败`);
      } else {
        message.success(`成功导入 ${result.created} 条记录`);
      }
      handleReset();
      onSuccess();
      onCancel();
    } catch (err: any) {
      message.error(`导入失败: ${err?.message || err}`);
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    setPreviewData([]);
    setParseResult(null);
    setFileSelected(false);
  };

  const handleCancel = () => {
    handleReset();
    onCancel();
  };

  const previewColumns = [
    { title: '日期', dataIndex: 'date', width: 100 },
    { title: '客户', dataIndex: 'customer', width: 100, ellipsis: true },
    { title: '项目', dataIndex: 'project', width: 100, ellipsis: true },
    { title: '类型', dataIndex: 'work_type', width: 70, render: (t: string) => <Tag>{t}</Tag> },
    { title: '内容', dataIndex: 'content', width: 150, ellipsis: true },
    { title: '进度', dataIndex: 'progress', width: 80, render: (p: string) => <Tag>{p}</Tag> },
  ];

  return (
    <Modal
      title="导入数据"
      open={visible}
      onCancel={handleCancel}
      width={800}
      footer={
        fileSelected && previewData.length > 0
          ? [
              <Button key="cancel" onClick={handleCancel}>取消</Button>,
              <Button key="reset" onClick={handleReset}>重新选择</Button>,
              <Button key="confirm" type="primary" loading={loading} onClick={handleConfirmImport}>
                确认导入 ({previewData.length} 条)
              </Button>,
            ]
          : [<Button key="cancel" onClick={handleCancel}>取消</Button>]
      }
    >
      {!fileSelected ? (
        <Upload.Dragger
          accept=".xlsx,.json"
          showUploadList={false}
          beforeUpload={handleFileSelect}
          maxCount={1}
          style={{ padding: '24px 0' }}
        >
          <p className="ant-upload-drag-icon">
            <InboxOutlined />
          </p>
          <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
          <p className="ant-upload-hint">支持 .xlsx 和 .json 格式</p>
        </Upload.Dragger>
      ) : (
        <div>
          {parseResult && parseResult.errors.length > 0 && (
            <Alert
              type="warning"
              showIcon
              closable
              message={`解析: ${parseResult.created} 条有效, ${parseResult.skipped} 条无效`}
              description={parseResult.errors.slice(0, 10).map((e: string, i: number) => (
                <div key={i}>{e}</div>
              ))}
              style={{ marginBottom: 16 }}
            />
          )}
          {previewData.length > 0 && (
            <div>
              <Text type="secondary" style={{ marginBottom: 8, display: 'block' }}>
                预览数据（共 {previewData.length} 条）
              </Text>
              <Table
                columns={previewColumns}
                dataSource={previewData}
                rowKey={(_, idx) => String(idx)}
                size="small"
                scroll={{ x: 600, y: 300 }}
                pagination={false}
              />
            </div>
          )}
        </div>
      )}
    </Modal>
  );
};

export default ImportModal;
