'use client';

import React from 'react';
import { Tab, TabList } from '@headlessui/react';
import { 
  TagIcon, 
  ShoppingBagIcon, 
  ClipboardDocumentListIcon 
} from '@heroicons/react/24/outline';

interface TabNavigationProps {
  activeTab: string;
  setActiveTab: React.Dispatch<React.SetStateAction<string>>;
}

export const TabNavigation: React.FC<TabNavigationProps> = ({ 
  activeTab, 
  setActiveTab 
}) => {
  return (
    <TabList className='flex border-b mb-2'>
      <Tab
        onClick={() => setActiveTab('details')}
        className={`px-4 py-2 text-sm font-medium focus:outline-none ${
          activeTab === 'details'
            ? 'border-b-2 border-primary text-primary'
            : 'text-gray-500 hover:text-gray-700'
        }`}
      >
        <span className='flex items-center'>
          <TagIcon className='h-4 w-4 mr-2' />
          Details
        </span>
      </Tab>
      <Tab
        onClick={() => setActiveTab('products')}
        className={`px-4 py-2 text-sm font-medium focus:outline-none ${
          activeTab === 'products'
            ? 'border-b-2 border-primary text-primary'
            : 'text-gray-500 hover:text-gray-700'
        }`}
      >
        <span className='flex items-center'>
          <ShoppingBagIcon className='h-4 w-4 mr-2' />
          Products
        </span>
      </Tab>
      <Tab
        onClick={() => setActiveTab('categories')}
        className={`px-4 py-2 text-sm font-medium focus:outline-none ${
          activeTab === 'categories'
            ? 'border-b-2 border-primary text-primary'
            : 'text-gray-500 hover:text-gray-700'
        }`}
      >
        <span className='flex items-center'>
          <ClipboardDocumentListIcon className='h-4 w-4 mr-2' />
          Categories
        </span>
      </Tab>
    </TabList>
  );
};
