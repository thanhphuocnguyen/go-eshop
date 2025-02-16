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
import { useActionState } from 'react';

export default function LoginPage() {
  const [state, action, pending] = useActionState(login, undefined);

  return (
    <div className='w-screen h-[90vh] flex items-center justify-center'>
      <div className='w-full max-w-lg px-4'>
        <Fieldset
          as='form'
          action={action}
          className='space-y-6 rounded-xl bg-gray-600 shadow-md p-6 sm:p-10'
        >
          <Legend className='text-base/7 font-semibold text-white'>
            Login
          </Legend>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Username</Label>
            <Input
              disabled={false}
              id='username'
              name='username'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
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
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {
              <Label className='text-red-500 text-sm/6 mt-1'>
                {state?.errors?.password?.join(', ')}
              </Label>
            }
          </Field>
          <Button className='btn secondary' type='submit'>
            {pending ? 'Logging in' : 'Login'}
          </Button>
        </Fieldset>
      </div>
    </div>
  );
}
