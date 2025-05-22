'use client';

import React from 'react';
import { Tab, TabList } from '@headlessui/react';
import {
  TagIcon,
  ShoppingBagIcon,
  ClipboardDocumentListIcon,
} from '@heroicons/react/24/outline';
import { motion } from 'framer-motion';

interface TabNavigationProps {
  activeTab: string;
  setActiveTab: React.Dispatch<React.SetStateAction<string>>;
  isEditMode?: boolean;
}

export const TabNavigation: React.FC<TabNavigationProps> = ({
  activeTab,
  setActiveTab,
  isEditMode = false,
}) => {
  return (
    <TabList className='flex border-b mb-2 relative'>
      <Tab
        onClick={() => setActiveTab('details')}
        className={`px-4 py-2 text-sm font-medium focus:outline-none relative ${
          activeTab === 'details'
            ? 'text-primary'
            : 'text-gray-500 hover:text-gray-700'
        }`}
      >
        <span className='flex items-center'>
          <TagIcon className='h-4 w-4 mr-2' />
          Details
        </span>
        {activeTab === 'details' && (
          <motion.div
            className='absolute bottom-0 left-0 right-0 h-0.5 bg-primary'
            layoutId='activeTab'
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
          />
        )}
      </Tab>
      {isEditMode && (
        <>
          <Tab
            onClick={() => setActiveTab('products')}
            className={`px-4 py-2 text-sm font-medium focus:outline-none relative ${
              activeTab === 'products'
                ? 'text-primary'
                : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            <span className='flex items-center'>
              <ShoppingBagIcon className='h-4 w-4 mr-2' />
              Products
            </span>
            {activeTab === 'products' && (
              <motion.div
                className='absolute bottom-0 left-0 right-0 h-0.5 bg-primary'
                layoutId='activeTab'
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3 }}
              />
            )}
          </Tab>
          <Tab
            onClick={() => setActiveTab('categories')}
            className={`px-4 py-2 text-sm font-medium focus:outline-none relative ${
              activeTab === 'categories'
                ? 'text-primary'
                : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            <span className='flex items-center'>
              <ClipboardDocumentListIcon className='h-4 w-4 mr-2' />
              Categories
            </span>
            {activeTab === 'categories' && (
              <motion.div
                className='absolute bottom-0 left-0 right-0 h-0.5 bg-primary'
                layoutId='activeTab'
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3 }}
              />
            )}
          </Tab>
        </>
      )}
    </TabList>
  );
};
