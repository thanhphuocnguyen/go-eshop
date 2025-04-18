'use client';

import { login } from '@/app/actions/auth';
import {
  Button,
  Field,
  Fieldset,
  Input,
  Label,
  Legend,
} from '@headlessui/react';
import clsx from 'clsx';
import { useActionState, useEffect } from 'react';
import { toast } from 'react-toastify';

export default function LoginPage() {
  const [state, action, pending] = useActionState(login, undefined);

  useEffect(() => {
    if (state?.message) {
      toast(state.message, { type: 'error' });
    }
  }, [state]);

  return (
    <div className='w-screen h-[90vh] flex items-center justify-center'>
      <div className='w-full max-w-lg px-4'>
        <Fieldset
          as='form'
          autoComplete='on'
          action={action}
          className='space-y-6 rounded-xl flex flex-col justify-center bg-secondary shadow-md p-6 sm:p-10'
        >
          <Legend className='font-semibold text-xl text-white'>Login</Legend>
          <Field>
            <Label className='text-sm/4 font-medium text-white'>Username</Label>
            <Input
              disabled={false}
              id='username'
              name='username'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white h-12 py-1.5 px-3 text-sm/6 text-gray-500',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-primary'
              )}
            />
            {state?.errors?.username && (
              <Label className='text-red-500 text-sm/6 mt-1'>
                {state?.errors?.username.join(', ')}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Password</Label>
            <Input
              type='password'
              id='password'
              name='password'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white h-12 py-1.5 px-3 text-sm/6 text-gray-500',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-primary',
              )}
            />
            {
              <Label className='text-red-500 text-sm/6 mt-1'>
                {state?.errors?.password?.join(', ')}
              </Label>
            }
          </Field>
          <Button
            className='btn btn-primary btn-lg mx-auto'
            disabled={pending}
            type='submit'
          >
            {pending ? 'Submitting' : 'Submit'}
          </Button>
        </Fieldset>
      </div>
    </div>
  );
}
