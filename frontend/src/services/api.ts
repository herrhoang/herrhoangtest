import axios from 'axios';
import { Account, Transaction, Category } from '../types';

interface ApiResponse<T> {
  data: T;
}

const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

interface AccountAPI {
  getAll: () => Promise<ApiResponse<Account[]>>;
  create: (data: Partial<Account>) => Promise<ApiResponse<Account>>;
  update: (id: number, data: Partial<Account>) => Promise<ApiResponse<Account>>;
}

interface TransactionAPI {
  getAll: () => Promise<ApiResponse<Transaction[]>>;
  create: (data: Partial<Transaction>) => Promise<ApiResponse<Transaction>>;
}

interface CategoryAPI {
  getAll: () => Promise<ApiResponse<Category[]>>;
  create: (data: Partial<Category>) => Promise<ApiResponse<Category>>;
}

interface APIService {
  accountApi: AccountAPI;
  transactionApi: TransactionAPI;
  categoryApi: CategoryAPI;
}

export const accountApi: AccountAPI = {
  getAll: async () => {
    try {
      const response = await api.get<ApiResponse<Account[]>>('/accounts');
      return response.data;
    } catch (error) {
      console.error('Error fetching accounts:', error);
      throw error;
    }
  },
  create: async (data) => {
    try {
      const response = await api.post<ApiResponse<Account>>('/accounts', data);
      return response.data;
    } catch (error) {
      console.error('Error creating account:', error);
      throw error;
    }
  },
  update: async (id, data) => {
    try {
      const response = await api.put<ApiResponse<Account>>(`/accounts/${id}`, data);
      return response.data;
    } catch (error) {
      console.error('Error updating account:', error);
      throw error;
    }
  }
};

export const transactionApi: TransactionAPI = {
  getAll: async () => {
    try {
      const response = await api.get<ApiResponse<Transaction[]>>('/transactions');
      return response.data;
    } catch (error) {
      console.error('Error fetching transactions:', error);
      throw error;
    }
  },
  create: async (data) => {
    try {
      const response = await api.post<ApiResponse<Transaction>>('/transactions', data);
      return response.data;
    } catch (error) {
      console.error('Error creating transaction:', error);
      throw error;
    }
  }
};

export const categoryApi: CategoryAPI = {
  getAll: async () => {
    try {
      const response = await api.get<ApiResponse<Category[]>>('/categories');
      return response.data;
    } catch (error) {
      console.error('Error fetching categories:', error);
      throw error;
    }
  },
  create: async (data) => {
    try {
      const response = await api.post<ApiResponse<Category>>('/categories', data);
      return response.data;
    } catch (error) {
      console.error('Error creating category:', error);
      throw error;
    }
  }
};

const apiService: APIService = {
  accountApi,
  transactionApi,
  categoryApi
};

export default apiService;

