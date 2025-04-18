import { UserModel } from '../types/user';

export const extractUserFromCookie = (
  cookie: string | undefined
): UserModel | null => {
  if (!cookie) return null;
  const user = JSON.parse(cookie);
  return user as UserModel;
};
