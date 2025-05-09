import { GeneralCategoryModel } from '@/lib/definitions';
import {
  Button,
  Field,
  Fieldset,
  Input,
  Label,
  Legend,
} from '@headlessui/react';
import clsx from 'clsx';
import Image from 'next/image';
import React, { useEffect } from 'react';
import { z } from 'zod';

interface CategoryEditFormProps {
  data?: GeneralCategoryModel;
  title: string;
  handleSave: (data: FormData) => Promise<void>;
}

const UpdateCategoryFormSchema = z.object({
  name: z.string().nonempty(),
  description: z.string().optional(),
  display_order: z.number().optional(),
  slug: z.string().nonempty(),
});

export const CategoryEditForm: React.FC<CategoryEditFormProps> = ({
  title,
  data,
  handleSave,
}) => {
  const [base64, setBase64] = React.useState<string | null>(null);
  const [file, setFile] = React.useState<File | null>(null);
  const [isLoading, setIsLoading] = React.useState(false);
  const [state, setState] = React.useState<{
    name?: string[];
    slug?: string[];
  }>({});

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const parseResult = UpdateCategoryFormSchema.safeParse({
      name: formData.get('name') as string,
      slug: formData.get('slug') as string,
    });
    if (!parseResult.success) {
      const errors = parseResult.error.flatten().fieldErrors;
      setState(errors);
      return;
    }
    setIsLoading(true);
    if (!file) {
      formData.delete('image');
    }
    await handleSave(formData);
    setIsLoading(false);
  };

  useEffect(() => {
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setBase64(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  }, [file]);

  return (
    <Fieldset
      onSubmit={handleSubmit}
      as='form'
      className='space-y-6 rounded-xl flex flex-col justify-center p-0 sm:p-10'
    >
      <div className='flex justify-between'>
        <Legend className='font-bold text-2xl text-gray-600'>{title}</Legend>
        <Button
          disabled={isLoading}
          type='submit'
          className={clsx(
            'btn btn-lg btn-primary btn-elevated',
            isLoading
              ? 'cursor-not-allowed btn-secondary'
              : 'cursor-pointer btn-green'
          )}
        >
          {isLoading ? 'Saving' : 'Save'}
        </Button>
      </div>
      <div className='flex space-x-6'>
        <div className={'flex-1 flex flex-col space-y-5'}>
          <div className='w-full flex space-x-2'>
            <Field as='div' className='flex-1'>
              <Label className='text-sm/3 font-medium text-gray-600'>
                Name
              </Label>
              <Input
                disabled={false}
                id='name'
                name='name'
                placeholder='Enter name...'
                defaultValue={data?.name ?? ''}
                className={clsx(
                  'mt-1 block w-full rounded-lg border border-blue-400 bg-white h-12 py-1.5 px-3 text-sm/6 text-gray-600',
                  'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-green-500'
                )}
              />
              {state?.name && (
                <Label className='text-red-500 text-sm/6 mt-1'>
                  {state.name.join(', ')}
                </Label>
              )}
            </Field>
            <Field as='div'>
              <Label className='text-sm/3 font-medium text-gray-600'>
                Slug
              </Label>
              <Input
                type='text'
                placeholder='Enter slug...'
                id='slug'
                defaultValue={data?.slug ?? ''}
                name='slug'
                className={clsx(
                  'mt-1 block w-full rounded-lg border border-blue-400 bg-white h-12 py-1.5 px-3 text-sm/6 text-gray-600',
                  'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-green-500'
                )}
              />
              {state?.slug && (
                <Label className='text-red-500 text-sm/6 mt-1'>
                  {state.slug.join(', ')}
                </Label>
              )}
            </Field>
          </div>
          <Field as='div'>
            <Label className='text-sm/3 font-medium text-gray-600'>
              Description
            </Label>
            <Input
              type='text'
              placeholder='Enter description'
              id='description'
              name='description'
              defaultValue={data?.description ?? ''}
              className={clsx(
                'mt-1 block w-full rounded-lg border border-blue-400 bg-white h-12 py-1.5 px-3 text-sm/6 text-gray-600',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-green-500'
              )}
            />
          </Field>
        </div>

        <div className='flex-2'>
          <Label
            title='Upload image'
            className='relative h-full w-40 shadow-md shadow-lime-500 text-sm/6 font-medium text-gray-600'
          >
            <Image
              className='rounded-sm object-contain h-auto w-auto'
              width={150}
              height={150}
              alt='Category Image'
              priority
              src={
                base64 || data?.image_url || '/images/product-placeholder.webp'
              }
            />
            <Input
              type='file'
              className={
                'absolute cursor-pointer inset-0 w-full h-full opacity-0'
              }
              id='image'
              name='image'
              accept='image/*'
              onChange={(e) => {
                setFile(e.target?.files?.[0] ?? null);
              }}
            />
          </Label>
        </div>
      </div>
    </Fieldset>
  );
};
