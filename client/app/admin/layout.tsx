'use client';

import React from 'react';
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
  UserCircleIcon,
  ArrowRightOnRectangleIcon,
} from '@heroicons/react/24/outline';
import './admin.css';
import Image from 'next/image';
import Link from 'next/link';
import clsx from 'clsx';
import { redirect, usePathname } from 'next/navigation';
import { ChevronUpIcon } from '@heroicons/react/16/solid';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { deleteCookie } from 'cookies-next';
import { useAppUser } from '@/lib/contexts/AppUserContext';

export default function Layout({ children }: { children: React.ReactNode }) {
  const { user } = useAppUser();
  const logout = async () => {
    deleteCookie('access_token');
    deleteCookie('refresh_token');
    redirect('/login');
  };

  const pathname = usePathname();

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
            <div className='relative'>
              <button className='p-2 hover:bg-gray-100 rounded-full transition-colors'>
                <BellIcon height={20} width={20} className='text-gray-600' />
              </button>
              <span className='absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-4 h-4 flex items-center justify-center'>2</span>
            </div>
            <div className='border-l h-10 border-gray-300' />
            <Menu as={'section'} className="relative">
              <MenuButton className="account-button">
                {({ open }) => (
                  <>
                    <div className="account-avatar">
                      {user?.fullname ? (
                        <span className="font-medium">
                          {user.fullname.charAt(0).toUpperCase()}
                        </span>
                      ) : (
                        <UserCircleIcon height={16} width={16} />
                      )}
                    </div>
                    <div className="flex flex-col items-start">
                      <span className="text-sm font-semibold truncate max-w-[120px]">
                        {user?.fullname || 'User'}
                      </span>
                      <span className="text-xs text-gray-500 truncate max-w-[120px]">
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
              <MenuItems className="menu-dropdown absolute right-0 w-56 mt-2">
                <MenuItem>
                  {({ active }) => (
                    <Link
                      href="/profile"
                      className={clsx('menu-item group', active && 'bg-green-100 text-green-700')}
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
                      className={clsx('menu-item group', active && 'bg-green-100 text-green-700')}
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
