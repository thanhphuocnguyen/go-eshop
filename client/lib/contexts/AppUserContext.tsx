'use client';

import { createContext, useContext, useEffect, useState } from 'react';
import { getCookie } from 'cookies-next/client';
import { UserModel } from '@/lib/definitions';
import { useUser } from '@/lib/hooks/useUser';
import { KeyedMutator } from 'swr';

// Define types for cart items and cart

interface AppUserContextType {
  isLoading: boolean;
  user: UserModel | undefined;
  mutateUser: KeyedMutator<UserModel>;
  setIsLoggedIn: (isLoggedIn: boolean) => void;
  isLoggedIn: boolean;
}

const AppUserContext = createContext<AppUserContextType | undefined>(undefined);
export const AppUserContextConsumer = AppUserContext.Consumer;

export function AppUserProvider({ children }: { children: React.ReactNode }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const { user, mutateUser, isLoading } = useUser(isLoggedIn);

  const value = {
    user,
    isLoading,
    mutateUser,
    isLoggedIn,
    setIsLoggedIn,
  } satisfies AppUserContextType;
  useEffect(() => {
    if (getCookie('access_token')) {
      setIsLoggedIn(true);
    }
  }, []);
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
