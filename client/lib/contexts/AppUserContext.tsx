'use client';

import { createContext, useContext } from 'react';
import { getCookie } from 'cookies-next/client';
import { UserModel } from '@/lib/definitions';
import { useUser } from '@/lib/hooks/useUser';
import { KeyedMutator } from 'swr';

// Define types for cart items and cart

interface AppUserContextType {
  isLoading: boolean;
  user: UserModel | undefined;
  mutateUser: KeyedMutator<UserModel>;
}

const AppUserContext = createContext<AppUserContextType | undefined>(undefined);
export const AppUserContextConsumer = AppUserContext.Consumer;

export function AppUserProvider({ children }: { children: React.ReactNode }) {
  const { user, mutateUser, isLoading } = useUser(getCookie('access_token'));
  const value = {
    user,
    isLoading,
    mutateUser,
  } satisfies AppUserContextType;

  return (
    <AppUserContext.Provider value={value}>{children}</AppUserContext.Provider>
  );
}

export const useAppUser = (): AppUserContextType => {
  const context = useContext(AppUserContext);
  if (context === undefined) {
    throw new Error('useAppUser must be used within an AppUserProvider');
  }
  return context;
};
