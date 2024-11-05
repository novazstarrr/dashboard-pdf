import { Card, Spinner, Alert } from 'react-bootstrap';
import { FaCloudUploadAlt, FaFileUpload } from 'react-icons/fa';
import { useState, useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import axios from '../lib/axios';
import 'react-toastify/dist/ReactToastify.css';
import PropTypes from 'prop-types';

export const UploadPDF = ({ fetchUserFiles }) => {
  const [uploadStatus, setUploadStatus] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const onDrop = useCallback(async (acceptedFiles) => {
    const uploadedFile = acceptedFiles[0];
    if (!uploadedFile) return;

    if (uploadedFile.type !== 'application/pdf') {
      setUploadStatus('Error: Please upload a PDF file');
      return;
    }

    try {
      setIsLoading(true);
      const formData = new FormData();
      formData.append('file', uploadedFile);

      const token = localStorage.getItem('token');
      await axios.post(`/api/files/upload`, formData, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'multipart/form-data',
        },
      });

      setUploadStatus('File uploaded successfully!');
      fetchUserFiles();
    } catch (error) {
      setUploadStatus('Error uploading file: ' + error.message);
    } finally {
      setIsLoading(false);
    }
  }, [fetchUserFiles]);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'application/pdf': ['.pdf']
    },
    multiple: false
  });

  return (
    <Card className="shadow-sm h-100">
      <Card.Body>
        <Card.Title>
          <FaCloudUploadAlt className="me-2" />
          Upload PDF
        </Card.Title>
        <div
          {...getRootProps()}
          className={`border border-2 rounded p-4 mb-3 text-center ${
            isDragActive ? 'border-primary bg-light' : 'border-dashed'
          }`}
          style={{ cursor: 'pointer' }}
        >
          <input {...getInputProps()} />
          {isLoading ? (
            <Spinner animation="border" variant="primary" />
          ) : isDragActive ? (
            <>
              <FaFileUpload className="text-primary mb-2" size={32} />
              <p className="mb-0">Drop the PDF here...</p>
            </>
          ) : (
            <>
              <FaCloudUploadAlt className="text-primary mb-2" size={32} />
              <p className="mb-0">Drag & drop a PDF here, or click to select</p>
            </>
          )}
        </div>
        {uploadStatus && (
          <Alert variant={uploadStatus.startsWith('Error') ? 'danger' : 'success'}>
            {uploadStatus}
          </Alert>
        )}
      </Card.Body>
    </Card>
  );
};

UploadPDF.propTypes = {
  fetchUserFiles: PropTypes.func.isRequired,
};
