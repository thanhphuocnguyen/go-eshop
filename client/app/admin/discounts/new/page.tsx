'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { FormProvider, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import dayjs from 'dayjs';
import { Button, TabGroup, TabPanels } from '@headlessui/react';

import { DetailsPanel, TabNavigation } from '../_components';
import {
  createDiscountSchema,
  CreateDiscountFormData,
  CreateDiscountOutputData,
} from '../_types';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';

export default function NewDiscountPage() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState('details');

  // Initialize react-hook-form with Zod resolver
  const discountForm = useForm<
    CreateDiscountFormData,
    unknown,
    CreateDiscountOutputData
  >({
    reValidateMode: 'onBlur',
    mode: 'onBlur',
    resolver: zodResolver(createDiscountSchema),
    defaultValues: {
      code: 'FREESHIP5',
      discountType: {
        id: 'percentage',
        name: 'Percentage',
      },
      discountValue: 5,
      isActive: true,
      startsAt: dayjs().format('YYYY-MM-DDTHH:mm'),
      expiresAt: dayjs().add(1, 'month').format('YYYY-MM-DDTHH:mm'),
      description: 'Free ship for orders over $50',
    },
  });

  const {
    handleSubmit,
    formState: { isSubmitting, isDirty },
  } = discountForm;

  // Submit handler
  const onSubmit = async (formData: CreateDiscountOutputData) => {
    try {
      const { error } = await clientSideFetch<string>(
        ADMIN_API_PATHS.DISCOUNTS,
        {
          method: 'POST',
          body: formData,
        }
      );
      // In a real implementation, this would send the data to an API

      if (error) {
        toast.error(
          <div>
            Failed to create discount:
            <div>{JSON.stringify(error)}</div>
          </div>
        );
        return;
      }

      // Redirect to the discounts list page
      router.push('/admin/discounts');
    } catch (err) {
      toast.error(
        <div>
          Failed to create discount:
          <div>{JSON.stringify(err)}</div>
        </div>
      );
      console.error('Error creating discount:', err);
    }
  };

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
          <h1 className='text-2xl font-semibold'>Create New Discount</h1>
        </div>
      </div>

      <FormProvider {...discountForm}>
        <form onSubmit={handleSubmit(onSubmit, (err) => console.log(err))}>
          <TabGroup>
            <TabNavigation
              isEditMode={false}
              activeTab={activeTab}
              setActiveTab={setActiveTab}
            />
            <TabPanels>
              <DetailsPanel />
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
              disabled={isSubmitting || !isDirty}
              className={`px-4 py-2 cursor-pointer bg-primary text-white rounded-md hover:bg-primary/80 ${
                isSubmitting ? 'opacity-75 cursor-not-allowed' : ''
              }`}
            >
              {isSubmitting ? 'Creating...' : 'Create Discount'}
            </Button>
          </div>
        </form>
      </FormProvider>
    </div>
  );
}
