import {
  Description,
  Dialog,
  DialogBackdrop,
  DialogPanel,
  DialogTitle,
} from '@headlessui/react';
import React from 'react';

interface StyledModalProps {
  title: string;
  children: React.ReactNode;
  open: boolean;
  onClose: () => void;
}
const StyledModal: React.FC<StyledModalProps> = (props) => {
  return (
    <Dialog
      as='div'
      className='relative z-10 focus:outline-none'
      onClose={props.onClose}
      open={props.open}
    >
      <DialogBackdrop className='fixed inset-0 bg-black/30' />

      <div className='fixed inset-0 z-10 w-screen overflow-y-auto'>
        <div className='flex min-h-full items-center justify-center p-4'>
          <DialogPanel className='p-4 bg-white'>
            <DialogTitle className='text-2xl font-bold text-center text-gray-800'>
              {props.title}
            </DialogTitle>
            <Description as='div'>{props.children}</Description>
          </DialogPanel>
        </div>
      </div>
    </Dialog>
  );
};

export default StyledModal;
