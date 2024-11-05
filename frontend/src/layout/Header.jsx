import { Navbar, Container, Nav } from 'react-bootstrap';
import { useAuth } from '../hooks/useAuth';
import { Link, useNavigate } from 'react-router-dom';
import { 
  FaFileAlt, 
  FaSignInAlt, 
  FaUserPlus, 
  FaSignOutAlt, 
  FaUsersCog,
  FaFile
} from 'react-icons/fa';

export const Header = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = (e) => {
    e.preventDefault();
    try {
      logout();
      localStorage.removeItem('token');
      navigate('/login');
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  return (
    <Navbar bg="white" expand="lg" className="shadow-sm">
      <Container>
        <Navbar.Brand as={Link} to="/" className="text-primary">
          <FaFileAlt className="me-2" />
          PDF Manager
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="ms-auto">
            {user ? (
              <>
                <Nav.Link as={Link} to="/user-management" className="text-dark cursor-pointer">
                  <FaUsersCog className="me-1" />
                  User Management
                </Nav.Link>
                <Nav.Link as={Link} to="/dashboard" className="text-dark cursor-pointer">
                  <FaFile className="me-1" />
                  Dashboard
                </Nav.Link>
                <button 
                  onClick={handleLogout}
                  className="btn nav-link text-dark cursor-pointer border-0 bg-transparent"
                >
                  <FaSignOutAlt className="me-1" />
                  Logout
                </button>
              </>
            ) : (
              <>
                <Nav.Link as={Link} to="/login" className="text-dark">
                  <FaSignInAlt className="me-1" />
                  Login
                </Nav.Link>
                <Nav.Link as={Link} to="/register" className="text-dark">
                  <FaUserPlus className="me-1" />
                  Register
                </Nav.Link>
              </>
            )}
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};
