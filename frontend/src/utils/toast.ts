import toast from 'react-hot-toast';

// Success toast
export const showSuccess = (message: string): void => {
  toast.success(message, {
    duration: 3000,
    position: 'top-right',
  });
};

// Error toast
export const showError = (message: string): void => {
  toast.error(message, {
    duration: 4000,
    position: 'top-right',
  });
};

// Info toast
export const showInfo = (message: string): void => {
  toast(message, {
    duration: 3000,
    position: 'top-right',
    icon: 'ℹ️',
  });
};

// Loading toast (returns toast id for dismissal)
export const showLoading = (message: string): string => {
  return toast.loading(message, {
    position: 'top-right',
  });
};

// Dismiss specific toast
export const dismissToast = (toastId: string): void => {
  toast.dismiss(toastId);
};

// Dismiss all toasts
export const dismissAllToasts = (): void => {
  toast.dismiss();
};

// Promise toast - automatically shows loading/success/error
export const toastPromise = <T>(
  promise: Promise<T>,
  messages: {
    loading: string;
    success: string;
    error: string;
  }
): Promise<T> => {
  return toast.promise(
    promise,
    {
      loading: messages.loading,
      success: messages.success,
      error: messages.error,
    },
    {
      position: 'top-right',
    }
  );
};
