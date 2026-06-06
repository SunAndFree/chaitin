import React from 'react';
import { Button, DatePicker, Input, Select, Space, Row, Col } from 'antd';
import { SearchOutlined, ClearOutlined } from '@ant-design/icons';
import type { FilterParams } from '../../types/workArrangement';
import { WORK_TYPES, PROGRESSES } from '../../types/workArrangement';

const { RangePicker } = DatePicker;

interface FilterBarProps {
  filters: FilterParams;
  onChange: (filters: FilterParams) => void;
  onSearch: () => void;
  onReset: () => void;
}

const FilterBar: React.FC<FilterBarProps> = ({ filters, onChange, onSearch, onReset }) => {
  const updateFilter = (key: keyof FilterParams, value: any) => {
    onChange({ ...filters, [key]: value || '' });
  };

  return (
    <div
      style={{
        padding: '16px 0',
        borderBottom: '1px solid #f0f0f0',
        marginBottom: 16,
      }}
    >
      <Row gutter={[12, 12]} align="middle">
        <Col xs={24} sm={12} md={8} lg={5}>
          <RangePicker
            value={
              filters.date_from && filters.date_to
                ? [filters.date_from as any, filters.date_to as any]
                : null
            }
            onChange={(dates: any) => {
              if (dates && dates[0] && dates[1]) {
                onChange({
                  ...filters,
                  date_from: dates[0].format('YYYY-MM-DD'),
                  date_to: dates[1].format('YYYY-MM-DD'),
                });
              } else {
                const { date_from, date_to, ...rest } = filters;
                onChange(rest);
              }
            }}
            placeholder={['开始日期', '结束日期']}
            style={{ width: '100%' }}
            allowClear
          />
        </Col>
        <Col xs={24} sm={12} md={8} lg={4}>
          <Input
            placeholder="客户名称"
            value={filters.customer || ''}
            onChange={(e) => updateFilter('customer', e.target.value)}
            allowClear
            prefix={<SearchOutlined style={{ color: '#bfbfbf' }} />}
          />
        </Col>
        <Col xs={24} sm={12} md={8} lg={4}>
          <Input
            placeholder="项目名称"
            value={filters.project || ''}
            onChange={(e) => updateFilter('project', e.target.value)}
            allowClear
            prefix={<SearchOutlined style={{ color: '#bfbfbf' }} />}
          />
        </Col>
        <Col xs={12} sm={6} md={4} lg={3}>
          <Select
            placeholder="工作类型"
            value={filters.work_type || undefined}
            onChange={(v) => updateFilter('work_type', v)}
            allowClear
            style={{ width: '100%' }}
            options={WORK_TYPES.map((t) => ({ label: t, value: t }))}
          />
        </Col>
        <Col xs={12} sm={6} md={4} lg={3}>
          <Select
            placeholder="工作进度"
            value={filters.progress || undefined}
            onChange={(v) => updateFilter('progress', v)}
            allowClear
            style={{ width: '100%' }}
            options={PROGRESSES.map((p) => ({ label: p, value: p }))}
          />
        </Col>
        <Col xs={24} sm={12} md={8} lg={5}>
          <Space>
            <Button type="primary" icon={<SearchOutlined />} onClick={onSearch}>
              查询
            </Button>
            <Button icon={<ClearOutlined />} onClick={onReset}>
              重置
            </Button>
          </Space>
        </Col>
      </Row>
    </div>
  );
};

export default FilterBar;
