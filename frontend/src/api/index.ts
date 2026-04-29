import axios from "axios";
import type {
  Asset, AssetHistory, AssetModel, AuditSession, DashboardStats,
  Department, Location, LoginResponse, Notification, PaginatedResponse, User,
} from "../types";

const api = axios.create({ baseURL: "/api/v1" });

// Attach JWT to every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

// Auto-logout on 401
api.interceptors.response.use(
  (r) => r,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      window.location.href = "/login";
    }
    return Promise.reject(err);
  }
);

// ── Auth ──
export const login = (username: string, password: string) =>
  api.post<LoginResponse>("/auth/login", { username, password });

export const getMe = () => api.get<User>("/auth/me");

export const changePassword = (currentPassword: string, newPassword: string) =>
  api.put("/auth/password", { currentPassword, newPassword });

// ── Dashboard ──
export const getDashboard = () => api.get<DashboardStats>("/dashboard");

// ── Assets ──
export const listAssets = (params?: Record<string, unknown>) =>
  api.get<PaginatedResponse<Asset>>("/assets", { params });

export const getAsset = (id: string) => api.get<Asset>(`/assets/${id}`);

export const createAsset = (data: Record<string, unknown>) =>
  api.post<Asset>("/assets", data);

export const updateAsset = (id: string, data: Record<string, unknown>) =>
  api.put<Asset>(`/assets/${id}`, data);

export const deleteAsset = (id: string) => api.delete(`/assets/${id}`);

export const assignAsset = (id: string, data: { userId: string; locationId?: string; notes?: string }) =>
  api.post(`/assets/${id}/assign`, data);

export const transferAsset = (id: string, data: { toUserId: string; locationId?: string; notes?: string }) =>
  api.post(`/assets/${id}/transfer`, data);

export const disposeAsset = (id: string, data: { reason: string; residualValue?: number; notes?: string }) =>
  api.post(`/assets/${id}/dispose`, data);

export const getAssetHistory = (id: string) =>
  api.get<AssetHistory[]>(`/assets/${id}/history`);

export const getMyAssets = () => api.get<Asset[]>("/my-assets");

// ── Asset Models ──
export const listAssetModels = () => api.get<AssetModel[]>("/asset-models");
export const createAssetModel = (data: Record<string, unknown>) =>
  api.post<AssetModel>("/asset-models", data);

// ── Locations ──
export const listLocations = () => api.get<Location[]>("/locations-list");
export const createLocation = (data: Record<string, unknown>) =>
  api.post<Location>("/locations", data);
export const updateLocation = (id: string, data: Record<string, unknown>) =>
  api.put<Location>(`/locations/${id}`, data);
export const deleteLocation = (id: string) => api.delete(`/locations/${id}`);

// ── Departments ──
export const listDepartments = () => api.get<Department[]>("/departments-list");
export const createDepartment = (data: Record<string, unknown>) =>
  api.post<Department>("/departments", data);
export const updateDepartment = (id: string, data: Record<string, unknown>) =>
  api.put<Department>(`/departments/${id}`, data);
export const deleteDepartment = (id: string) => api.delete(`/departments/${id}`);

// ── Users ──
export const listUsers = (params?: Record<string, unknown>) =>
  api.get<PaginatedResponse<User>>("/users", { params });
export const createUser = (data: Record<string, unknown>) =>
  api.post<User>("/users", data);
export const updateUser = (id: string, data: Record<string, unknown>) =>
  api.put<User>(`/users/${id}`, data);
export const deactivateUser = (id: string) =>
  api.put(`/users/${id}/deactivate`);

// ── Audits ──
export const listAudits = (params?: Record<string, unknown>) =>
  api.get<PaginatedResponse<AuditSession>>("/audits", { params });
export const getAudit = (id: string) => api.get<AuditSession>(`/audits/${id}`);
export const startAudit = (locationId: string, notes?: string) =>
  api.post<AuditSession>("/audits", { locationId, notes });
export const runCVAudit = (sessionId: string, file: File) => {
  const fd = new FormData();
  fd.append("file", file);
  return api.post<AuditSession>(`/audits/${sessionId}/cv`, fd, {
    headers: { "Content-Type": "multipart/form-data" },
  });
};

// ── Depreciation ──
export const calculateDepreciation = (assetId: string, method: string) =>
  api.post("/depreciation/calculate", { assetId, method });
export const revalueAsset = (assetId: string, newValue: number, reason: string) =>
  api.post("/depreciation/revalue", { assetId, newValue, reason });

// ── Reports ──
export const generateReport = (reportType: string, format: string) =>
  api.post("/reports/generate", { reportType, format });

// ── Notifications ──
export const listNotifications = () => api.get<Notification[]>("/notifications");
export const markNotificationRead = (id: string) => api.put(`/notifications/${id}/read`);
export const markAllRead = () => api.put("/notifications/read-all");

export default api;
