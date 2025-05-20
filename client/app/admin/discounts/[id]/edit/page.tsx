'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import dayjs from 'dayjs';
import { TabGroup, TabPanels, Button } from '@headlessui/react';

import {
  DetailsPanel,
  ProductsPanel,
  CategoriesPanel,
  TabNavigation
} from '../../_components';
import { discountSchema, DiscountFormData, ProductType, CategoryType } from '../../_types';

// Mock data for a specific discount
const mockDiscount = {
  id: '1',
  code: 'SUMMER25',
  description: 'Summer sale discount',
  discountType: 'percentage',
  discountValue: 25,
  minPurchaseAmount: 50,
  maxDiscountAmount: 100,
  usageLimit: 1000,
  usedCount: 342,
  isActive: true,
  startsAt: '2025-05-01T00:00:00Z',
  expiresAt: '2025-08-31T23:59:59Z',
  createdAt: '2025-04-15T08:30:00Z',
  updatedAt: '2025-04-20T14:45:30Z',
  products: [
    { id: 'prod1', name: 'Summer T-Shirt', price: 29.99 },
    { id: 'prod2', name: 'Beach Shorts', price: 39.99 },
    { id: 'prod3', name: 'Sunglasses', price: 25.0 },
  ],
  categories: [
    { id: 'cat1', name: 'Summer Collection' },
    { id: 'cat2', name: 'Beachwear' },
  ],
  usageHistory: [
    {
      id: 'order1',
      orderId: 'ORD-12345',
      customerName: 'John Doe',
      amount: 89.97,
      discountAmount: 22.49,
      date: '2025-05-10T15:30:00Z',
    },
    {
      id: 'order2',
      orderId: 'ORD-12346',
      customerName: 'Jane Smith',
      amount: 129.99,
      discountAmount: 32.5,
      date: '2025-05-12T09:15:00Z',
    },
    {
      id: 'order3',
      orderId: 'ORD-12350',
      customerName: 'Robert Johnson',
      amount: 75.5,
      discountAmount: 18.88,
      date: '2025-05-14T17:22:00Z',
    },
  ],
};

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

export default function EditDiscountPage({
  params,
}: {
  params: { id: string };
}) {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
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
      startsAt: '',
      expiresAt: new Date(),
      products: [],
      categories: [],
    },
  });

  useEffect(() => {
    // In a real implementation, this would fetch the discount data and available products/categories from an API
    const fetchData = async () => {
      try {
        // Simulate API call with a delay
        await new Promise((resolve) => setTimeout(resolve, 500));

        // Set form values from the fetched discount
        setValue('code', mockDiscount.code);
        setValue('description', mockDiscount.description);
        setValue('discountType', mockDiscount.discountType as 'percentage' | 'fixed_amount');
        setValue('discountValue', mockDiscount.discountValue);
        setValue('minPurchaseAmount', mockDiscount.minPurchaseAmount);
        setValue('maxDiscountAmount', mockDiscount.maxDiscountAmount);
        setValue('usageLimit', mockDiscount.usageLimit);
        setValue('isActive', mockDiscount.isActive);
        setValue(
          'startsAt',
          dayjs(mockDiscount.startsAt).format('YYYY-MM-DDTHH:mm')
        );
        setValue(
          'expiresAt',
          dayjs(mockDiscount.expiresAt).toDate()
        );
        setValue('products', mockDiscount.products);
        setValue('categories', mockDiscount.categories);

        // Set available products/categories (excluding already selected ones)
        setAvailableProducts(
          mockAllProducts.filter(
            (p) => !mockDiscount.products.some((sp) => sp.id === p.id)
          )
        );

        setAvailableCategories(
          mockAllCategories.filter(
            (c) => !mockDiscount.categories.some((sc) => sc.id === c.id)
          )
        );
      } catch (error) {
        console.error('Error fetching discount data:', error);
        setError('Failed to load discount data. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [setValue, params.id]);

  // Submit handler
  const onSubmit = async (data: DiscountFormData) => {
    setSubmitting(true);
    setError(null);
    
    try {
      // Mock API call with a delay
      await new Promise((resolve) => setTimeout(resolve, 800));
      console.log('Submitted data:', data);

      // Redirect to the discount details page after successful update
      router.push(`/admin/discounts/${params.id}`);
    } catch (error) {
      console.error('Error updating discount:', error);
      setError('Failed to update discount. Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className='p-4'>
        <div className='flex justify-center items-center h-64'>
          <div className='animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary'></div>
        </div>
      </div>
    );
  }

  return (
    <div className='p-4'>
      <div className='mb-4 flex items-center justify-between'>
        <div className="flex items-center">
          <Button
            onClick={() => router.back()}
            className='mr-4 p-2 rounded-full hover:bg-gray-100'
          >
            <ArrowLeftIcon className='h-5 w-5' />
          </Button>
          <h1 className='text-2xl font-semibold'>Edit Discount</h1>
        </div>
      </div>

      {error && (
        <div className='mb-4 p-3 bg-red-100 text-red-800 rounded-md'>
          {error}
        </div>
      )}

      <div className='bg-white p-6 rounded-lg shadow-md'>
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
            <Link href={`/admin/discounts/${params.id}`}>
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
              {submitting ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
