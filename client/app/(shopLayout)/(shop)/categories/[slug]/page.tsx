import { apiFetch } from '@/lib/apis/api';
import { PUBLIC_API_PATHS } from '@/lib/constants/api';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import { Metadata } from 'next';
import CategoryDetailClient from './_components/CategoryDetailClient';
import { cache } from 'react';

interface CategoryPageProps {
  params: Promise<{ slug: string }>;
}

export const getCategory = cache(async (slug: string) => {
  const { error, data } = await apiFetch<GenericResponse<GeneralCategoryModel>>(
    PUBLIC_API_PATHS.CATEGORY.replace(':slug', slug),
    {
      nextOptions: {
        next: { revalidate: 3600, tags: [`category-${slug}`] }, // Cache for 1 hour
      },
    }
  );
  if (error) {
    throw new Error(error.details, {
      cause: error,
    });
  }
  return data;
});

// Generate metadata for SEO
export async function generateMetadata({
  params,
}: CategoryPageProps): Promise<Metadata> {
  const { slug } = await params;
  // Fetch the category by slug directly
  try {
    const category = await getCategory(slug);

    if (!category) {
      return {
        title: 'Category Not Found',
        description: 'The requested category could not be found',
      };
    }

    return {
      title: `${category.name} - Shop by Category`,
      description:
        category.description ||
        `Browse our collection of ${category.name} products`,
      openGraph: {
        title: `${category.name} - Shop by Category`,
        description:
          category.description ||
          `Browse our collection of ${category.name} products`,
        images: category.image_url ? [{ url: category.image_url }] : [],
      },
    };
  } catch (error) {
    console.error('Error generating metadata:', error);
    return {
      title: 'Shop by Category',
      description: 'Browse our products by category',
    };
  }
}

export default async function CategoryDetailPage({
  params,
}: CategoryPageProps) {
  const { slug } = await params;
  const category = await getCategory(slug);

  // Pass the fetched data to the client component
  return <CategoryDetailClient category={category} />;
}
