import { ChangeHandler, RefCallBack } from 'react-hook-form';

export type HookFormProps = {
  ref?: RefCallBack;
  type?: string;
  onChange: ChangeHandler;
  onBlur: ChangeHandler;
  min?: string | number;
  name: string;
  max?: string | number;
  maxLength?: number;
  minLength?: number;
  pattern?: string;
  required?: boolean;
  disabled?: boolean;
};
