'use client';

import clsx from 'clsx';

type TabType = 'profile' | 'addresses' | 'orders' | 'security';

interface ProfileTabsProps {
  activeTab: TabType;
  setActiveTab: (tab: TabType) => void;
}

export default function ProfileTabs({ activeTab, setActiveTab }: ProfileTabsProps) {
  const tabs = [
    { id: 'profile', label: 'Personal Info' },
    { id: 'addresses', label: 'Addresses' },
    { id: 'orders', label: 'Orders' },
    { id: 'security', label: 'Security' },
  ] as const;

  return (
    <div className="border-b border-gray-200">
      <nav className="flex -mb-px">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as TabType)}
            className={clsx(
              'w-1/4 py-4 px-1 text-center border-b-2 font-medium text-sm',
              activeTab === tab.id
                ? 'border-indigo-500 text-indigo-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            )}
          >
            {tab.label}
          </button>
        ))}
      </nav>
    </div>
  );
}