import { useState } from 'react';
import PropTypes from 'prop-types';
import { Button, Modal, InputGroup, Form } from 'react-bootstrap';
import { FaShare, FaCopy } from 'react-icons/fa';
import axios from '../lib/axios';
import { toast } from 'react-toastify';

const ShareButton = ({ fileId }) => {
  const [showShareModal, setShowShareModal] = useState(false);
  const [shareableLink, setShareableLink] = useState('');
  const [isSharing, setIsSharing] = useState(false);

  const handleShare = async () => {
    try {
      setIsSharing(true);
      const token = localStorage.getItem('token');
      if (!token) {
        toast.error('Authentication token not found.');
        return;
      }

      const response = await axios.post(
        `/api/files/${fileId}/share`,
        {},
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const { shareableLink } = response.data;
      if (!shareableLink) {
        throw new Error('Invalid share response from server');
      }

      setShareableLink(shareableLink);
      setShowShareModal(true);
      toast.success('Share link generated successfully');
    } catch (err) {
      console.error('Share error:', err);
      toast.error('Error generating share link: ' + (err.response?.data?.message || err.message));
    } finally {
      setIsSharing(false);
    }
  };

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(shareableLink);
      toast.success('Link copied to clipboard!');
    } catch (err) {
      console.error('Copy error:', err);
      toast.error('Failed to copy link');
    }
  };

  return (
    <>
      <Button
        variant="outline-primary"
        size="sm"
        onClick={handleShare}
        disabled={isSharing}
        className="cursor-pointer"
      >
        <FaShare /> {isSharing ? 'Sharing...' : 'Share'}
      </Button>

      <Modal show={showShareModal} onHide={() => setShowShareModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Share File</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <InputGroup>
            <Form.Control type="text" value={shareableLink} readOnly />
            <Button variant="outline-secondary" onClick={handleCopy}>
              <FaCopy /> Copy
            </Button>
          </InputGroup>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => setShowShareModal(false)}>
            Close
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};

ShareButton.propTypes = {
  fileId: PropTypes.oneOfType([PropTypes.string, PropTypes.number]).isRequired,
};

export default ShareButton;
