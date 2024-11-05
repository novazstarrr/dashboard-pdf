import axios from '../lib/axios';

export const authService = {
  login: async (email, password) => {
    try {
      const response = await axios.post(`/api/login`, {
        email,
        password,
      });
      if (response.data.token) {
        localStorage.setItem('token', response.data.token);
      }
      return response.data;
    } catch (error) {
      if (error.response) {
        throw new Error(error.response.data.message || 'An error occurred');
      } else if (error.request) {
        throw new Error('No response from server');
      } else {
        throw new Error('An error occurred');
      }
    }
  },

  register: async (userData) => {
    try {
      const response = await axios.post(`/api/register`, userData);
      return response.data;
    } catch (error) {
      if (error.response) {
        throw new Error(error.response.data.message || 'An error occurred');
      } else if (error.request) {
        throw new Error('No response from server');
      } else {
        throw new Error('An error occurred');
      }
    }
  },

  getCurrentUser: async () => {
    try {
      const token = localStorage.getItem('token');
      if (!token) return null;

      const response = await axios.get('/api/users/me', {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      return response.data;
    } catch (error) {
      if (error.response?.status === 401 || error.response?.status === 400) {
        localStorage.removeItem('token');
        return null;
      }
      throw error;
    }
  },

  logout: () => {
    localStorage.removeItem('token');
  },

  getToken: () => {
    return localStorage.getItem('token');
  }
};
