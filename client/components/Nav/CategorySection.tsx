'use client';
import {
  Popover,
  PopoverButton,
  PopoverGroup,
  PopoverPanel,
} from '@headlessui/react';
import Image from 'next/image';
import React from 'react';

const navigation = {
  categories: [
    {
      id: 'women',
      name: 'Women',
      featured: [
        {
          name: 'New Arrivals',
          href: '#',
          imageSrc:
            'https://tailwindui.com/plus/img/ecommerce-images/mega-menu-category-01.jpg',
          imageAlt:
            'Models sitting back to back, wearing Basic Tee in black and bone.',
        },
        {
          name: 'Basic Tees',
          href: '#',
          imageSrc:
            'https://tailwindui.com/plus/img/ecommerce-images/mega-menu-category-02.jpg',
          imageAlt:
            'Close up of Basic Tee fall bundle with off-white, ochre, olive, and black tees.',
        },
      ],
      sections: [
        {
          id: 'clothing',
          name: 'Clothing',
          items: [
            { name: 'Tops', href: '#' },
            { name: 'Dresses', href: '#' },
            { name: 'Pants', href: '#' },
            { name: 'Denim', href: '#' },
            { name: 'Sweaters', href: '#' },
            { name: 'T-Shirts', href: '#' },
            { name: 'Jackets', href: '#' },
            { name: 'Activewear', href: '#' },
            { name: 'Browse All', href: '#' },
          ],
        },
        {
          id: 'accessories',
          name: 'Accessories',
          items: [
            { name: 'Watches', href: '#' },
            { name: 'Wallets', href: '#' },
            { name: 'Bags', href: '#' },
            { name: 'Sunglasses', href: '#' },
            { name: 'Hats', href: '#' },
            { name: 'Belts', href: '#' },
          ],
        },
        {
          id: 'brands',
          name: 'Brands',
          items: [
            { name: 'Full Nelson', href: '#' },
            { name: 'My Way', href: '#' },
            { name: 'Re-Arranged', href: '#' },
            { name: 'Counterfeit', href: '#' },
            { name: 'Significant Other', href: '#' },
          ],
        },
      ],
    },
    {
      id: 'men',
      name: 'Men',
      featured: [
        {
          name: 'New Arrivals',
          href: '#',
          imageSrc:
            'https://tailwindui.com/plus/img/ecommerce-images/product-page-04-detail-product-shot-01.jpg',
          imageAlt:
            'Drawstring top with elastic loop closure and textured interior padding.',
        },
        {
          name: 'Artwork Tees',
          href: '#',
          imageSrc:
            'https://tailwindui.com/plus/img/ecommerce-images/category-page-02-image-card-06.jpg',
          imageAlt:
            'Three shirts in gray, white, and blue arranged on table with same line drawing of hands and shapes overlapping on front of shirt.',
        },
      ],
      sections: [
        {
          id: 'clothing',
          name: 'Clothing',
          items: [
            { name: 'Tops', href: '#' },
            { name: 'Pants', href: '#' },
            { name: 'Sweaters', href: '#' },
            { name: 'T-Shirts', href: '#' },
            { name: 'Jackets', href: '#' },
            { name: '', href: '#' },
            { name: 'Browse All', href: '#' },
          ],
        },
        {
          id: 'accessories',
          name: 'Accessories',
          items: [
            { name: 'Watches', href: '#' },
            { name: 'Wallets', href: '#' },
            { name: 'Bags', href: '#' },
            { name: 'Sunglasses', href: '#' },
            { name: 'Hats', href: '#' },
            { name: 'Belts', href: '#' },
          ],
        },
        {
          id: 'brands',
          name: 'Brands',
          items: [
            { name: 'Re-Arranged', href: '#' },
            { name: 'Counterfeit', href: '#' },
            { name: 'Full Nelson', href: '#' },
            { name: 'My Way', href: '#' },
          ],
        },
      ],
    },
  ],
  pages: [
    { name: 'Company', href: '#' },
    { name: 'Stores', href: '#' },
  ],
};
const CategorySection = () => {
  return (
    <>
      {/* Flyout menus */}
      <PopoverGroup className='hidden lg:ml-8 lg:block lg:self-stretch'>
        <div className='flex h-full space-x-8'>
          {navigation.categories.map((category) => (
            <Popover key={category.name} className='flex'>
              <div className='relative flex'>
                <PopoverButton className='relative z-10 -mb-px flex items-center border-b-2 border-transparent pt-px text-sm font-medium text-gray-700 transition-colors duration-200 ease-out hover:text-gray-800 data-open:border-indigo-600 data-open:text-indigo-600'>
                  {category.name}
                </PopoverButton>
              </div>

              <PopoverPanel
                transition
                className='absolute inset-x-0 top-full text-sm text-gray-500 transition data-closed:opacity-0 data-enter:duration-200 data-enter:ease-out data-leave:duration-150 data-leave:ease-in'
              >
                {/* Presentational element used to render the bottom shadow, if we put the shadow on the actual panel it pokes out the top, so we use this shorter element to hide the top of the shadow */}
                <div
                  aria-hidden='true'
                  className='absolute inset-0 top-1/2 bg-white shadow-sm'
                />

                <div className='relative bg-white'>
                  <div className='mx-auto max-w-7xl px-8'>
                    <div className='grid grid-cols-2 gap-x-8 gap-y-10 py-16'>
                      <div className='col-start-2 grid grid-cols-2 gap-x-8'>
                        {category.featured.map((item) => (
                          <div
                            key={item.name}
                            className='group relative text-base sm:text-sm'
                          >
                            <Image
                            width={200}
                            height={200}
                              alt={item.imageAlt}
                              src={item.imageSrc}
                              className='aspect-square w-full rounded-lg bg-gray-100 object-cover group-hover:opacity-75'
                            />
                            <a
                              href={item.href}
                              className='mt-6 block font-medium text-gray-900'
                            >
                              <span
                                aria-hidden='true'
                                className='absolute inset-0 z-10'
                              />
                              {item.name}
                            </a>
                            <p aria-hidden='true' className='mt-1'>
                              Shop now
                            </p>
                          </div>
                        ))}
                      </div>
                      <div className='row-start-1 grid grid-cols-3 gap-x-8 gap-y-10 text-sm'>
                        {category.sections.map((section) => (
                          <div key={section.name}>
                            <p
                              id={`${section.name}-heading`}
                              className='font-medium text-gray-900'
                            >
                              {section.name}
                            </p>
                            <ul
                              role='list'
                              aria-labelledby={`${section.name}-heading`}
                              className='mt-6 space-y-6 sm:mt-4 sm:space-y-4'
                            >
                              {section.items.map((item) => (
                                <li key={item.name} className='flex'>
                                  <a
                                    href={item.href}
                                    className='hover:text-gray-800'
                                  >
                                    {item.name}
                                  </a>
                                </li>
                              ))}
                            </ul>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                </div>
              </PopoverPanel>
            </Popover>
          ))}

          {navigation.pages.map((page) => (
            <a
              key={page.name}
              href={page.href}
              className='flex items-center text-sm font-medium text-gray-700 hover:text-gray-800'
            >
              {page.name}
            </a>
          ))}
        </div>
      </PopoverGroup>
    </>
  );
};

export default CategorySection;
