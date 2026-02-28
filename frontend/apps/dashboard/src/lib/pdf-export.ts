/**
 * PDF Export Utilities for Orders
 *
 * REQUIRES INSTALLATION:
 * npm install jspdf html2canvas
 *
 * Usage:
 * import { exportOrderToPDF } from '@/lib/pdf-export'
 * exportOrderToPDF(orderData, 'order-123.pdf')
 */

import { debugLog, debugError } from '@/lib/debug'

type Order = {
  id: number;
  order_number: string;
  customer_name: string;
  customer_email: string;
  customer_phone: string;
  delivery_address: string;
  delivery_notes?: string;
  status: string;
  payment_status: string;
  payment_method: string;
  subtotal: number;
  tax: number;
  delivery_fee: number;
  discount: number;
  total: number;
  created_at: string;
  items: {
    id: number;
    product_name: string;
    quantity: number;
    unit_price: number;
    subtotal: number;
    special_instructions?: string;
  }[];
  status_history?: {
    id: number;
    status: string;
    notes?: string;
    changed_at: string;
    changed_by?: string;
  }[];
};

const formatAmount = (value: number | null | undefined) => {
  const numericValue = typeof value === 'number' ? value : 0;
  return Number.isFinite(numericValue) ? numericValue.toFixed(2) : '0.00';
};

/**
 * Export a single order as PDF invoice
 */
export async function exportOrderToPDF(order: Order, filename: string = `order-${order.order_number}.pdf`) {
  try {
    // Dynamic import to avoid build-time errors if jspdf not installed
    const { jsPDF } = await import('jspdf');
    const html2canvas = (await import('html2canvas')).default;

    // Create HTML content
    const htmlContent = generateOrderHTML(order);

    // Create temporary container
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = htmlContent;
    tempDiv.style.position = 'absolute';
    tempDiv.style.left = '-9999px';
    tempDiv.style.width = '800px';
    tempDiv.style.padding = '20px';
    tempDiv.style.backgroundColor = 'white';
    tempDiv.style.fontFamily = 'Arial, sans-serif';
    document.body.appendChild(tempDiv);

    try {
      // Convert HTML to canvas
      const canvas = await html2canvas(tempDiv, {
        scale: 2,
        backgroundColor: '#ffffff',
      });

      // Create PDF
      const pdf = new jsPDF({
        orientation: 'portrait',
        unit: 'mm',
        format: 'a4',
      });

      const imgData = canvas.toDataURL('image/png');
      const pageWidth = pdf.internal.pageSize.getWidth();
      const pageHeight = pdf.internal.pageSize.getHeight();
      const imgWidth = pageWidth - 20; // 10mm margins
      const imgHeight = (canvas.height * imgWidth) / canvas.width;

      let heightLeft = imgHeight;
      let position = 10;

      // Add first page
      pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
      heightLeft -= pageHeight - 20;

      // Add additional pages if needed
      while (heightLeft >= 0) {
        position = heightLeft - imgHeight;
        pdf.addPage();
        pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
        heightLeft -= pageHeight;
      }

      // Download PDF
      pdf.save(filename);
      debugLog(`✅ PDF exported: ${filename}`);
    } finally {
      // Clean up
      document.body.removeChild(tempDiv);
    }
  } catch (error) {
    debugError('❌ PDF export error:', error);

    if (error instanceof Error && error.message.includes('jspdf')) {
      throw new Error('jsPDF library not installed. Run: npm install jspdf html2canvas');
    }

    throw error;
  }
}

/**
 * Export multiple orders as CSV
 */
export function exportOrdersToCSV(
  orders: Order[],
  filename: string = 'orders-export.csv',
  locale: 'en' | 'ar' = 'en'
) {
  try {
    const isArabic = locale === 'ar'
    // Prepare CSV headers
    const headers = isArabic
      ? [
          'رقم الطلب',
          'اسم العميل',
          'البريد الإلكتروني',
          'الهاتف',
          'عنوان التوصيل',
          'عدد العناصر',
          'الإجمالي الفرعي',
          'الضريبة',
          'رسوم التوصيل',
          'الخصم',
          'الإجمالي',
          'الحالة',
          'حالة الدفع',
          'طريقة الدفع',
          'تاريخ الإنشاء',
        ]
      : [
          'Order Number',
          'Customer Name',
          'Email',
          'Phone',
          'Delivery Address',
          'Items Count',
          'Subtotal',
          'Tax',
          'Delivery Fee',
          'Discount',
          'Total',
          'Status',
          'Payment Status',
          'Payment Method',
          'Created Date',
        ];

    // Prepare CSV rows
      const rows = orders.map((order) => [
        order.order_number,
        `"${order.customer_name}"`, // Quote strings with special chars
        `"${order.customer_email}"`,
        `"${order.customer_phone}"`,
        `"${order.delivery_address}"`,
        order.items?.length || 0,
        formatAmount(order.subtotal),
        formatAmount(order.tax),
        formatAmount(order.delivery_fee),
        formatAmount(order.discount),
        formatAmount(order.total),
        localizeOrderStatus(order.status, locale),
        localizePaymentStatus(order.payment_status, locale),
        order.payment_method,
        new Date(order.created_at).toLocaleString(isArabic ? 'ar-EG' : 'en-US'),
      ]);

    // Build CSV content
    const csvContent = [
      headers.join(','),
      ...rows.map((row) => row.join(',')),
    ].join('\n');

    // Download CSV
    const bom = '\uFEFF'
    const blob = new Blob([bom + csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);

    link.setAttribute('href', url);
    link.setAttribute('download', filename);
    link.style.visibility = 'hidden';

    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    debugLog(`✅ CSV exported: ${filename}`);
  } catch (error) {
    debugError('❌ CSV export error:', error);
    throw error;
  }
}

function localizeOrderStatus(status: string, locale: 'en' | 'ar') {
  if (locale !== 'ar') return status
  const map: Record<string, string> = {
    pending: 'قيد الانتظار',
    confirmed: 'مؤكد',
    preparing: 'قيد التحضير',
    ready: 'جاهز',
    out_for_delivery: 'قيد التوصيل',
    delivered: 'تم التوصيل',
    cancelled: 'ملغى',
  }
  return map[status] || status
}

function localizePaymentStatus(status: string, locale: 'en' | 'ar') {
  if (locale !== 'ar') return status
  const map: Record<string, string> = {
    pending: 'قيد الانتظار',
    paid: 'مدفوع',
    failed: 'فشل',
    refunded: 'مسترد',
  }
  return map[status] || status
}

/**
 * Generate HTML content for PDF invoice
 */
function generateOrderHTML(order: Order): string {
  const createdDate = new Date(order.created_at).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
  const paymentStatusClass = order.payment_status === 'paid' ? 'invoice-badge-paid' : 'invoice-badge-unpaid';

  const itemsHTML = order.items
    .map(
      (item) => `
    <tr>
      <td class="invoice-cell">${item.product_name}</td>
      <td class="invoice-cell invoice-cell-center">${item.quantity}</td>
      <td class="invoice-cell invoice-cell-right">$${formatAmount(item.unit_price)}</td>
      <td class="invoice-cell invoice-cell-right">$${formatAmount(item.subtotal)}</td>
    </tr>
    ${
      item.special_instructions
        ? `<tr>
          <td colspan="4" class="invoice-item-note">
            Note: ${item.special_instructions}
          </td>
        </tr>`
        : ''
    }
  `
    )
    .join('');

  return `
    <div class="invoice-root">
      <style>
        .invoice-root {
          font-family: Arial, sans-serif;
          color: #333;
        }

        .invoice-header {
          text-align: center;
          margin-bottom: 30px;
          border-bottom: 2px solid #007bff;
          padding-bottom: 20px;
        }

        .invoice-title {
          margin: 0;
          color: #007bff;
          font-size: 28px;
        }

        .invoice-order-number {
          margin: 5px 0;
          color: #666;
        }

        .invoice-meta {
          display: flex;
          justify-content: space-between;
          margin-bottom: 30px;
        }

        .invoice-column {
          flex: 1;
        }

        .invoice-column-right {
          text-align: right;
        }

        .invoice-section-title {
          margin: 0 0 10px 0;
          color: #333;
        }

        .invoice-text {
          margin: 5px 0;
        }

        .invoice-text-muted {
          margin: 5px 0;
          color: #666;
        }

        .invoice-text-italic {
          font-style: italic;
        }

        .invoice-badge {
          display: inline-block;
          padding: 4px 8px;
          border-radius: 4px;
        }

        .invoice-badge-status {
          background-color: #e8f5e9;
          color: #2e7d32;
        }

        .invoice-badge-paid {
          background-color: #e8f5e9;
          color: #2e7d32;
        }

        .invoice-badge-unpaid {
          background-color: #fff3e0;
          color: #e65100;
        }

        .invoice-table {
          width: 100%;
          margin-bottom: 30px;
          border-collapse: collapse;
        }

        .invoice-table-head {
          background-color: #f5f5f5;
          border-bottom: 2px solid #ddd;
        }

        .invoice-heading {
          padding: 10px;
          font-weight: bold;
          text-align: left;
        }

        .invoice-heading-center {
          text-align: center;
        }

        .invoice-heading-right {
          text-align: right;
        }

        .invoice-cell {
          padding: 8px;
          border-bottom: 1px solid #e0e0e0;
          text-align: left;
        }

        .invoice-cell-center {
          text-align: center;
        }

        .invoice-cell-right {
          text-align: right;
        }

        .invoice-item-note {
          padding: 4px 8px;
          font-size: 12px;
          color: #666;
          font-style: italic;
        }

        .invoice-totals-wrap {
          display: flex;
          justify-content: flex-end;
          margin-bottom: 30px;
        }

        .invoice-totals {
          width: 300px;
        }

        .invoice-total-row {
          display: flex;
          justify-content: space-between;
          padding: 8px 0;
          border-bottom: 1px solid #ddd;
        }

        .invoice-total-row-discount {
          color: #2e7d32;
        }

        .invoice-total-row-final {
          display: flex;
          justify-content: space-between;
          padding: 12px 0;
          font-size: 16px;
          font-weight: bold;
          border-top: 2px solid #333;
        }

        .invoice-footer {
          text-align: center;
          color: #666;
          font-size: 12px;
          padding-top: 20px;
          border-top: 1px solid #ddd;
        }
      </style>
      <!-- Header -->
      <div class="invoice-header">
        <h1 class="invoice-title">INVOICE</h1>
        <p class="invoice-order-number">Order #${order.order_number}</p>
      </div>

      <div class="invoice-meta">
        <!-- Customer Info -->
        <div class="invoice-column">
          <h3 class="invoice-section-title">Bill To:</h3>
          <p class="invoice-text"><strong>${order.customer_name}</strong></p>
          <p class="invoice-text-muted">${order.customer_email}</p>
          <p class="invoice-text-muted">${order.customer_phone}</p>
          <p class="invoice-text-muted">${order.delivery_address}</p>
          ${
            order.delivery_notes
              ? `<p class="invoice-text-muted invoice-text-italic">Note: ${order.delivery_notes}</p>`
              : ''
          }
        </div>

        <!-- Order Info -->
        <div class="invoice-column invoice-column-right">
          <p class="invoice-text"><strong>Date:</strong> ${createdDate}</p>
          <p class="invoice-text"><strong>Status:</strong> <span class="invoice-badge invoice-badge-status">${order.status}</span></p>
          <p class="invoice-text"><strong>Payment Method:</strong> ${order.payment_method}</p>
          <p class="invoice-text"><strong>Payment Status:</strong> <span class="invoice-badge ${paymentStatusClass}">${order.payment_status}</span></p>
        </div>
      </div>

      <!-- Items Table -->
      <table class="invoice-table">
        <thead>
          <tr class="invoice-table-head">
            <th class="invoice-heading">Item</th>
            <th class="invoice-heading invoice-heading-center">Qty</th>
            <th class="invoice-heading invoice-heading-right">Unit Price</th>
            <th class="invoice-heading invoice-heading-right">Total</th>
          </tr>
        </thead>
        <tbody>
          ${itemsHTML}
        </tbody>
      </table>

      <!-- Totals -->
      <div class="invoice-totals-wrap">
        <div class="invoice-totals">
          <div class="invoice-total-row">
            <span>Subtotal:</span>
            <span>$${formatAmount(order.subtotal)}</span>
          </div>
          ${
            order.tax > 0
              ? `
            <div class="invoice-total-row">
              <span>Tax:</span>
              <span>$${formatAmount(order.tax)}</span>
            </div>
          `
              : ''
          }
          <div class="invoice-total-row">
            <span>Delivery Fee:</span>
            <span>$${formatAmount(order.delivery_fee)}</span>
          </div>
          ${
            order.discount > 0
              ? `
            <div class="invoice-total-row invoice-total-row-discount">
              <span>Discount:</span>
              <span>-$${formatAmount(order.discount)}</span>
            </div>
          `
              : ''
          }
          <div class="invoice-total-row-final">
            <span>Total:</span>
            <span>$${formatAmount(order.total)}</span>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="invoice-footer">
        <p>Thank you for your order!</p>
        <p class="invoice-text">Generated on ${new Date().toLocaleString()}</p>
      </div>
    </div>
  `;
}
