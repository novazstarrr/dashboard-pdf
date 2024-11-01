import { Navigate } from 'react-router-dom';
import { LoginPage } from './pages/login/Login.page.jsx';
import { RegisterPage } from './pages/register/Register.page.jsx';
import { DashboardPage } from './pages/dashboard/Dashboard.page.jsx';
import { SharedFilePage } from './pages/shared/[shareId]/SharedFile.page.jsx';
import { PrivateRoute } from './components/PrivateRoute.jsx'; 

export const routes = [
  {
    path: '/login',
    element: <LoginPage />
  },
  {
    path: '/register',
    element: <RegisterPage />
  },
  {
    path: '/dashboard',
    element: (
      <PrivateRoute>
        <DashboardPage />
      </PrivateRoute>
    )
  },
  {
    path: '/shared/:shareId',
    element: <SharedFilePage />
  },
  {
    path: '*',
    element: <Navigate to="/login" />
  }
];
