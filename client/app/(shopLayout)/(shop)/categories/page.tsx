'use client';

import { apiFetch } from '@/lib/apis/api';
import { API_PATHS } from '@/lib/constants/api';
import { CategoryProductModel, GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import Image from 'next/image';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import CategoryProductSkeleton from '@/components/Product/CategoryProductSkeleton';
import ProductCard from '@/components/Product/ProductCard';

export default function CategoryPage() {
  const [categories, setCategories] = useState<GeneralCategoryModel[]>([]);
  const [categoryProducts, setCategoryProducts] = useState<{ [key: string]: CategoryProductModel[] }>({});
  const [loading, setLoading] = useState(true);
  const [loadingProducts, setLoadingProducts] = useState<{ [key: string]: boolean }>({});

  // Fetch all categories
  useEffect(() => {
    async function fetchCategories() {
      try {
        const { data } = await apiFetch<GenericResponse<GeneralCategoryModel[]>>(
          `${API_PATHS.CATEGORIES}?page=1&page_size=10`
        );
        
        if (data) {
          setCategories(data);
          
          // Initialize loading state for each category
          const initialLoadingState: { [key: string]: boolean } = {};
          data.forEach(category => {
            initialLoadingState[category.id] = true;
          });
          setLoadingProducts(initialLoadingState);
          
          // Fetch products for each category
          data.forEach(category => {
            fetchProductsByCategory(category.id);
          });
        }
      } catch (error) {
        console.error('Failed to fetch categories:', error);
      } finally {
        setLoading(false);
      }
    }
    
    fetchCategories();
  }, []);

  // Fetch products for a specific category
  async function fetchProductsByCategory(categoryId: string) {
    try {
      const { data } = await apiFetch<GenericResponse<CategoryProductModel[]>>(
        API_PATHS.CATEGORY_PRODUCTS.replace(':id', categoryId)
      );
      
      if (data) {
        setCategoryProducts(prev => ({
          ...prev,
          [categoryId]: data
        }));
      }
    } catch (error) {
      console.error(`Failed to fetch products for category ${categoryId}:`, error);
    } finally {
      setLoadingProducts(prev => ({
        ...prev,
        [categoryId]: false
      }));
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className='text-3xl font-bold mb-8'>Categories</h1>
        <div className="space-y-12">
          {[1, 2, 3].map(i => (
            <div key={i} className="animate-pulse">
              <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
              <CategoryProductSkeleton />
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className='container mx-auto px-4 py-8'>
      <div className='mb-8'>
        <h1 className='text-3xl font-bold'>Categories</h1>
        <p className='text-gray-500 mt-2'>Explore our product categories</p>
      </div>
      
      <div className='space-y-12'>
        {categories.map((category) => (
          <div key={category.id} className='pb-8 border-b border-gray-200 last:border-0'>
            <div className='flex items-center mb-6'>
              {category.image_url && (
                <div className='h-16 w-16 mr-4 relative overflow-hidden rounded-lg'>
                  <Image
                    src={category.image_url}
                    alt={category.name}
                    fill
                    className='object-cover'
                  />
                </div>
              )}
              <div>
                <h2 className='text-2xl font-bold text-gray-800'>{category.name}</h2>
                <p className='text-gray-600'>{category.description}</p>
                <Link 
                  href={`/categories/${category.slug}`}
                  className='text-indigo-600 text-sm hover:underline mt-1 inline-block'
                >
                  View all products in {category.name} →
                </Link>
              </div>
            </div>

            {/* Products section */}
            {loadingProducts[category.id] ? (
              <CategoryProductSkeleton />
            ) : (
              <>
                {categoryProducts[category.id] && categoryProducts[category.id].length > 0 ? (
                  <div className='grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6'>
                    {categoryProducts[category.id].slice(0, 4).map((product) => (
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
                  <p className='text-gray-500 italic'>No products available in this category.</p>
                )}
              </>
            )}
          </div>
        ))}

        {categories.length === 0 && (
          <div className='text-center py-12'>
            <p className='text-xl text-gray-600'>No categories found.</p>
          </div>
        )}
      </div>
    </div>
  );
}
