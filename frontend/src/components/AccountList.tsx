import React, { useEffect, useState } from 'react';
import { Table, Card, Button, Modal, Form, Input, message } from 'antd';
import { Account } from '../types';
import apiService from '../services/api';
const { accountApi } = apiService;

const AccountList: React.FC = () => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();

  const fetchAccounts = async () => {
    try {
      const response = await accountApi.getAll();
      setAccounts(response.data);
    } catch (error) {
      message.error('获取账户列表失败');
    }
  };

  useEffect(() => {
    fetchAccounts();
  }, []);

  const handleCreate = async (values: any) => {
    try {
      await accountApi.create(values);
      message.success('创建账户成功');
      setIsModalVisible(false);
      form.resetFields();
      fetchAccounts();
    } catch (error) {
      message.error('创建账户失败');
    }
  };

  const columns = [
    {
      title: '账户名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '余额',
      dataIndex: 'balance',
      key: 'balance',
      render: (balance: number) => `¥${balance.toFixed(2)}`,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
  ];

  return (
    <Card
      title="账户列表"
      extra={<Button type="primary" onClick={() => setIsModalVisible(true)}>新建账户</Button>}
    >
      <Table
        dataSource={accounts}
        columns={columns}
        rowKey="id"
      />
      
      <Modal
        title="新建账户"
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        onOk={form.submit}
      >
        <Form form={form} onFinish={handleCreate}>
          <Form.Item
            name="name"
            label="账户名称"
            rules={[{ required: true, message: '请输入账户名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="balance"
            label="初始余额"
            rules={[{ required: true, message: '请输入初始余额' }]}
          >
            <Input type="number" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default AccountList;
