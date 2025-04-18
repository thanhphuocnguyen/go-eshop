'use client';
import { API_PATHS } from '@/lib/constants/api';
import {
  Button,
  Dialog,
  DialogPanel,
  DialogTitle,
  Field,
  Fieldset,
  Input,
  Label,
} from '@headlessui/react';
import clsx from 'clsx';
import React, { useState } from 'react';
import { getCookie } from 'cookies-next';
import { toast } from 'react-toastify';
import {
  GenericResponse,
  AttributeFormModel,
  AttributeFormSchema,
  AttributeValueFormModel,
} from '@/lib/definitions';
import { XCircleIcon } from '@heroicons/react/24/outline';

interface AddNewDialogProps {
  selectedAttribute: AttributeFormModel | null;
  open: boolean;
  onClose: () => void;
  handleSubmitted: (attribute: AttributeFormModel) => void;
}

export const AddNewDialog: React.FC<AddNewDialogProps> = ({
  selectedAttribute,
  onClose,
  open,
  handleSubmitted,
}) => {
  const [isLoading, setIsLoading] = React.useState(false);
  const [newValues, setNewValues] = useState<AttributeValueFormModel[]>([]);

  const [state, setState] = React.useState<{
    name?: string[];
    values?: string[];
  }>({
    name: [],
  });

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);
    const formData = new FormData(e.currentTarget);
    const name = formData.get('name') as string;

    const validatedFields = AttributeFormSchema.safeParse({
      name,
    });

    if (!validatedFields.success) {
      setState(validatedFields.error.flatten().fieldErrors);
    }
    const body: AttributeFormModel = {
      name,
      values: newValues,
    };
    if (selectedAttribute?.id) {
      body.id = selectedAttribute.id;
    }
    if (selectedAttribute?.values) {
      body.values = [...selectedAttribute.values, ...newValues];
    }
    const token = await getCookie('token');
    const resp = await apiFetch(
      body.id
        ? API_PATHS.ATTRIBUTE.replace(':id', body.id.toString())
        : API_PATHS.ATTRIBUTES,
      {
        method: body.id ? 'PUT' : 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      }
    );

    setIsLoading(false);
    if (resp.status !== 200 && resp.status !== 201) {
      toast('Failed to create attribute', { type: 'error' });
      return;
    }
    const data: GenericResponse<AttributeFormModel> = await resp.json();
    toast(
      selectedAttribute?.id ? 'Updated successfully' : 'Created successfully',
      { type: 'success' }
    );
    handleSubmitted(data.data);
    setNewValues([]);
    onClose();
  };

  return (
    <Dialog
      open={open}
      as='div'
      className='relative z-10 focus:outline-none'
      onClose={onClose}
    >
      <div className='fixed inset-0 z-10 w-screen overflow-y-auto'>
        <div className='flex min-h-full items-center justify-center p-3'>
          <DialogPanel
            transition
            className='w-full max-w-2xl rounded-xl bg-form-field-bg p-6 backdrop-blur-2xl duration-300 ease-out data-[closed]:transform-[scale(95%)] data-[closed]:opacity-0'
          >
            <DialogTitle
              as='h3'
              className='text-center text-xl/7 font-bold text-form-field-contrast-text'
            >
              {selectedAttribute ? 'Edit' : 'Add new'} attribute
            </DialogTitle>
            <div>
              <Fieldset
                as='form'
                onSubmit={handleSubmit}
                className='flex my-3 w-full flex-col justify-center space-y-5'
              >
                <Field>
                  <Label className='text-sm/6 font-medium text-form-field-label-contrast-text'>
                    Name
                  </Label>
                  <Input
                    disabled={false}
                    id='name'
                    name='name'
                    defaultValue={selectedAttribute?.name ?? ''}
                    placeholder='Enter attribute name'
                    className={clsx(
                      'mt-1 block w-full rounded-lg border border-form-field-outline bg-white h-12 py-1.5 px-3 text-sm/6 text-form-field-contrast-text',
                      'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-form-field-outline-hover'
                    )}
                  />
                  {state?.name && (
                    <Label className='text-red-500 text-sm/6 mt-1'>
                      {state?.name.join(', ')}
                    </Label>
                  )}
                </Field>
                <div>
                  <div className='text-sm/6 font-medium text-form-field-label-contrast-text'>
                    Values
                  </div>
                  <div className='p-2 flex gap-2 flex-wrap bg-white rounded-lg border border-form-field-outline'>
                    {selectedAttribute?.values?.map((value, index) => (
                      <span
                        className='rounded-2xl text-sm text-white px-3 py-2 h-auto bg-secondary'
                        key={value.id ?? index}
                      >
                        {value.value}
                      </span>
                    ))}
                    {newValues.map((value) => (
                      <span
                        className='rounded-2xl flex gap-2 items-center text-sm text-white px-3 py-2 bg-tertiary'
                        key={value.value}
                      >
                        {value.value}
                        <button
                          title='Remove'
                          onClick={() => {
                            setNewValues(
                              newValues.filter((e) => e.value !== value.value)
                            );
                          }}
                          className='text-button-danger cursor-pointer'
                        >
                          <XCircleIcon className='size-6 ' />
                        </button>
                      </span>
                    ))}
                    <Field>
                      <Input
                        disabled={false}
                        id={`new-value`}
                        name={`new-value`}
                        placeholder='Enter to add'
                        onKeyDown={(e) => {
                          if (e.key === 'Enter') {
                            e.preventDefault();
                            const value = e.currentTarget.value;
                            if (value) {
                              setNewValues([...newValues, { value: value }]);
                              e.currentTarget.value = '';
                            }
                          }
                        }}
                        className={clsx(
                          ' block w-36 rounded-lg bg-white py-1.5 px-3 text-sm/6 text-form-field-contrast-text',
                          'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2'
                        )}
                      />
                    </Field>
                  </div>
                </div>
                <Button
                  disabled={isLoading}
                  className='btn bg-sky-300 hover:bg-sky-500 mx-auto'
                  type='submit'
                >
                  {isLoading ? 'Submitting' : 'Submit'}
                </Button>
              </Fieldset>
            </div>
          </DialogPanel>
        </div>
      </div>
    </Dialog>
  );
};
