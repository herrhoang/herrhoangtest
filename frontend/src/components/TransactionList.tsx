import React, { useEffect, useState } from 'react';
import { Table, Card, Button, Modal, Form, Input, Select, message } from 'antd';
import { Transaction, Account, Category } from '../types';
import { transactionApi, accountApi, categoryApi } from '../services/api';

const { Option } = Select;

const TransactionList: React.FC = () => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();

  const fetchData = async () => {
    try {
      const [transactionsRes, accountsRes, categoriesRes] = await Promise.all([
        transactionApi.getAll(),
        accountApi.getAll(),
        categoryApi.getAll(),
      ]);
      setTransactions(transactionsRes.data);
      setAccounts(accountsRes.data);
      setCategories(categoriesRes.data);
    } catch (error) {
      message.error('获取数据失败');
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleCreate = async (values: any) => {
    try {
      await transactionApi.create(values);
      message.success('创建交易记录成功');
      setIsModalVisible(false);
      form.resetFields();
      fetchData();
    } catch (error) {
      message.error('创建交易记录失败');
    }
  };

  const columns = [
    {
      title: '账户',
      dataIndex: 'account_id',
      key: 'account_id',
      render: (accountId: number) => 
        accounts.find(a => a.id === accountId)?.name || '未知账户',
    },
    {
      title: '类别',
      dataIndex: 'category_id',
      key: 'category_id',
      render: (categoryId: number) =>
        categories.find(c => c.id === categoryId)?.name || '未知类别',
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number, record: Transaction) => 
        `${record.type === 'expense' ? '-' : '+'}¥${amount.toFixed(2)}`,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString(),
    },
  ];

  return (
    <Card
      title="交易记录"
      extra={<Button type="primary" onClick={() => setIsModalVisible(true)}>新建交易</Button>}
    >
      <Table
        dataSource={transactions}
        columns={columns}
        rowKey="id"
      />
      
      <Modal
        title="新建交易"
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        onOk={form.submit}
      >
        <Form form={form} onFinish={handleCreate}>
          <Form.Item
            name="account_id"
            label="账户"
            rules={[{ required: true, message: '请选择账户' }]}
          >
            <Select>
              {accounts.map(account => (
                <Option key={account.id} value={account.id}>{account.name}</Option>
              ))}
            </Select>
          </Form.Item>
          
          <Form.Item
            name="category_id"
            label="类别"
            rules={[{ required: true, message: '请选择类别' }]}
          >
            <Select>
              {categories.map(category => (
                <Option key={category.id} value={category.id}>{category.name}</Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="type"
            label="类型"
            rules={[{ required: true, message: '请选择类型' }]}
          >
            <Select>
              <Option value="income">收入</Option>
              <Option value="expense">支出</Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="amount"
            label="金额"
            rules={[{ required: true, message: '请输入金额' }]}
          >
            <Input type="number" step="0.01" />
          </Form.Item>

          <Form.Item
            name="description"
            label="描述"
          >
            <Input.TextArea />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default TransactionList;
