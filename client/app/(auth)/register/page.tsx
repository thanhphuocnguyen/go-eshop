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
import { useActionState } from 'react';

export default function RegisterPage() {
  const [state, action, pending] = useActionState(register, undefined);

  return (
    <div className='w-screen h-[90vh] flex items-center justify-center'>
      <div className='w-full max-w-lg px-4'>
        <Fieldset
          as='form'
          action={action}
          className='space-y-6 rounded-xl bg-gray-700 bg-light-green shadow-md p-6 sm:p-10'
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
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Password</Label>
            <Input
              type='password'
              name='password'
              id='password'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
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
              <Label className='text-red-500 text-sm/6 mt-1'>
                {state.errors.confirmPassword.join(', ')}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Email</Label>
            <Input
              type='email'
              id='email'
              name='email'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {state?.errors?.email && (
              <Label className='text-red-500 text-sm/6 mt-1'>
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
              id='phone'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {state?.errors?.phone && (
              <Label className='text-red-500 text-sm/6 mt-1'>
                {state?.errors.phone.join(', ')}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>FullName</Label>
            <Input
              type='tel'
              id='fullname'
              name='fullname'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {state?.errors?.fullname && (
              <Label className='text-red-500 text-sm/6 mt-1'>
                {state?.errors.fullname.join(', ')}
              </Label>
            )}
          </Field>
          <Button className='btn mr-0 ml-auto' disabled={pending} type='submit'>
            {pending ? 'Submitting' : 'Register'}
          </Button>
        </Fieldset>
      </div>
    </div>
  );
}
