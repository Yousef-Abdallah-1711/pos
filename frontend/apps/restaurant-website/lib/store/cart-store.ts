import { create } from "zustand";
import { persist } from "zustand/middleware";

export interface CartItem {
  id: number;
  productId: number;
  name: string;
  price: number;
  quantity: number;
  trackInventory?: boolean;
  quantityInStock?: number;
  image?: string;
  notes?: string;
  total: number;
}

export interface AddToCartProduct {
  id: number;
  name: string;
  price: number;
  image?: string;
  track_inventory?: boolean;
  quantity_in_stock?: number;
}

interface CartStore {
  items: CartItem[];
  addItem: (product: AddToCartProduct) => void;
  removeItem: (productId: number) => void;
  updateQuantity: (productId: number, quantity: number) => void;
  updateNotes: (productId: number, notes: string) => void;
  clearCart: () => void;
  getTotal: () => number;
  getTotalItems: () => number;
}

/**
 * Zustand cart store with localStorage persistence
 */
export const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      items: [],

      addItem: (product) => {
        set((state) => {
          const existingItem = state.items.find((item) => item.productId === product.id);
          const quantityInStock = Number(
            product.quantity_in_stock ?? existingItem?.quantityInStock ?? 0
          );
          const trackInventory =
            typeof product.track_inventory === 'boolean'
              ? product.track_inventory
              : existingItem?.trackInventory ?? false;
          if (trackInventory && quantityInStock <= 0) {
            return state;
          }

          if (existingItem) {
            // Increase quantity - create new array with updated item
            const nextQuantity = existingItem.quantity + 1;
            const cappedQuantity =
              trackInventory && quantityInStock > 0
                ? Math.min(nextQuantity, quantityInStock)
                : nextQuantity;

            if (cappedQuantity === existingItem.quantity) {
              return state;
            }

            return {
              items: state.items.map((item) =>
                item.productId === product.id
                  ? {
                    ...item,
                    quantity: cappedQuantity,
                    total: cappedQuantity * item.price,
                  }
                  : item
              ),
            };
          } else {
            // Add new item - return new state
            return {
              items: [
                ...state.items,
                {
                  id: Date.now(),
                  productId: product.id,
                  name: product.name,
                  price: product.price,
                  quantity: 1,
                  trackInventory,
                  quantityInStock,
                  image: product.image,
                  total: product.price,
                },
              ],
            };
          }
        });
      },

      removeItem: (productId) => {
        set((state) => ({
          items: state.items.filter((item) => item.productId !== productId),
        }));
      },

      updateQuantity: (productId, quantity) => {
        set((state) => ({
          items: state.items
            .map((item) => {
              if (item.productId !== productId) return item;

              if (quantity <= 0) {
                return null;
              }

              const available = Number(item.quantityInStock ?? 0);
              const trackInventory = item.trackInventory ?? false;
              const targetQuantity =
                trackInventory && available > 0 ? Math.min(quantity, available) : quantity;

              if (targetQuantity <= 0) {
                return null;
              }

              return {
                ...item,
                quantity: targetQuantity,
                total: targetQuantity * item.price,
              };
            })
            .filter((item): item is CartItem => item !== null),
        }));
      },

      updateNotes: (productId, notes) => {
        set((state) => ({
          items: state.items.map((item) =>
            item.productId === productId
              ? {
                ...item,
                notes,
              }
              : item
          ),
        }));
      },

      clearCart: () => {
        set({ items: [] });
      },

      getTotal: () => {
        return get().items.reduce((total, item) => total + item.total, 0);
      },

      getTotalItems: () => {
        return get().items.reduce((total, item) => total + item.quantity, 0);
      },
    }),
    {
      name: "cart-store",
      partialize: (state) => ({
        items: state.items,
      }),
    }
  )
);

/**
 * Optimized selectors to prevent unnecessary re-renders
 * Use these selectors instead of the entire store in components
 */
export const selectCartItems = (state: CartStore) => state.items;
export const selectCartTotal = (state: CartStore) => state.getTotal();
export const selectCartItemsCount = (state: CartStore) => state.getTotalItems();
export const selectCartItemByProductId = (productId: number) => (state: CartStore) =>
  state.items.find((item) => item.productId === productId);
