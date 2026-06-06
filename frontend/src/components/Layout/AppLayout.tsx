import React from 'react';
import { Layout, Menu, Typography } from 'antd';
import {
  ScheduleOutlined,
  SettingOutlined,
} from '@ant-design/icons';

const { Sider, Content, Header } = Layout;
const { Title } = Typography;

interface AppLayoutProps {
  children: React.ReactNode;
  activeKey: string;
  onMenuClick: (key: string) => void;
}

const AppLayout: React.FC<AppLayoutProps> = ({ children, activeKey, onMenuClick }) => {
  const menuItems = [
    {
      key: 'work-arrangements',
      icon: <ScheduleOutlined />,
      label: '工作安排',
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置',
    },
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        breakpoint="lg"
        collapsedWidth="60"
        theme="light"
        style={{
          borderRight: '1px solid #f0f0f0',
          boxShadow: '2px 0 8px rgba(0,0,0,0.04)',
        }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            borderBottom: '1px solid #f0f0f0',
          }}
        >
          <Title
            level={5}
            style={{ margin: 0, color: '#1677ff', fontWeight: 700, letterSpacing: 1 }}
          >
            📋 工作管理
          </Title>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[activeKey]}
          items={menuItems}
          onClick={({ key }) => onMenuClick(key)}
          style={{ borderRight: 0, marginTop: 8 }}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            borderBottom: '1px solid #f0f0f0',
            height: 64,
            boxShadow: '0 1px 4px rgba(0,0,0,0.04)',
          }}
        >
          <Title level={4} style={{ margin: 0, fontWeight: 500 }}>
            工作安排管理系统
          </Title>
        </Header>
        <Content
          style={{
            margin: 24,
            padding: 24,
            background: '#fff',
            borderRadius: 8,
            minHeight: 280,
            boxShadow: '0 1px 4px rgba(0,0,0,0.04)',
          }}
        >
          {children}
        </Content>
      </Layout>
    </Layout>
  );
};

export default AppLayout;
