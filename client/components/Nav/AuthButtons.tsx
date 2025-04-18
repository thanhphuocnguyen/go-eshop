'use client';

import { logout } from '@/app/actions/auth';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { ChevronDownIcon, UserIcon } from '@heroicons/react/16/solid';
import { ArrowUpTrayIcon } from '@heroicons/react/24/outline';
import clsx from 'clsx';
import Link from 'next/link';
import React from 'react';

const AuthButtons: React.FC<{ role?: string; name?: string }> = ({
  role,
  name,
}) => {
  return name ? (
    <Menu>
      <MenuButton className='inline-flex items-center gap-2 rounded-md bg-gray-800 py-1.5 px-3 text-sm/6 font-semibold text-white shadow-inner shadow-white/10 focus:outline-none data-[hover]:bg-gray-700 data-[open]:bg-gray-700 data-[focus]:outline-1 data-[focus]:outline-white'>
        {({ active }) => (
          <>
            <span className='font-bold text-lg'>{name}</span>
            <span className={clsx('ml-2', active ? 'rotate-180' : 'rotate-0')}>
              <ChevronDownIcon height={24} width={24} />
            </span>
          </>
        )}
      </MenuButton>
      <MenuItems
        transition
        anchor='bottom end'
        className='w-52 z-50 origin-top-right rounded-xl border border-white/5 bg-gray-300 p-1 text-sm/6 text-white transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0'
      >
        {role === 'admin' && (
          <MenuItem>
            <Link
              href='/admin'
              className='group flex w-full items-center gap-2 rounded-lg py-1.5 px-3 data-[focus]:bg-white/10'
            >
              <UserIcon className='size-4 fill-gray-500' />
              Admin
              <kbd className='ml-auto hidden font-sans text-xs text-white/50 group-data-[focus]:inline'>
                ⌘E
              </kbd>
            </Link>
          </MenuItem>
        )}
        <MenuItem>
          <Link
            href='/profile'
            className='group flex w-full items-center gap-2 rounded-lg py-1.5 px-3 data-[focus]:bg-white/10'
          >
            <UserIcon className='size-4 fill-white/30' />
            Profile
            <kbd className='ml-auto hidden font-sans text-xs text-white/50 group-data-[focus]:inline'>
              ⌘E
            </kbd>
          </Link>
        </MenuItem>
        <MenuItem>
          <button
            className='group flex w-full items-center gap-2 rounded-lg py-1.5 px-3 data-[focus]:bg-white/10'
            onClick={logout}
          >
            <ArrowUpTrayIcon className='size-4' />
            Logout
          </button>
        </MenuItem>
      </MenuItems>
    </Menu>
  ) : (
    <>
      <Link
        href='/login'
        className='text-sm font-medium text-gray-700 hover:text-gray-800'
      >
        Sign in
      </Link>
      <span aria-hidden='true' className='h-6 w-px bg-gray-200' />
      <Link
        href='/register'
        className='text-sm font-medium text-gray-700 hover:text-gray-800'
      >
        Create account
      </Link>
    </>
  );
};

export default AuthButtons;
