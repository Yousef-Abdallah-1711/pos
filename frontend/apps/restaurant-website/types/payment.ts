export type PaymentMethodType = "cash" | "paymob_online" | "screenshot_transfer";

export interface PaymentMethod {
  type: PaymentMethodType;
  label: string;
  enabled: boolean;
}

export interface PaymobCheckoutResponse {
  success: boolean;
  paymentId: string;
  iframeUrl?: string;
  status: "pending" | "paid" | "failed";
  message?: string;
}

export interface PaymentScreenshot {
  fileName: string;
  fileUrl: string;
  uploadedAt: string;
  note?: string;
}
