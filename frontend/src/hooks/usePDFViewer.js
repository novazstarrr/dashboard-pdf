import { useState, useEffect } from 'react';
import axios from '../lib/axios';

export const usePDFViewer = (file, show) => {
  const [state, setState] = useState({
    pdfUrl: null,
    error: null,
    loading: false,
  });

  useEffect(() => {
    let isMounted = true;
    let objectUrl = null;
    const token = localStorage.getItem('token');

    const fetchPDF = async () => {
      if (!file?.id || !show) return;

      setState(prev => ({ ...prev, loading: true, error: null }));

      try {
        const response = await axios({
          method: 'get',
          url: `/api/files/${file.id}/view`,
          headers: { Authorization: `Bearer ${token}` },
          responseType: 'blob'
        });

        if (!isMounted) return;

        objectUrl = URL.createObjectURL(response.data);
        setState({ loading: false, error: null, pdfUrl: objectUrl });
      } catch (error) {
        if (!isMounted) return;

        const errorMessage = error.response?.data?.error 
          || error.response?.data?.message 
          || 'Failed to load PDF';
        
        setState({ loading: false, error: errorMessage, pdfUrl: null });
        console.error('Error loading PDF:', error);
      }
    };

    fetchPDF();

    return () => {
      isMounted = false;
      if (objectUrl) {
        URL.revokeObjectURL(objectUrl);
      }
    };
  }, [file?.id, show]); 

  return state;
}; 