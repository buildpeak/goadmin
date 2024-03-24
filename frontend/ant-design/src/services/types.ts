export interface SignUpRequest {
  username: string;
  password: string;
  first_name: string;
  last_name: string;
  agreement: boolean;
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
