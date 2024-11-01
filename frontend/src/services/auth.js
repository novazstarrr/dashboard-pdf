import axios from '../lib/axios';

export const authService = {
  login: async (email, password) => {
    try {
      const response = await axios.post(`/login`, {
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
      const response = await axios.post(`/register`, userData);
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

  logout: () => {
    localStorage.removeItem('token');
  },

  getToken: () => {
    return localStorage.getItem('token');
  }
};
