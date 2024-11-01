// src/hooks/useAuth.js
import { useContext } from 'react';
import AuthContext from '../context/AuthContext';
import axios from '../lib/axios';  

export const useAuth = () => {
  const { user, setUser } = useContext(AuthContext);
  
  if (!setUser) {
    throw new Error('useAuth must be used within an AuthProvider');
  }

  const register = async (userData) => {
    try {
      console.log('Sending registration data:', userData);

      const response = await axios.post('/api/register', {
        email: userData.email,
        password: userData.password,
        firstName: userData.firstName,
        surname: userData.surname,
        dob: userData.dob,
      }, {
        withCredentials: true
      });
      
      return response.data;
    } catch (error) {
      console.error('Registration error:', error);
      throw new Error(error.response?.data?.message || 'Registration failed');
    }
  };

  const login = async ({ email, password }) => {
    try {
      const response = await axios.post(`/api/login`, {
        email,
        password
      });
      
      if (response.data) {
        localStorage.setItem('token', response.data.token);
        setUser(response.data.user);
      }
      
      return response.data;
    } catch (error) {
      console.error('Auth hook login error:', error.response || error);
      throw error;
    }
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem('token');
  };

  return { 
    user,
    login,
    logout,
    register,
  };
};

