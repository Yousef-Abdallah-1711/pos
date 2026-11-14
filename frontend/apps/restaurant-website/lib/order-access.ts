/**
 * Public order access tokens.
 *
 * The API returns an unguessable `tracking_token` once, when an order is created. It is required
 * (header `X-Order-Token`) to read or cancel that order through the public endpoints, so order ids
 * and numbers alone reveal nothing. Tokens are kept in localStorage on the device that placed the
 * order; other devices can still look an order up by number plus the checkout phone number.
 */

const STORAGE_KEY = 'order-access-tokens';
const MAX_REMEMBERED_ORDERS = 25;

interface TokenStore {
  byId: Record<string, string>;
  byNumber: Record<string, string>;
  order: string[]; // ids, oldest first, to cap storage
}

function read(): TokenStore {
  const empty: TokenStore = { byId: {}, byNumber: {}, order: [] };
  if (typeof window === 'undefined') return empty;
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) return empty;
    const parsed = JSON.parse(raw);
    return {
      byId: parsed?.byId ?? {},
      byNumber: parsed?.byNumber ?? {},
      order: Array.isArray(parsed?.order) ? parsed.order : [],
    };
  } catch {
    return empty;
  }
}

function write(store: TokenStore): void {
  if (typeof window === 'undefined') return;
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(store));
  } catch {
    // storage unavailable (private mode / quota): tracking falls back to the phone number
  }
}

/** Remember the token returned by order creation. Accepts the raw create response or its `data`. */
export function rememberOrderToken(response: unknown): void {
  const data = (response as any)?.data ?? response;
  const token: unknown = data?.tracking_token;
  if (typeof token !== 'string' || token === '') return;

  const store = read();
  const id = data?.id != null ? String(data.id) : '';
  const number = typeof data?.order_number === 'string' ? data.order_number : '';
  if (id) {
    store.byId[id] = token;
    store.order = [...store.order.filter((existing) => existing !== id), id];
  }
  if (number) store.byNumber[number] = token;

  while (store.order.length > MAX_REMEMBERED_ORDERS) {
    const dropped = store.order.shift();
    if (dropped) delete store.byId[dropped];
  }
  write(store);
}

/** Token for a public order URL such as `/public/orders/12/track` or `/public/orders/number/ORD-1/lookup`. */
export function tokenForOrderUrl(url: string | undefined): string | null {
  if (!url) return null;
  const store = read();
  const byNumber = url.match(/\/public\/orders\/number\/([^/?#]+)/);
  if (byNumber) {
    return store.byNumber[decodeURIComponent(byNumber[1])] ?? null;
  }
  const byId = url.match(/\/public\/orders\/(\d+)(?:[/?#]|$)/);
  if (byId) {
    return store.byId[byId[1]] ?? null;
  }
  return null;
}

/** True when the request creates an order. */
export function isCreateOrderRequest(method: string | undefined, url: string | undefined): boolean {
  return (method ?? '').toLowerCase() === 'post' && /\/public\/orders\/?(?:\?.*)?$/.test(url ?? '');
}

/** Bounded history contains only numeric order IDs with a device-held capability. */
export function rememberedOrderIds(): string[] {
  const store = read();
  return store.order.slice(-MAX_REMEMBERED_ORDERS).reverse().filter(id =>
    typeof id === 'string' && /^[1-9][0-9]{0,18}$/.test(id) &&
    typeof store.byId[id] === 'string' && store.byId[id].length > 0 && store.byId[id].length <= 512
  );
}
