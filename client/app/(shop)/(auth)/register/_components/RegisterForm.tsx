'use client';

import { Button, Fieldset } from '@headlessui/react';
import clsx from 'clsx';
import { FaUser, FaLock, FaEnvelope, FaPhone, FaIdCard } from 'react-icons/fa';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import {
  GenericResponse,
  RegisterForm,
  registerSchema,
} from '@/lib/definitions';
import { zodResolver } from '@hookform/resolvers/zod';
import { TextField } from '@/components/FormFields';
import { useRouter } from 'next/navigation';
import { apiFetch } from '@/lib/api/api';
import { API_PATHS } from '@/lib/constants/api';
import { toast } from 'react-toastify';
import { signIn } from 'next-auth/react';

export default function RegisterFormComponent() {
  const router = useRouter();
  const form = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
  });
  const { register, handleSubmit, formState: { isSubmitting, errors } } = form;
  
  const onSubmit = async (body: RegisterForm) => {
    try {
      const { data, error } = await apiFetch<GenericResponse<unknown>>(
        API_PATHS.REGISTER,
        {
          method: 'POST',
          body,
        }
      );
      
      if (error) {
        toast.error(error.details || 'Registration failed');
        return;
      }
      
      toast.success('Registration successful! Please sign in.');
      
      // Optionally auto-login after registration
      const loginResult = await signIn('credentials', {
        redirect: false,
        email: body.email,
        password: body.password,
      });
      
      if (loginResult?.error) {
        // If auto-login fails, redirect to login page
        router.push('/login');
      } else {
        // If auto-login succeeds, redirect to homepage
        router.refresh();
        router.push('/');
      }
    } catch (err) {
      console.error('Registration error:', err);
      toast.error('An unexpected error occurred during registration');
    }
  };
  
  return (
    <Fieldset
      as='form'
      onSubmit={handleSubmit(onSubmit)}
      aria-label='register form'
      className={clsx(
        'my-auto border-2 p-6 border-gray-200 rounded-md shadow-md flex flex-col'
      )}
    >
      <h2 className='text-xl mb-1 font-bold'>Create your Account</h2>
      <div className='text-sm mb-6'>
        <span>Start your order in seconds. Already have an account? </span>
        <Link href='/login' className=' text-blue-500'>
          Login here.
        </Link>
      </div>
      <div className='grid grid-cols-2 gap-4'>
        <TextField
          {...register('email')}
          placeholder='john.doe@example.com'
          type='email'
          icon={<FaEnvelope />}
          label='Email'
          error={errors.email?.message}
        />
        <TextField
          {...register('fullname')}
          name='fullname'
          type='text'
          icon={<FaIdCard />}
          placeholder='John Doe'
          label='Full Name'
          error={errors.fullname?.message}
        />
        <TextField
          {...register('phone')}
          name='phone'
          type='text'
          icon={<FaPhone />}
          placeholder='+1 234 567 890'
          label='Phone'
          error={errors.phone?.message}
        />
        <TextField
          label='Password'
          icon={<FaLock />}
          type='password'
          placeholder='********'
          {...register('password')}
          error={errors.password?.message}
        />
        <TextField
          label='Confirm Password'
          icon={<FaLock />}
          placeholder='********'
          type='password'
          {...register('confirmPassword')}
          error={errors.confirmPassword?.message}
        />
      </div>
      <hr className='my-8' />
      <div>
        <h3 className='text-lg font-medium mb-3'>Address section</h3>
        <div className='grid grid-cols-2 gap-4'>
          <TextField
            label='Street Address'
            icon={<FaUser />}
            type='text'
            placeholder='123 Main St'
            {...register('address.street')}
            error={errors.address?.street?.message}
          />
          <TextField
            label='City'
            icon={<FaUser />}
            type='text'
            placeholder='New York'
            {...register('address.city')}
            error={errors.address?.city?.message}
          />
          <TextField
            label='State'
            icon={<FaUser />}
            type='text'
            placeholder='NY'
            {...register('address.state')}
            error={errors.address?.state?.message}
          />
          <TextField
            label='Zip Code'
            icon={<FaUser />}
            type='text'
            placeholder='10001'
            {...register('address.zipCode')}
            error={errors.address?.zipCode?.message}
          />
        </div>
      </div>
      <Button
        className={'mt-12 w-full btn btn-primary btn-lg'}
        type='submit'
        disabled={isSubmitting}
      >
        <span className='text-lg'>{isSubmitting ? 'Creating Account...' : 'Create Account'}</span>
      </Button>
    </Fieldset>
  );
}