import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import axios from '../lib/axios';
import { Container, Alert, Button, Spinner } from 'react-bootstrap';

export const SharedFile = () => {
    const { shareId } = useParams();
    const [file, setFile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchSharedFile = async () => {
            if (!shareId) {
                setError('Invalid share link');
                setLoading(false);
                return;
            }

            try {
                console.log('Fetching shared file:', shareId);
                const response = await axios.get(`/api/shared/${shareId}`);
                console.log('Shared file response:', response.data);
                setFile(response.data);
            } catch (err) {
                console.error('Error fetching shared file:', err);
                setError(err.response?.data?.message || 'File not found or link has expired');
            } finally {
                setLoading(false);
            }
        };

        fetchSharedFile();
    }, [shareId]);

    const handleDownload = async () => {
        try {
            const response = await axios.get(`/api/shared/${shareId}/download`, {
                responseType: 'blob'
            });
            
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', file.name);
            document.body.appendChild(link);
            link.click();
            link.remove();
        } catch (err) {
            console.error('Download error:', err);
            setError('Error downloading file');
        }
    };

    if (loading) {
        return (
            <Container className="mt-5 text-center">
                <Spinner animation="border" role="status">
                    <span className="visually-hidden">Loading...</span>
                </Spinner>
            </Container>
        );
    }

    if (error) {
        return (
            <Container className="mt-5">
                <Alert variant="danger">{error}</Alert>
            </Container>
        );
    }

    if (!file) {
        return (
            <Container className="mt-5">
                <Alert variant="warning">File not found</Alert>
            </Container>
        );
    }

    return (
        <Container className="mt-5">
            <div className="card">
                <div className="card-body">
                    <h2 className="card-title">{file.name}</h2>
                    <p className="text-muted">
                        Size: {(file.size / 1024 / 1024).toFixed(2)} MB
                    </p>
                    <Button 
                        variant="primary"
                        onClick={handleDownload}
                    >
                        Download File
                    </Button>
                </div>
            </div>
        </Container>
    );
};
