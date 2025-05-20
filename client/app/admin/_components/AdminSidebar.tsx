'use client';

import React from 'react';
import {
  CircleStackIcon,
  CubeIcon,
  FireIcon,
  HomeIcon,
  NewspaperIcon,
  RectangleStackIcon,
  ShoppingCartIcon,
  UsersIcon,
  StarIcon,
  TagIcon,
} from '@heroicons/react/24/outline';
import Image from 'next/image';
import Link from 'next/link';
import clsx from 'clsx';
import { usePathname } from 'next/navigation';

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
  {
    name: 'Discounts',
    icon: TagIcon,
    path: '/admin/discounts',
    href: '/admin/discounts',
  },
  {
    name: 'Ratings',
    icon: StarIcon,
    path: '/admin/ratings',
    href: '/admin/ratings',
  },
];

export default function AdminSidebar() {
  const pathname = usePathname();
  
  return (
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
  );
}