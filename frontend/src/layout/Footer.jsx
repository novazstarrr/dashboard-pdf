import { Container } from 'react-bootstrap';
import { 
  FaGithub, 
  FaTwitter, 
  FaLinkedin, 
  FaHeart 
} from 'react-icons/fa';

export const Footer = () => {
  return (
    <footer className="bg-light py-4 mt-auto border-top">
      <Container>
        <div className="d-flex flex-wrap justify-content-between align-items-center">
          <div className="col-md-4 d-flex align-items-center">
            <span className="text-muted">
              Made with <FaHeart className="text-danger mx-1" /> by Your Team
            </span>
          </div>
          <ul className="nav col-md-4 justify-content-end list-unstyled d-flex">
            <li className="ms-3">
              <a className="text-muted" href="#" target="_blank" rel="noopener noreferrer">
                <FaGithub size={24} />
              </a>
            </li>
            <li className="ms-3">
              <a className="text-muted" href="#" target="_blank" rel="noopener noreferrer">
                <FaTwitter size={24} />
              </a>
            </li>
            <li className="ms-3">
              <a className="text-muted" href="#" target="_blank" rel="noopener noreferrer">
                <FaLinkedin size={24} />
              </a>
            </li>
          </ul>
        </div>
      </Container>
    </footer>
  );
};
