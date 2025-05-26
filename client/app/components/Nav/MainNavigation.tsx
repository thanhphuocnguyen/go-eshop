'use client';

import { Fragment, useEffect, useState } from 'react';
import { Popover, Transition } from '@headlessui/react';
import {
  ChevronDownIcon,
  Square3Stack3DIcon,
  TagIcon,
  BuildingStorefrontIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import Image from 'next/image';

interface Category {
  id: string;
  name: string;
  slug: string;
  imageUrl?: string;
}

interface Collection {
  id: string;
  name: string;
  slug: string;
  imageUrl?: string;
}

interface Brand {
  id: string;
  name: string;
  slug: string;
  logo_url?: string;
}

export default function MainNavigation() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [collections, setCollections] = useState<Collection[]>([]);
  const [brands, setBrands] = useState<Brand[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    async function fetchData() {
      setIsLoading(true);
      
      try {
        // Fetch categories
        const categoriesResponse = await clientSideFetch<Category[]>(
          `${PUBLIC_API_PATHS.CATEGORIES}?limit=10`
        );
        if (categoriesResponse.data) {
          setCategories(categoriesResponse.data);
        }
        
        // Fetch collections
        const collectionsResponse = await clientSideFetch<Collection[]>(
          `${PUBLIC_API_PATHS.COLLECTIONS}?limit=8`
        );
        if (collectionsResponse.data) {
          setCollections(collectionsResponse.data);
        }
        
        // Fetch brands
        const brandsResponse = await clientSideFetch<Brand[]>(
          `${PUBLIC_API_PATHS.BRANDS}?limit=8`
        );
        if (brandsResponse.data) {
          setBrands(brandsResponse.data);
        }
      } catch (error) {
        console.error('Failed to fetch navigation data', error);
      } finally {
        setIsLoading(false);
      }
    }
    
    fetchData();
  }, []);
  
  return (
    <div className="ml-8 flex space-x-8">
      {/* Categories Dropdown */}
      <NavDropdown 
        title="Categories" 
        items={categories}
        viewAllLink="/categories"
        icon={<Square3Stack3DIcon className="h-5 w-5 text-indigo-600" />}
        isLoading={isLoading}
      />
      
      {/* Collections Dropdown */}
      <NavDropdown 
        title="Collections" 
        items={collections}
        viewAllLink="/collections"
        icon={<TagIcon className="h-5 w-5 text-indigo-600" />}
        isLoading={isLoading}
      />
      
      {/* Brands Dropdown */}
      <NavDropdown 
        title="Brands" 
        items={brands}
        viewAllLink="/brands"
        icon={<BuildingStorefrontIcon className="h-5 w-5 text-indigo-600" />}
        isLoading={isLoading}
      />
      
      {/* Direct link to all products */}
      <Link 
        href="/products" 
        className="flex items-center text-sm font-medium text-gray-700 hover:text-gray-800"
      >
        All Products
      </Link>
    </div>
  );
}

interface NavDropdownProps {
  title: string;
  items: { id: string; name: string; slug: string; imageUrl?: string; logo_url?: string }[];
  viewAllLink: string;
  icon: React.ReactNode;
  isLoading: boolean;
}

function NavDropdown({ title, items, viewAllLink, icon, isLoading }: NavDropdownProps) {
  return (
    <Popover className="relative">
      {({ open }) => (
        <>
          <Popover.Button 
            className={`
              ${open ? 'text-indigo-600' : 'text-gray-700'}
              group inline-flex items-center rounded-md text-sm font-medium hover:text-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2
            `}
          >
            <span>{title}</span>
            <ChevronDownIcon
              className={`
                ${open ? 'text-indigo-600 rotate-180' : 'text-gray-400'}
                ml-2 h-5 w-5 transition-transform duration-200
              `}
              aria-hidden="true"
            />
          </Popover.Button>

          <Transition
            as={Fragment}
            enter="transition ease-out duration-200"
            enterFrom="opacity-0 translate-y-1"
            enterTo="opacity-100 translate-y-0"
            leave="transition ease-in duration-150"
            leaveFrom="opacity-100 translate-y-0"
            leaveTo="opacity-0 translate-y-1"
          >
            <Popover.Panel className="absolute left-1/2 z-10 mt-3 w-64 max-w-sm -translate-x-1/2 transform px-2 sm:px-0">
              <div className="overflow-hidden rounded-lg shadow-lg ring-1 ring-black ring-opacity-5">
                <div className="relative grid gap-6 bg-white p-6">
                  {isLoading ? (
                    <div className="space-y-2">
                      {[...Array(4)].map((_, idx) => (
                        <div key={idx} className="h-6 bg-gray-200 rounded animate-pulse"></div>
                      ))}
                    </div>
                  ) : items.length > 0 ? (
                    items.map((item) => (
                      <Link
                        key={item.id}
                        href={`${viewAllLink}/${item.slug}`}
                        className="flex items-center p-2 -m-3 transition duration-150 ease-in-out rounded-lg hover:bg-gray-50"
                      >
                        {(item.imageUrl || item.logo_url) && (
                          <div className="flex-shrink-0 flex items-center justify-center h-10 w-10 rounded-md bg-indigo-50 text-indigo-600">
                            <Image
                              src={item.imageUrl || item.logo_url || ''}
                              alt={item.name}
                              width={30}
                              height={30}
                              className="object-contain"
                            />
                          </div>
                        )}
                        <div className="ml-4">
                          <p className="text-sm font-medium text-gray-900">{item.name}</p>
                        </div>
                      </Link>
                    ))
                  ) : (
                    <p className="text-sm text-gray-500 text-center py-4">No items available</p>
                  )}
                </div>
                <div className="bg-gray-50 p-4">
                  <Link
                    href={viewAllLink}
                    className="flow-root px-2 py-2 transition duration-150 ease-in-out rounded-md hover:bg-gray-100"
                  >
                    <span className="flex items-center">
                      {icon}
                      <span className="ml-3 text-sm font-medium text-gray-900">View all {title}</span>
                    </span>
                  </Link>
                </div>
              </div>
            </Popover.Panel>
          </Transition>
        </>
      )}
    </Popover>
  );
}
