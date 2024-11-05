import { useState, useEffect, useMemo } from 'react';
import PropTypes from 'prop-types';
import { authService } from '../services/auth';
import { AuthContext } from './authContextValue';

export { AuthContext };

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const refreshUser = async () => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        setUser(null);
        setLoading(false);
        return false;
      }

      const userData = await authService.getCurrentUser();
      if (userData) {
        setUser(userData);
        return true;
      }
      
      setUser(null);
      return false;
    } catch (error) {
      console.error('Failed to refresh user:', error);
      setUser(null);
      return false;
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
  };

  useEffect(() => {
    refreshUser();
  }, []);

  useEffect(() => {
    const interval = setInterval(refreshUser, 5 * 60 * 1000);
    return () => clearInterval(interval);
  }, []);

  const value = useMemo(() => ({
    user,
    setUser,
    loading,
    refreshUser,
    logout,
    isAuthenticated: !!user
  }), [user, loading]);

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

AuthProvider.propTypes = {
  children: PropTypes.node.isRequired,
};
