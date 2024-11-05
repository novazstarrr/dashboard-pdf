// src/hooks/useAuth.js
import { useContext } from 'react';
import { AuthContext } from '../context/authContextValue';
import axios from '../lib/axios';

export const useAuth = () => {
  const context = useContext(AuthContext);
  
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }

  const { user, setUser, loading, refreshUser, logout } = context;

  const register = async (userData) => {
    try {
      const response = await axios.post('/api/register', {
        email: userData.email,
        password: userData.password,
        firstName: userData.firstName,
        surname: userData.surname,
        dob: userData.dob,
      });
      return response.data;
    } catch (error) {
      console.error('Registration error:', error);
      throw new Error(error.response?.data?.message || 'Registration failed');
    }
  };

  const login = async ({ email, password }) => {
    try {
      const response = await axios.post('/api/login', {
        email,
        password
      });
      
      if (response.data) {
        localStorage.setItem('token', response.data.token);
        setUser(response.data.user);
      }
      
      return response.data;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  };

  return { 
    user,
    login,
    logout,
    register,
    loading,
    refreshUser,
    isAuthenticated: !!user
  };
};

