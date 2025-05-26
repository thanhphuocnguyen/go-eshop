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
} from '@/app/lib/definitions';
import { zodResolver } from '@hookform/resolvers/zod';
import { TextField } from '@/app/components/FormFields';
import { useRouter } from 'next/navigation';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';

export default function RegisterFormComponent() {
  const router = useRouter();
  const form = useForm<RegisterForm>({
    reValidateMode: 'onBlur',
    resolver: zodResolver(registerSchema),
  });

  const {
    register,
    handleSubmit,
    formState: { isSubmitting, errors },
  } = form;

  const onSubmit = async (body: RegisterForm) => {
    try {
      const { data, error } = await clientSideFetch<GenericResponse<unknown>>(
        PUBLIC_API_PATHS.REGISTER,
        {
          method: 'POST',
          body,
        }
      );

      if (error) {
        toast.error(error.details || 'Registration failed');
        return;
      }

      if (data) {
        toast.success('Registration successful! Please sign in.');
        router.refresh();
        router.push('/login');
      }

      // Optionally auto-login after registration
    } catch (err) {
      console.error('Registration error:', err);
      toast.error('An unexpected error occurred during registration');
    }
  };

  return (
    <Fieldset
      as='form'
      onSubmit={handleSubmit(onSubmit, (err) => {
        console.log(err);
      })}
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
          {...register('username')}
          placeholder='john_doe_123'
          type='text'
          icon={<FaEnvelope />}
          label='Username'
          error={errors.username?.message}
        />
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
            label='Phone'
            icon={<FaUser />}
            type='phone'
            placeholder='+1 234 567 890'
            {...register('address.phone')}
            error={errors.address?.phone?.message}
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
            label='District'
            icon={<FaUser />}
            type='text'
            placeholder='Los Angeles'
            {...register('address.district')}
            error={errors.address?.district?.message}
          />
          <TextField
            label='Ward'
            icon={<FaUser />}
            type='text'
            placeholder='Downtown'
            {...register('address.ward')}
            error={errors.address?.ward?.message}
          />
        </div>
      </div>
      <Button
        className={'mt-12 w-full btn btn-primary btn-lg'}
        type='submit'
        disabled={isSubmitting}
      >
        <span className='text-lg'>
          {isSubmitting ? 'Creating Account...' : 'Create Account'}
        </span>
      </Button>
    </Fieldset>
  );
}
