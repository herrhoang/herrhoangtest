import React, { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Input, message, Card, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { Account } from '../types';
import { accountApi } from '../services/api';

const AccountPage: React.FC = () => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [editingAccount, setEditingAccount] = useState<Account | null>(null);
  const [form] = Form.useForm();

  const fetchAccounts = async () => {
    try {
      const data = await accountApi.getAll();
      setAccounts(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Failed to fetch accounts:', error);
      message.error('获取账户列表失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await accountApi.delete(id);
      message.success('账户删除成功');
      fetchAccounts();
    } catch (error: any) {
      message.error(error?.response?.data?.error || '删除失败');
    }
  };

  useEffect(() => {
    fetchAccounts();
  }, []);

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      console.log('Form values:', values);

      // 确保余额是数字类型
      const balance = parseFloat(values.balance);
      if (isNaN(balance)) {
        throw new Error('请输入有效的余额');
      }

      if (editingAccount) {
        const updateResponse = await accountApi.update(editingAccount.id, {
          ...values,
          balance: balance
        });
        console.log('Update response:', updateResponse);
        message.success('账户更新成功');
      } else {
        const createResponse = await accountApi.create({
          name: values.name,
          balance: balance
        });
        console.log('Create response:', createResponse);
        message.success('账户创建成功');
      }
      setIsModalVisible(false);
      form.resetFields();
      setEditingAccount(null);
      fetchAccounts();
    } catch (error) {
      console.error('Operation failed:', error);
      message.error(error instanceof Error ? error.message : '操作失败');
    }
  };

  const handleEdit = (record: Account) => {
    setEditingAccount(record);
    form.setFieldsValue(record);
    setIsModalVisible(true);
  };

  const handleAdd = () => {
    setEditingAccount(null);
    form.resetFields();
    setIsModalVisible(true);
  };

  const columns = [
    {
      title: '账户名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '账户余额',
      dataIndex: 'balance',
      key: 'balance',
      render: (balance: number) => `¥${balance.toFixed(2)}`,
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Account) => (
        <>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除此账户吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </>
      ),
    },
  ];

  return (
    <div>
      <Card
        title="账户管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            添加账户
          </Button>
        }
      >
        <Table
          dataSource={accounts}
          columns={columns}
          rowKey="id"
          pagination={false}
        />
      </Card>

      <Modal
        title={editingAccount ? '编辑账户' : '添加账户'}
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={() => {
          setIsModalVisible(false);
          form.resetFields();
          setEditingAccount(null);
        }}
      >
        <Form form={form} layout="vertical" initialValues={{ type: 'bank' }}>
          <Form.Item
            name="name"
            label="账户名称"
            rules={[{ required: true, message: '请输入账户名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="balance"
            label="账户余额"
            rules={[
              { required: true, message: '请输入账户余额' },
              {
                validator: async (_, value) => {
                  if (isNaN(value) || value === '') {
                    throw new Error('请输入有效的数字');
                  }
                  const numValue = parseFloat(value);
                  if (numValue < 0) {
                    throw new Error('余额不能为负数');
                  }
                }
              }
            ]}
          >
            <Input 
              type="number" 
              prefix="¥" 
              step="0.01"
            />
          </Form.Item>
          <Form.Item
            name="type"
            label="账户类型"
            hidden
          >
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default AccountPage;
