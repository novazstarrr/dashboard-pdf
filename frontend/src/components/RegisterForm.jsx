import { useReducer, useState } from 'react';
import { Form, Button, Alert, Container, Row, Col, Card } from 'react-bootstrap';
import { useAuth } from '../hooks/useAuth';
import { useNavigate, Link } from 'react-router-dom';
import { z } from 'zod';

const initialFormState = {
  firstName: '',
  surname: '',
  dob: '',
  email: '',
  password: '',
  confirmPassword: '',
};

const formReducer = (state, action) => {
  switch (action.type) {
    case 'SET_FIELD':
      return {
        ...state,
        [action.field]: action.value
      };
    case 'RESET_FORM':
      return initialFormState;
    default:
      return state;
  }
};

const registerSchema = z.object({
  firstName: z.string()
    .min(1, 'First name is required')
    .regex(/^[A-Za-z\s'-]+$/, {
      message: 'First name can only contain letters, spaces, hyphens, and apostrophes'
    }),
  surname: z.string()
    .min(1, 'Surname is required')
    .regex(/^[A-Za-z\s'-]+$/, {
      message: 'Surname can only contain letters, spaces, hyphens, and apostrophes'
    }),
  dob: z.string()
    .min(1, 'Date of birth is required')
    .refine((date) => {
      const birthDate = new Date(date);
      const today = new Date();
      const age = today.getFullYear() - birthDate.getFullYear();
      return age >= 13;
    }, 'You must be at least 13 years old'),
  email: z.string().email('Invalid email address'),
  password: z.string()
    .min(8, 'Password must be at least 8 characters')
    .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
    .regex(/[a-z]/, 'Password must contain at least one lowercase letter')
    .regex(/[0-9]/, 'Password must contain at least one number')
    .regex(/[^A-Za-z0-9]/, 'Password must contain at least one special character'),
  confirmPassword: z.string()
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
});

const formFields = [
  {
    name: 'firstName',
    label: 'First Name',
    type: 'text',
    placeholder: 'Enter your first name',
  },
  {
    name: 'surname',
    label: 'Surname',
    type: 'text',
    placeholder: 'Enter your surname',
  },
  {
    name: 'dob',
    label: 'Date of Birth',
    type: 'date',
  },
  {
    name: 'email',
    label: 'Email',
    type: 'email',
    placeholder: 'Enter your email',
  },
  {
    name: 'password',
    label: 'Password',
    type: 'password',
    placeholder: 'Enter your password',
  },
  {
    name: 'confirmPassword',
    label: 'Confirm Password',
    type: 'password',
    placeholder: 'Confirm your password',
  },
];

export const RegisterForm = () => {
  const [formState, dispatch] = useReducer(formReducer, initialFormState);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [validationErrors, setValidationErrors] = useState({});
  const { register } = useAuth();
  const navigate = useNavigate();

  const validateForm = () => {
    try {
      registerSchema.parse(formState);
      setValidationErrors({});
      return true;
    } catch (error) {
      const errors = {};
      error.errors.forEach((err) => {
        errors[err.path[0]] = [err.message];
      });
      setValidationErrors(errors);
      return false;
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    // eslint-disable-next-line no-unused-vars
    const { confirmPassword, ...userData } = formState;
    console.log('Sending registration data:', userData);

    try {
      setError('');
      setLoading(true);
      await register(userData);
      navigate('/login');
    } catch (err) {
      console.error('Registration error:', err);
      setError(err.message || 'Failed to register');
    } finally {
      setLoading(false);
    }
  };

  const handleFieldChange = (field, value) => {
    dispatch({ type: 'SET_FIELD', field, value });
  };

  const renderFormField = ({ name, label, type, placeholder }) => {
    const checkPasswordRequirement = (regex) => {
      if (name !== 'password' || !formState.password) return true;
      return regex.test(formState.password);
    };

    return (
      <Form.Group className="mb-3" key={name}>
        <Form.Label>{label}</Form.Label>
        <Form.Control
          type={type}
          value={formState[name]}
          onChange={(e) => handleFieldChange(name, e.target.value)}
          required
          placeholder={placeholder}
          size="lg"
          isInvalid={validationErrors[name]}
        />
        {name === 'password' && (
          <Form.Text>
            Password must contain:
            <ul className="mb-0 ps-3 mt-1">
              <li className={!checkPasswordRequirement(/.{8,}/) ? 'text-danger' : 'text-muted'}>
                At least 8 characters
              </li>
              <li className={!checkPasswordRequirement(/[A-Z]/) ? 'text-danger' : 'text-muted'}>
                One uppercase letter
              </li>
              <li className={!checkPasswordRequirement(/[a-z]/) ? 'text-danger' : 'text-muted'}>
                One lowercase letter
              </li>
              <li className={!checkPasswordRequirement(/[0-9]/) ? 'text-danger' : 'text-muted'}>
                One number
              </li>
              <li className={!checkPasswordRequirement(/[^A-Za-z0-9]/) ? 'text-danger' : 'text-muted'}>
                One special character
              </li>
            </ul>
          </Form.Text>
        )}
        {validationErrors[name] && (
          <Form.Control.Feedback type="invalid">
            {validationErrors[name].map((error, index) => (
              <div key={index}>{error}</div>
            ))}
          </Form.Control.Feedback>
        )}
      </Form.Group>
    );
  };

  return (
    <div className="min-vh-100 d-flex align-items-center justify-content-center bg-light py-4">
      <Container>
        <Row className="justify-content-center">
          <Col xs={12} md={8} lg={5}>
            <Card className="shadow-lg mx-auto" style={{ maxWidth: '400px' }}>
              <Card.Body className="p-4">
                <h2 className="text-center mb-4">Register</h2>
                {error && <Alert variant="danger">{error}</Alert>}
                <Form onSubmit={handleSubmit}>
                  {formFields.map(renderFormField)}
                  <Button
                    className="w-100 mb-3"
                    type="submit"
                    disabled={loading}
                    size="lg"
                  >
                    {loading ? 'Loading...' : 'Register'}
                  </Button>
                  <div className="text-center">
                    <small className="text-muted">
                      Already have an account? <Link to="/login">Login</Link>
                    </small>
                  </div>
                </Form>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </div>
  );
};
