import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import AccountPage from './pages/AccountPage';
import TransactionPage from './pages/TransactionPage';
import './App.css';

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/accounts" element={<AccountPage />} />
          <Route path="/transactions" element={<TransactionPage />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;
