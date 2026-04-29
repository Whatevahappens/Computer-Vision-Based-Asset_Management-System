export type Role = "ADMIN" | "ASSET_CUSTODIAN" | "ACCOUNTANT" | "EMPLOYEE";
export type UserStatus = "ACTIVE" | "INACTIVE" | "SUSPENDED";
export type AssetStatus = "ACTIVE" | "INACTIVE" | "DISPOSED" | "LOST" | "UNDER_MAINTENANCE";
export type AuditStatus = "PLANNED" | "IN_PROGRESS" | "COMPLETED" | "CANCELLED";
export type FindingType = "MATCHED" | "MISSING" | "UNREGISTERED" | "DAMAGED" | "MISPLACED";
export type DepreciationMethod = "STRAIGHT_LINE" | "DECLINING_BALANCE";

export interface User {
  id: string; firstName: string; lastName: string; email: string;
  username: string; phone: string; status: UserStatus; role: Role;
  departmentId?: string; department?: Department;
}
export interface Department { id: string; name: string; description: string; }
export interface Location { id: string; name: string; building: string; floor: string; room: string; capacity: number; }
export interface AssetModel {
  id: string; brand: string; modelName: string; assetType: string;
  category: string; defaultUnitPrice: number; defaultUsefulLifeMonths: number;
  depreciationMethod: DepreciationMethod;
}
export interface Asset {
  id: string; barcode: string; serialNumber: string; assetName: string;
  acquisitionPrice: number; acquisitionDate: string; usefulLifeMonths: number;
  currentValue: number; status: AssetStatus; nature: string; description: string;
  assetModelId?: string; assetModel?: AssetModel;
  departmentId?: string; department?: Department;
  locationId?: string; location?: Location;
  assignedUserId?: string; assignedUser?: User; createdAt: string;
}
export interface AssetHistory { id: string; changeType: string; changedAt: string; description: string; assetId: string; userId: string; }
export interface AuditSession {
  id: string; startedAt: string; endedAt?: string; status: AuditStatus; notes: string;
  locationId: string; location?: Location; performedBy: string; performer?: User;
  findings?: AuditFinding[]; summaries?: AuditSummary[];
}
export interface AuditFinding { id: string; type: FindingType; confidence: number; notes: string; auditSessionId: string; expectedAssetId?: string; detectedAssetId?: string; }
export interface AuditSummary { id: string; category: string; registeredCount: number; detectedCount: number; difference: number; auditSessionId: string; }
export interface Notification { id: string; title: string; message: string; type: string; isRead: boolean; createdAt: string; }
export interface DashboardStats { totalAssets: number; activeAssets: number; totalValue: number; totalAudits: number; }
export interface PaginatedResponse<T> { data: T[]; total: number; page: number; limit: number; totalPages: number; }
export interface LoginResponse { token: string; user: User; }
