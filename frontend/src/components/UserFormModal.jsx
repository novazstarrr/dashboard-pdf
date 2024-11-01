import React from 'react';
import PropTypes from 'prop-types';
import { Modal, Button, Form } from 'react-bootstrap';

/**
 * Component: UserFormModal
 * 
 * A modal dialog for adding or editing a user.
 *
 * @param {Object} props - Component properties.
 * @param {boolean} props.show - Controls the visibility of the modal.
 * @param {Function} props.handleClose - Function to close the modal.
 * @param {Function} props.handleSave - Function to save the user data.
 * @param {Object} props.userForm - The form data for the user.
 * @param {Function} props.handleChange - Function to handle form input changes.
 * @param {string} props.title - The title of the modal.
 */
const UserFormModal = ({
  show,
  handleClose,
  handleSave,
  userForm,
  handleChange,
  title,
}) => (
  <Modal show={show} onHide={handleClose}>
    <Modal.Header closeButton>
      <Modal.Title>{title}</Modal.Title>
    </Modal.Header>
    <Modal.Body>
      <Form>
        <Form.Group controlId="formFirstName">
          <Form.Label>First Name</Form.Label>
          <Form.Control
            type="text"
            name="firstName"
            value={userForm.firstName}
            onChange={handleChange}
            placeholder="Enter first name"
          />
        </Form.Group>
        <Form.Group controlId="formSurname" className="mt-2">
          <Form.Label>Surname</Form.Label>
          <Form.Control
            type="text"
            name="surname"
            value={userForm.surname}
            onChange={handleChange}
            placeholder="Enter surname"
          />
        </Form.Group>
        <Form.Group controlId="formEmail" className="mt-2">
          <Form.Label>Email</Form.Label>
          <Form.Control
            type="email"
            name="email"
            value={userForm.email}
            onChange={handleChange}
            placeholder="Enter email"
          />
        </Form.Group>
        <Form.Group controlId="formDob" className="mt-2">
          <Form.Label>Date of Birth</Form.Label>
          <Form.Control
            type="date"
            name="dob"
            value={userForm.dob}
            onChange={handleChange}
            placeholder="Enter date of birth"
          />
        </Form.Group>
      </Form>
    </Modal.Body>
    <Modal.Footer>
      <Button variant="secondary" onClick={handleClose}>Close</Button>
      <Button variant="primary" onClick={handleSave}>Save</Button>
    </Modal.Footer>
  </Modal>
);


UserFormModal.propTypes = {
  show: PropTypes.bool.isRequired,
  handleClose: PropTypes.func.isRequired,
  handleSave: PropTypes.func.isRequired,
  userForm: PropTypes.shape({
    firstName: PropTypes.string.isRequired,
    surname: PropTypes.string.isRequired,
    email: PropTypes.string.isRequired,
    dob: PropTypes.string.isRequired,
  }).isRequired,
  handleChange: PropTypes.func.isRequired,
  title: PropTypes.string.isRequired,
};

export default React.memo(UserFormModal);
