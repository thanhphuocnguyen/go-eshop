'use client';

import React from 'react';
import { Control, UseFormRegister, UseFormWatch, FieldErrors } from 'react-hook-form';
import { TagIcon } from '@heroicons/react/24/outline';
import { TabPanel, Textarea } from '@headlessui/react';
import { TextField } from '@/components/FormFields';
import { StyledComboBoxController } from '@/components/FormFields/StyledComboBoxController';
import { ControlledStyledCheckbox } from '@/components/FormFields/ControlledStyledCheckbox';
import { DiscountFormData } from '../_types';

interface DetailsPanelProps {
  register: UseFormRegister<DiscountFormData>;
  control: Control<DiscountFormData>;
  errors: FieldErrors<DiscountFormData>;
  watch: UseFormWatch<DiscountFormData>;
}

export const DetailsPanel: React.FC<DetailsPanelProps> = ({
  register,
  control,
  errors,
  watch,
}) => {
  return (
    <TabPanel className='grid grid-cols-1 md:grid-cols-2 gap-6 pt-6'>
      <div className='col-span-1'>
        <div className='relative'>
          <div className='absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none'>
            <TagIcon className='h-5 w-5 text-gray-400' />
          </div>
          <TextField
            label={
              <label className='block mb-2 text-sm font-medium'>
                Discount Code <span className='text-red-500'>*</span>
              </label>
            }
            type='text'
            {...register('code', {
              required: 'Discount code is required',
            })}
            placeholder='e.g. SUMMER25'
            error={errors.code?.message}
          />
        </div>
        <p className='mt-1 text-xs text-gray-500'>
          Code used by customers to apply the discount
        </p>
      </div>

      <div className='col-span-1'>
        <StyledComboBoxController
          control={control}
          label={
            <label className='block mb-2 text-sm font-medium'>
              Discount Type <span className='text-red-500'>*</span>
            </label>
          }
          name='discountType'
          options={[
            { id: 'percentage', name: 'Percentage' },
            { id: 'fixed_amount', name: 'Fixed Amount' },
          ]}
          error={errors.discountType?.message}
        />

        <p className='mt-1 text-xs text-gray-500'>
          How the discount will be calculated
        </p>
      </div>

      <div className='col-span-1'>
        <div className='relative'>
          <div className='absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none'>
            {watch('discountType') === 'percentage' ? '%' : '$'}
          </div>
          <TextField
            label={
              <label className='block mb-2 text-sm font-medium'>
                Discount Value <span className='text-red-500'>*</span>
              </label>
            }
            type='number'
            step='0.01'
            {...register('discountValue', {
              required: 'Discount value is required',
              min: { value: 0, message: 'Value must be positive' },
              valueAsNumber: true,
            })}
            placeholder={watch('discountType') === 'percentage' ? '10' : '10.00'}
            error={errors.discountValue?.message}
          />
        </div>
        <p className='mt-1 text-xs text-gray-500'>
          {watch('discountType') === 'percentage'
            ? 'Percentage off (0-100)'
            : 'Fixed amount to deduct from order'}
        </p>
      </div>

      <div className='col-span-1'>
        <div className='relative'>
          <div className='absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none'>
            <span className='text-gray-500'>$</span>
          </div>
          <TextField
            label={
              <label className='block mb-2 text-sm font-medium'>
                Minimum Purchase Amount
              </label>
            }
            type='number'
            step='0.01'
            {...register('minPurchaseAmount', {
              min: { value: 0, message: 'Amount must be positive' },
              valueAsNumber: true,
            })}
            placeholder='0.00'
            error={errors.minPurchaseAmount?.message}
          />
        </div>
        <p className='mt-1 text-xs text-gray-500'>
          Minimum order amount required (optional)
        </p>
      </div>

      <div className='col-span-1'>
        <div className='relative'>
          <div className='absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none'>
            <span className='text-gray-500'>$</span>
          </div>
          <TextField
            label={
              <label className='block mb-2 text-sm font-medium'>
                Maximum Discount Amount
              </label>
            }
            type='number'
            step='0.01'
            {...register('maxDiscountAmount', {
              min: { value: 0, message: 'Amount must be positive' },
              valueAsNumber: true,
            })}
            placeholder='0.00'
            error={errors.maxDiscountAmount?.message}
          />
        </div>

        <p className='mt-1 text-xs text-gray-500'>
          Maximum amount to discount (optional)
        </p>
      </div>

      <div className='col-span-1'>
        <TextField
          label={
            <label className='block mb-2 text-sm font-medium'>
              Usage Limit
            </label>
          }
          type='number'
          {...register('usageLimit', {
            min: {
              value: 1,
              message: 'Usage limit must be at least 1',
            },
            valueAsNumber: true,
          })}
          placeholder='No limit'
          error={errors.usageLimit?.message}
        />

        <p className='mt-1 text-xs text-gray-500'>
          Maximum number of times this discount can be used (optional)
        </p>
      </div>

      <div className='col-span-1'>
        <label className='block mb-2 text-sm font-medium'>Status</label>
        <div className='flex items-center'>
          <ControlledStyledCheckbox
            control={control}
            name='isActive'
            label='Active (can be used by customers)'
          />
        </div>
      </div>

      <div className='col-span-1'>
        <TextField
          label={
            <label className='block mb-2 text-sm font-medium'>
              Start Date <span className='text-red-500'>*</span>
            </label>
          }
          type='datetime-local'
          {...register('startsAt', {
            required: 'Start date is required',
          })}
          error={errors.startsAt?.message}
        />
      </div>

      <div className='col-span-1'>
        <TextField
          label={
            <label className='block mb-2 text-sm font-medium'>
              Expiry Date <span className='text-red-500'>*</span>
            </label>
          }
          type='datetime-local'
          {...register('expiresAt', {
            required: 'Expiry date is required',
          })}
          error={errors.expiresAt?.message}
        />
      </div>

      <div className='col-span-2'>
        <label className='block mb-2 text-sm font-medium'>
          Description
        </label>
        <Textarea
          className={`block w-full border ${
            errors.description ? 'border-red-500' : 'border-gray-300'
          } rounded-md px-3 py-2`}
          {...register('description')}
          placeholder='Description of this discount'
          rows={3}
        ></Textarea>
        {errors.description && (
          <p className='mt-1 text-xs text-red-500'>
            {errors.description.message}
          </p>
        )}
      </div>
    </TabPanel>
  );
};
