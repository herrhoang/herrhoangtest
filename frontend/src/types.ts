export interface Account {
  id: number;
  name: string;
  balance: number;
  type: string;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  id: number;
  account_id: number;
  category_id: number;
  amount: number;
  type: 'income' | 'expense';
  description: string;
  date: string;
  created_at: string;
}

export interface Category {
  id: number;
  name: string;
  type: 'income' | 'expense';
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}
