export interface SignUpRequest {
  username: string;
  password: string;
  first_name: string;
  last_name: string;
  agreement: boolean;
  recaptcha_token?: string;
}

export interface User {
  id: string;
  username: string;
  first_name: string;
  last_name: string;
  email: string;
  picture: string;
  active: boolean;
  deleted_at: Date | null;
}

export interface JwtToken {
  access_token: string;
  refresh_token: string;
}

export interface Role {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Permission {
  id: string;
  name: string;
  rule_type: string;
  rule: string;
  created_at: string;
  updated_at: string;
}

export interface CreateRoleRequest {
  name: string;
}

export interface CreatePermissionRequest {
  name: string;
  rule_type: string;
  rule: string;
}

export interface AuditLog {
  id: string;
  actor_id: string;
  actor_username: string;
  action: string;
  resource_type: string;
  resource_id: string;
  detail: string;
  ip_address: string;
  created_at: string;
}
