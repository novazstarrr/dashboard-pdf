import { useState } from 'react';
import PropTypes from 'prop-types';
import { Form, InputGroup, Button } from 'react-bootstrap';
import { FaSearch } from 'react-icons/fa';

const SearchBar = ({ searchTerm, onSearch }) => {
  const [showAdvanced, setShowAdvanced] = useState(false);

  return (
    <Form.Group className="mb-3">
      <InputGroup>
        <Form.Control
          type="text"
          placeholder="Search files..."
          value={searchTerm}
          onChange={(e) => onSearch(e.target.value)}
        />
        
        <Button 
          variant="outline-secondary"
          onClick={() => setShowAdvanced(!showAdvanced)}
        >
          <FaSearch />
        </Button>
      </InputGroup>

      {showAdvanced && (
        <div className="mt-2 p-2 border rounded">
          <small className="text-muted">
            Search tips:
            <ul className="mb-0">
              <li>Search by file name, size (e.g., {'500KB'}, {'1.5MB'}), or type</li>
              <li>Use dates in YYYY-MM-DD format or relative terms (e.g., {'today'}, {'last week'})</li>
            </ul>
          </small>
        </div>
      )}
    </Form.Group>
  );
};

SearchBar.propTypes = {
  searchTerm: PropTypes.string.isRequired,
  onSearch: PropTypes.func.isRequired,
};

export default SearchBar;
