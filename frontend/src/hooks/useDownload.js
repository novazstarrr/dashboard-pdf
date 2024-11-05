import { useState } from 'react';
import axios from '../lib/axios';
import { toast } from 'react-toastify';

const useDownload = () => {
  const [isDownloading, setIsDownloading] = useState(false);

  const downloadFile = async (fileId, fileName) => {
    const token = localStorage.getItem('token');
    if (!token) {
      toast.error('Authentication token not found.');
      return;
    }

    setIsDownloading(true);
    try {
      const response = await axios.get(`/api/files/${fileId}/download`, {
        headers: { Authorization: `Bearer ${token}` },
        responseType: 'blob',
      });

      const blob = new Blob([response.data]);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = fileName;
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);

      toast.success('File downloaded successfully.');
    } catch (error) {
      console.error('Download error:', error);
      toast.error(`Error downloading file: ${error.message}`);
    } finally {
      setIsDownloading(false);
    }
  };

  return { downloadFile, isDownloading };
};

export default useDownload; 