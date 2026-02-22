"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { debugError } from "@/lib/debug";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button, Input } from "@pos-saas/ui";
import { CheckCircle } from "lucide-react";
import { useThemeColors } from "@/lib/hooks/use-theme";

// Validation schemas
const checkoutSchema = z.object({
  customerName: z.string().min(2, "Name must be at least 2 characters"),
  customerEmail: z.string().email("Invalid email address"),
  customerPhone: z.string().min(10, "Phone must be at least 10 digits"),
  deliveryAddress: z.string().min(5, "Address must be at least 5 characters"),
  deliveryNotes: z.string().optional(),
  paymentMethod: z.enum(["credit_card", "debit_card", "cash", "paypal"]),
});

type CheckoutFormData = z.infer<typeof checkoutSchema>;

interface CheckoutFormProps {
  onSubmit?: (data: CheckoutFormData) => Promise<void>;
  locale?: "en" | "ar";
  isLoading?: boolean;
  cartTotal?: number;
}

/**
 * Checkout form component with validation
 * Uses theme CSS variables for colors
 */
export function CheckoutForm({
  onSubmit,
  locale = "en",
  isLoading = false,
  cartTotal = 0,
}: CheckoutFormProps) {
  const [submitSuccess, setSubmitSuccess] = useState(false);
  const isRTL = locale === "ar";
  const themeColors = useThemeColors() as any;

  // Fallback colors if theme not loaded
  const bgColor = (themeColors as any)?.background || "#ffffff";
  const textColor = (themeColors as any)?.text || "#1f2937";
  const borderColor = (themeColors as any)?.border || "#e5e7eb";
  const primaryColor = (themeColors as any)?.primary || "#f97316";
  const mutedTextColor = "hsl(var(--theme-text) / 0.7)";

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CheckoutFormData>({
    resolver: zodResolver(checkoutSchema),
  });

  const handleFormSubmit = async (data: CheckoutFormData) => {
    try {
      await onSubmit?.(data);
      setSubmitSuccess(true);
      reset();
      setTimeout(() => setSubmitSuccess(false), 5000);
    } catch (error) {
      debugError("Checkout error:", error);
    }
  };

  const labels = {
    en: {
      customerInfo: "Customer Information",
      name: "Full Name",
      email: "Email Address",
      phone: "Phone Number",
      delivery: "Delivery Information",
      address: "Delivery Address",
      notes: "Special Delivery Instructions",
      payment: "Payment Method",
      creditCard: "Credit Card",
      debitCard: "Debit Card",
      cash: "Cash on Delivery",
      paypal: "PayPal",
      orderSummary: "Order Summary",
      subtotal: "Subtotal",
      deliveryFee: "Delivery Fee",
      total: "Total Amount",
      placeOrder: "Place Order",
      processing: "Processing...",
      successMessage: "Order placed successfully!",
    },
    ar: {
      customerInfo: "معلومات العميل",
      name: "الاسم الكامل",
      email: "عنوان البريد الإلكتروني",
      phone: "رقم الهاتف",
      delivery: "معلومات التوصيل",
      address: "عنوان التسليم",
      notes: "تعليمات توصيل خاصة",
      payment: "طريقة الدفع",
      creditCard: "بطاقة ائتمان",
      debitCard: "بطاقة خصم",
      cash: "الدفع عند الاستلام",
      paypal: "باي بال",
      orderSummary: "ملخص الطلب",
      subtotal: "المجموع الفرعي",
      deliveryFee: "رسوم التوصيل",
      total: "المبلغ الإجمالي",
      placeOrder: "تأكيد الطلب",
      processing: "جاري المعالجة...",
      successMessage: "تم تأكيد الطلب بنجاح!",
    },
  };

  const t = labels[locale as keyof typeof labels] || labels.en;

  return (
    <form
      onSubmit={handleSubmit(handleFormSubmit)}
      className={`checkout-form space-y-8 ${isRTL ? "rtl" : "ltr"}`}
    >
      {/* Success Message */}
      {submitSuccess && (
        <div className="success-message border rounded-lg p-4 flex items-center gap-3">
          <CheckCircle className="success-icon h-5 w-5 flex-shrink-0" />
          <p className="success-text">{t.successMessage}</p>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Left Column - Form */}
        <div className="lg:col-span-2 space-y-6">
          {/* Customer Information */}
          <div className={`checkout-card rounded-lg p-6 ${isRTL ? "text-right" : ""}`}>
            <h3 className="section-title text-lg font-semibold mb-4">
              {t.customerInfo}
            </h3>
            <div className="space-y-4">
              <div>
                <label className="field-label block text-sm font-medium mb-2">
                  {t.name}
                </label>
                <Input
                  {...register("customerName")}
                  placeholder="John Doe"
                  className={`checkout-input ${
                    errors.customerName ? "is-error" : ""
                  }`}
                />
                {errors.customerName && (
                  <p className="field-error text-sm mt-1">
                    {errors.customerName.message}
                  </p>
                )}
              </div>

              <div>
                <label className="field-label block text-sm font-medium mb-2">
                  {t.email}
                </label>
                <Input
                  {...register("customerEmail")}
                  type="email"
                  placeholder="john@example.com"
                  className={`checkout-input ${
                    errors.customerEmail ? "is-error" : ""
                  }`}
                />
                {errors.customerEmail && (
                  <p className="field-error text-sm mt-1">
                    {errors.customerEmail.message}
                  </p>
                )}
              </div>

              <div>
                <label className="field-label block text-sm font-medium mb-2">
                  {t.phone}
                </label>
                <Input
                  {...register("customerPhone")}
                  type="tel"
                  placeholder="+1 (555) 123-4567"
                  className={`checkout-input ${
                    errors.customerPhone ? "is-error" : ""
                  }`}
                />
                {errors.customerPhone && (
                  <p className="field-error text-sm mt-1">
                    {errors.customerPhone.message}
                  </p>
                )}
              </div>
            </div>
          </div>

          {/* Delivery Information */}
          <div className={`checkout-card rounded-lg p-6 ${isRTL ? "text-right" : ""}`}>
            <h3 className="section-title text-lg font-semibold mb-4">
              {t.delivery}
            </h3>
            <div className="space-y-4">
              <div>
                <label className="field-label block text-sm font-medium mb-2">
                  {t.address}
                </label>
                <textarea
                  {...register("deliveryAddress")}
                  placeholder="123 Main Street, Apt 4B, City, State 12345"
                  rows={3}
                  className={`textarea-field w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 resize-none ${
                    errors.deliveryAddress ? "is-error" : ""
                  }`}
                />
                {errors.deliveryAddress && (
                  <p className="field-error text-sm mt-1">
                    {errors.deliveryAddress.message}
                  </p>
                )}
              </div>

              <div>
                <label className="field-label block text-sm font-medium mb-2">
                  {t.notes}
                </label>
                <textarea
                  {...register("deliveryNotes")}
                  placeholder="e.g., Ring doorbell, Leave at door, etc."
                  rows={2}
                  className="textarea-field w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 resize-none"
                />
              </div>
            </div>
          </div>

          {/* Payment Method */}
          <div className={`checkout-card rounded-lg p-6 ${isRTL ? "text-right" : ""}`}>
            <h3 className="section-title text-lg font-semibold mb-4">
              {t.payment}
            </h3>
            <div className="space-y-3">
              {[
                { value: "credit_card", label: t.creditCard },
                { value: "debit_card", label: t.debitCard },
                { value: "cash", label: t.cash },
                { value: "paypal", label: t.paypal },
              ].map((option) => (
                <label
                  key={option.value}
                  className="payment-option flex items-center p-3 border rounded-lg cursor-pointer transition-colors"
                >
                  <input
                    type="radio"
                    {...register("paymentMethod")}
                    value={option.value}
                    className="mr-3"
                  />
                  <span className="option-label font-medium">
                    {option.label}
                  </span>
                </label>
              ))}
              {errors.paymentMethod && (
                <p className="field-error text-sm mt-1">
                  {errors.paymentMethod.message}
                </p>
              )}
            </div>
          </div>
        </div>

        {/* Right Column - Order Summary */}
        <div className="lg:col-span-1">
          <div className="checkout-card rounded-lg p-6 sticky top-6">
            <h3 className="section-title text-lg font-semibold mb-4">
              {t.orderSummary}
            </h3>
            <div className="summary-rows space-y-3 mb-6 pb-6">
              <div className="summary-row flex justify-between">
                <span>{t.subtotal}</span>
                <span>${cartTotal.toFixed(2)}</span>
              </div>
              <div className="summary-row flex justify-between">
                <span>{t.deliveryFee}</span>
                <span>${(cartTotal > 50 ? 0 : 5).toFixed(2)}</span>
              </div>
            </div>

            <div className="summary-total flex justify-between text-lg font-bold mb-6">
              <span>{t.total}</span>
              <span className="summary-total-value">
                ${(cartTotal + (cartTotal > 50 ? 0 : 5)).toFixed(2)}
              </span>
            </div>

            <Button
              type="submit"
              className="w-full"
              size="lg"
              disabled={isLoading}
            >
              {isLoading ? t.processing : t.placeOrder}
            </Button>
          </div>
        </div>
      </div>
      <style jsx>{`
        .checkout-form {
          --checkout-bg: ${bgColor};
          --checkout-text: ${textColor};
          --checkout-border: ${borderColor};
          --checkout-primary: ${primaryColor};
          --checkout-muted: ${mutedTextColor};
        }

        .checkout-form.rtl {
          direction: rtl;
        }

        .checkout-card {
          background-color: var(--checkout-bg);
        }

        .section-title,
        .field-label,
        .option-label,
        .summary-total {
          color: var(--checkout-text);
        }

        .field-error {
          color: #ef4444;
        }

        .textarea-field {
          border-color: var(--checkout-border);
          color: var(--checkout-text);
          background-color: hsl(var(--theme-background) / 0.5);
        }

        .textarea-field.is-error {
          border-color: #ef4444;
        }

        .payment-option {
          border-color: var(--checkout-border);
        }

        .payment-option:hover {
          background-color: hsl(var(--theme-text) / 0.02);
        }

        .summary-rows {
          border-bottom: 1px solid var(--checkout-border);
        }

        .summary-row {
          color: var(--checkout-muted);
        }

        .summary-total-value {
          color: var(--checkout-primary);
        }

        .success-message {
          background-color: #dcfce7;
          border-color: #bbf7d0;
        }

        .success-icon {
          color: #16a34a;
        }

        .success-text {
          color: #166534;
        }

        :global(.checkout-input) {
          border-color: var(--checkout-border) !important;
          color: var(--checkout-text) !important;
          background-color: hsl(var(--theme-background) / 0.5) !important;
        }

        :global(.checkout-input.is-error) {
          border-color: #ef4444 !important;
        }
      `}</style>
    </form>
  );
}
