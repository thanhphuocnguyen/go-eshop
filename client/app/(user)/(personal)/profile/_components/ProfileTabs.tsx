'use client';

import clsx from 'clsx';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

export default function ProfileTabs() {
  const pathname = usePathname();
  
  const tabs = [
    { id: 'profile', path: '/profile', label: 'Personal Info' },
    { id: 'addresses', path: '/profile/addresses', label: 'Addresses' },
    { id: 'security', path: '/profile/security', label: 'Security' },
  ] as const;

  const isTabActive = (tabPath: string) => {
    if (tabPath === '/profile' && pathname === '/profile') {
      return true;
    }
    if (tabPath !== '/profile' && pathname.startsWith(tabPath)) {
      return true;
    }
    return false;
  };

  return (
    <div className="border-b border-gray-200">
      <nav className="flex -mb-px">
        {tabs.map((tab) => (
          <Link
            key={tab.id}
            href={tab.path}
            className={clsx(
              'w-1/3 py-4 px-1 text-center border-b-2 font-medium text-sm',
              isTabActive(tab.path)
                ? 'border-indigo-500 text-indigo-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            )}
          >
            {tab.label}
          </Link>
        ))}
      </nav>
    </div>
  );
}