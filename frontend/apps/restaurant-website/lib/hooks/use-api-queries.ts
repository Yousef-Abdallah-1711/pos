import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../api-client";

// Types
export interface Category {
  id: number;
  name: string;
  description?: string;
}

export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  categoryId: number;
  image?: string;
  isAvailable: boolean;
  rating?: number;
}

export interface OrderItem {
  id?: number;
  productId: number;
  name: string;
  price: number;
  quantity: number;
  notes?: string;
}

export interface CreateOrderItemPayload {
  product_id: number;
  quantity: number;
  special_instructions?: string;
}

export interface CreateOrderPayload {
  customer_name: string;
  customer_email: string;
  customer_phone: string;
  delivery_address?: string;
  items: CreateOrderItemPayload[];
  payment_method: string;
  notes?: string;
  table_id?: number;
  order_source?: 'website' | 'dine_in';
}

export interface Order {
  id: number;
  orderNumber: string;
  status: "pending" | "confirmed" | "preparing" | "ready" | "out_for_delivery" | "delivered" | "cancelled";
  totalAmount: number;
  createdAt: string;
  estimatedDelivery?: string;
  items?: OrderItem[];
  customerName?: string;
  customerEmail?: string;
  customerPhone?: string;
  deliveryAddress?: string;
}

export interface OrderResponse {
  success: boolean;
  data: Order;
  message?: string;
}

export interface OrderListResponse {
  success: boolean;
  data: Order[];
  pagination?: {
    page: number;
    limit: number;
    total: number;
    pages: number;
  };
  message?: string;
}

// Query Keys
const queryKeys = {
  categories: (slug: string) => ["categories", slug],
  products: (slug: string) => ["products", slug],
  productsByCategory: (categoryId: number, slug: string) => ["products", categoryId, slug],
  orders: () => ["orders"],
  orderById: (id: number) => ["orders", id],
  orderByNumber: (number: string) => ["orders", "number", number],
  searchProducts: (query: string, slug: string) => ["search", slug, query],
};

/**
 * Fetch all categories
 */
export function useCategories(restaurantSlug: string) {
  return useQuery({
    queryKey: queryKeys.categories(restaurantSlug),
    queryFn: async () => {
      if (!restaurantSlug) return [];
      const response = await apiClient.get<any>(
        `/public/restaurants/${restaurantSlug}/categories?lang=en`
      );
      return response?.categories || [];
    },
    enabled: !!restaurantSlug,
    staleTime: 5 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
    retry: 1,
    retryDelay: 1000,
  });
}

/**
 * Fetch all products
 */
export function useProducts(restaurantSlug: string) {
  return useQuery({
    queryKey: queryKeys.products(restaurantSlug),
    queryFn: async () => {
      if (!restaurantSlug) return [];
      const response = await apiClient.get<any>(
        `/public/restaurants/${restaurantSlug}/products?lang=en`
      );
      return response?.products || [];
    },
    enabled: !!restaurantSlug,
    staleTime: 5 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
    retry: 1,
    retryDelay: 1000,
  });
}

/**
 * Fetch products by category
 */
export function useProductsByCategory(categoryId: number, restaurantSlug: string) {
  return useQuery({
    queryKey: queryKeys.productsByCategory(categoryId, restaurantSlug),
    queryFn: async () => {
      if (!restaurantSlug) return [];
      const response = await apiClient.get<any>(
        `/public/restaurants/${restaurantSlug}/categories/${categoryId}/products?lang=en`
      );
      return response?.products || [];
    },
    enabled: !!restaurantSlug && !!categoryId,
    staleTime: 5 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
    retry: 1,
    retryDelay: 1000,
  });
}

/**
 * Search products with AbortController support
 * Previous requests are automatically cancelled when new search queries come in
 */
export function useSearchProducts(query: string, restaurantSlug: string) {
  return useQuery({
    queryKey: queryKeys.searchProducts(query, restaurantSlug),
    queryFn: async ({ signal }) => {
      if (!restaurantSlug) return [];
      const response = await apiClient.getCancellable<any>(
        `/public/restaurants/${restaurantSlug}/search?q=${encodeURIComponent(query)}&lang=en`,
        signal
      );
      return response?.products || [];
    },
    enabled: !!restaurantSlug && query.length > 0,
    staleTime: 2 * 60 * 1000,
    gcTime: 15 * 60 * 1000,
    retry: 0,
  });
}

/**
 * Fetch all orders (customer orders)
 */
export function useOrders() {
  return useQuery({
    queryKey: queryKeys.orders(),
    queryFn: async () => {
      const response = await apiClient.get<OrderListResponse>("/public/orders");
      return response;
    },
    staleTime: 1 * 60 * 1000,
    gcTime: 10 * 60 * 1000,
    retry: 1,
    retryDelay: 1000,
  });
}

/**
 * Fetch order by ID
 */
export function useOrder(orderId: number) {
  return useQuery({
    queryKey: queryKeys.orderById(orderId),
    queryFn: async () => {
      const response = await apiClient.get<OrderResponse>(`/public/orders/${orderId}`);
      return response;
    },
    enabled: !!orderId,
    staleTime: 30 * 1000,
    gcTime: 5 * 60 * 1000,
    retry: 1,
    retryDelay: 1000,
  });
}

/**
 * Fetch order by order number
 */
export function useOrderByNumber(orderNumber: string) {
  return useQuery({
    queryKey: queryKeys.orderByNumber(orderNumber),
    queryFn: async () => {
      const response = await apiClient.get<OrderResponse>(
        `/public/orders/number/${orderNumber}/lookup`,
      );
      return response;
    },
    enabled: !!orderNumber && orderNumber.length > 0,
    staleTime: 30 * 1000,
    gcTime: 5 * 60 * 1000,
    retry: 1,
    retryDelay: 1000,
  });
}

/**
 * Create new order mutation
 */
export function useCreateOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (orderData: CreateOrderPayload) => {
      const response = await apiClient.post<OrderResponse>("/public/orders", orderData);
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.orders() });
    },
  });
}

/**
 * Get order status
 */
export function useOrderStatus(orderId: number) {
  return useQuery({
    queryKey: ["orderStatus", orderId],
    queryFn: async () => {
      const response = await apiClient.get<{
        success: boolean;
        data: { status: string; updatedAt: string };
      }>(`/public/orders/${orderId}/status`);
      return response;
    },
    enabled: !!orderId,
    staleTime: 15 * 1000,
    gcTime: 1 * 60 * 1000,
    refetchInterval: 10 * 1000,
    retry: 1,
    retryDelay: 1000,
  });
}

/**
 * Track order with full history
 */
export function useTrackOrder(orderId: number) {
  return useQuery({
    queryKey: ["trackOrder", orderId],
    queryFn: async () => {
      const response = await apiClient.get<OrderResponse>(`/public/orders/${orderId}/track`);
      return response;
    },
    enabled: !!orderId,
    staleTime: 30 * 1000,
    gcTime: 5 * 60 * 1000,
    refetchInterval: 15 * 1000,
    retry: 1,
    retryDelay: 1000,
  });
}

/**
 * Cancel order mutation
 */
export function useCancelOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (orderId: number) => {
      const response = await apiClient.delete<OrderResponse>(`/public/orders/${orderId}`);
      return response;
    },
    onSuccess: (_, orderId) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.orderById(orderId) });
      queryClient.invalidateQueries({ queryKey: queryKeys.orders() });
    },
  });
}

/**
 * Validate order before creation
 */
export function useValidateOrder() {
  return useMutation({
    mutationFn: async (orderData: any) => {
      const response = await apiClient.post<{
        success: boolean;
        data: { valid: boolean; errors?: string[] };
      }>("/public/orders/validate", orderData);
      return response;
    },
  });
}
