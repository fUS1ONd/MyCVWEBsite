import toastLib from 'react-hot-toast';

// Re-export toast for direct usage
export const toast = toastLib;

// Success toast
export const showSuccess = (message: string): void => {
  toastLib.success(message, {
    duration: 3000,
    position: 'top-right',
  });
};

// Error toast
export const showError = (message: string): void => {
  toastLib.error(message, {
    duration: 4000,
    position: 'top-right',
  });
};

// Info toast
export const showInfo = (message: string): void => {
  toastLib(message, {
    duration: 3000,
    position: 'top-right',
    icon: 'ℹ️',
  });
};

// Loading toast (returns toast id for dismissal)
export const showLoading = (message: string): string => {
  return toastLib.loading(message, {
    position: 'top-right',
  });
};

// Dismiss specific toast
export const dismissToast = (toastId: string): void => {
  toastLib.dismiss(toastId);
};

// Dismiss all toasts
export const dismissAllToasts = (): void => {
  toastLib.dismiss();
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
  return toastLib.promise(
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
