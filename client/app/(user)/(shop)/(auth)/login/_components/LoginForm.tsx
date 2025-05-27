'use client';

import { Button, Fieldset } from '@headlessui/react';
import clsx from 'clsx';
import { FaLock, FaEnvelope, FaUser } from 'react-icons/fa';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { TextField } from '@/app/components/FormFields';
import { toast } from 'react-toastify';
import { useRouter } from 'next/navigation';
import { z } from 'zod';
import { useState } from 'react';
import { loginAction } from '@/app/actions/auth';

// Login form schema
const loginSchema = z
  .object({
    email: z.string().email({ message: 'Invalid email address' }).optional(),
    username: z
      .string()
      .min(5, { message: 'Username must be at least 5 characters' })
      .optional(),
    password: z
      .string()
      .min(6, { message: 'Password must be at least 6 characters' }),
  })
  .refine(
    (data) => {
      // Ensure either email or username is provided
      return data.email !== undefined || data.username !== undefined;
    },
    {
      message: 'Either email or username must be provided',
      path: ['email'],
    }
  );

type LoginForm = z.infer<typeof loginSchema>;

type LoginMethod = 'email' | 'username';

export default function LoginFormComponent() {
  const router = useRouter();
  const [loginMethod, setLoginMethod] = useState<LoginMethod>('email');

  const form = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: loginMethod === 'email' ? 'admin@eshop.com' : undefined,
      username: loginMethod === 'username' ? 'admin' : undefined,
      password: 'secret123',
    },
    mode: 'onChange',
  });

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    resetField,
    setValue,
  } = form;

  const toggleLoginMethod = () => {
    const newMethod = loginMethod === 'email' ? 'username' : 'email';
    setLoginMethod(newMethod);

    // Clear both fields first
    resetField('email');
    resetField('username');

    // Set default for the active method
    if (newMethod === 'email') {
      setValue('email', 'admin@eshop.com');
    } else {
      setValue('username', 'admin');
    }
  };

  const onSubmit = async (data: LoginForm) => {
    const formData = new FormData();
    if (loginMethod === 'email') {
      formData.append('email', data.email || '');
    } else {
      formData.append('username', data.username || '');
    }
    formData.append('password', data.password);

    const rs = await loginAction(formData);
    if (rs?.error) {
      toast.error(
        <div>
          Login failed. Please check your credentials.
          <div>{JSON.stringify(rs.error.message)}</div>
        </div>
      );
      return;
    }

    if (rs.accessToken) {
      localStorage.setItem('accessToken', rs.accessToken);
      localStorage.setItem('refreshToken', rs.refreshToken);
      toast.success('Login successful!');
      router.refresh();
      router.push('/');
    } else {
      toast.error('Login failed. Please check your credentials.');
    }
  };

  return (
    <Fieldset
      as='form'
      onSubmit={handleSubmit(onSubmit)}
      aria-label='login form'
      className={clsx(
        'my-auto border-2 p-6 border-gray-200 rounded-md shadow-md flex flex-col'
      )}
    >
      <h2 className='text-xl mb-1 font-bold'>Sign In</h2>
      <div className='text-sm mb-6'>
        <span>New to our shop? </span>
        <Link href='/register' className='text-blue-500'>
          Create an account here.
        </Link>
      </div>

      <div className='flex justify-center mb-4'>
        <div className='inline-flex rounded-md shadow-sm' role='group'>
          <button
            type='button'
            className={clsx(
              'px-4 py-2 text-sm font-medium rounded-l-lg border',
              loginMethod === 'email'
                ? 'bg-blue-500 text-white border-blue-500'
                : 'bg-white text-gray-700 border-gray-200 hover:bg-gray-100'
            )}
            onClick={(e) => {
              e.preventDefault();
              if (loginMethod !== 'email') toggleLoginMethod();
            }}
          >
            Email
          </button>
          <button
            type='button'
            className={clsx(
              'px-4 py-2 text-sm font-medium rounded-r-lg border',
              loginMethod === 'username'
                ? 'bg-blue-500 text-white border-blue-500'
                : 'bg-white text-gray-700 border-gray-200 hover:bg-gray-100'
            )}
            onClick={(e) => {
              e.preventDefault();
              if (loginMethod !== 'username') toggleLoginMethod();
            }}
          >
            Username
          </button>
        </div>
      </div>

      <div className='flex flex-col gap-4'>
        {loginMethod === 'email' ? (
          <TextField
            {...register('email')}
            placeholder='john.doe@example.com'
            type='email'
            icon={<FaEnvelope />}
            label='Email'
            error={errors.email?.message}
          />
        ) : (
          <TextField
            {...register('username')}
            placeholder='johndoe'
            type='text'
            icon={<FaUser />}
            label='Username'
            error={errors.username?.message}
          />
        )}
        <TextField
          {...register('password')}
          label='Password'
          icon={<FaLock />}
          type='password'
          placeholder='********'
          error={errors.password?.message}
        />
      </div>
      <div className='mt-3 mb-6 text-right'>
        <Link href='/forgot-password' className='text-sm text-blue-500'>
          Forgot password?
        </Link>
      </div>
      <Button
        className={'w-full btn btn-primary btn-lg'}
        type='submit'
        disabled={isSubmitting}
      >
        <span className='text-lg'>
          {isSubmitting ? 'Signing in...' : 'Sign In'}
        </span>
      </Button>
    </Fieldset>
  );
}
