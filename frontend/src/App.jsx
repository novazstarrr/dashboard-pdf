import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Header } from './layout/Header.jsx';
import { Footer } from './layout/Footer.jsx';
import { AuthProvider } from './context/AuthContext.jsx';
import { routes } from './Routes.jsx';
import { SharedFile } from './components/SharedFile.jsx';
import { UserManagement } from './components/UserManagement.jsx';
import { Dashboard } from './components/Dashboard.jsx';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

function App() {
  return (
    <Router>
      <AuthProvider>
        <div className="d-flex flex-column min-vh-100">
          <Header />
          <main className="flex-grow-1">
            <Routes>
              <Route path="/" element={<Navigate to="/login" replace />} />
              <Route path="/shared/:shareId" element={<SharedFile />} />
              <Route path="/user-management" element={<UserManagement />} />
              <Route path="/dashboard" element={<Dashboard />} />
              {routes.map((route) => (
                <Route
                  key={route.path}
                  path={route.path}
                  element={route.element}
                />
              ))}
              <Route path="*" element={<Navigate to="/login" replace />} />
            </Routes>
          </main>
          <Footer />
          <ToastContainer
            position="bottom-right"
            autoClose={3000}
            hideProgressBar={false}
            newestOnTop={false}
            closeOnClick
            rtl={false}
            pauseOnFocusLoss
            draggable
            pauseOnHover
          />
        </div>
      </AuthProvider>
    </Router>
  );
}

export default App;
