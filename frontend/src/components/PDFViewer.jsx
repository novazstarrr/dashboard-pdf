import { memo } from 'react';
import PropTypes from 'prop-types';
import { Modal, Alert, Spinner } from 'react-bootstrap';
import { usePDFViewer } from '../hooks/usePDFViewer';

export const PDFViewer = memo(({ file, show, onHide }) => {
  const { pdfUrl, error, loading } = usePDFViewer(file, show);

  return (
    <Modal 
      show={show} 
      onHide={onHide} 
      size="lg"
      aria-labelledby="pdf-viewer-modal"
    >
      <Modal.Header closeButton>
        <Modal.Title id="pdf-viewer-modal">
          {file?.name || 'PDF Viewer'}
        </Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {loading && (
          <div className="d-flex justify-content-center p-4">
            <Spinner animation="border" variant="primary" />
          </div>
        )}
        {error && (
          <Alert variant="danger">
            {error}
          </Alert>
        )}
        {pdfUrl && !loading && !error && (
          <iframe
            src={pdfUrl}
            style={{ width: '100%', height: '500px' }}
            title={`PDF Viewer - ${file?.name}`}
            aria-label={`PDF document ${file?.name}`}
          />
        )}
      </Modal.Body>
    </Modal>
  );
});

PDFViewer.displayName = 'PDFViewer';

PDFViewer.propTypes = {
  file: PropTypes.shape({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
  }),
  show: PropTypes.bool.isRequired,
  onHide: PropTypes.func.isRequired,
};
