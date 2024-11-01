
import PropTypes from 'prop-types';
import { Form } from 'react-bootstrap';


const SearchBar = ({ searchTerm, onSearch }) => (
  <Form.Group className="mb-3">
    <Form.Control
      type="text"
      placeholder="Search files by name, type, or date..."
      value={searchTerm}
      onChange={(e) => onSearch(e.target.value)}
    />
  </Form.Group>
);

SearchBar.propTypes = {
  searchTerm: PropTypes.string.isRequired,
  onSearch: PropTypes.func.isRequired,
};

export default SearchBar;
