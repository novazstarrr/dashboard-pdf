import { useReducer } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from './useAuth';
import { handleLoginError } from '../utils/errorHandling';

const initialState = {
  email: '',
  password: '',
  error: '',
  loading: false
};

const formReducer = (state, action) => {
  switch (action.type) {
    case 'SET_FIELD':
      return { ...state, [action.field]: action.value };
    case 'SET_ERROR':
      return { ...state, error: action.payload };
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'RESET_FORM':
      return initialState;
    default:
      return state;
  }
};

export const useLoginForm = () => {
  const [state, dispatch] = useReducer(formReducer, initialState);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    const { email, password } = state;

    try {
      dispatch({ type: 'SET_ERROR', payload: '' });
      dispatch({ type: 'SET_LOADING', payload: true });
      
      if (!email || !password) {
        dispatch({ type: 'SET_ERROR', payload: 'Please enter both email and password' });
        return;
      }

      await login({ email, password });
      navigate('/dashboard');
    } catch (err) {
      handleLoginError(err, dispatch);
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  };

  return {
    state,
    dispatch,
    handleSubmit
  };
}; 