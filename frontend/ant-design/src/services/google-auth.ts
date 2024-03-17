import axios, { isAxiosError } from "axios";

const backendUrl = process.env.REACT_APP_BACKEND_URL;

export async function verifyGoogleIdToken(idToken: string): Promise<void> {
  try {
    const response = await axios.post(`${backendUrl}/auth/signin-with-google`, {
      id_token: idToken,
    });
    return response.data;
  } catch (error) {
    // if error is 404, redirect to signup page
    if (isAxiosError(error) && error.response?.status === 404) {
      window.location.href = "/signup";
    }
    console.error(error);
  }
}
