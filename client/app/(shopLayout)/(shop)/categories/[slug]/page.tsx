'use client';

import { apiFetch } from '@/lib/apis/api';
import { API_PATHS } from '@/lib/constants/api';
import { CategoryProductModel, GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import CategoryProductSkeleton from '@/components/Product/CategoryProductSkeleton';
import ProductCard from '@/components/Product/ProductCard';

export default function CategoryDetailPage({ 
  params 
}: { 
  params: { slug: string } 
}) {
  const { slug } = params;
  const [category, setCategory] = useState<GeneralCategoryModel | null>(null);
  const [products, setProducts] = useState<CategoryProductModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingProducts, setLoadingProducts] = useState(true);
  
  useEffect(() => {
    async function fetchCategoryAndProducts() {
      setLoading(true);
      setLoadingProducts(true);
      
      try {
        // First fetch the category to get its ID
        const categoriesResponse = await apiFetch<GenericResponse<GeneralCategoryModel[]>>(
          `${API_PATHS.CATEGORIES}?page=1&page_size=100`
        );
        
        if (categoriesResponse.data) {
          const foundCategory = categoriesResponse.data.find(cat => cat.slug === slug);
          
          if (foundCategory) {
            setCategory(foundCategory);
            
            // Then fetch products for this category
            const productsResponse = await apiFetch<GenericResponse<CategoryProductModel[]>>(
              API_PATHS.CATEGORY_PRODUCTS.replace(':id', foundCategory.id)
            );
            
            if (productsResponse.data) {
              setProducts(productsResponse.data);
            }
          }
        }
      } catch (error) {
        console.error('Failed to fetch category or products:', error);
      } finally {
        setLoading(false);
        setLoadingProducts(false);
      }
    }
    
    fetchCategoryAndProducts();
  }, [slug]);

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="animate-pulse">
          <div className="h-10 bg-gray-200 rounded w-1/3 mb-6"></div>
          <div className="h-6 bg-gray-200 rounded w-1/2 mb-12"></div>
          <CategoryProductSkeleton />
        </div>
      </div>
    );
  }

  if (!category) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Link href="/categories" className="inline-flex items-center text-indigo-600 hover:underline mb-6">
          <ArrowLeftIcon className="h-4 w-4 mr-2" />
          Back to all categories
        </Link>
        <div className="text-center py-12">
          <h1 className="text-2xl font-bold text-gray-800 mb-2">Category Not Found</h1>
          <p className="text-gray-600">The category you are looking for does not exist or has been removed.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <Link href="/categories" className="inline-flex items-center text-indigo-600 hover:underline mb-6">
        <ArrowLeftIcon className="h-4 w-4 mr-2" />
        Back to all categories
      </Link>
      
      <div className="flex items-center mb-8">
        {category.image_url && (
          <div className="h-24 w-24 mr-6 relative overflow-hidden rounded-lg shadow-md">
            <Image
              src={category.image_url}
              alt={category.name}
              fill
              className="object-cover"
            />
          </div>
        )}
        <div>
          <h1 className="text-3xl font-bold text-gray-800">{category.name}</h1>
          {category.description && (
            <p className="text-gray-600 mt-2 max-w-2xl">{category.description}</p>
          )}
        </div>
      </div>
      
      <div>
        <h2 className="text-xl font-semibold text-gray-800 mb-6">Products in this category</h2>
        
        {loadingProducts ? (
          <CategoryProductSkeleton />
        ) : (
          <>
            {products.length > 0 ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                {products.map((product) => (
                  <ProductCard
                    key={product.id}
                    ID={product.id}
                    name={product.name}
                    image={product.image_url}
                    priceFrom={product.price_from}
                    priceTo={product.price_to}
                    rating={4.5}
                  />
                ))}
              </div>
            ) : (
              <div className="text-center py-12 bg-gray-50 rounded-lg">
                <p className="text-lg text-gray-600">No products available in this category.</p>
                <Link href="/categories" className="text-indigo-600 hover:underline mt-2 inline-block">
                  Browse other categories
                </Link>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
