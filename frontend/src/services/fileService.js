import axios from '../lib/axios';

export const fileService = {
    getFiles: () => axios.get('/api/files/my'),
    shareFile: (fileId) => axios.post(`/api/files/${fileId}/share`),
    deleteFile: (fileId) => axios.delete(`/api/files/${fileId}`),
    downloadFile: (fileId) => axios.get(`/api/files/${fileId}/download`, { responseType: 'blob' }),
}; 