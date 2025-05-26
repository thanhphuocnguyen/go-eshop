import Image from 'next/image';
import Link from 'next/link';
import CategoriesSection from './components/CategoriesSection';
import CollectionsSection from './components/CollectionsSection';
import { Metadata } from 'next';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { GeneralCategoryModel } from '@/app/lib/definitions';
import { ArrowRightIcon } from '@heroicons/react/24/outline';
import { serverSideFetch } from '@/app/lib/api/apiServer';

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
  const { data, error } = await serverSideFetch<{
    categories: GeneralCategoryModel[];
    collections: GeneralCategoryModel[];
  }>(PUBLIC_API_PATHS.ADVERTISE_CATEGORIES, {
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
      <section className='new-arrival-ads relative overflow-hidden min-h-[700px] flex items-center justify-center'>
        <div className='absolute inset-0 bg-gradient-to-b from-black/40 via-black/20 to-black/60'></div>
        <div className='relative w-full max-w-4xl mx-auto flex flex-col justify-center items-center z-10 px-4'>
          <div className='overflow-hidden mb-4'>
            <span className='block text-amber-300 text-xl font-medium uppercase tracking-wider mb-3 animate-fadeIn'>Just Launched</span>
            <h2 className='text-white text-6xl md:text-7xl font-bold text-center leading-tight'>
              New Arrivals <span className='text-blue-400'>Are Here</span>
            </h2>
          </div>
          <div className='h-1 w-24 bg-white my-8'></div>
          <p className='relative pb-10 pt-4 z-1 text-white text-xl md:text-2xl text-center max-w-2xl opacity-90'>
            The new arrivals have, well, newly arrived. Check out the latest
            options from our summer small-batch release while they&apos;re still
            in stock.
          </p>
          <Link
            href={'/collections/new-arrivals'}
            className='group relative mx-auto bg-white hover:bg-blue-600 text-xl text-black hover:text-white p-5 px-8 rounded-lg transition-all duration-300 shadow-lg hover:shadow-xl flex items-center space-x-2'
          >
            <span>Shop New Arrivals</span>
            <ArrowRightIcon className='w-5 h-5 group-hover:translate-x-1 transition-transform' />
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
        <section aria-label='banner' className='relative pt-24 h-[800px] overflow-hidden'>
          <div className='relative h-full'>
            <div className='relative w-full h-full'>
              <Image
                className='rounded-2xl object-cover relative scale-105 hover:scale-100 transition-all duration-700 ease-in-out'
                alt='shop workspace'
                src='/images/banners/home-page-01-feature-section-01.jpg'
                fill
                priority
                sizes="100vw"
                quality={90}
              />
              <div className='absolute inset-0 bg-gradient-to-r from-black/60 via-black/30 to-transparent rounded-2xl'></div>
            </div>
            <div className='absolute inset-0 flex flex-col justify-center'>
              <div className='w-1/2 ml-16 mr-auto flex flex-col content-start text-left max-w-2xl'>
                <div className='overflow-hidden'>
                  <h2 className='text-white text-6xl font-bold leading-tight transform transition-transform duration-500 ease-out'>
                    Level up <span className="text-blue-300">your desk</span>
                  </h2>
                </div>
                <div className='h-1 w-20 bg-blue-400 my-6'></div>
                <p className='pb-10 pt-2 text-white text-xl max-w-lg text-opacity-90'>
                  Make your desk beautiful and organized. Post a picture to
                  social media and watch it get more likes than life-changing
                  announcements. Reflect on the shallow nature of existence. At
                  least you have a really nice desk setup.
                </p>
                <Link
                  href={'/collections/workspace'}
                  className='group w-fit flex items-center space-x-2 bg-white hover:bg-blue-600 text-xl text-black hover:text-white p-5 rounded-lg transition-all duration-300 shadow-lg hover:shadow-xl'
                >
                  <span>Shop Workspace</span>
                  <ArrowRightIcon className='w-5 h-5 group-hover:translate-x-1 transition-transform' />
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
        <section aria-label='banner' className='relative my-32 mb-48 h-[700px] overflow-hidden'>
          <div className='h-full'>
            <div className='relative h-full rounded-2xl'>
              <Image
                className='rounded-2xl object-cover relative scale-105 hover:scale-100 transition-all duration-700 ease-in-out'
                src={'/images/banners/home-page-01-feature-section-02.jpg'}
                alt='productivity tools and focus system'
                fill
                sizes="100vw"
                quality={90}
              />
              <div className='absolute inset-0 bg-gradient-to-l from-black/70 via-black/40 to-transparent rounded-2xl'></div>
            </div>
            <div className='absolute h-full inset-0 flex items-center'>
              <div className='text-right flex flex-col content-center justify-center text-white mr-16 ml-auto gap-6 w-1/2 max-w-2xl'>
                <div className='overflow-hidden'>
                  <h2 className='text-5xl font-bold leading-tight transform transition-transform duration-500 ease-out'>
                    <span className="text-amber-200">Simple</span> productivity
                  </h2>
                </div>
                <div className='h-1 w-20 bg-amber-300 my-2 ml-auto'></div>
                <p className='text-xl text-white/90 max-w-lg ml-auto'>
                  Endless tasks, limited hours, a single piece of paper. Not
                  really a haiku, but we&apos;re doing our best here. No kanban
                  boards, burn-down charts, or tangled flowcharts with our Focus
                  system. Just the undeniable urge to fill empty circles.
                </p>
                <div className='mt-4 ml-auto'>
                  <Link 
                    className='group inline-flex items-center space-x-2 bg-amber-100 hover:bg-amber-500 text-amber-900 hover:text-white font-medium py-5 px-8 rounded-lg transition-all duration-300 shadow-lg hover:shadow-xl' 
                    href={'/collections/focus'}
                  >
                    <span>Shop Productivity</span>
                    <ArrowRightIcon className='w-5 h-5 group-hover:translate-x-1 transition-transform' />
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
