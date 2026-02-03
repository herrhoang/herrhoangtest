import React, { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, message, Card, DatePicker } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { Transaction, Account, Category } from '../types';
import { transactionApi, accountApi, categoryApi } from '../services/api';

const { Option } = Select;

const TransactionPage: React.FC = () => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();
  const transactionType = Form.useWatch('type', form);

  const fetchData = async () => {
    try {
      const [transactionsRes, accountsRes, categoriesRes] = await Promise.all([
        transactionApi.getAll(),
        accountApi.getAll(),
        categoryApi.getAll(),
      ]);
      setTransactions(Array.isArray(transactionsRes) ? transactionsRes : []);
      setAccounts(Array.isArray(accountsRes) ? accountsRes : []);
      setCategories(Array.isArray(categoriesRes) ? categoriesRes : []);
    } catch (error) {
      message.error('获取数据失败');
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      const transactionData = {
        account_id: values.account_id,
        category_id: values.category_id,
        amount: parseFloat(values.amount),
        type: values.type,
        description: values.note || '',
        date: values.date ? values.date.format('YYYY-MM-DD') : undefined,
      };
      await transactionApi.create(transactionData);
      message.success('交易记录创建成功');
      setIsModalVisible(false);
      form.resetFields();
      fetchData();
    } catch (error) {
      message.error('创建失败');
    }
  };

  const columns = [
    {
      title: '日期',
      dataIndex: 'date',
      key: 'date',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (type === 'income' ? '收入' : '支出'),
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number, record: Transaction) => {
        const color = record.type === 'income' ? '#52c41a' : '#f5222d';
        return <span style={{ color }}>{`¥${amount.toFixed(2)}`}</span>;
      },
    },
    {
      title: '账户',
      dataIndex: 'account_id',
      key: 'account_id',
      render: (accountId: number) => {
        const account = accounts.find(a => a.id === accountId);
        return account?.name || '-';
      },
    },
    {
      title: '分类',
      dataIndex: 'category_id',
      key: 'category_id',
      render: (categoryId: number) => {
        const category = categories.find(c => c.id === categoryId);
        return category?.name || '-';
      },
    },
    {
      title: '备注',
      dataIndex: 'description',
      key: 'description',
    },
  ];

  return (
    <div>
      <Card
        title="交易记录"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setIsModalVisible(true)}>
            添加交易
          </Button>
        }
      >
        <Table
          dataSource={transactions}
          columns={columns}
          rowKey="id"
          pagination={{ pageSize: 10 }}
        />
      </Card>

      <Modal
        title="添加交易"
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={() => {
          setIsModalVisible(false);
          form.resetFields();
        }}
      >
        <Form
          form={form}
          layout="vertical"
          onValuesChange={(changedValues) => {
            if ('type' in changedValues) {
              form.setFieldValue('category_id', undefined);
            }
          }}
        >
          <Form.Item
            name="type"
            label="交易类型"
            rules={[{ required: true, message: '请选择交易类型' }]}
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
            <Input type="number" prefix="¥" />
          </Form.Item>
          <Form.Item
            name="account_id"
            label="账户"
            rules={[{ required: true, message: '请选择账户' }]}
          >
            <Select>
              {accounts.map(account => (
                <Option key={account.id} value={account.id}>
                  {account.name}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="category_id"
            label="分类"
            rules={[{ required: true, message: '请选择分类' }]}
          >
            <Select placeholder={transactionType ? '请选择分类' : '请先选择交易类型'}>
              {categories
                .filter(c => !transactionType || c.type === transactionType)
                .map(category => (
                  <Option key={category.id} value={category.id}>
                    {category.name}
                  </Option>
                ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="date"
            label="日期"
            initialValue={dayjs()}
            rules={[{ required: true, message: '请选择日期' }]}
          >
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="note" label="备注">
            <Input.TextArea />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default TransactionPage;
