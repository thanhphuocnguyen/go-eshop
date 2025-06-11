'use client';

import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { AttributeFormModel, AttributeFormSchema } from '@/app/lib/definitions';

import { Button, Field, Fieldset, Input, Label } from '@headlessui/react';
import { zodResolver } from '@hookform/resolvers/zod';
import clsx from 'clsx';
import { use, useEffect } from 'react';
import { FormProvider, useFieldArray, useForm } from 'react-hook-form';
import { toast } from 'react-toastify';
import useSWR from 'swr';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from '@dnd-kit/core';
import {
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import ValueItem from './_components/ValueItem';
import { clientSideFetch } from '@/app/lib/api/apiClient';

export default function Page({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const form = useForm<AttributeFormModel>({
    resolver: zodResolver(AttributeFormSchema),
    defaultValues: {
      name: '',
      values: [
        {
          code: '',
          name: '',
          displayOrder: 1,
          isActive: true,
        },
      ],
    },
  });

  const { register, control, reset, handleSubmit } = form;

  const { fields, append, remove, replace } = useFieldArray({
    control,
    name: 'values',
    keyName: 'key',
  });

  const { data: attribute } = useSWR(
    id ? ADMIN_API_PATHS.ATTRIBUTE.replace(':id', id) : null,
    async (url) => {
      const response = await clientSideFetch<AttributeFormModel>(url, {});
      if (response.error) {
        toast('Failed to fetch attribute', { type: 'error' });
        return;
      }
      return response.data;
    },
    {
      refreshInterval: 0,
      dedupingInterval: 60000, // 1 minute deduplication
    }
  );

  async function submitHandler(data: AttributeFormModel) {
    data.values = data.values?.map((item, i) => ({
      ...item,
      displayOrder: i + 1,
    }));

    const response = await clientSideFetch<AttributeFormModel>(
      ADMIN_API_PATHS.ATTRIBUTE.replace(':id', id),
      {
        method: 'PUT',
        body: data,
      }
    );
    if (response.error) {
      toast('Failed to update attribute', { type: 'error' });
      return;
    }
    toast('Attribute updated successfully', { type: 'success' });
    reset(response.data);
  }

  useEffect(() => {
    if (attribute) {
      reset(attribute);
    }
  }, [attribute, reset]);

  return (
    <div className='p-5 bg-white overflow-y-auto'>
      <Link
        className='flex text-sky-500 hover:underline underline-offset-1'
        href='/admin/attributes'
      >
        <ArrowLeftIcon className='h-5 w-5 mr-2' />
        <span className='hidden md:inline'>Back to </span>
        Attributes
      </Link>
      <Fieldset
        as='form'
        className='w-full'
        onSubmit={handleSubmit(submitHandler, (err) => {
          console.log(err);
        })}
      >
        <FormProvider {...form}>
          <div className='flex justify-between mb-5'>
            <h1 className='text-lg font-bold text-primary'>Attribute Edit</h1>
            <Button type='submit' className={clsx('btn btn-primary')}>
              Update
            </Button>
          </div>

          <Field>
            <Label className='text-sm/6 font-semibold'>Name</Label>
            <Input
              {...register('name')}
              className={clsx(
                'mt-1 block w-full rounded-lg border border-form-field-outline bg-white h-12 py-1.5 px-3 text-sm/6 text-form-field-contrast-text',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-form-field-outline-hover'
              )}
            />
          </Field>

          <Field className='mt-4'>
            <div className='flex items-center justify-between'>
              <Label className='text-sm/6 font-semibold'>Values</Label>
              <Button
                type='button'
                onClick={() =>
                  append({
                    code: '',
                    name: '',
                    displayOrder: fields.length,
                    isActive: true,
                  })
                }
                className='btn btn-primary'
              >
                Add Value
              </Button>
            </div>
            <DndContext
              sensors={sensors}
              collisionDetection={closestCenter}
              onDragEnd={handleDragEnd}
            >
              <Field as='ul' className='mt-2 gap-4 w-full'>
                <SortableContext
                  items={fields.map((item) => item.key)}
                  id='key'
                  strategy={verticalListSortingStrategy}
                >
                  {fields.map((item, idx) => (
                    <ValueItem
                      item={item}
                      remove={remove}
                      id={item.key}
                      idx={idx}
                      key={item.key}
                    />
                  ))}
                </SortableContext>
              </Field>
            </DndContext>
          </Field>
        </FormProvider>
      </Fieldset>
    </div>
  );

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;

    if (active.id !== over?.id) {
      const oldIndex = fields.findIndex((item) => item.key === active.id);
      const newIndex = fields.findIndex((item) => item.key === over?.id);

      if (oldIndex !== -1 && newIndex !== -1) {
        // Create a new array with the moved item in the correct position
        const updatedFields = [...fields];
        const [movedItem] = updatedFields.splice(oldIndex, 1);
        updatedFields.splice(newIndex, 0, movedItem);

        // Update the field array values using replace function
        // This is more efficient than using setValue for field arrays
        const newValues = updatedFields.map((field, index) => ({
          ...field,
          displayOrder: index + 1, // Update display order based on new position
        }));

        replace(newValues);
      }
    }
  }
}
