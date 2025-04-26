import { API_PATHS } from '@/lib/constants/api';
import {
  GenericResponse,
  BrandFormSchema,
  GeneralCategoryModel,
  BrandRequest,
} from '@/lib/definitions';
import {
  Button,
  Dialog,
  DialogPanel,
  DialogTitle,
  Field,
  Input,
  Fieldset,
  Label,
} from '@headlessui/react';
import clsx from 'clsx';
import React from 'react';
import { getCookie } from 'cookies-next';
import { toast } from 'react-toastify';
import { apiFetch } from '@/lib/apis/api';

interface BrandFormDialogProps {
  open: boolean;
  onClose: () => void;
  handleSubmitted: (brand: GeneralCategoryModel) => void;
  selectedBrand: GeneralCategoryModel | null;
}

export const BrandFormDialog: React.FC<BrandFormDialogProps> = ({
  onClose,
  open,
  handleSubmitted,
  selectedBrand,
}) => {
  const [isLoading, setIsLoading] = React.useState(false);
  const [state, setState] = React.useState<{
    name?: string[];
  }>({
    name: [],
  });
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);
    const formData = new FormData(e.currentTarget);
    const name = formData.get('name') as string;
    const description = formData.get('description') as string;
    const image_url = formData.get('image_url') as string;
    const slug = formData.get('slug') as string;
    const validatedFields = BrandFormSchema.safeParse({
      name,
      description,
      image_url,
    });
    if (!validatedFields.success) {
      setState(validatedFields.error.flatten().fieldErrors);
      setIsLoading(false);
      return;
    }
    const body: BrandRequest = {
      name,
      description,
      image_url,
      slug,
    };
    const response = await apiFetch(API_PATHS.BRANDS, {
      method: 'POST',
      body,
    });
    if (response.status !== 201) {
      toast('Failed to create brand', { type: 'error' });
      setIsLoading(false);
      return;
    }
    const data: GenericResponse<GeneralCategoryModel> = await response.json();
    toast(data.message, { type: 'success' });
    handleSubmitted(data.data);
    onClose();
  };
  return (
    <Dialog
      open={open}
      as='div'
      className='relative z-10 focus:outline-none bg-green-200'
      onClose={onClose}
    >
      <div className='fixed inset-0 z-10 w-screen overflow-y-auto'>
        <div className='flex min-h-full items-center justify-center p-3'>
          <DialogPanel
            transition
            className='w-full max-w-md rounded-xl bg-green-200 p-6 backdrop-blur-2xl duration-300 ease-out data-[closed]:transform-[scale(95%)] data-[closed]:opacity-0'
          >
            <DialogTitle
              as='h3'
              className='text-center text-xl/7 font-bold text-gray-500'
            >
              Add new attribute
            </DialogTitle>
            <div>
              <Fieldset
                as='form'
                onSubmit={handleSubmit}
                className='flex my-3 w-full flex-col justify-center space-y-5'
              >
                <Field>
                  <Label className='text-sm/6 font-medium text-gray-800'>
                    Name
                  </Label>
                  <Input
                    disabled={false}
                    id='name'
                    name='name'
                    defaultValue={selectedBrand?.name ?? ''}
                    placeholder='Enter attribute name'
                    className={clsx(
                      'mt-1  block w-full rounded-lg border-none bg-yellow-300 h-12 py-1.5 px-3 text-sm/6 text-gray-800',
                      'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white'
                    )}
                  />
                  {state?.name && (
                    <Label className='text-red-500 text-sm/6 mt-1'>
                      {state?.name.join(', ')}
                    </Label>
                  )}
                </Field>
                <Field>
                  <Label className='text-sm/6 font-medium text-gray-800'>
                    Description
                  </Label>
                  <Input
                    disabled={false}
                    id='description'
                    name='description'
                    defaultValue={selectedBrand?.description ?? ''}
                    placeholder='Enter description'
                    className={clsx(
                      'mt-1  block w-full rounded-lg border-none bg-yellow-300 h-12 py-1.5 px-3 text-sm/6 text-gray-800',
                      'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white'
                    )}
                  />
                  {state?.name && (
                    <Label className='text-red-500 text-sm/6 mt-1'>
                      {state?.name.join(', ')}
                    </Label>
                  )}
                </Field>
                <Field>
                  <Label className='text-sm/6 font-medium text-gray-800'>
                    Slug
                  </Label>
                  <Input
                    disabled={false}
                    id='slug'
                    name='slug'
                    defaultValue={selectedBrand?.slug ?? ''}
                    placeholder='Enter slug'
                    className={clsx(
                      'mt-1  block w-full rounded-lg border-none bg-yellow-300 h-12 py-1.5 px-3 text-sm/6 text-gray-800',
                      'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white'
                    )}
                  />
                  {state?.name && (
                    <Label className='text-red-500 text-sm/6 mt-1'>
                      {state?.name.join(', ')}
                    </Label>
                  )}
                </Field>
                <Button
                  disabled={isLoading}
                  className='btn bg-blue-300 hover:bg-blue-500 mx-auto'
                  type='submit'
                >
                  {isLoading ? 'Creating' : 'Create'}
                </Button>
              </Fieldset>
            </div>
          </DialogPanel>
        </div>
      </div>
    </Dialog>
  );
};
