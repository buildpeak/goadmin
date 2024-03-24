import axios, { isAxiosError } from "axios";
import { JwtToken, SignUpRequest, User } from "./types";

const backendUrl = process.env.REACT_APP_BACKEND_URL;

export function doesStatusMatch(error: unknown, status: number) {
  if (isAxiosError(error)) {
    return error.response?.status === status;
  }
  return false;
}

function handleAxiosError(error: unknown) {
  if (isAxiosError(error)) {
    if (doesStatusMatch(error, 401)) {
      window.location.href = "/login";
    }
    if (doesStatusMatch(error, 400)) {
      throw new Error("Invalid input");
    }
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
    if (isAxiosError(error) && error.response?.status === 404) {
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
