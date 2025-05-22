'use client';

import clsx from 'clsx';
import React from 'react';
import { CSS } from '@dnd-kit/utilities';
import { useSortable } from '@dnd-kit/sortable';
import { Button, Field, Input, Label, Switch } from '@headlessui/react';
import { Controller, useFormContext } from 'react-hook-form';
import { AttributeFormModel, AttributeValueFormModel } from '@/app/lib/definitions';
import { Bars3Icon, XCircleIcon } from '@heroicons/react/24/outline';

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
      className={clsx(
        'border w-full border-form-field-outline rounded-lg py-2 px-6 mb-2'
      )}
      style={style}
    >
      <div className={clsx('flex gap-4 items-center')}>
        <Button {...listeners} {...attributes} className={'col-span-1'}>
          <Bars3Icon className='size-5 text-form-field-contrast-text' />
        </Button>
        <div className='w-56'>
          <Field>
            <Label className='text-sm/6 font-semibold'>Code</Label>
            <Input
              {...register(`values.${idx}.code`)}
              className={clsx(
                'mt-1 block w-full rounded-lg border border-form-field-outline bg-white h-12 py-1.5 px-3 text-sm/6 text-form-field-contrast-text',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-form-field-outline-hover'
              )}
            />
          </Field>
        </div>
        <div className='w-56'>
          <Field>
            <Label className='text-sm/6 font-semibold'>Name</Label>
            <Input
              {...register(`values.${idx}.name`)}
              className={clsx(
                'mt-1 block w-full rounded-lg border border-form-field-outline bg-white h-12 py-1.5 px-3 text-sm/6 text-form-field-contrast-text',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-form-field-outline-hover'
              )}
            />
          </Field>
        </div>
        <div className='w-40'>
          <Field className='flex items-center gap-2'>
            <Controller
              control={control}
              name={`values.${idx}.isActive`}
              render={({ field: { value, onChange, ...rest } }) => (
                <Switch
                  {...rest}
                  checked={value}
                  onChange={onChange}
                  className={`${
                    value ? 'bg-blue-500' : 'bg-gray-300'
                  } relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2`}
                >
                  <span className='sr-only'>Active</span>
                  <span
                    className={`${
                      value ? 'translate-x-6' : 'translate-x-1'
                    } inline-block h-4 w-4 transform rounded-full bg-white transition-transform`}
                  />
                </Switch>
              )}
            />
            <Label>Active</Label>
          </Field>
        </div>
        <div className='flex flex-1 items-center justify-end w-20'>
          <button type='button' onClick={() => remove(idx)}>
            <XCircleIcon className='size-6 text-white bg-danger rounded-full' />
          </button>
        </div>
      </div>
    </li>
  );
};

export default ValueItem;
