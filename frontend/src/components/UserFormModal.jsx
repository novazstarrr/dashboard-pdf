import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Modal, Button, Form } from 'react-bootstrap';
import { z } from 'zod';

const userSchema = z.object({
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
  email: z.string()
    .min(1, 'Email is required')
    .email('Invalid email address'),
  dob: z.string()
    .min(1, 'Date of birth is required')
    .refine((date) => {
      const birthDate = new Date(date);
      return !isNaN(birthDate.getTime());
    }, 'Invalid date format')
});

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
}) => {
  const [validationErrors, setValidationErrors] = useState({});

  const validateField = (name, value) => {
    try {
      const fieldSchema = userSchema.shape[name];
      fieldSchema.parse(value);
      setValidationErrors(prev => ({ ...prev, [name]: null }));
    } catch (error) {
      setValidationErrors(prev => ({
        ...prev,
        [name]: error.errors[0].message
      }));
    }
  };

  const handleFieldChange = (e) => {
    const { name, value } = e.target;
    handleChange(e);
    validateField(name, value);
  };

  const handleBlur = (e) => {
    const { name, value } = e.target;
    validateField(name, value);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      userSchema.parse(userForm);
      await handleSave();
      handleClose();
    } catch (err) {
      const errors = {};
      err.errors.forEach((error) => {
        errors[error.path[0]] = error.message;
      });
      setValidationErrors(errors);
    }
  };

  return (
    <Modal show={show} onHide={handleClose} size="md">
      <Modal.Header closeButton>
        <Modal.Title>{title}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form onSubmit={handleSubmit}>
          <Form.Group controlId="formFirstName">
            <Form.Label>First Name</Form.Label>
            <Form.Control
              type="text"
              name="firstName"
              value={userForm.firstName}
              onChange={handleFieldChange}
              onBlur={handleBlur}
              placeholder="Enter first name"
              isInvalid={!!validationErrors.firstName}
            />
            {validationErrors.firstName && (
              <Form.Text className="text-danger">
                {validationErrors.firstName}
              </Form.Text>
            )}
          </Form.Group>

          <Form.Group controlId="formSurname" className="mt-2">
            <Form.Label>Surname</Form.Label>
            <Form.Control
              type="text"
              name="surname"
              value={userForm.surname}
              onChange={handleFieldChange}
              onBlur={handleBlur}
              placeholder="Enter surname"
              isInvalid={!!validationErrors.surname}
            />
            {validationErrors.surname && (
              <Form.Text className="text-danger">
                {validationErrors.surname}
              </Form.Text>
            )}
          </Form.Group>

          <Form.Group controlId="formEmail" className="mt-2">
            <Form.Label>Email</Form.Label>
            <Form.Control
              type="email"
              name="email"
              value={userForm.email}
              onChange={handleFieldChange}
              onBlur={handleBlur}
              placeholder="Enter email"
              isInvalid={!!validationErrors.email}
            />
            {validationErrors.email && (
              <Form.Text className="text-danger">
                {validationErrors.email}
              </Form.Text>
            )}
          </Form.Group>

          <Form.Group controlId="formDob" className="mt-2">
            <Form.Label>Date of Birth</Form.Label>
            <Form.Control
              type="date"
              name="dob"
              value={userForm.dob}
              onChange={handleFieldChange}
              onBlur={handleBlur}
              isInvalid={!!validationErrors.dob}
            />
            {validationErrors.dob && (
              <Form.Text className="text-danger">
                {validationErrors.dob}
              </Form.Text>
            )}
          </Form.Group>
        </Form>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleClose}>Close</Button>
        <Button 
          variant="primary" 
          onClick={handleSubmit}
          disabled={Object.keys(validationErrors).some(key => validationErrors[key])}
        >
          Save
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

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
