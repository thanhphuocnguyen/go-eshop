import { GeneralCategoryModel } from '@/lib/definitions';
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
        <div className='h-[500px] grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4'>
          {categories.map((e) => (
            <Link href={`/categories/${e.slug}`} key={e.id} className="block w-full h-full">
              <div className='relative bg-white rounded-md shadow-md h-full'>
                <h2 className='text-xl font-bold absolute bottom-5 text-white text-center left-0 right-0 mx-auto z-10 '>
                  {e.name}
                </h2>
                {
                  <Image
                    fill
                    alt='product-image'
                    className='object-cover rounded-md'
                    src={e.image_url ?? '/images/product-placeholder.webp'}
                  />
                }
                <div className='absolute z-0 h-1/2 opacity-45 inset-x-0 bottom-0 bg-gradient-to-t from-slate-600 via-white to-transparent rounded-md'></div>
              </div>
            </Link>
          ))}
        </div>
      ) : null}
    </section>
  );
}
