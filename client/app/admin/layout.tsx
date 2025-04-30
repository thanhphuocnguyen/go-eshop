'use client';

import React, { useEffect } from 'react';
import {
  BellIcon,
  CircleStackIcon,
  CubeIcon,
  FireIcon,
  HomeIcon,
  NewspaperIcon,
  RectangleStackIcon,
  ShoppingCartIcon,
  UsersIcon,
} from '@heroicons/react/24/outline';
import './admin.css';
import Image from 'next/image';
import Link from 'next/link';
import clsx from 'clsx';
import { redirect, usePathname } from 'next/navigation';
import { ChevronUpIcon } from '@heroicons/react/16/solid';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { CookieValueTypes, deleteCookie, getCookie } from 'cookies-next';

export default function Layout({ children }: { children: React.ReactNode }) {
  const logout = async () => {
    deleteCookie('access_token');
    deleteCookie('refresh_token');
    deleteCookie('session_id');
    deleteCookie('user_role');
    deleteCookie('user_id');
    deleteCookie('user_name');
    redirect('/login');
  };

  const pathname = usePathname();
  const [fullName, setFullName] = React.useState<string | null>(null);

  useEffect(() => {
    setFullName((getCookie('user_name') as CookieValueTypes) ?? '');
  }, []);

  return (
    <div className='flex h-screen'>
      <nav className='p-3 bg-primary flex w-1/6 flex-col shadow-sm'>
        <Link href={'/'}>
          <div className='mb-4 p-3 rounded-md relative flex items-center gap-3 hover:bg-secondary'>
            <Image
              className='object-cover rounded-lg'
              src='/images/logos/logo.webp'
              alt='logo'
              width={40}
              height={40}
            />
            <div className='text-xl font-bold text-white'>Simple Life</div>
          </div>
        </Link>

        <div className='text-lg font-bold my-1 text-white'>Settings</div>
        <ul className='flex flex-col gap-1'>
          {NavBarItems.map((item) => (
            <Link href={item.href} key={item.path}>
              <li
                className={clsx(
                  'side-bar-item',
                  pathname === item.path && 'side-bar-item-active'
                )}
              >
                <item.icon height={20} width={20} />
                {item.name}
              </li>
            </Link>
          ))}
        </ul>
      </nav>
      <main className='w-5/6 block bg-white'>
        <section className='flex items-center h-[60px] p-3 border-b shadow-sm text-black justify-end'>
          <div className='flex items-center gap-3'>
            <div>
              <BellIcon height={20} width={20} />
            </div>
            <div className='border-l h-10 border-gray-600' />
            <Menu as={'section'}>
              <MenuButton
                className={
                  'w-52 justify-between flex gap-2 hover:bg-gray-200 rounded-md items-center p-2'
                }
              >
                {({ active }) => (
                  <>
                    <div className='inline-flex items-center gap-1'>
                      <div className='h-8 w-8 rounded-full bg-gray-600'></div>
                      <span className='text-md text-ellipsis font-semibold'>
                        {fullName ?? ''}
                      </span>
                    </div>
                    <span
                      className={clsx(
                        'transition-all',
                        active ? 'transform rotate-180' : 'transform rotate-0'
                      )}
                    >
                      <ChevronUpIcon height={20} width={20} />
                    </span>
                  </>
                )}
              </MenuButton>
              <MenuItems
                anchor='bottom'
                className={'bg-white shadow-lg rounded-lg w-52'}
              >
                <MenuItem>
                  <Link
                    href={'/profile'}
                    className='block w-full text-left font-semibold p-3 cursor-pointer data-[focus]:bg-green-500 border-b'
                  >
                    Profile
                  </Link>
                </MenuItem>
                <MenuItem>
                  <button
                    onClick={logout}
                    className='block w-full text-left font-semibold p-3 cursor-pointer data-[focus]:bg-green-500'
                  >
                    Logout
                  </button>
                </MenuItem>
              </MenuItems>
            </Menu>
          </div>
        </section>
        <section className='p-4 m-5 border-2 border-gray-200 shadow-sm flex flex-col h-[90%] rounded-md'>
          {children}
        </section>
      </main>
    </div>
  );
}

const NavBarItems = [
  {
    name: 'Dashboard',
    icon: HomeIcon,
    path: '/admin/dashboard',
    href: '/admin/dashboard',
  },
  {
    name: 'Users',
    icon: UsersIcon,
    path: '/admin/users',
    href: '/admin/users',
  },
  {
    name: 'Categories',
    icon: CircleStackIcon,
    path: '/admin/categories',
    href: '/admin/categories',
  },
  {
    name: 'Products',
    icon: FireIcon,
    path: '/admin/products',
    href: '/admin/products',
  },
  {
    name: 'Collections',
    icon: RectangleStackIcon,
    path: '/admin/collections',
    href: '/admin/collections',
  },
  {
    name: 'Brands',
    icon: NewspaperIcon,
    path: '/admin/brands',
    href: '/admin/brands',
  },
  {
    name: 'Attributes',
    icon: CubeIcon,
    path: '/admin/attributes',
    href: '/admin/attributes',
  },
  {
    name: 'Orders',
    icon: ShoppingCartIcon,
    path: '/admin/orders',
    href: '/admin/orders',
  },
];
