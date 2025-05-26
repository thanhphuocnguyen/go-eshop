'use client';
import React from 'react';
import clsx from 'clsx';
import { Field, Switch } from '@headlessui/react';
import { useFormContext, useWatch } from 'react-hook-form';
import { ProductModelForm, VariantModelForm } from '@/app/lib/definitions';
import { TextField } from '@/app/components/FormFields';
import { StyledComboBoxController } from '@/app/components/FormFields/StyledComboBoxController';
import { useAttributes } from '../../hooks/useAttributes';

interface AttributeFormProps {
  index: number;
  item: VariantModelForm;
  selectedAttributes: string[];
  onRemove: (index: number) => void;
}
export const VariantForm: React.FC<AttributeFormProps> = ({
  index,
  item,
  selectedAttributes,
  onRemove,
}) => {
  const {
    control,
    register,
    setValue,
    formState: { errors },
  } = useFormContext<ProductModelForm>();
  const { attributes } = useAttributes();

  const isActive = useWatch({
    control,
    name: `variants.${index}.isActive`,
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
                setValue(`variants.${index}.isActive`, checked)
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
                setValue(`variants.${index}.isActive`, !isActive);
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
          {selectedAttributes?.map((attribute, idx) => {
            const attr = attributes?.find((e) => e.id === attribute);
            return attribute ? (
              <StyledComboBoxController
                control={control}
                key={attribute}
                error={
                  errors.variants?.[index]?.attributes?.[idx]?.valueObject?.id
                    ?.message
                }
                getDisplayValue={(attrOpt) =>
                  attr && attr.name.toLowerCase().startsWith('color')
                    ? attrOpt?.name
                    : attrOpt.code
                }
                name={`variants.${index}.attributes.${idx}.valueObject`}
                label={attr?.name ?? ''}
                options={
                  attributes?.find((e) => e.id === attr?.id)?.values ?? []
                }
              />
            ) : null;
          })}
          <TextField
            {...register(`variants.${index}.price`)}
            label='Price'
            placeholder='Enter Price'
            error={errors.variants?.[index]?.price?.message}
            type='number'
            disabled={false}
          />
          <TextField
            {...register(`variants.${index}.stockQty`)}
            label='Stock'
            placeholder='Enter Stock'
            error={errors.variants?.[index]?.stockQty?.message}
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
