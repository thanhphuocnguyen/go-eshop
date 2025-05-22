'use client';

import React, { useState, useEffect, use } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { FormProvider, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import dayjs from 'dayjs';
import { TabGroup, TabPanels, Button } from '@headlessui/react';

import {
  DetailsPanel,
  ProductsPanel,
  CategoriesPanel,
  TabNavigation,
} from '../../_components';
import { createDiscountSchema, CreateDiscountFormData } from '../../_types';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import useSWR from 'swr';
import { clientSideFetch } from '@/app/lib/apis/apiClient';
import { Discount } from '../page';
import { toast } from 'react-toastify';

// Mock data for a specific discount

// Mock data for products and categories selection

export default function EditDiscountPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const router = useRouter();
  const [activeTab, setActiveTab] = useState('details');

  const { data: discount, isLoading } = useSWR(
    ADMIN_API_PATHS.DISCOUNT.replace(':id', id),
    (url) =>
      clientSideFetch<Discount>(url).then((res) => {
        if (res.error) {
          throw new Error(res.error.details);
        }
        return res.data;
      }),
    {
      onError: (err) => {
        toast.error(
          <div>
            Failed to fetch discount details.
            <div>{JSON.stringify(err)}</div>
          </div>
        );
      },
    }
  );

  // Initialize react-hook-form with Zod resolver
  const editForm = useForm<CreateDiscountFormData>({
    resolver: zodResolver(createDiscountSchema),
    defaultValues: {
      code: '',
      description: '',
      discountType: { id: 'percentage', name: 'Percentage' },
      discountValue: 0,
      minPurchaseAmount: null,
      maxDiscountAmount: null,
      usageLimit: null,
      isActive: true,
      startsAt: '',
      expiresAt: dayjs().add(1, 'month').format('YYYY-MM-DDTHH:mm'),
    },
  });
  const {
    handleSubmit,
    reset,
    formState: { isSubmitting, isDirty },
  } = editForm;

  useEffect(() => {
    if (discount) {
      reset({
        code: discount.code,
        description: discount.description,
        discountType: {
          id: discount.discountType,
          name:
            discount.discountType === 'percentage'
              ? 'Percentage'
              : 'Fixed Amount',
        },
        discountValue: discount.discountValue,
        minPurchaseAmount: discount.minPurchase,
        maxDiscountAmount: discount.maxDiscount,
        usageLimit: discount.usageLimit,
        isActive: discount.isActive,
        startsAt: dayjs(discount.startsAt).format('YYYY-MM-DDTHH:mm'),
        expiresAt: dayjs(discount.expiresAt).format('YYYY-MM-DDTHH:mm'),
      });
    }
    // In a real implementation, this would fetch the discount data and available products/categories from an API
  }, [id, discount]);

  // Submit handler
  const onSubmit = async (data: CreateDiscountFormData) => {
    try {
      // Mock API call with a delay
      await new Promise((resolve) => setTimeout(resolve, 800));
      console.log('Submitted data:', data);

      // Redirect to the discount details page after successful update
      router.push(`/admin/discounts/${id}`);
    } catch (error) {
      console.error('Error updating discount:', error);
    } finally {
    }
  };

  if (isLoading) {
    return (
      <div className='p-4'>
        <div className='flex justify-center items-center h-64'>
          <div className='animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary'></div>
        </div>
      </div>
    );
  }

  return (
    <div className='p-4 h-full'>
      <div className='mb-4 flex items-center justify-between'>
        <div className='flex items-center'>
          <Button
            onClick={() => router.back()}
            className='mr-4 p-2 rounded-full hover:bg-gray-100'
          >
            <ArrowLeftIcon className='h-5 w-5' />
          </Button>
          <h1 className='text-2xl font-semibold'>Edit Discount</h1>
        </div>
      </div>

      <div className=''>
        <FormProvider {...editForm}>
          <form onSubmit={handleSubmit(onSubmit)}>
            <TabGroup>
              <TabNavigation
                activeTab={activeTab}
                setActiveTab={setActiveTab}
              />
              <TabPanels>
                <DetailsPanel />

                <ProductsPanel />

                <CategoriesPanel />
              </TabPanels>
            </TabGroup>

            <div className='mt-8 flex justify-end border-t pt-6'>
              <Link href={`/admin/discounts/${id}`}>
                <Button
                  type='button'
                  className='mr-3 px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50'
                >
                  Cancel
                </Button>
              </Link>
              <Button
                type='submit'
                disabled={isSubmitting || !isDirty}
                className={`px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 ${
                  isSubmitting ? 'opacity-75 cursor-not-allowed' : ''
                }`}
              >
                {isSubmitting ? 'Saving...' : 'Save Changes'}
              </Button>
            </div>
          </form>
        </FormProvider>
      </div>
    </div>
  );
}
