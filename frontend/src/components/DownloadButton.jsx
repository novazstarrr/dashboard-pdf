import { useState } from 'react';
import PropTypes from 'prop-types';
import { Button } from 'react-bootstrap';
import { FaDownload } from 'react-icons/fa';
import axios from '../lib/axios';
import { toast } from 'react-toastify';

const DownloadButton = ({ fileId, fileName, disabled }) => {
  const [isDownloading, setIsDownloading] = useState(false);

  const handleDownload = async () => {
    try {
      setIsDownloading(true);
      const token = localStorage.getItem('token');
      if (!token) {
        toast.error('Authentication token not found.');
        return;
      }

      const response = await axios.get(`/api/files/${fileId}/download`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        responseType: 'blob',
      });

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', fileName);
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);

      toast.success('File downloaded successfully');
    } catch (err) {
      console.error('Download error:', err);
      toast.error('Error downloading file: ' + (err.response?.data?.message || err.message));
    } finally {
      setIsDownloading(false);
    }
  };

  return (
    <Button
      variant="outline-success"
      size="sm"
      onClick={handleDownload}
      disabled={disabled || isDownloading}
      className="me-2 cursor-pointer"
    >
      <FaDownload className="me-1" />
      {isDownloading ? 'Downloading...' : 'Download'}
    </Button>
  );
};

DownloadButton.propTypes = {
  fileId: PropTypes.oneOfType([PropTypes.string, PropTypes.number]).isRequired,
  fileName: PropTypes.string.isRequired,
  disabled: PropTypes.bool,
};

DownloadButton.defaultProps = {
  disabled: false,
};

export default DownloadButton;
