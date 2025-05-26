import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { serverSideFetch } from '@/app/lib/api/apiServer';
import { Metadata } from 'next';
import Image from 'next/image';
import { notFound } from 'next/navigation';
import { cache } from 'react';
import { GetCollectionBySlugModel } from '../lib/definitions';

// Define the params type
type Params = {
  slug: string;
};

// Cached collection fetch function
async function getCollectionBySlug(
  slug: string
): Promise<GetCollectionBySlugModel> {
  const result = await serverSideFetch<GetCollectionBySlugModel>(
    PUBLIC_API_PATHS.COLLECTION.replace(':slug', slug)
  );

  if (result.error) {
    throw new Error(result.error.details, {
      cause: result.error,
    });
  }
  if (!result.data || result.error) {
    notFound();
  }

  return result.data;
}
export const collectionBySlugCache = cache(getCollectionBySlug);

// Dynamic metadata generation
export async function generateMetadata({
  params,
}: {
  params: Params;
}): Promise<Metadata> {
  try {
    const { slug } = await params;
    const { collection } = await collectionBySlugCache(slug);

    return {
      title: `${collection.name} Collection | eShop`,
      description:
        collection.description ||
        `Browse the ${collection.name} collection at eShop.`,
      openGraph: {
        title: `${collection.name} Collection | eShop`,
        description:
          collection.description ||
          `Browse the ${collection.name} collection at eShop.`,
        images: collection.imageUrl
          ? [
              {
                url: collection.imageUrl,
                width: 1200,
                height: 630,
                alt: `${collection.name} Collection`,
              },
            ]
          : [],
      },
    };
  } catch (error) {
    return {
      title: 'Collection | eShop',
      description: 'Explore our curated collections',
    };
  }
}

export default async function CollectionLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<Params>;
}) {
  const { slug } = await params;
  const { collection } = await collectionBySlugCache(slug);

  return (
    <div className='container mx-auto px-4 py-8'>
      {/* Collection Header Section - Avatar style */}
      <div className='mb-8 bg-white shadow-sm rounded-xl overflow-hidden'>
        <div className='flex flex-col md:flex-row md:h-48 lg:h-56'>
          {/* Avatar Image */}
          <div className='w-full md:w-48 lg:w-56 relative'>
            {collection.imageUrl ? (
              <div className='relative h-48 md:h-full w-full'>
                <Image
                  src={collection.imageUrl}
                  alt={collection.name}
                  fill
                  priority
                  className='object-cover'
                  sizes='(max-width: 768px) 100vw, 200px'
                />
              </div>
            ) : (
              <div className='h-48 md:h-full w-full bg-gradient-to-r from-blue-600 to-indigo-800 flex items-center justify-center'>
                <span className='text-4xl font-bold text-white'>
                  {collection.name.charAt(0)}
                </span>
              </div>
            )}
          </div>

          {/* Content Box */}
          <div className='flex-1 p-4 md:p-6 flex flex-col justify-center'>
            <div className='max-w-3xl'>
              <div className='mb-1 text-xs text-indigo-600 font-medium'>
                Collection
              </div>
              <h1 className='text-2xl md:text-3xl font-bold text-gray-900 mb-2 tracking-tight'>
                {collection.name}
              </h1>
              <p className='text-gray-600 text-sm md:text-base mb-2'>
                {collection.description
                  ? collection.description.length > 120
                    ? collection.description.substring(0, 120) + '...'
                    : collection.description
                  : `Explore our ${collection.name} collection.`}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Children will render the page content */}
      {children}
    </div>
  );
}
