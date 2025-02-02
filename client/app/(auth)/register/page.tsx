import {
  Button,
  Field,
  Fieldset,
  Input,
  Label,
  Legend,
} from '@headlessui/react';
import { zodResolver } from '@hookform/resolvers/zod';
import clsx from 'clsx';
import { useForm } from 'react-hook-form';
import { RegisterForm, registerSchema } from './_lib/type';

export default function RegisterPage() {
  const {
    handleSubmit,
    register,
    formState: { isSubmitting, isValid, isDirty, errors },
  } = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      confirmPassword: '',
      email: '',
      fullname: '',
      password: '',
      phone: '',
      username: '',
    },
  });

  const onSubmit = async (data: RegisterForm) => {
    console.log(data);
  };

  return (
    <div className='w-full h-[90vh] flex items-center justify-center'>
      <div className='w-full max-w-lg px-4'>
        <Fieldset
          as='form'
          onSubmit={handleSubmit(onSubmit)}
          className='space-y-6 rounded-xl bg-light-green shadow-md p-6 sm:p-10'
        >
          <Legend className='text-base/7 font-semibold text-white'>
            Login
          </Legend>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Username</Label>
            <Input
              {...register('username')}
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
              {...register('password')}
              type='password'
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
              {...register('confirmPassword')}
              type='password'
              id='confirm-password'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {errors.confirmPassword && (
              <Label className='text-red-500 text-sm/6 mt-1'>
                {errors.confirmPassword?.message}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Email</Label>
            <Input
              {...register('email')}
              type='email'
              id='email'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {errors.email && (
              <Label className='text-red-500 text-sm/6 mt-1'>
                {errors.email?.message}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>
              Phone Number
            </Label>
            <Input
              type='tel'
              {...register('phone')}
              id='phone-number'
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {errors.phone && (
              <Label className='text-red-500 text-sm/6 mt-1'>
                {errors.phone?.message}
              </Label>
            )}
          </Field>
          <Button disabled={!isDirty && isValid} type='submit'>
            {isSubmitting ? 'Submitting' : 'Register'}
          </Button>
        </Fieldset>
      </div>
    </div>
  );
}
