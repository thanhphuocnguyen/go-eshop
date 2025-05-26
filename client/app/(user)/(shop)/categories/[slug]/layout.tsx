import { cache, Suspense } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import CategoryProductSkeleton from '@/app/components/Product/CategoryProductSkeleton';
import { GeneralCategoryModel } from '@/app/lib/definitions';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { serverSideFetch } from '@/app/lib/api/apiServer';
import { Metadata } from 'next';
import { notFound } from 'next/navigation';

// Create a cache for the category data

export async function generateMetadata({
  params,
}: {
  params: { slug: string };
}): Promise<Metadata> {
  const category = await getCategoryBySlug(params.slug);

  if (!category) {
    return {
      title: 'Category Not Found',
      description: 'The requested category could not be found.',
    };
  }

  return {
    title: `${category.name} | Shop by Category`,
    description:
      category.description ||
      `Browse our selection of ${category.name} products.`,
    openGraph: {
      title: `${category.name} | Shop by Category`,
      description:
        category.description ||
        `Browse our selection of ${category.name} products.`,
      images: category.imageUrl ? [{ url: category.imageUrl }] : [],
    },
    twitter: {
      card: 'summary_large_image',
      title: category.name,
      description:
        category.description ||
        `Browse our selection of ${category.name} products.`,
      images: category.imageUrl ? [category.imageUrl] : [],
    },
  };
}

export async function getCategoryBySlug(
  slug: string
): Promise<GeneralCategoryModel | null> {
  const { data, error } = await serverSideFetch<GeneralCategoryModel>(
    `${PUBLIC_API_PATHS.CATEGORIES}/${slug}`
  );

  if (error || !data) {
    console.error('Failed to fetch category:', error);
    notFound();
    return null;
  }

  return data;
}
export const categoryCache = cache(getCategoryBySlug);

export default async function CategoryLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const category = await categoryCache(slug);

  if (!category) {
    return (
      <div className='container mx-auto px-4 py-8'>
        <Link
          href='/categories'
          className='inline-flex items-center text-indigo-600 hover:underline mb-6'
        >
          <ArrowLeftIcon className='h-4 w-4 mr-2' />
          Back to all categories
        </Link>
        <div className='text-center py-12'>
          <h1 className='text-2xl font-bold text-gray-800 mb-2'>
            Category Not Found
          </h1>
          <p className='text-gray-600'>
            The category you are looking for does not exist or has been removed.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className='container h-max mx-auto px-4 py-8'>
      <Link
        href='/categories'
        className='inline-flex items-center text-indigo-600 hover:underline mb-6'
      >
        <ArrowLeftIcon className='h-4 w-4 mr-2' />
        Back to all categories
      </Link>

      <div className='flex items-center mb-8'>
        {category.imageUrl && (
          <div className='h-24 w-24 mr-6 relative overflow-hidden rounded-lg shadow-md'>
            <Image
              src={category.imageUrl}
              alt={category.name}
              fill
              className='object-cover'
            />
          </div>
        )}
        <div>
          <h1 className='text-3xl font-bold text-gray-800'>{category.name}</h1>
          {category.description && (
            <p className='text-gray-600 mt-2 max-w-2xl'>
              {category.description}
            </p>
          )}
        </div>
      </div>

      <Suspense fallback={<CategoryProductSkeleton />}>{children}</Suspense>
    </div>
  );
}
