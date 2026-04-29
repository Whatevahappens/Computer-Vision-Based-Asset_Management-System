import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { ConfigProvider } from "antd";
import { AuthProvider, useAuth } from "./context/AuthContext";
import AppLayout from "./layouts/AppLayout";
import LoginPage from "./pages/LoginPage";
import DashboardPage from "./pages/DashboardPage";
import AssetsPage from "./pages/AssetsPage";
import AuditPage from "./pages/AuditPage";
import UsersPage from "./pages/UsersPage";
import DepreciationPage from "./pages/DepreciationPage";
import MyAssetsPage from "./pages/MyAssetsPage";
import LocationsPage from "./pages/LocationsPage";
import DepartmentsPage from "./pages/DepartmentsPage";
import ReportsPage from "./pages/ReportsPage";
import NotificationsPage from "./pages/NotificationsPage";
import AssignmentsPage from "./pages/AssignmentsPage";
import AssetModelsPage from "./pages/AssetModelsPage";

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { token } = useAuth();
  return token ? <>{children}</> : <Navigate to="/login" replace />;
};

const App: React.FC = () => (
  <ConfigProvider
    theme={{
      token: {
        colorPrimary: "#1677ff",
        fontFamily: "'IBM Plex Sans', -apple-system, sans-serif",
        borderRadius: 6,
      },
    }}
  >
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/" element={<ProtectedRoute><AppLayout /></ProtectedRoute>}>
            <Route index element={<DashboardPage />} />
            <Route path="assets" element={<AssetsPage />} />
            <Route path="audits" element={<AuditPage />} />
            <Route path="users" element={<UsersPage />} />
            <Route path="depreciation" element={<DepreciationPage />} />
            <Route path="my-assets" element={<MyAssetsPage />} />
            <Route path="locations" element={<LocationsPage />} />
            <Route path="departments" element={<DepartmentsPage />} />
            <Route path="reports" element={<ReportsPage />} />
            <Route path="notifications" element={<NotificationsPage />} />
            <Route path="assignments" element={<AssignmentsPage />} />
            <Route path="asset-models" element={<AssetModelsPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  </ConfigProvider>
);

export default App;
