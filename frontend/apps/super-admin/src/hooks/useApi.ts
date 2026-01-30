import { useCallback, useState } from "react";
import { AxiosError } from "axios";

interface UseApiOptions {
  onSuccess?: (data: unknown) => void;
  onError?: (error: AxiosError) => void;
}

export function useApi<T>(
  apiFunction: (...args: unknown[]) => Promise<T>,
  options?: UseApiOptions
) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<AxiosError | null>(null);

  const execute = useCallback(
    async (...args: unknown[]) => {
      setLoading(true);
      setError(null);
      try {
        const result = await apiFunction(...args);
        setData(result);
        options?.onSuccess?.(result);
        return result;
      } catch (err) {
        const axiosError = err as AxiosError;
        setError(axiosError);
        options?.onError?.(axiosError);
        throw axiosError;
      } finally {
        setLoading(false);
      }
    },
    [apiFunction, options]
  );

  return { data, loading, error, execute };
}
