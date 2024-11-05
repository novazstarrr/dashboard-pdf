import { useState, useCallback } from 'react';
import { Card, Button, Form, ListGroup, Pagination } from 'react-bootstrap';
import { toast } from 'react-toastify';
import useUsers from '../hooks/useUsers';
import UserFormModal from './UserFormModal.jsx';

export const UserManagement = () => {
  const {
    paginatedUsers,
    loading,
    error,
    currentPage,
    totalPages,
    addUser,
    editUser,
    deleteUser,
    searchTerm,
    setSearchTerm,
    onPageChange,
    fetchUsers,
  } = useUsers();

  const [selectedUser, setSelectedUser] = useState(null);
  const [showUserModal, setShowUserModal] = useState(false);
  const [userForm, setUserForm] = useState({ firstName: '', surname: '', email: '', dob: '' });

  
  const handleUserFormChange = useCallback((e) => {
    const { name, value } = e.target;
    setUserForm((prev) => ({ ...prev, [name]: value }));
  }, []);

  
  const handleUserSave = useCallback(async () => {
    try {
      if (selectedUser) {
        await editUser(selectedUser.id, userForm);
      } else {
        await addUser(userForm);
        await fetchUsers();
      }
      setShowUserModal(false);
      toast.success('User saved successfully');
    } catch (error) {
      toast.error('Error saving user: ' + error.message);
    }
  }, [selectedUser, userForm, addUser, editUser, fetchUsers]);

  
  const handleCloseModal = useCallback(() => {
    setShowUserModal(false);
  }, []);

  
  const handleOpenAddUserModal = useCallback(() => {
    setSelectedUser(null);
    setUserForm({ firstName: '', surname: '', email: '', dob: '' }); 
    setShowUserModal(true);
  }, []);

  
  const handleOpenEditUserModal = useCallback((user) => {
    setSelectedUser(user);
    setUserForm({
      firstName: user.firstName,
      surname: user.surname,
      email: user.email,
      dob: user.dob ? user.dob.split('T')[0] : '', 
    });
    setShowUserModal(true);
  }, []);

  if (loading) {
    return <div>Loading users...</div>;
  }

  if (error) {
    return <div>Error loading users.</div>;
  }

  return (
    <div className="flex-grow-1" style={{ background: '#f5f5f5' }}>
      <div className="py-4 container" style={{ maxWidth: '1000px' }}>
        <Card className="shadow-sm">
          <Card.Body>
            <Button onClick={handleOpenAddUserModal}>Add User</Button>
            
            <Form.Group className="mt-3">
              <Form.Control
                type="text"
                placeholder="Search users by name, email, or DOB..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </Form.Group>

            <ListGroup className="mt-3">
              {paginatedUsers.map((user) => (
                <ListGroup.Item key={user.id} className="d-flex justify-content-between align-items-center">
                  <div>
                    {user.firstName} {user.surname} - {user.email} - DOB: {user.dob && !isNaN(new Date(user.dob)) ? new Date(user.dob).toLocaleDateString('en-GB') : 'N/A'}
                  </div>
                  <div>
                    <Button variant="outline-primary" size="sm" onClick={() => handleOpenEditUserModal(user)}>Edit</Button>
                    <Button variant="outline-danger" size="sm" onClick={() => deleteUser(user.id)}>Delete</Button>
                  </div>
                </ListGroup.Item>
              ))}
              {paginatedUsers.length === 0 && (
                <ListGroup.Item>No users found</ListGroup.Item>
              )}
            </ListGroup>

            <Pagination className="mt-3">
              <Pagination.First onClick={() => onPageChange(1)} disabled={currentPage === 1} />
              <Pagination.Prev onClick={() => onPageChange(currentPage - 1)} disabled={currentPage === 1} />
              {[...Array(totalPages)].map((_, idx) => (
                <Pagination.Item
                  key={idx + 1}
                  active={idx + 1 === currentPage}
                  onClick={() => onPageChange(idx + 1)}
                >
                  {idx + 1}
                </Pagination.Item>
              ))}
              <Pagination.Next onClick={() => onPageChange(currentPage + 1)} disabled={currentPage === totalPages} />
              <Pagination.Last onClick={() => onPageChange(totalPages)} disabled={currentPage === totalPages} />
            </Pagination>
          </Card.Body>
        </Card>

        <UserFormModal
          show={showUserModal}
          handleClose={handleCloseModal}
          handleSave={handleUserSave}
          userForm={userForm}
          handleChange={handleUserFormChange}
          title={selectedUser ? 'Edit User' : 'Add User'}
        />
      </div>
    </div>
  );
};
