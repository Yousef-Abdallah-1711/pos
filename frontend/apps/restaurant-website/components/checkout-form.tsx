"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { debugError } from "@/lib/debug";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button, Input } from "@pos-saas/ui";
import { AlertCircle, CheckCircle } from "lucide-react";
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
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-8">
      {/* Success Message */}
      {submitSuccess && (
        <div
          className="border rounded-lg p-4 flex items-center gap-3"
          style={{
            backgroundColor: "#dcfce7",
            borderColor: "#bbf7d0",
          }}
        >
          <CheckCircle
            className="h-5 w-5 flex-shrink-0"
            style={{ color: "#16a34a" }}
          />
          <p style={{ color: "#166534" }}>{t.successMessage}</p>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Left Column - Form */}
        <div className="lg:col-span-2 space-y-6">
          {/* Customer Information */}
          <div
            className={`rounded-lg p-6 ${isRTL ? "text-right" : ""}`}
            style={{ backgroundColor: bgColor }}
          >
            <h3
              className="text-lg font-semibold mb-4"
              style={{ color: textColor }}
            >
              {t.customerInfo}
            </h3>
            <div className="space-y-4">
              <div>
                <label
                  className="block text-sm font-medium mb-2"
                  style={{ color: textColor }}
                >
                  {t.name}
                </label>
                <Input
                  {...register("customerName")}
                  placeholder="John Doe"
                  style={{
                    borderColor: errors.customerName ? "#ef4444" : borderColor,
                    color: textColor,
                    backgroundColor: `hsl(var(--theme-background) / 0.5)`,
                  }}
                />
                {errors.customerName && (
                  <p style={{ color: "#ef4444" }} className="text-sm mt-1">
                    {errors.customerName.message}
                  </p>
                )}
              </div>

              <div>
                <label
                  className="block text-sm font-medium mb-2"
                  style={{ color: textColor }}
                >
                  {t.email}
                </label>
                <Input
                  {...register("customerEmail")}
                  type="email"
                  placeholder="john@example.com"
                  style={{
                    borderColor: errors.customerEmail ? "#ef4444" : borderColor,
                    color: textColor,
                    backgroundColor: `hsl(var(--theme-background) / 0.5)`,
                  }}
                />
                {errors.customerEmail && (
                  <p style={{ color: "#ef4444" }} className="text-sm mt-1">
                    {errors.customerEmail.message}
                  </p>
                )}
              </div>

              <div>
                <label
                  className="block text-sm font-medium mb-2"
                  style={{ color: textColor }}
                >
                  {t.phone}
                </label>
                <Input
                  {...register("customerPhone")}
                  type="tel"
                  placeholder="+1 (555) 123-4567"
                  style={{
                    borderColor: errors.customerPhone ? "#ef4444" : borderColor,
                    color: textColor,
                    backgroundColor: `hsl(var(--theme-background) / 0.5)`,
                  }}
                />
                {errors.customerPhone && (
                  <p style={{ color: "#ef4444" }} className="text-sm mt-1">
                    {errors.customerPhone.message}
                  </p>
                )}
              </div>
            </div>
          </div>

          {/* Delivery Information */}
          <div
            className={`rounded-lg p-6 ${isRTL ? "text-right" : ""}`}
            style={{ backgroundColor: bgColor }}
          >
            <h3
              className="text-lg font-semibold mb-4"
              style={{ color: textColor }}
            >
              {t.delivery}
            </h3>
            <div className="space-y-4">
              <div>
                <label
                  className="block text-sm font-medium mb-2"
                  style={{ color: textColor }}
                >
                  {t.address}
                </label>
                <textarea
                  {...register("deliveryAddress")}
                  placeholder="123 Main Street, Apt 4B, City, State 12345"
                  rows={3}
                  style={{
                    borderColor: errors.deliveryAddress ? "#ef4444" : borderColor,
                    color: textColor,
                    backgroundColor: `hsl(var(--theme-background) / 0.5)`,
                  }}
                  className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 resize-none"
                />
                {errors.deliveryAddress && (
                  <p style={{ color: "#ef4444" }} className="text-sm mt-1">
                    {errors.deliveryAddress.message}
                  </p>
                )}
              </div>

              <div>
                <label
                  className="block text-sm font-medium mb-2"
                  style={{ color: textColor }}
                >
                  {t.notes}
                </label>
                <textarea
                  {...register("deliveryNotes")}
                  placeholder="e.g., Ring doorbell, Leave at door, etc."
                  rows={2}
                  style={{
                    borderColor: borderColor,
                    color: textColor,
                    backgroundColor: `hsl(var(--theme-background) / 0.5)`,
                  }}
                  className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 resize-none"
                />
              </div>
            </div>
          </div>

          {/* Payment Method */}
          <div
            className={`rounded-lg p-6 ${isRTL ? "text-right" : ""}`}
            style={{ backgroundColor: bgColor }}
          >
            <h3
              className="text-lg font-semibold mb-4"
              style={{ color: textColor }}
            >
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
                  className="flex items-center p-3 border rounded-lg cursor-pointer transition-colors"
                  style={{
                    borderColor: borderColor,
                  }}
                  onMouseEnter={(e) =>
                    (e.currentTarget.style.backgroundColor = `hsl(var(--theme-text) / 0.02)`)
                  }
                  onMouseLeave={(e) =>
                    (e.currentTarget.style.backgroundColor = "transparent")
                  }
                >
                  <input
                    type="radio"
                    {...register("paymentMethod")}
                    value={option.value}
                    className="mr-3"
                  />
                  <span className="font-medium" style={{ color: textColor }}>
                    {option.label}
                  </span>
                </label>
              ))}
              {errors.paymentMethod && (
                <p style={{ color: "#ef4444" }} className="text-sm mt-1">
                  {errors.paymentMethod.message}
                </p>
              )}
            </div>
          </div>
        </div>

        {/* Right Column - Order Summary */}
        <div className="lg:col-span-1">
          <div
            className="rounded-lg p-6 sticky top-6"
            style={{ backgroundColor: bgColor }}
          >
            <h3
              className="text-lg font-semibold mb-4"
              style={{ color: textColor }}
            >
              {t.orderSummary}
            </h3>
            <div
              className="space-y-3 mb-6 pb-6"
              style={{
                borderBottomColor: borderColor,
                borderBottomWidth: "1px",
              }}
            >
              <div
                className="flex justify-between"
                style={{ color: `hsl(var(--theme-text) / 0.7)` }}
              >
                <span>{t.subtotal}</span>
                <span>${cartTotal.toFixed(2)}</span>
              </div>
              <div
                className="flex justify-between"
                style={{ color: `hsl(var(--theme-text) / 0.7)` }}
              >
                <span>{t.deliveryFee}</span>
                <span>${(cartTotal > 50 ? 0 : 5).toFixed(2)}</span>
              </div>
            </div>

            <div
              className="flex justify-between text-lg font-bold mb-6"
              style={{ color: textColor }}
            >
              <span>{t.total}</span>
              <span style={{ color: primaryColor }}>
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
    </form>
  );
}
