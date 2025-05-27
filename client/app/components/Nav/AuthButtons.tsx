'use client';

import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { ChevronDownIcon, UserIcon } from '@heroicons/react/16/solid';
import { ArrowUpTrayIcon, ShoppingBagIcon } from '@heroicons/react/24/outline';
import clsx from 'clsx';
import Link from 'next/link';
import { redirect, useRouter } from 'next/navigation';
import React from 'react';
import { logoutAction } from '@/app/actions/auth';
import { useUser } from '@/app/hooks';

const AuthButtons: React.FC = ({}) => {
  const { user, isLoading, mutateUser } = useUser();
  const router = useRouter();

  const logout = async () => {
    await logoutAction();
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('sessionId');
    mutateUser(undefined, { revalidate: false });
    router.refresh();
    redirect('/login');
  };
  if (isLoading) {
    return (
      <div className='flex items-center gap-2'>
        <span className='h-6 w-6 animate-spin rounded-full border-4 border-t-transparent border-gray-200' />
      </div>
    );
  }

  return (
    <>
      {user ? (
        <Menu>
          <MenuButton className='inline-flex items-center gap-2 rounded-md bg-indigo-600 py-2 px-4 text-sm font-medium text-white shadow-md hover:bg-indigo-700 transition-all duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2'>
            {({ active }) => (
              <>
                <span className='font-medium'>{user?.fullname}</span>
                <span
                  className={clsx(
                    'transition-transform duration-200',
                    active ? 'rotate-180' : 'rotate-0'
                  )}
                >
                  <ChevronDownIcon height={20} width={20} />
                </span>
              </>
            )}
          </MenuButton>
          <MenuItems
            transition
            anchor='bottom end'
            className='w-56 z-50 mt-2 origin-top-right rounded-lg border border-gray-200 bg-white p-1.5 shadow-lg ring-1 ring-black ring-opacity-5 transition duration-150 ease-out focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0'
          >
            {user?.role === 'admin' && (
              <MenuItem>
                <Link
                  href='/admin'
                  className='group flex w-full items-center gap-2.5 rounded-md py-2 px-3 text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-all duration-150'
                >
                  <UserIcon className='size-5 text-indigo-500' />
                  Admin
                  <kbd className='ml-auto hidden font-sans text-xs text-gray-400 group-hover:inline'>
                    ⌘E
                  </kbd>
                </Link>
              </MenuItem>
            )}
            <MenuItem>
              <Link
                href='/profile'
                className='group flex w-full items-center gap-2.5 rounded-md py-2 px-3 text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-all duration-150'
              >
                <UserIcon className='size-5 text-indigo-500' />
                Profile
                <kbd className='ml-auto hidden font-sans text-xs text-gray-400 group-hover:inline'>
                  ⌘E
                </kbd>
              </Link>
            </MenuItem>
            <MenuItem>
              <Link
                href='/orders'
                className='group flex w-full items-center gap-2.5 rounded-md py-2 px-3 text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-all duration-150'
              >
                <ShoppingBagIcon className='size-5 text-indigo-500' />
                Orders
              </Link>
            </MenuItem>
            <MenuItem>
              <button
                className='group flex w-full items-center gap-2.5 rounded-md py-2 px-3 text-gray-700 hover:bg-red-300 hover:text-red-600 transition-all duration-150'
                onClick={logout}
              >
                <ArrowUpTrayIcon className='size-5 text-white' />
                Logout
              </button>
            </MenuItem>
          </MenuItems>
        </Menu>
      ) : (
        <>
          <Link
            href='/login'
            className='text-sm font-medium text-gray-700 hover:text-indigo-600 hover:underline transition-colors duration-150'
          >
            Sign in
          </Link>
          <span aria-hidden='true' className='h-6 w-px bg-gray-200' />
          <Link
            href='/register'
            className='text-sm font-medium px-3 py-1.5 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors duration-150 shadow-sm'
          >
            Create account
          </Link>
        </>
      )}
    </>
  );
};

export default AuthButtons;
