import { Navigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import PropTypes from 'prop-types';
import { useEffect } from 'react';

export const PrivateRoute = ({ children }) => {
  const { user, loading, refreshUser } = useAuth();

  useEffect(() => {
    if (!user && !loading) {
      refreshUser();
    }
  }, [user, loading, refreshUser]);

  if (loading) {
    return (
      <div className="d-flex justify-content-center align-items-center min-vh-100">
        <div className="spinner-border text-primary" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" />;
  }

  return children;
};

PrivateRoute.propTypes = {
  children: PropTypes.node.isRequired,
};
