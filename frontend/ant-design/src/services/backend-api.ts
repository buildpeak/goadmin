import axios, { isAxiosError } from "axios";
import {
  JwtToken,
  SignUpRequest,
  User,
  Role,
  Permission,
  AuditLog,
  CreateRoleRequest,
  CreatePermissionRequest,
} from "./types";

const backendUrl = process.env.REACT_APP_BACKEND_URL || "";

function getAuthHeader() {
  const token = localStorage.getItem("accessToken");
  return token ? { Authorization: `Bearer ${token}` } : {};
}

const api = axios.create({ baseURL: backendUrl });

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("accessToken");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export interface ApiError {
  type: string;
  title: string;
  detail: string;
  status: number;
  instance: string;
  errors?: {
    pointer: string;
    detail: string;
  }[];
}

export class ApiError extends Error {
  constructor(apiErr: ApiError) {
    super(apiErr.detail);
  }
}

export function statusMatch(error: unknown, status: number) {
  if (isAxiosError(error)) {
    return error.response?.status === status;
  }
  return false;
}

function handleAxiosError(error: unknown) {
  if (isAxiosError(error)) {
    const status = error.response?.status;

    switch (status) {
      case 400:
        throw new Error("Invalid input");
      default:
        throw new ApiError(error.response?.data);
    }
  }
}

let isRefreshing = false;
let refreshQueue: Array<{
  resolve: (token: string) => void;
  reject: (err: unknown) => void;
}> = [];

async function tryRefreshToken(): Promise<string | null> {
  const refreshToken = localStorage.getItem("refreshToken");
  if (!refreshToken) return null;

  try {
    const response = await axios.post(`${backendUrl}/auth/refresh`, {
      refresh_token: refreshToken,
    });
    const data = response.data as JwtToken;
    localStorage.setItem("accessToken", data.access_token);
    localStorage.setItem("refreshToken", data.refresh_token);
    return data.access_token;
  } catch {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    return null;
  }
}

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (!isAxiosError(error) || error.response?.status !== 401) {
      return Promise.reject(error);
    }

    const originalRequest = error.config;
    if (!originalRequest) return Promise.reject(error);

    // Don't try to refresh for the refresh endpoint itself
    if (originalRequest.url?.includes("/auth/refresh")) {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      window.location.href = "/login";
      return Promise.reject(error);
    }

    // Don't try for auth endpoints (login, signup, etc.)
    if (originalRequest.url?.includes("/auth/")) {
      return Promise.reject(error);
    }

    if (!isRefreshing) {
      isRefreshing = true;
      const newToken = await tryRefreshToken();
      isRefreshing = false;

      if (newToken) {
        refreshQueue.forEach(({ resolve }) => resolve(newToken));
        refreshQueue = [];
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        return api(originalRequest);
      }

      refreshQueue.forEach(({ reject }) => reject(new Error("refresh failed")));
      refreshQueue = [];
      window.location.href = "/login";
      return Promise.reject(error);
    }

    return new Promise((resolve, reject) => {
      refreshQueue.push({
        resolve: (token: string) => {
          originalRequest.headers.Authorization = `Bearer ${token}`;
          resolve(api(originalRequest));
        },
        reject,
      });
    });
  }
);

export async function login(
  username: string,
  password: string
): Promise<JwtToken> {
  try {
    const response = await axios.post(`${backendUrl}/auth/login`, {
      username,
      password,
    });

    return response.data as JwtToken;
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function verifyGoogleIdToken(idToken: string): Promise<JwtToken> {
  try {
    const response = await axios.post(`${backendUrl}/auth/signin-with-google`, {
      id_token: idToken,
    });

    return response.data as JwtToken;
  } catch (error) {
    if (statusMatch(error, 404)) {
      window.location.href = "/signup";
    }

    handleAxiosError(error);
    throw error;
  }
}

export async function signUp(req: SignUpRequest): Promise<User> {
  try {
    const response = await axios.post(`${backendUrl}/auth/signup`, req);
    return response.data as User;
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function logout(accessToken: string) {
  try {
    await axios.post(`${backendUrl}/auth/logout`, null, {
      headers: { Authorization: `Bearer ${accessToken}` },
    });
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function getUsers(): Promise<User[]> {
  try {
    const response = await api.get("/v1/users");
    return response.data as User[];
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function getUser(id: string): Promise<User> {
  try {
    const response = await api.get(`/v1/users/${id}`);
    return response.data as User;
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function updateUser(
  id: string,
  data: Partial<User>
): Promise<User> {
  try {
    const response = await api.patch(`/v1/users/${id}`, data);
    return response.data as User;
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function refreshToken(
  refresh_token: string
): Promise<JwtToken> {
  try {
    const response = await axios.post(`${backendUrl}/auth/refresh`, {
      refresh_token,
    });
    return response.data as JwtToken;
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function forgotPassword(email: string): Promise<void> {
  try {
    await axios.post(`${backendUrl}/auth/forgot-password`, { email });
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function resetPassword(
  token: string,
  password: string
): Promise<void> {
  try {
    await axios.post(`${backendUrl}/auth/reset-password`, { token, password });
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function getProfile(): Promise<User> {
  try {
    const response = await api.get("/auth/profile");
    return response.data as User;
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

// ── Roles ──────────────────────────────────────────

export async function listRoles(): Promise<Role[]> {
  try {
    const response = await api.get("/v1/roles");
    return response.data as Role[];
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function createRole(data: CreateRoleRequest): Promise<Role> {
  try {
    const response = await api.post("/v1/roles", data);
    return response.data as Role;
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function updateRole(
  id: string,
  data: Partial<Role>
): Promise<Role> {
  try {
    const response = await api.patch(`/v1/roles/${id}`, data);
    return response.data as Role;
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function deleteRole(id: string): Promise<void> {
  try {
    await api.delete(`/v1/roles/${id}`);
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

// ── Permissions ────────────────────────────────────

export async function listPermissions(): Promise<Permission[]> {
  try {
    const response = await api.get("/v1/permissions");
    return response.data as Permission[];
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function createPermission(
  data: CreatePermissionRequest
): Promise<Permission> {
  try {
    const response = await api.post("/v1/permissions", data);
    return response.data as Permission;
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function updatePermission(
  id: string,
  data: Partial<Permission>
): Promise<Permission> {
  try {
    const response = await api.patch(`/v1/permissions/${id}`, data);
    return response.data as Permission;
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function deletePermission(id: string): Promise<void> {
  try {
    await api.delete(`/v1/permissions/${id}`);
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

// ── Role-Permission assignments ────────────────────

export async function listPermissionsByRole(
  roleId: string
): Promise<Permission[]> {
  try {
    const response = await api.get(`/v1/roles/${roleId}/permissions`);
    return response.data as Permission[];
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function addPermissionToRole(
  roleId: string,
  permissionId: string
): Promise<void> {
  try {
    await api.post(`/v1/roles/${roleId}/permissions/${permissionId}`);
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function removePermissionFromRole(
  roleId: string,
  permissionId: string
): Promise<void> {
  try {
    await api.delete(`/v1/roles/${roleId}/permissions/${permissionId}`);
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

// ── User-Role assignments ──────────────────────────

export async function listRolesByUser(userId: string): Promise<Role[]> {
  try {
    const response = await api.get(`/v1/users/${userId}/roles`);
    return response.data as Role[];
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function assignRoleToUser(
  userId: string,
  roleId: string
): Promise<void> {
  try {
    await api.post(`/v1/users/${userId}/roles/${roleId}`);
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function unassignRoleFromUser(
  userId: string,
  roleId: string
): Promise<void> {
  try {
    await api.delete(`/v1/users/${userId}/roles/${roleId}`);
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

// ── User-Permission assignments (direct) ───────────

export async function listPermissionsByUser(
  userId: string
): Promise<Permission[]> {
  try {
    const response = await api.get(`/v1/users/${userId}/permissions`);
    return response.data as Permission[];
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function assignPermissionToUser(
  userId: string,
  permissionId: string
): Promise<void> {
  try {
    await api.post(`/v1/users/${userId}/permissions/${permissionId}`);
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

export async function unassignPermissionFromUser(
  userId: string,
  permissionId: string
): Promise<void> {
  try {
    await api.delete(`/v1/users/${userId}/permissions/${permissionId}`);
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}

// ── Audit Logs ────────────────────────────────────

export async function listAuditLogs(params?: {
  actor?: string;
  action?: string;
  resource?: string;
  limit?: number;
  offset?: number;
}): Promise<AuditLog[]> {
  try {
    const query = new URLSearchParams();
    if (params?.actor) query.set("actor", params.actor);
    if (params?.action) query.set("action", params.action);
    if (params?.resource) query.set("resource", params.resource);
    if (params?.limit) query.set("limit", String(params.limit));
    if (params?.offset) query.set("offset", String(params.offset));
    const response = await api.get(`/v1/audit-logs?${query.toString()}`);
    return response.data as AuditLog[];
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}
