import { Button, Dialog, DialogPanel, DialogTitle } from '@headlessui/react';
import clsx from 'clsx';
import React from 'react';

interface ConfirmDialogProps {
  title: string;
  message: string;
  open: boolean;
  confirmStyle?: string;
  onClose: () => void;
  onConfirm: () => void;
}
export const ConfirmDialog: React.FC<ConfirmDialogProps> = ({
  message,
  onClose,
  confirmStyle,
  onConfirm,
  open,
  title,
}) => {
  return (
    <Dialog
      open={open}
      as='div'
      className='relative z-10 focus:outline-none shadow-lg'
      onClose={onClose}
    >
      <div className='fixed inset-0 z-10 w-screen overflow-y-auto'>
        <div className='flex min-h-full items-center justify-center p-4'>
          <DialogPanel
            transition
            className='w-full max-w-md rounded-xl bg-white border border-green-600 shadow-green-500 shadow-sm p-6 backdrop-blur-2xl duration-300 ease-out data-[closed]:transform-[scale(95%)] data-[closed]:opacity-0'
          >
            <DialogTitle as='h3' className='text-xl/7 font-medium text-red-600'>
              {title}
            </DialogTitle>
            <p className='mt-2 text-base/6 text-gray-800'>{message}</p>
            <div className='mt-10 flex justify-end'>
              <Button className='btn btn-secondary' onClick={onClose}>
                Cancel
              </Button>
              <Button
                onClick={onConfirm}
                className={clsx('ml-2 btn bg-button-danger', confirmStyle)}
              >
                Confirm
              </Button>
            </div>
          </DialogPanel>
        </div>
      </div>
    </Dialog>
  );
};
