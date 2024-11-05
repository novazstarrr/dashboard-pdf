import { Container, Row, Col, Card, ListGroup, Spinner, Pagination } from 'react-bootstrap';
import { FaFile } from 'react-icons/fa';
import { useAuth } from '../hooks/useAuth';
import { useState, useEffect, useCallback } from 'react';
import axios from '../lib/axios';
import { PDFViewer } from './PDFViewer';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { UploadPDF } from './UploadPDF';
import SearchBar from './SearchBar'; 
import DownloadButton from './DownloadButton';
import ShareButton from './ShareButton';
import DeleteButton from './DeleteButton';

export const Dashboard = () => {
  const { user } = useAuth();
  const [userFiles, setUserFiles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(5);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0); 
  const [selectedFile, setSelectedFile] = useState(null);
  const [showViewer, setShowViewer] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  const fetchUserFiles = useCallback(async () => {
    try {
      setIsLoading(true);
      const token = localStorage.getItem('token');
      if (!token) {
        toast.error('Authentication token not found.');
        setIsLoading(false);
        return;
      }

      const response = await axios.get(
        `/api/files/my?page=${page}&page_size=${pageSize}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      setUserFiles(response.data.data);
      setTotalItems(response.data.total); 
      setTotalPages(Math.ceil(response.data.total / pageSize));

      console.log('Total files:', response.data.total);
      console.log('Total pages:', Math.ceil(response.data.total / pageSize));
    } catch (error) {
      console.error('Error fetching files:', error);
      toast.error('Error fetching files.');
    } finally {
      setIsLoading(false);
    }
  }, [page, pageSize]);

  useEffect(() => {
    fetchUserFiles();
  }, [fetchUserFiles]);

  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString();
  };

  const handleDeleteSuccess = () => {
    const updatedTotalItems = totalItems - 1; 
    const updatedMaxPage = Math.max(1, Math.ceil(updatedTotalItems / pageSize));

    if (page > updatedMaxPage) {
      setPage(updatedMaxPage);
    } else {
      const isCurrentPageEmpty = (updatedTotalItems % pageSize) === 0 && updatedMaxPage < page;
      if (isCurrentPageEmpty && page > 1) {
        setPage(page - 1);
      } else {
        fetchUserFiles();
      }
    }

    setTotalItems(updatedTotalItems); 
  };

  const fetchSearchResults = useCallback(async (term) => {
    try {
      if (!term) {
        fetchUserFiles();
        return;
      }

      const token = localStorage.getItem('token');
      if (!token) {
        toast.error('Authentication token not found.');
        return;
      }

      const response = await axios.get(
        `/api/files/search?q=${encodeURIComponent(term)}&page=${page}&page_size=${pageSize}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      setUserFiles(response.data.data);
      setTotalItems(response.data.total);
      setTotalPages(Math.max(1, Math.ceil(response.data.total / pageSize)));
    } catch (error) {
      console.error('Error searching files:', error);
      toast.error('Error searching files.');
      setUserFiles([]);
      setTotalItems(0);
      setTotalPages(1);
    }
  }, [page, pageSize, fetchUserFiles]);

  useEffect(() => {
    const handler = setTimeout(() => {
      if (searchTerm) {
        fetchSearchResults(searchTerm);
      } else {
        fetchUserFiles();
      }
    }, 300);

    return () => clearTimeout(handler);
  }, [searchTerm, fetchSearchResults, fetchUserFiles]);

  const handleSearch = (term) => {
    setSearchTerm(term);
    setPage(1); 
  };

  return (
    <div className="flex-grow-1" style={{ background: '#f5f5f5' }}>
      <Container className="py-4">
        <Row className="mb-4">
          <Col>
            <h1 className="fw-bold">Dashboard</h1>
            <p className="text-muted">Welcome back{user?.email ? `, ${user.email}` : ''}</p>
          </Col>
        </Row>

        <Row className="g-4">
          <Col xs={12} md={6} lg={4}>
            <UploadPDF fetchUserFiles={fetchUserFiles} />
          </Col>

          <Col xs={12} md={6} lg={8}>
            <Card className="shadow-sm h-100">
              <Card.Body>
                <Card.Title>
                  <FaFile className="me-2" />
                  Your Files
                </Card.Title>

                
                <SearchBar 
                  searchTerm={searchTerm} 
                  onSearch={handleSearch}
                />

                {isLoading ? (
                  <div className="text-center py-3">
                    <Spinner animation="border" variant="primary" />
                  </div>
                ) : (
                  <>
                    <ListGroup>
                      {userFiles.map((file) => (
                        <ListGroup.Item
                          key={file.id}
                          className="d-flex justify-content-between align-items-center"
                        >
                          <div
                            onClick={() => {
                              setSelectedFile(file);
                              setShowViewer(true);
                            }}
                            className="cursor-pointer"
                            style={{ cursor: 'pointer' }}
                          >
                            <FaFile className="me-2" />
                            {file.name}
                            <br />
                            <small className="text-muted">
                              Size: {formatFileSize(file.size)}
                              <br />
                              Type: {file.mimeType}
                              <br />
                              Created: {formatDate(file.createdAt)}
                              <br />
                              Updated: {formatDate(file.updatedAt)}
                            </small>
                          </div>
                          <div className="btn-group">
                            <DownloadButton
                              fileId={file.id}
                              fileName={file.name}
                            />
                            <DeleteButton
                              fileId={file.id}
                              onDelete={handleDeleteSuccess}
                            />
                            <ShareButton fileId={file.id} />
                          </div>
                        </ListGroup.Item>
                      ))}
                      {userFiles.length === 0 && (
                        <ListGroup.Item>
                          {searchTerm ? 'No matching files found' : 'No files uploaded yet'}
                        </ListGroup.Item>
                      )}
                    </ListGroup>

                    {!searchTerm && totalPages > 0 && (
                      <div className="d-flex justify-content-between align-items-center mt-3">
                        <small className="text-muted">
                          Showing {userFiles.length} file(s)
                        </small>
                        <Pagination>
                          <Pagination.First
                            onClick={() => setPage(1)}
                            disabled={page === 1}
                          />
                          <Pagination.Prev
                            onClick={() => setPage((p) => Math.max(1, p - 1))}
                            disabled={page === 1}
                          />

                          {[...Array(totalPages)].map((_, idx) => (
                            <Pagination.Item
                              key={idx + 1}
                              active={idx + 1 === page}
                              onClick={() => setPage(idx + 1)}
                            >
                              {idx + 1}
                            </Pagination.Item>
                          ))}

                          <Pagination.Next
                            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                            disabled={page === totalPages}
                          />
                          <Pagination.Last
                            onClick={() => setPage(totalPages)}
                            disabled={page === totalPages}
                          />
                        </Pagination>
                      </div>
                    )}
                  </>
                )}
              </Card.Body>
            </Card>
          </Col>
        </Row>

        <PDFViewer
          file={selectedFile}
          show={showViewer}
          onHide={() => {
            setShowViewer(false);
            setSelectedFile(null);
          }}
        />
      </Container>
    </div>
  );
};
