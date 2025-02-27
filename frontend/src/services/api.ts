import axios from 'axios';
import { Account, Transaction, Category } from '../types';

const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

interface AccountAPI {
  getAll: () => Promise<Account[]>;
  create: (data: Partial<Account>) => Promise<Account>;
  update: (id: number, data: Partial<Account>) => Promise<Account>;
}

interface TransactionAPI {
  getAll: () => Promise<Transaction[]>;
  create: (data: Partial<Transaction>) => Promise<Transaction>;
}

interface CategoryAPI {
  getAll: () => Promise<Category[]>;
  create: (data: Partial<Category>) => Promise<Category>;
}

interface APIService {
  accountApi: AccountAPI;
  transactionApi: TransactionAPI;
  categoryApi: CategoryAPI;
}

export const accountApi: AccountAPI = {
  getAll: () => api.get('/accounts').then(res => res.data),
  create: (data) => api.post('/accounts', data).then(res => res.data),
  update: (id, data) => api.put(`/accounts/${id}`, data).then(res => res.data)
};

export const transactionApi: TransactionAPI = {
  getAll: () => api.get('/transactions').then(res => res.data),
  create: (data) => api.post('/transactions', data).then(res => res.data)
};

export const categoryApi: CategoryAPI = {
  getAll: () => api.get('/categories').then(res => res.data),
  create: (data) => api.post('/categories', data).then(res => res.data)
};

const apiService: APIService = {
  accountApi,
  transactionApi,
  categoryApi
};

export default apiService;

