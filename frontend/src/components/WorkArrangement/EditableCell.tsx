import React from 'react';
import { Input, InputNumber, DatePicker, Select } from 'antd';
import dayjs from 'dayjs';
import type { WorkArrangement } from '../../types/workArrangement';
import { WORK_TYPES, LOCATIONS, PROGRESSES } from '../../types/workArrangement';

export type EditableField =
  | 'id' | 'date' | 'customer' | 'project' | 'work_type'
  | 'location' | 'partner' | 'content' | 'duration' | 'progress' | 'notes';

interface EditableCellProps {
  field: EditableField;
  value: any;
  editing: boolean;
  onChange: (value: any) => void;
  style?: React.CSSProperties;
}

const EditableCell: React.FC<EditableCellProps> = ({ field, value, editing, onChange, style }) => {
  if (!editing) {
    // Render display value when not editing
    if (field === 'date') {
      return <span style={style}>{value || '-'}</span>;
    }
    if (field === 'duration') {
      return <span style={style}>{value > 0 ? value : '-'}</span>;
    }
    return <span style={style}>{value ?? '-'}</span>;
  }

  switch (field) {
    case 'date':
      return (
        <DatePicker
          value={value ? dayjs(value) : dayjs()}
          onChange={(d) => onChange(d ? d.format('YYYY-MM-DD') : '')}
          format="YYYY-MM-DD"
          style={{ width: '100%' }}
          size="small"
        />
      );
    case 'id':
      return (
        <Input
          value={value ?? ''}
          onChange={(e) => onChange(e.target.value)}
          size="small"
          style={{ width: '100%' }}
        />
      );
    case 'work_type':
      return (
        <Select
          value={value || '售前'}
          onChange={onChange}
          options={WORK_TYPES.map((t) => ({ label: t, value: t }))}
          size="small"
          style={{ width: '100%' }}
        />
      );
    case 'location':
      return (
        <Select
          value={value || '远程'}
          onChange={onChange}
          options={LOCATIONS.map((l) => ({ label: l, value: l }))}
          size="small"
          style={{ width: '100%' }}
        />
      );
    case 'partner':
      return (
        <Select
          value={value || '否'}
          onChange={onChange}
          options={[{ label: '是', value: '是' }, { label: '否', value: '否' }]}
          size="small"
          style={{ width: '100%' }}
        />
      );
    case 'progress':
      return (
        <Select
          value={value || '未开始'}
          onChange={onChange}
          options={PROGRESSES.map((p) => ({ label: p, value: p }))}
          size="small"
          style={{ width: '100%' }}
        />
      );
    case 'duration':
      return (
        <InputNumber
          value={value || 0}
          onChange={(v) => onChange(v ?? 0)}
          min={0}
          max={24}
          step={0.5}
          size="small"
          style={{ width: '100%' }}
        />
      );
    case 'customer':
    case 'project':
    case 'content':
    case 'notes':
      return (
        <Input
          value={value ?? ''}
          onChange={(e) => onChange(e.target.value)}
          size="small"
          style={{ width: '100%' }}
        />
      );
    default:
      return <span>{value}</span>;
  }
};

export default EditableCell;
