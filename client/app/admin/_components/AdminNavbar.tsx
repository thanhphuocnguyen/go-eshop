'use client';
import React from 'react';
import {
  BellIcon,
  UserCircleIcon,
  ArrowRightOnRectangleIcon,
} from '@heroicons/react/24/outline';
import { ChevronUpIcon } from '@heroicons/react/16/solid';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { redirect, useRouter } from 'next/navigation';
import clsx from 'clsx';
import Link from 'next/link';
import { logoutAction } from '@/app/actions/auth';
import { useUser } from '@/app/hooks';

export default function AdminNavbar() {
  const { user, mutateUser } = useUser();
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

  return (
    <section className='flex items-center h-[60px] p-3 border-b shadow-sm text-black justify-end'>
      <div className='flex items-center gap-3'>
        <div className='relative'>
          <button className='p-2 hover:bg-gray-100 rounded-full transition-colors'>
            <BellIcon height={20} width={20} className='text-gray-600' />
          </button>
          <span className='absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-4 h-4 flex items-center justify-center'>
            2
          </span>
        </div>
        <div className='border-l h-10 border-gray-300' />
        <Menu as={'section'} className='relative'>
          <MenuButton className='account-button'>
            {({ open }) => (
              <>
                <div className='account-avatar'>
                  {user?.fullname ? (
                    <span className='font-medium'>
                      {user.fullname.charAt(0).toUpperCase()}
                    </span>
                  ) : (
                    <UserCircleIcon height={16} width={16} />
                  )}
                </div>
                <div className='flex flex-col items-start'>
                  <span className='text-sm font-semibold truncate max-w-[120px]'>
                    {user?.fullname || 'User'}
                  </span>
                  <span className='text-xs text-gray-500 truncate max-w-[120px]'>
                    {user?.email || 'user@example.com'}
                  </span>
                </div>
                <ChevronUpIcon
                  height={16}
                  width={16}
                  className={clsx(
                    'text-gray-500 transition-transform duration-200 menu-transition',
                    open ? 'rotate-180' : 'rotate-0'
                  )}
                />
              </>
            )}
          </MenuButton>
          <MenuItems className='menu-dropdown absolute right-0 w-56 mt-2'>
            <MenuItem>
              {({ active }) => (
                <Link
                  href='/profile'
                  className={clsx(
                    'menu-item group',
                    active && 'bg-green-100 text-green-700'
                  )}
                >
                  <UserCircleIcon height={18} width={18} />
                  My Profile
                </Link>
              )}
            </MenuItem>
            <MenuItem>
              {({ active }) => (
                <button
                  onClick={logout}
                  className={clsx(
                    'menu-item group',
                    active && 'bg-green-100 text-green-700'
                  )}
                >
                  <ArrowRightOnRectangleIcon height={18} width={18} />
                  Sign Out
                </button>
              )}
            </MenuItem>
          </MenuItems>
        </Menu>
      </div>
    </section>
  );
}
