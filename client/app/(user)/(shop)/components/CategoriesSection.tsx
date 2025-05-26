import { GeneralCategoryModel } from '@/app/lib/definitions';
import { ArrowRightIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';
import Link from 'next/link';

export default function CategoriesSection({
  categories,
}: {
  categories?: GeneralCategoryModel[];
}) {
  return (
    <section className='pt-24'>
      <div className=' flex justify-between mb-2'>
        <h4 className='font-bold text-2xl'>Shop by Category</h4>
        <Link
          href={'/categories'}
          className='text-blue-500 flex space-x-2 items-center'
        >
          <span>Browse All Categories</span>
          <span>
            <ArrowRightIcon className='size-4' />
          </span>
        </Link>
      </div>
      {categories?.length ? (
        <div className='grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6'>
          {categories.map((e) => (
            <Link href={`/categories/${e.slug}`} key={e.id} className="group block">
              <div className='bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300'>
                <div className='relative h-[30rem] w-full'>
                  <Image
                    fill
                    alt={`${e.name} category`}
                    className='object-cover'
                    src={e.imageUrl ?? '/images/product-placeholder.webp'}
                  />
                  <div className='absolute inset-0 bg-gradient-to-t from-black/70 via-black/30 to-transparent rounded-b-lg opacity-50 group-hover:opacity-80 transition-opacity duration-300'></div>
                  <div className='absolute bottom-0 left-0 right-0 p-4 text-center z-10'>
                    <h2 className='text-lg font-bold text-white group-hover:text-blue-100 transition-colors drop-shadow-md'>
                      {e.name}
                    </h2>
                  </div>
                </div>
              </div>
            </Link>
          ))}
        </div>
      ) : null}
    </section>
  );
}
