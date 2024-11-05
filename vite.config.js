import { defineConfig } from 'vite';

export default defineConfig({
  // ... other config
  optimizeDeps: {
    exclude: ['react-beautiful-dnd']
  }
}); 