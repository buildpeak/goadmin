import { JwtPayload } from "jwt-decode";

export interface GoogleJwtPayload extends JwtPayload {
  email: string;
  email_verified: boolean;
  exp: number;
  family_name: string;
  given_name: string;
  iat: number;
  iss: string;
  locale: string;
  name: string;
  picture: string;
  sub: string;
}
