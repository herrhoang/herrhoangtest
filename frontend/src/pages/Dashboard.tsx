import React, { useEffect, useState } from 'react';
import { Card, Statistic, Row, Col } from 'antd';
import { Account, Transaction } from '../types';
import { accountApi, transactionApi } from '../services/api';

const Dashboard: React.FC = () => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [transactions, setTransactions] = useState<Transaction[]>([]);

  const fetchData = async () => {
    try {
      const [accountsRes, transactionsRes] = await Promise.all([
        accountApi.getAll(),
        transactionApi.getAll(),
      ]);
      setAccounts(Array.isArray(accountsRes) ? accountsRes : []);
      setTransactions(Array.isArray(transactionsRes) ? transactionsRes : []);
    } catch (error) {
      console.error('获取数据失败:', error);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const totalBalance = accounts.reduce((sum, account) => sum + account.balance, 0);
  const totalIncome = transactions
    .filter(t => t.type === 'income')
    .reduce((sum, t) => sum + t.amount, 0);
  const totalExpense = transactions
    .filter(t => t.type === 'expense')
    .reduce((sum, t) => sum + t.amount, 0);

  return (
    <div>
      <h2>财务概览</h2>
      <Row gutter={16}>
        <Col span={8}>
          <Card>
            <Statistic
              title="总资产"
              value={totalBalance}
              precision={2}
              prefix="¥"
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="总收入"
              value={totalIncome}
              precision={2}
              prefix="¥"
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="总支出"
              value={totalExpense}
              precision={2}
              prefix="¥"
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
      </Row>
      {accounts.length > 0 && (
        <Card title="金融账户" style={{ marginTop: 24 }}>
          <Row gutter={16}>
            {accounts.map((account) => (
              <Col key={account.id} span={8}>
                <Card size="small">
                  <Statistic
                    title={account.name}
                    value={account.balance}
                    precision={2}
                    prefix="¥"
                  />
                </Card>
              </Col>
            ))}
          </Row>
        </Card>
      )}
    </div>
  );
};

export default Dashboard;
