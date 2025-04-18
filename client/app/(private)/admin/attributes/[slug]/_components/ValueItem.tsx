'use client';
import clsx from 'clsx';
import React from 'react';
import { CSS } from '@dnd-kit/utilities';
import { useSortable } from '@dnd-kit/sortable';
import { Checkbox, Field, Input, Label } from '@headlessui/react';
import { Controller, useFormContext } from 'react-hook-form';
import { AttributeFormModel, AttributeValueFormModel } from '@/lib/definitions';
import { XCircleIcon } from '@heroicons/react/24/outline';

interface ValueItemProps {
  idx: number;
  id: string;
  remove: (index: number) => void;
  item: AttributeValueFormModel;
}
const ValueItem: React.FC<ValueItemProps> = ({ idx, id, remove }) => {
  const { register, control } = useFormContext<AttributeFormModel>();
  const { attributes, listeners, setNodeRef, transform, transition } =
    useSortable({ id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };
  return (
    <li
      ref={setNodeRef}
      className={clsx('border border-form-field-outline rounded-lg p-4 mb-2')}
      {...attributes}
      {...listeners}
      style={style}
    >
      <div className={clsx('grid grid-cols-6 gap-4 items-center')}>
        <div className='col-span-2'>
          <Field>
            <Label className='text-sm/6 font-semibold'>Value</Label>
            <Input
              {...register(`values.${idx}.value`)}
              className={clsx(
                'mt-1 block w-full rounded-lg border border-form-field-outline bg-white h-12 py-1.5 px-3 text-sm/6 text-form-field-contrast-text',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-form-field-outline-hover'
              )}
            />
          </Field>
        </div>
        <div className='col-span-2'>
          <Field>
            <Label className='text-sm/6 font-semibold'>Display Value</Label>
            <Input
              {...register(`values.${idx}.display_value`)}
              className={clsx(
                'mt-1 block w-full rounded-lg border border-form-field-outline bg-white h-12 py-1.5 px-3 text-sm/6 text-form-field-contrast-text',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-form-field-outline-hover'
              )}
            />
          </Field>
        </div>
        <div className='col-span-1'>
          <Field className='flex items-center gap-2'>
            <Controller
              control={control}
              name={`values.${idx}.is_active`}
              render={({ field: { value, onChange, ...rest } }) => (
                <Checkbox
                  {...rest}
                  checked={value}
                  onChange={(checked) => {
                    onChange(checked);
                  }}
                  className={clsx(
                    'group block size-5 rounded border bg-white data-[checked]:bg-blue-500'
                  )}
                  id={`values.${idx}.is_default`}
                >
                  <svg
                    className='stroke-white opacity-0 group-data-[checked]:opacity-100'
                    viewBox='0 0 14 14'
                    fill='none'
                  >
                    <path
                      d='M3 8L6 11L11 3.5'
                      strokeWidth={2}
                      strokeLinecap='round'
                      strokeLinejoin='round'
                    />
                  </svg>
                </Checkbox>
              )}
            />
            <Label>Active</Label>
          </Field>
        </div>
        <div className='col-span-1 flex justify-end'>
          <button type='button' onClick={() => remove(idx)}>
            <XCircleIcon className='size-6 text-white bg-danger rounded-full' />
          </button>
        </div>
      </div>
    </li>
  );
};

export default ValueItem;
