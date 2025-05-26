'use client';

import React from 'react';
import { useFormContext, useWatch } from 'react-hook-form';
import { TagIcon } from '@heroicons/react/24/outline';
import { TabPanel } from '@headlessui/react';
import { TextAreaField, TextField } from '@/app/components/FormFields';
import { StyledComboBoxController } from '@/app/components/FormFields/StyledComboBoxController';
import { ControlledStyledCheckbox } from '@/app/components/FormFields/ControlledStyledCheckbox';
import { CreateDiscountFormData, discountTypeOptions } from '../_types';
import { AnimatePresence, motion } from 'framer-motion';

interface DetailsPanelProps {}

export const DetailsPanel: React.FC<DetailsPanelProps> = ({}) => {
  const {
    register,
    control,
    formState: { errors },
  } = useFormContext<CreateDiscountFormData>();
  const discountType = useWatch({ control, name: 'discountType' });

  return (
    <TabPanel as={AnimatePresence} mode='wait'>
      <motion.div
        key='details'
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        exit={{ opacity: 0, x: 20 }}
        transition={{ duration: 0.3 }}
      >
        <div className='grid grid-cols-1 md:grid-cols-2 gap-6 pt-6'>
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
              placeholder='Select discount type'
              name='discountType'
              options={discountTypeOptions}
              error={errors.discountType?.message}
            />

            <p className='mt-1 text-xs text-gray-500'>
              How the discount will be calculated
            </p>
          </div>

          <div className='col-span-1'>
            <TextField
              label={
                <div className='relative'>
                  <label className='block ml-4 mb-2 text-sm font-medium'>
                    Discount Value <span className='text-red-500'>*</span>
                  </label>
                  <div className='absolute inset-y-0 left-0 flex items-center pointer-events-none'>
                    {discountType?.id === 'percentage' ? '%' : '$'}
                  </div>
                </div>
              }
              type='number'
              step='0.01'
              {...register('discountValue', {
                required: 'Discount value is required',
                valueAsNumber: true,
              })}
              placeholder={discountType?.id === 'percentage' ? '10' : '10.00'}
              error={errors.discountValue?.message}
            />
            <p className='mt-1 text-xs text-gray-500'>
              {discountType?.id === 'percentage'
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
                {...register('minPurchaseAmount', {})}
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
                {...register('maxDiscountAmount', {})}
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
              {...register('usageLimit', {})}
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
              {...register('startsAt')}
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
            <TextAreaField
              {...register('description')}
              placeholder='Description of this discount'
              rows={3}
            ></TextAreaField>
            {errors.description && (
              <p className='mt-1 text-xs text-red-500'>
                {errors.description.message}
              </p>
            )}
          </div>
        </div>
      </motion.div>
    </TabPanel>
  );
};
