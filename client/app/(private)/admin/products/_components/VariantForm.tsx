'use client';
import { ProductModelForm, VariantModelForm } from '@/lib/definitions';
import React from 'react';
import { useFormContext, useWatch } from 'react-hook-form';
import { TextField } from '@/components/FormFields';

import clsx from 'clsx';
import { Field, Switch } from '@headlessui/react';
import { StyledComboBoxController } from '@/components/FormFields/StyledComboBoxController';
import { useAttributes } from '../../_lib/hooks/useAttributes';

interface AttributeFormProps {
  index: number;
  item: VariantModelForm;
  onRemove: (index: number) => void;
}
export const VariantForm: React.FC<AttributeFormProps> = ({
  index,
  onRemove,
  item,
}) => {
  const {
    control,
    register,
    getValues,
    setValue,
    formState: { errors },
  } = useFormContext<ProductModelForm>();
  const { attributes } = useAttributes();
  const isActive = useWatch({
    control,
    name: `variants.${index}.is_active`,
  });

  return (
    <div className='w-full'>
      <div className='flex w-full justify-between px-2 py-1 text-left mb-0'>
        <div className='flex items-center gap-5'>
          <span className='text-lg font-semibold'>
            {item.sku ?? `Variant ${index + 1}`}
          </span>
          <Field className='flex items-center gap-2'>
            <Switch
              checked={!!isActive}
              onChange={(checked) =>
                setValue(`variants.${index}.is_active`, checked)
              }
              className={clsx(
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2',
                isActive ? 'bg-primary' : 'bg-gray-200'
              )}
            >
              <span className='sr-only'>Active Variant</span>
              <span
                className={clsx(
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  isActive ? 'translate-x-6' : 'translate-x-1'
                )}
              />
            </Switch>
            <span
              className='font-semibold cursor-pointer'
              onClick={() => {
                setValue(`variants.${index}.is_active`, !isActive);
              }}
            >
              Active
            </span>
          </Field>
        </div>
        <div className='flex items-center'>
          <button
            type='button'
            className={clsx(
              'btn btn-danger btn-sm transition-all duration-300 hover:scale-105'
            )}
            onClick={() => onRemove(index)}
          >
            Remove
          </button>
        </div>
      </div>

      <div className='px-2 pt-2 pb-0'>
        <div className='grid grid-cols-5 gap-x-3'>
          {getValues(`variants.${index}.attributes`)?.map((attribute, idx) =>
            attribute ? (
              <StyledComboBoxController
                control={control}
                key={attribute.id}
                error={!!errors.variants?.[index]?.attributes?.[idx]?.value?.id}
                message={
                  errors.variants?.[index]?.attributes?.[idx]?.value?.id
                    ?.message
                }
                getDisplayValue={(attribute) =>
                  attribute.display_value || attribute.value
                }
                name={`variants.${index}.attributes.${idx}.value`}
                label={attribute.name}
                options={
                  attributes?.find((e) => e.id === attribute.id)?.values ?? []
                }
              />
            ) : null
          )}
          <TextField
            {...register(`variants.${index}.price`)}
            label='Price'
            placeholder='Enter Price'
            error={!!errors.variants?.[index]?.price?.message}
            message={errors.variants?.[index]?.price?.message}
            type='number'
            disabled={false}
          />
          <TextField
            {...register(`variants.${index}.stock`)}
            label='Stock'
            placeholder='Enter Stock'
            error={!!errors.variants?.[index]?.stock?.message}
            message={errors.variants?.[index]?.stock?.message}
            type='number'
            disabled={false}
          />
          <TextField
            {...register(`variants.${index}.weight`)}
            label='Weight'
            placeholder='Enter Weight'
            type='number'
            disabled={false}
          />
        </div>
      </div>
    </div>
  );
};
