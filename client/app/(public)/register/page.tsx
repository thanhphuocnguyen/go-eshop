'use client';

import { register } from '@/app/actions/auth';
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

export default function RegisterPage() {
  const [state, action, pending] = useActionState(register, undefined);

  useEffect(() => {
    if (state?.message) {
      toast(state.message, { type: 'error' });
    }
  }, [state]);

  return (
    <div className='h-[90vh] flex items-center justify-center'>
      <div className='px-4 w-1/3'>
        <Fieldset
          as='form'
          action={action}
          className='space-y-6 rounded-xl bg-gray-700 shadow-md p-6 sm:p-10'
        >
          <Legend className='text-base/7 font-semibold text-white'>
            Login
          </Legend>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Username</Label>
            <Input
              id='username'
              name='username'
              value={state?.data?.username}
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {state?.errors?.username && (
              <Label className='text-button-danger text-sm/6 mt-1'>
                {state.errors.username.join(', ')}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Password</Label>
            <Input
              type='password'
              name='password'
              value={state?.data?.password}
              id='password'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {state?.errors?.password && (
              <Label className='text-button-danger text-sm/6 mt-1'>
                {state.errors.password.join(', ')}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>
              Confirm Password
            </Label>
            <Input
              type='password'
              id='confirmPassword'
              name='confirmPassword'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {state?.errors?.confirmPassword && (
              <Label className='text-button-danger text-sm/6 mt-1'>
                {state.errors.confirmPassword.join(', ')}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Email</Label>
            <Input
              type='email'
              value={state?.data?.email}
              id='email'
              name='email'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {state?.errors?.email && (
              <Label className='text-button-danger text-sm/6 mt-1'>
                {state.errors.email?.join(', ')}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>
              Phone Number
            </Label>
            <Input
              type='tel'
              name='phone'
              value={state?.data?.phone}
              id='phone'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {state?.errors?.phone && (
              <Label className='text-button-danger text-sm/6 mt-1'>
                {state?.errors.phone.join(', ')}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>
              Full Name
            </Label>
            <Input
              type='text'
              value={state?.data?.fullname}
              id='fullname'
              name='fullname'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {state?.errors?.fullname && (
              <Label className='text-button-danger text-sm/6 mt-1'>
                {state?.errors.fullname.join(', ')}
              </Label>
            )}
          </Field>
          <Button
            className='btn btn-lg btn-primary mr-0 ml-auto'
            disabled={pending}
            type='submit'
          >
            {pending ? 'Registering' : 'Register'}
          </Button>
        </Fieldset>
      </div>
    </div>
  );
}
