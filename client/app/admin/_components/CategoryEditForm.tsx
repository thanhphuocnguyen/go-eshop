import { GeneralCategoryModel } from '@/app/lib/definitions';
import {
  Button,
  Field,
  Fieldset,
  Input,
  Label,
  Legend,
  Switch,
} from '@headlessui/react';
import clsx from 'clsx';
import React from 'react';
import { z } from 'zod';
import ImageUploader from '@/app/components/ImageUploader';

interface CategoryEditFormProps {
  data?: GeneralCategoryModel;
  title: string;
  handleSave: (data: FormData) => Promise<void>;
}

const UpdateCategoryFormSchema = z.object({
  name: z.string().nonempty(),
  description: z.string().optional(),
  displayOrder: z.number().optional(),
  slug: z.string().nonempty(),
  published: z.boolean().optional(),
  remarkable: z.boolean().optional(),
});

export const CategoryEditForm: React.FC<CategoryEditFormProps> = ({
  title,
  data,
  handleSave,
}) => {
  const [file, setFile] = React.useState<File | null>(null);
  const [isLoading, setIsLoading] = React.useState(false);
  const [published, setPublished] = React.useState<boolean>(
    data?.published ?? true
  );
  const [remarkable, setRemarkable] = React.useState<boolean>(
    data?.remarkable ?? false
  );
  const [state, setState] = React.useState<{
    name?: string[];
    slug?: string[];
  }>({});

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);

    // Add toggle values to form data
    formData.append('published', published.toString());
    formData.append('remarkable', remarkable.toString());

    const parseResult = UpdateCategoryFormSchema.safeParse({
      name: formData.get('name') as string,
      slug: formData.get('slug') as string,
      published: published,
      remarkable: remarkable,
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

          <div className='flex space-x-8 mt-6'>
            <Field as='div' className='flex items-center gap-3'>
              <Label className='text-sm/3 font-medium text-gray-600'>
                Published
              </Label>
              <Switch
                checked={published}
                onChange={setPublished}
                className={`${
                  published ? 'bg-green-500' : 'bg-gray-300'
                } relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2`}
              >
                <span
                  className={`${
                    published ? 'translate-x-6' : 'translate-x-1'
                  } inline-block h-4 w-4 transform rounded-full bg-white transition-transform`}
                />
              </Switch>
            </Field>

            <Field as='div' className='flex items-center gap-3'>
              <Label className='text-sm/3 font-medium text-gray-600'>
                Remarkable
              </Label>
              <Switch
                checked={remarkable}
                onChange={setRemarkable}
                className={`${
                  remarkable ? 'bg-blue-500' : 'bg-gray-300'
                } relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2`}
              >
                <span
                  className={`${
                    remarkable ? 'translate-x-6' : 'translate-x-1'
                  } inline-block h-4 w-4 transform rounded-full bg-white transition-transform`}
                />
              </Switch>
            </Field>
          </div>
        </div>

        <div className='flex-2'>
          <ImageUploader
            defaultImage={data?.imageUrl}
            name='image'
            label='Upload image'
            onChange={(newFile) => {
              setFile(newFile);
            }}
            className='w-40 h-auto'
            width={150}
            height={150}
            maxFileSizeMB={2.5}
            aspectRatio={1}
          />
        </div>
      </div>
    </Fieldset>
  );
};
