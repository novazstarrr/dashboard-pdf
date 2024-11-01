export const handleLoginError = (err, dispatch) => {
  console.error('Login error details:', {
    status: err.response?.status,
    data: err.response?.data,
    message: err.response?.data?.message,
    error: err.message
  });
  
  const errorMessage = err.response?.data?.message || err.message || 'An unknown error occurred';
  
  if (!err.response) {
    dispatch({ type: 'SET_ERROR', payload: 'Network error. Please check your connection.' });
    return;
  }

  switch (err.response.status) {
    case 400:
      dispatch({ type: 'SET_ERROR', payload: errorMessage || 'Invalid email or password format' });
      break;
    case 401:
      dispatch({ type: 'SET_ERROR', payload: 'Invalid email or password' });
      break;
    case 404:
      dispatch({ type: 'SET_ERROR', payload: 'No account exists with this email. Would you like to register?' });
      break;
    case 429:
      dispatch({ type: 'SET_ERROR', payload: 'Too many failed attempts. Please try again in a few minutes.' });
      break;
    case 500:
      dispatch({ type: 'SET_ERROR', payload: 'Server error. Please try again later.' });
      break;
    default:
      dispatch({ type: 'SET_ERROR', payload: `Login failed: ${errorMessage}` });
  }
}; 