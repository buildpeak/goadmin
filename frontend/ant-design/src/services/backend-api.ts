import axios, { isAxiosError } from "axios";
import { JwtToken, SignUpRequest, User } from "./types";

const backendUrl = process.env.REACT_APP_BACKEND_URL;

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
      case 401:
        window.location.href = "/login";
        break;
      case 400:
        throw new Error("Invalid input");
      default:
        throw new ApiError(error.response?.data);
    }
  }
}

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
    // if error is 404, redirect to signup page
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
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
  } catch (error) {
    handleAxiosError(error);
    throw error;
  }
}
