import React from 'react';
import { Layout as AntLayout, Menu } from 'antd';
import { Link, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  WalletOutlined,
  TransactionOutlined,
} from '@ant-design/icons';

const { Header, Content, Sider } = AntLayout;

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const location = useLocation();

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: <Link to="/">仪表盘</Link>,
    },
    {
      key: '/accounts',
      icon: <WalletOutlined />,
      label: <Link to="/accounts">账户管理</Link>,
    },
    {
      key: '/transactions',
      icon: <TransactionOutlined />,
      label: <Link to="/transactions">交易记录</Link>,
    },
  ];

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Header style={{ padding: 0, background: '#fff' }}>
        <div style={{ float: 'left', width: 200, height: 31, margin: '16px 24px 16px 0', background: 'rgba(255, 255, 255, 0.2)' }}>
          <h1 style={{ margin: 0, color: '#1890ff', textAlign: 'center' }}>个人记账系统</h1>
        </div>
      </Header>
      <AntLayout>
        <Sider width={200} style={{ background: '#fff' }}>
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            style={{ height: '100%', borderRight: 0 }}
            items={menuItems}
          />
        </Sider>
        <AntLayout style={{ padding: '24px' }}>
          <Content style={{ background: '#fff', padding: 24, margin: 0, minHeight: 280 }}>
            {children}
          </Content>
        </AntLayout>
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;
