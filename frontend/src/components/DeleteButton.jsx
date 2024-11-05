import PropTypes from 'prop-types';
import { Button } from 'react-bootstrap';
import { FaTrashAlt } from 'react-icons/fa';
import useDelete from '../hooks/useDelete';

const DeleteButton = ({ fileId, onDelete }) => {
  const { deleteFile, isDeleting } = useDelete();

  const handleDelete = () => {
    deleteFile(fileId, onDelete);
  };

  return (
    <Button
      variant="outline-danger"
      size="sm"
      onClick={handleDelete}
      disabled={isDeleting}
      className="cursor-pointer"
    >
      <FaTrashAlt className="me-1" />
      {isDeleting ? 'Deleting...' : 'Delete'}
    </Button>
  );
};

DeleteButton.propTypes = {
  fileId: PropTypes.oneOfType([PropTypes.string, PropTypes.number]).isRequired,
  onDelete: PropTypes.func,
};

export default DeleteButton; 