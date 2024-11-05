import { Form, Button, Container, Row, Col, Card } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { useLoginForm } from '../hooks/useLoginForm';
import { ErrorAlert } from './common/ErrorAlert';

export const LoginForm = () => {
  const { state, dispatch, handleSubmit } = useLoginForm();
  const { email, password, error, loading } = state;

  return (
    <div className="min-vh-100 d-flex align-items-center justify-content-center" style={{ background: '#f5f5f5', padding: '20px' }}>
      <Container>
        <Row className="justify-content-center">
          <Col xs={12} md={8} lg={5}>
            <Card className="shadow-lg" style={{ maxWidth: '400px', margin: '0 auto' }}>
              <Card.Body className="p-4">
                <h2 className="text-center mb-4">Login</h2>
                <ErrorAlert error={error} />
                <Form onSubmit={handleSubmit}>
                  <Form.Group className="mb-3">
                    <Form.Label>Email</Form.Label>
                    <Form.Control
                      type="email"
                      value={email}
                      onChange={(e) => dispatch({ 
                        type: 'SET_FIELD', 
                        field: 'email', 
                        value: e.target.value 
                      })}
                      required
                      placeholder="Enter your email"
                      size="lg"
                      isInvalid={error && error.includes('email')}
                    />
                    <Form.Control.Feedback type="invalid">
                      Please enter a valid email address
                    </Form.Control.Feedback>
                  </Form.Group>
                  
                  <Form.Group className="mb-4">
                    <Form.Label>Password</Form.Label>
                    <Form.Control
                      type="password"
                      value={password}
                      onChange={(e) => dispatch({ 
                        type: 'SET_FIELD', 
                        field: 'password', 
                        value: e.target.value 
                      })}
                      required
                      placeholder="Enter your password"
                      size="lg"
                      isInvalid={error && error.includes('password')}
                    />
                  </Form.Group>

                  <Button
                    className="w-100 mb-3"
                    type="submit"
                    disabled={loading}
                    size="lg"
                  >
                    {loading ? (
                      <>
                        <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                        Logging in...
                      </>
                    ) : (
                      'Login'
                    )}
                  </Button>

                  <div className="text-center">
                    <small className="text-muted">
                      Don&apos;t have an account? <Link to="/register">Register</Link>
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
