'use client';
import {
  Dialog,
  DialogBackdrop,
  DialogPanel,
  Tab,
  TabGroup,
  TabList,
  TabPanel,
  TabPanels,
} from '@headlessui/react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';
import React, { Fragment } from 'react';
import clsx from 'clsx';

interface CategoriesDialogProps {
  open: boolean;
  setOpen: (open: boolean) => void;
  navigation: any;
}
const CategoriesDialog: React.FC<CategoriesDialogProps> = ({
  open,
  setOpen,
  navigation,
}) => {
  return (
    <div>
      {/* Mobile menu */}
      <Dialog open={open} onClose={setOpen} className='relative z-40 lg:hidden'>
        <DialogBackdrop
          transition
          className={clsx(
            'fixed inset-0 bg-black/25 transition-opacity duration-300 ease-linear data-closed:opacity-0'
          )}
        />

        <div className={clsx('fixed inset-0 z-40 flex')}>
          <DialogPanel
            transition
            className={clsx(
              'relative flex w-full max-w-xs transform flex-col overflow-y-auto bg-white pb-12 shadow-xl transition duration-300 ease-in-out data-closed:-translate-x-full'
            )}
          >
            <div className={clsx('flex px-4 pt-5 pb-2')}>
              <button
                type='button'
                onClick={() => setOpen(false)}
                className={clsx(
                  'relative -m-2 inline-flex items-center justify-center rounded-md p-2 text-gray-400'
                )}
              >
                <span className={clsx('absolute -inset-0.5')} />
                <span className={clsx('sr-only')}>Close menu</span>
                <XMarkIcon aria-hidden='true' className={clsx('size-6')} />
              </button>
            </div>

            {/* Links */}
            <TabGroup className={clsx('mt-2')}>
              <div className={clsx('border-b border-gray-200')}>
                <TabList className={clsx('-mb-px flex space-x-8 px-4')}>
                  {navigation.categories.map((category) => (
                    <Tab
                      key={category.name}
                      className={clsx(
                        'flex-1 border-b-2 border-transparent px-1 py-4 text-base font-medium whitespace-nowrap text-gray-900 data-selected:border-indigo-600 data-selected:text-indigo-600'
                      )}
                    >
                      {category.name}
                    </Tab>
                  ))}
                </TabList>
              </div>
              <TabPanels as={Fragment}>
                {navigation.categories.map((category) => (
                  <TabPanel
                    key={category.name}
                    className={clsx('space-y-10 px-4 pt-10 pb-8')}
                  >
                    <div className={clsx('grid grid-cols-2 gap-x-4')}>
                      {category.featured.map((item) => (
                        <div key={item.name} className={clsx('group relative text-sm')}>
                          <Image
                            alt={item.imageAlt}
                            src={item.imageSrc}
                            className={clsx(
                              'aspect-square w-full rounded-lg bg-gray-100 object-cover group-hover:opacity-75'
                            )}
                          />
                          <a
                            href={item.href}
                            className={clsx('mt-6 block font-medium text-gray-900')}
                          >
                            <span
                              aria-hidden='true'
                              className={clsx('absolute inset-0 z-10')}
                            />
                            {item.name}
                          </a>
                          <p aria-hidden='true' className={clsx('mt-1')}>
                            Shop now
                          </p>
                        </div>
                      ))}
                    </div>
                    {category.sections.map((section) => (
                      <div key={section.name}>
                        <p
                          id={`${category.id}-${section.id}-heading-mobile`}
                          className={clsx('font-medium text-gray-900')}
                        >
                          {section.name}
                        </p>
                        <ul
                          role='list'
                          aria-labelledby={`${category.id}-${section.id}-heading-mobile`}
                          className={clsx('mt-6 flex flex-col space-y-6')}
                        >
                          {section.items.map((item) => (
                            <li key={item.name} className={clsx('flow-root')}>
                              <a
                                href={item.href}
                                className={clsx('-m-2 block p-2 text-gray-500')}
                              >
                                {item.name}
                              </a>
                            </li>
                          ))}
                        </ul>
                      </div>
                    ))}
                  </TabPanel>
                ))}
              </TabPanels>
            </TabGroup>

            <div className={clsx('space-y-6 border-t border-gray-200 px-4 py-6')}>
              {navigation.pages.map((page) => (
                <div key={page.name} className={clsx('flow-root')}>
                  <a
                    href={page.href}
                    className={clsx('-m-2 block p-2 font-medium text-gray-900')}
                  >
                    {page.name}
                  </a>
                </div>
              ))}
            </div>

            <div className={clsx('space-y-6 border-t border-gray-200 px-4 py-6')}>
              <div className={clsx('flow-root')}>
                <a
                  href='#'
                  className={clsx('-m-2 block p-2 font-medium text-gray-900')}
                >
                  Sign in
                </a>
              </div>
              <div className={clsx('flow-root')}>
                <a
                  href='#'
                  className={clsx('-m-2 block p-2 font-medium text-gray-900')}
                >
                  Create account
                </a>
              </div>
            </div>

            <div className={clsx('border-t border-gray-200 px-4 py-6')}>
              <a href='#' className={clsx('-m-2 flex items-center p-2')}>
                <img
                  alt=''
                  src='https://tailwindui.com/plus/img/flags/flag-canada.svg'
                  className={clsx('block h-auto w-5 shrink-0')}
                />
                <span className={clsx('ml-3 block text-base font-medium text-gray-900')}>
                  CAD
                </span>
                <span className={clsx('sr-only')}>, change currency</span>
              </a>
            </div>
          </DialogPanel>
        </div>
      </Dialog>
    </div>
  );
};

export default CategoriesDialog;
