import { useState, useEffect, useCallback, useMemo } from 'react';
import axios from '../lib/axios';
import { toast } from 'react-toastify';
import { z } from 'zod';

/**
 * Custom Hook: useUsers
 * Manages user data, including fetching, adding, editing, deleting, and reordering.
 *
 * @returns {Object} - Users data and related handlers.
 */
const useUsers = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  
  const [currentPage, setCurrentPage] = useState(1);
  const ITEMS_PER_PAGE = 5;
  
  const [searchTerm, setSearchTerm] = useState('');

  
  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/users');
      setUsers(response.data);
      setError(null);
    } catch (err) {
      console.error('Error fetching users:', err);
      setError(err);
      toast.error('Failed to load users.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  // Move userSchema outside of the component/hook
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
        const today = new Date();
        const age = today.getFullYear() - birthDate.getFullYear();
        return age >= 13;
      }, 'You must be at least 13 years old')
  });

  /**
   * Adds a new user and navigates to the last page.
   * @param {Object} userData - The data of the user to add.
   */
  const addUser = useCallback(async (userData) => {
    try {
      // Validate the data first
      const validationResult = userSchema.safeParse(userData);
      
      if (!validationResult.success) {
        const errors = validationResult.error.errors.map(err => ({
          field: err.path[0],
          message: err.message
        }));
        throw new Error(JSON.stringify(errors));
      }

      // Check for duplicate email
      const existingUser = users.find(user => user.email === userData.email);
      if (existingUser) {
        throw new Error(JSON.stringify([{
          field: 'email',
          message: 'This email is already registered'
        }]));
      }

      const formattedData = { ...userData, dob: new Date(userData.dob).toISOString() };
      const { data } = await axios.post('/api/users', formattedData);
      
      const newUser = {
        ...data,
        dob: data.dob ? new Date(data.dob).toISOString() : "0001-01-01T00:00:00Z",
      };

      setUsers((prevUsers) => {
        const newUsers = [...prevUsers, newUser];
        const newTotalPages = Math.ceil(newUsers.length / ITEMS_PER_PAGE) || 1;
        setCurrentPage(newTotalPages); 
        return newUsers;
      });
      toast.success('User added successfully!');
    } catch (err) {
      console.error('Error adding user:', err);
      
      try {
        const validationErrors = JSON.parse(err.message);
        validationErrors.forEach(error => {
          toast.error(`${error.field}: ${error.message}`);
        });
      } catch {
        if (err.response?.data?.message) {
          toast.error(`Failed to add user: ${err.response.data.message}`);
        } else {
          toast.error('Failed to add user.');
        }
      }
      throw err;
    }
  }, [ITEMS_PER_PAGE, setCurrentPage, users, setUsers, userSchema]);

  /**
   * Edits an existing user.
   * @param {string|number} userId - The ID of the user to edit.
   * @param {Object} userData - The updated user data.
   */
  const editUser = useCallback(async (userId, userData) => {
    try {
      const validationResult = userSchema.safeParse(userData);
      
      if (!validationResult.success) {
        const errors = validationResult.error.errors.map(err => ({
          field: err.path[0],
          message: err.message
        }));
        throw new Error(JSON.stringify(errors));
      }

      const numericId = typeof userId === 'string' ? parseInt(userId, 10) : userId;
      const formattedData = { 
        ...userData, 
        dob: new Date(userData.dob).toISOString() 
      };

      const { data } = await axios.put(`/api/users/${numericId}`, formattedData);
      
      setUsers((prevUsers) => {
        const updatedUsers = prevUsers.map((user) =>
          user.id === numericId ? { ...user, ...data } : user
        );
        return updatedUsers;
      });

      toast.success('User updated successfully!');
    } catch (err) {
      console.error('Error editing user:', err);
      
      try {
        const validationErrors = JSON.parse(err.message);
        validationErrors.forEach(error => {
          toast.error(`${error.field}: ${error.message}`);
        });
      } catch {
        if (err.response?.data?.message) {
          toast.error(`Failed to update user: ${err.response.data.message}`);
        } else {
          toast.error('Failed to update user.');
        }
      }
      throw err;
    }
  }, [setUsers, userSchema]);

  /**
   * Deletes a user.
   * @param {string|number} userId - delete user id
   */
  const deleteUser = useCallback(async (userId) => {
    try {
      await axios.delete(`/api/users/${userId}`);
      setUsers((prevUsers) => prevUsers.filter((user) => user.id !== userId));
      toast.success('User deleted successfully!');
    } catch (err) {
      console.error('Error deleting user:', err);
      toast.error('Failed to delete user.');
      throw err;
    }
  }, []);

  /**
   * Reorders users locally after drag-and-drop and persists the new order.
   * @param {Array} reorderedUsers - The new order of users.
   */
  const reorderUsers = useCallback(async (reorderedUsers) => {
    try {
      await axios.put('/api/users/reorder', { users: reorderedUsers });
      setUsers(reorderedUsers);
      toast.success('Users reordered successfully!');
    } catch (err) {
      console.error('Error reordering users:', err);
      toast.error('Failed to reorder users.');
      throw err;
    }
  }, []);

  /**
   * Filters users based on the search term.
   * @returns {Array} - Filtered users.
   */
  const filteredUsers = useMemo(() => {
    if (!searchTerm) return users;
    const searchLower = searchTerm.toLowerCase();
    
    return users.filter((user) => {
      const dobString = user.dob ? new Date(user.dob).toLocaleDateString('en-GB') : '';
      
      return (
        user.firstName?.toLowerCase().includes(searchLower) ||
        user.surname?.toLowerCase().includes(searchLower) ||
        user.email?.toLowerCase().includes(searchLower) ||
        dobString.includes(searchTerm)
      );
    });
  }, [searchTerm, users]);

  /**
   * Gets the users for the current page.
   * @returns {Array} - Paginated users.
   */
  const paginatedUsers = useMemo(() => {
    const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
    return filteredUsers.slice(startIndex, startIndex + ITEMS_PER_PAGE);
  }, [currentPage, filteredUsers, ITEMS_PER_PAGE]);

  /**
   * Calculates the total number of pages.
   * @returns {number} - Total pages.
   */
  const totalPages = useMemo(() => {
    return Math.ceil(filteredUsers.length / ITEMS_PER_PAGE) || 1;
  }, [filteredUsers, ITEMS_PER_PAGE]);

  /**
   * Handles page changes.
   * @param {number} pageNumber - The page number to navigate to.
   */
  const onPageChange = useCallback((pageNumber) => {
    setCurrentPage(pageNumber);
  }, []);

  /**
   * Handles search term changes.
   * @param {string} term - The search term.
   */
  const handleSearch = useCallback((term) => {
    setSearchTerm(term);
    setCurrentPage(1); 
  }, []);

  return {
    users,
    paginatedUsers,
    loading,
    error,
    currentPage,
    totalPages,
    addUser,
    editUser,
    deleteUser,
    reorderUsers,
    searchTerm,
    setSearchTerm: handleSearch,
    onPageChange,
    fetchUsers,
  };
};

export default useUsers; 