import { Alert } from 'react-bootstrap';
import { Link } from 'react-router-dom';

export const ErrorAlert = ({ error }) => {
  if (!error) return null;
  
  return (
    <Alert variant="danger">
      <p>These credentials do not match any accounts on our system.</p>
    </Alert>
  );
}; 