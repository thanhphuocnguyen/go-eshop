import { CategoryProductModel } from '@/lib/definitions';
import { Button } from '@headlessui/react';
import Image from 'next/image';
import React from 'react';

interface CategoryProductList {
  products: CategoryProductModel[];
}
const CategoryProductList: React.FC<CategoryProductList> = ({ products }) => {
  return (
    <div className='px-10 h-full'>
      <h2 className='text-2xl font-semibold text-gray-600 my-4'>
        Products Linked
      </h2>
      <div className='flex bg-tableRow h-full w-full flex-col space-y-2'>
        {products?.map((e) => (
          <div
            className='bg-gray-50 flex justify-between gap-x-2 items-center w-full relative p-2 border-2 border-green-500 rounded-md'
            key={e.id}
          >
            <div className='flex p-4 gap-10 items-center'>
              <Image
                height={100}
                width={120}
                alt='product-image'
                objectFit='cover'
                className='rounded-md border border-lime-300'
                src={e.image_url || '/images/product-placeholder.webp'}
              />
              <div>
                <h3 className='text-gray-800 font-bold'>{e.name}</h3>
                <p
                  className='text-gray-600'
                  dangerouslySetInnerHTML={{ __html: e.description || '' }}
                ></p>
              </div>
            </div>
            <div className='flex mr-4'>
              <Button className='btn btn-danger'>Remove</Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default CategoryProductList;
