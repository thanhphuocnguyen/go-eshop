'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import dayjs from 'dayjs';
import {
  Button,
  TabGroup,
  TabPanels,
} from '@headlessui/react';

import { 
  DetailsPanel, 
  ProductsPanel, 
  CategoriesPanel, 
  TabNavigation 
} from '../_components';
import { discountSchema, DiscountFormData, ProductType, CategoryType } from '../_types';

// Mock data for products and categories selection
const mockAllProducts = [
  { id: 'prod1', name: 'Summer T-Shirt', price: 29.99 },
  { id: 'prod2', name: 'Beach Shorts', price: 39.99 },
  { id: 'prod3', name: 'Sunglasses', price: 25.0 },
  { id: 'prod4', name: 'Beach Hat', price: 19.99 },
  { id: 'prod5', name: 'Flip Flops', price: 15.99 },
  { id: 'prod6', name: 'Beach Towel', price: 22.99 },
];

const mockAllCategories = [
  { id: 'cat1', name: 'Summer Collection' },
  { id: 'cat2', name: 'Beachwear' },
  { id: 'cat3', name: 'Accessories' },
  { id: 'cat4', name: 'Footwear' },
  { id: 'cat5', name: 'Seasonal' },
];

export default function NewDiscountPage() {
  const router = useRouter();
  const [submitting, setSubmitting] = useState(false);
  const [activeTab, setActiveTab] = useState('details');
  const [error, setError] = useState<string | null>(null);

  const [availableProducts, setAvailableProducts] = useState<ProductType[]>([]);
  const [availableCategories, setAvailableCategories] = useState<CategoryType[]>([]);

  // Initialize react-hook-form with Zod resolver
  const {
    register,
    handleSubmit,
    control,
    setValue,
    watch,
    formState: { errors },
  } = useForm<DiscountFormData>({
    resolver: zodResolver(discountSchema),
    defaultValues: {
      code: '',
      description: '',
      discountType: 'percentage',
      discountValue: 0,
      minPurchaseAmount: null,
      maxDiscountAmount: null,
      usageLimit: null,
      isActive: true,
      startsAt: dayjs().format('YYYY-MM-DDTHH:mm'),
      expiresAt: dayjs().add(30, 'day').toDate(),
      products: [],
      categories: [],
    },
  });

  useEffect(() => {
    // In a real implementation, this would fetch available products/categories from an API
    const fetchData = async () => {
      try {
        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 500));

        setAvailableProducts(mockAllProducts);
        setAvailableCategories(mockAllCategories);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, []);

  // Submit handler
  const onSubmit = async (data: DiscountFormData) => {
    setSubmitting(true);
    setError(null);

    try {
      // In a real implementation, this would send the data to an API
      console.log('Submitting discount data:', data);

      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Redirect to the discounts list page
      router.push('/admin/discounts');
    } catch (err) {
      setError('Failed to create discount. Please try again.');
      console.error(err);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className='p-4 h-full'>
      <div className='mb-4 flex items-center justify-between'>
        <div className="flex items-center">
          <Button
            onClick={() => router.back()}
            className='mr-4 p-2 rounded-full hover:bg-gray-100'
          >
            <ArrowLeftIcon className='h-5 w-5' />
          </Button>
          <h1 className='text-2xl font-semibold'>Create New Discount</h1>
        </div>
      </div>

      {error && (
        <div className='mb-4 p-3 bg-red-100 text-red-800 rounded-md'>
          {error}
        </div>
      )}

      <div className=''>
        <form onSubmit={handleSubmit(onSubmit)}>
          <TabGroup>
            <TabNavigation activeTab={activeTab} setActiveTab={setActiveTab} />
            <TabPanels>
              <DetailsPanel 
                register={register} 
                control={control} 
                errors={errors} 
                watch={watch}
              />
              
              <ProductsPanel 
                watch={watch} 
                setValue={setValue} 
                availableProducts={availableProducts} 
                setAvailableProducts={setAvailableProducts}
              />
              
              <CategoriesPanel 
                watch={watch} 
                setValue={setValue} 
                availableCategories={availableCategories} 
                setAvailableCategories={setAvailableCategories}
              />
            </TabPanels>
          </TabGroup>

          <div className='mt-8 flex justify-end border-t pt-6'>
            <Link href='/admin/discounts'>
              <Button
                type='button'
                className='mr-3 px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50'
              >
                Cancel
              </Button>
            </Link>
            <Button
              type='submit'
              disabled={submitting}
              className={`px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 ${
                submitting ? 'opacity-75 cursor-not-allowed' : ''
              }`}
            >
              {submitting ? 'Creating...' : 'Create Discount'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
