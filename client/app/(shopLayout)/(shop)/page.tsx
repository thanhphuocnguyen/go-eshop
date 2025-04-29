import Image from 'next/image';
import Link from 'next/link';
import { Suspense } from 'react';
import Loading from '../../loading';
import CategoriesSection from './components/CategoriesSection';
import CollectionsSection from './components/CollectionsSection';
import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Homepage',
  description:
    'Welcome to our homepage! Explore our latest products and collections.',
};

export default async function Home() {
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
        <Suspense fallback={<Loading />}>
          <CategoriesSection />
        </Suspense>
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
        <Suspense fallback={<Loading />}>
          <CollectionsSection />
        </Suspense>

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
                  boards, burndown charts, or tangled flowcharts with our Focus
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
