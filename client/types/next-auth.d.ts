/* eslint-disable @typescript-eslint/no-unused-vars */
import { ErrorResponse } from '@/app/lib/definitions';
import NextAuth, { DefaultSession, DefaultJWT } from 'next-auth';
import { JWT } from 'next-auth/jwt';

declare module 'next-auth' {
  /**
   * Extend the built-in session types
   */
  interface Session {
    user: {
      id: string;
      role: string;
      accessToken: string;
    } & DefaultSession['user'];
  }

  /**
   * Extend the built-in user types
   */
  interface User {
    id: string;
    accessToken: string;
    accessTokenExpireAt: Date;
    refreshToken: string;
    refreshTokenExpireAt: Date;
    role: string;
    error?: ErrorResponse;
  }
}

declare module 'next-auth/jwt' {
  /**
   * Extend the built-in JWT types
   */
  interface JWT {
    role?: string;
    accessToken?: string;
  }
}
