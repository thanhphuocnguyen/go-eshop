import Image from 'next/image';
import Link from 'next/link';
import CategoriesSection from './components/CategoriesSection';
import CollectionsSection from './components/CollectionsSection';
import { Metadata } from 'next';
import { apiFetch } from '@/lib/apis/api';
import { PUBLIC_API_PATHS } from '@/lib/constants/api';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import { ArrowRightIcon } from '@heroicons/react/24/outline';

export const metadata: Metadata = {
  title: 'Homepage',
  description:
    'Welcome to our homepage! Explore our latest products and collections.',
};

// No data placeholder component
const NoDataPlaceholder = ({ type }: { type: string }) => {
  return (
    <div className='flex flex-col items-center justify-center py-10 bg-gray-50 rounded-lg border border-gray-200'>
      <Image
        src='/images/not-found.webp'
        alt='No data available'
        width={120}
        height={120}
        className='mb-4 opacity-60'
      />
      <h3 className='text-xl font-medium text-gray-700'>No {type} Available</h3>
      <p className='text-gray-500 text-center mt-2'>
        {type === 'Categories'
          ? 'Categories will be displayed here when they become available.'
          : 'Collections will be displayed here when they become available.'}
      </p>
    </div>
  );
};

export default async function Home() {
  const { data, error } = await apiFetch<
    GenericResponse<{
      categories: GeneralCategoryModel[];
      collections: GeneralCategoryModel[];
    }>
  >(PUBLIC_API_PATHS.HOME_PAGE_DATA, {
    nextOptions: {
      next: {
        tags: ['home'],
      },
    },
  });

  if (error) {
    return (
      <div className='flex justify-center items-center p-8'>
        <div className='animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600'></div>
      </div>
    );
  }

  return (
    <div className='block relative mb-10'>
      <section className='new-arrival-ads'>
        <div className='overlay'></div>
        <div className='relative w-[600] mx-auto flex flex-col justify-center items-center'>
          <h2 className='text-white relative text-5xl font-bold text-center'>
            New Arrivals Are Here
          </h2>
          <p className='relative pb-10 pt-4 z-1 text-white text-xl text-center'>
            The new arrivals have, well, newly arrived. Check out the latest
            options from our summer small-batch release while they&apos;re still
            in stock.
          </p>
          <Link
            href={'/categories/new-arrivals'}
            className='relative mx-auto bg-white hover:bg-white/60 text-xl text-black p-4 rounded-lg'
          >
            Shop New Arrivals
          </Link>
        </div>
      </section>
      <div className='px-10'>
        {data.categories && data.categories.length > 0 ? (
          <CategoriesSection categories={data.categories} />
        ) : (
          <div className='pt-24'>
            <div className='flex justify-between mb-2'>
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
            <NoDataPlaceholder type='Categories' />
          </div>
        )}
        <section className='relative pt-24 h-[700px]'>
          <div className='relative h-full'>
            <div className='relative w-full h-full'>
              <Image
                className='rounded-lg object-cover relative'
                alt='shop workspace'
                src='/images/banners/home-page-01-feature-section-01.jpg'
                fill
              />
              <div className='overlay rounded-md'></div>
            </div>
            <div className='absolute inset-0 flex flex-col m-auto'>
              <div className='w-1/2 m-auto flex flex-col content-center text-center'>
                <h2 className='text-white text-5xl font-bold text-center'>
                  Level up your desk
                </h2>
                <p className='pb-10 pt-4 z-1 text-white text-xl text-center'>
                  Make your desk beautiful and organized. Post a picture to
                  social media and watch it get more likes than life-changing
                  announcements. Reflect on the shallow nature of existence. At
                  least you have a really nice desk setup.
                </p>
                <Link
                  href={'/categories/workspaces'}
                  className='mx-auto bg-white hover:bg-white/60 text-xl text-black p-4 rounded-lg'
                >
                  Shop Workspace
                </Link>
              </div>
            </div>
          </div>
        </section>
        {data.collections && data.collections.length > 0 ? (
          <CollectionsSection collections={data.collections} />
        ) : (
          <div className='pt-20 pb-10'>
            <NoDataPlaceholder type='Collections' />
          </div>
        )}
        <section className='relative my-24 mb-40 h-[600px]'>
          <div className='h-full'>
            <div className='relative h-full rounded-md'>
              <Image
                className='rounded-lg object-cover relative'
                src={'/images/banners/home-page-01-feature-section-02.jpg'}
                alt='product-image'
                fill
              />
              <div className='overlay rounded-md'></div>
            </div>
            <div className='absolute h-full inset-0 m-auto'>
              <div className='text-center flex flex-col content-center justify-center text-white m-auto gap-4 w-1/2 h-full'>
                <h2 className='text-3xl font-bold'>Simple productivity</h2>
                <p>
                  Endless tasks, limited hours, a single piece of paper. Not
                  really a haiku, but we&apos;re doing our best here. No kanban
                  boards, burn-down charts, or tangled flowcharts with our Focus
                  system. Just the undeniable urge to fill empty circles.
                </p>
                <Link className='p-3' href={'/collections/focus'}>
                  Shop Productivity
                </Link>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
