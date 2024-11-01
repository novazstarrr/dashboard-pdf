import { useState } from 'react';
import axios from '../lib/axios';
import { toast } from 'react-toastify';

const useDelete = () => {
  const [isDeleting, setIsDeleting] = useState(false);

  const deleteFile = async (fileId, onDeleteSuccess) => {
    const confirmDelete = window.confirm('Are you sure you want to delete this file?');
    if (!confirmDelete) return;

    try {
      setIsDeleting(true);
      const token = localStorage.getItem('token');
      if (!token) {
        toast.error('Authentication token not found.');
        return;
      }

      await axios.delete(`/api/files/${fileId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      toast.success('File deleted successfully.');
      if (onDeleteSuccess) onDeleteSuccess(fileId);
    } catch (error) {
      console.error('Delete error:', error);
      toast.error(`Error deleting file: ${error.response?.data?.message || error.message}`);
    } finally {
      setIsDeleting(false);
    }
  };

  return { deleteFile, isDeleting };
};

export default useDelete;
