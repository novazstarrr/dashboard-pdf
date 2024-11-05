import { useState } from 'react';
import axios from '../lib/axios';
import { toast } from 'react-toastify';

const useShare = () => {
  const [isSharing, setIsSharing] = useState(false);
  const [shareableLink, setShareableLink] = useState('');
  const [showShareModal, setShowShareModal] = useState(false);

  const shareFile = async (fileId) => {
    const token = localStorage.getItem('token');
    if (!token) {
      toast.error('Authentication token not found.');
      return;
    }

    setIsSharing(true);
    try {
      const response = await axios.post(
        `/api/files/${fileId}/share`,
        {},
        { headers: { Authorization: `Bearer ${token}` } }
      );

      const { shareableLink } = response.data;
      if (!shareableLink) {
        throw new Error('Invalid share response from server.');
      }

      setShareableLink(shareableLink);
      setShowShareModal(true);
      toast.success('Share link generated successfully.');
    } catch (error) {
      console.error('Share error:', error);
      const errorMessage =
        error.response?.data?.message || error.message || 'Unknown error';
      toast.error(`Error generating share link: ${errorMessage}`);
    } finally {
      setIsSharing(false);
    }
  };

  const copyLink = async () => {
    try {
      await navigator.clipboard.writeText(shareableLink);
      toast.success('Link copied to clipboard!');
    } catch (err) {
      console.error('Copy error:', err);
      toast.error('Failed to copy link');
    }
  };

  const closeModal = () => setShowShareModal(false);

  return {
    shareFile,
    isSharing,
    shareableLink,
    showShareModal,
    copyLink,
    closeModal,
  };
};

export default useShare;
