'use client';

import { Button, Fieldset } from '@headlessui/react';
import clsx from 'clsx';
import { FaLock, FaEnvelope } from 'react-icons/fa';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { TextField } from '@/components/FormFields';
import { toast } from 'react-toastify';
import { useRouter } from 'next/navigation';
import { z } from 'zod';
import { signIn } from 'next-auth/react';

// Login form schema
const loginSchema = z.object({
  email: z.string().email({ message: 'Invalid email address' }),
  password: z
    .string()
    .min(6, { message: 'Password must be at least 6 characters' }),
});

type LoginForm = z.infer<typeof loginSchema>;

export default function LoginFormComponent() {
  const router = useRouter();
  const form = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: 'admin@simple-life.com',
      password: 'secret',
    },
  });

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = form;

  const onSubmit = async (data: LoginForm) => {
    try {
      const result = await signIn('credentials', {
        redirect: false,
        email: data.email,
        password: data.password,
      });

      if (result?.error) {
        toast.error(
          <div>
            Login failed. Please check your credentials.
            <div>{result.error}</div>
          </div>
        );
        return;
      }

      toast.success('Login successful!');
      router.refresh();
      router.push('/');
    } catch (error) {
      console.error(error);
      toast.error(
        <div>
          An unexpected error occurred.
          <div>{JSON.stringify(error)}</div>
        </div>
      );
      console.error(error);
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
      <div className='flex flex-col gap-4'>
        <TextField
          {...register('email')}
          placeholder='john.doe@example.com'
          type='email'
          icon={<FaEnvelope />}
          label='Email'
          error={errors.email?.message}
        />
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
