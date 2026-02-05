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
      console.log(`✅ PDF exported: ${filename}`);
    } finally {
      // Clean up
      document.body.removeChild(tempDiv);
    }
  } catch (error) {
    console.error('❌ PDF export error:', error);

    if (error instanceof Error && error.message.includes('jspdf')) {
      throw new Error('jsPDF library not installed. Run: npm install jspdf html2canvas');
    }

    throw error;
  }
}

/**
 * Export multiple orders as CSV
 */
export function exportOrdersToCSV(orders: Order[], filename: string = 'orders-export.csv') {
  try {
    // Prepare CSV headers
    const headers = [
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
        order.status,
        order.payment_status,
        order.payment_method,
        new Date(order.created_at).toLocaleString(),
      ]);

    // Build CSV content
    const csvContent = [
      headers.join(','),
      ...rows.map((row) => row.join(',')),
    ].join('\n');

    // Download CSV
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);

    link.setAttribute('href', url);
    link.setAttribute('download', filename);
    link.style.visibility = 'hidden';

    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    console.log(`✅ CSV exported: ${filename}`);
  } catch (error) {
    console.error('❌ CSV export error:', error);
    throw error;
  }
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

  const itemsHTML = order.items
    .map(
      (item) => `
    <tr>
      <td style="padding: 8px; border-bottom: 1px solid #e0e0e0;">${item.product_name}</td>
      <td style="padding: 8px; border-bottom: 1px solid #e0e0e0; text-align: center;">${item.quantity}</td>
      <td style="padding: 8px; border-bottom: 1px solid #e0e0e0; text-align: right;">$${formatAmount(item.unit_price)}</td>
      <td style="padding: 8px; border-bottom: 1px solid #e0e0e0; text-align: right;">$${formatAmount(item.subtotal)}</td>
    </tr>
    ${
      item.special_instructions
        ? `<tr>
          <td colspan="4" style="padding: 4px 8px; font-size: 12px; color: #666; font-style: italic;">
            Note: ${item.special_instructions}
          </td>
        </tr>`
        : ''
    }
  `
    )
    .join('');

  return `
    <div style="font-family: Arial, sans-serif; color: #333;">
      <!-- Header -->
      <div style="text-align: center; margin-bottom: 30px; border-bottom: 2px solid #007bff; padding-bottom: 20px;">
        <h1 style="margin: 0; color: #007bff; font-size: 28px;">INVOICE</h1>
        <p style="margin: 5px 0; color: #666;">Order #${order.order_number}</p>
      </div>

      <div style="display: flex; justify-content: space-between; margin-bottom: 30px;">
        <!-- Customer Info -->
        <div style="flex: 1;">
          <h3 style="margin: 0 0 10px 0; color: #333;">Bill To:</h3>
          <p style="margin: 5px 0;"><strong>${order.customer_name}</strong></p>
          <p style="margin: 5px 0; color: #666;">${order.customer_email}</p>
          <p style="margin: 5px 0; color: #666;">${order.customer_phone}</p>
          <p style="margin: 5px 0; color: #666;">${order.delivery_address}</p>
          ${
            order.delivery_notes
              ? `<p style="margin: 5px 0; color: #666; font-style: italic;">Note: ${order.delivery_notes}</p>`
              : ''
          }
        </div>

        <!-- Order Info -->
        <div style="flex: 1; text-align: right;">
          <p style="margin: 5px 0;"><strong>Date:</strong> ${createdDate}</p>
          <p style="margin: 5px 0;"><strong>Status:</strong> <span style="padding: 4px 8px; background-color: #e8f5e9; color: #2e7d32; border-radius: 4px;">${order.status}</span></p>
          <p style="margin: 5px 0;"><strong>Payment Method:</strong> ${order.payment_method}</p>
          <p style="margin: 5px 0;"><strong>Payment Status:</strong> <span style="padding: 4px 8px; background-color: ${
            order.payment_status === 'paid' ? '#e8f5e9' : '#fff3e0'
          }; color: ${order.payment_status === 'paid' ? '#2e7d32' : '#e65100'}; border-radius: 4px;">${order.payment_status}</span></p>
        </div>
      </div>

      <!-- Items Table -->
      <table style="width: 100%; margin-bottom: 30px; border-collapse: collapse;">
        <thead>
          <tr style="background-color: #f5f5f5; border-bottom: 2px solid #ddd;">
            <th style="padding: 10px; text-align: left; font-weight: bold;">Item</th>
            <th style="padding: 10px; text-align: center; font-weight: bold;">Qty</th>
            <th style="padding: 10px; text-align: right; font-weight: bold;">Unit Price</th>
            <th style="padding: 10px; text-align: right; font-weight: bold;">Total</th>
          </tr>
        </thead>
        <tbody>
          ${itemsHTML}
        </tbody>
      </table>

      <!-- Totals -->
      <div style="display: flex; justify-content: flex-end; margin-bottom: 30px;">
        <div style="width: 300px;">
          <div style="display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #ddd;">
            <span>Subtotal:</span>
            <span>$${formatAmount(order.subtotal)}</span>
          </div>
          ${
            order.tax > 0
              ? `
            <div style="display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #ddd;">
              <span>Tax:</span>
              <span>$${formatAmount(order.tax)}</span>
            </div>
          `
              : ''
          }
          <div style="display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #ddd;">
            <span>Delivery Fee:</span>
            <span>$${formatAmount(order.delivery_fee)}</span>
          </div>
          ${
            order.discount > 0
              ? `
            <div style="display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #ddd; color: #2e7d32;">
              <span>Discount:</span>
              <span>-$${formatAmount(order.discount)}</span>
            </div>
          `
              : ''
          }
          <div style="display: flex; justify-content: space-between; padding: 12px 0; font-size: 16px; font-weight: bold; border-top: 2px solid #333;">
            <span>Total:</span>
            <span>$${formatAmount(order.total)}</span>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div style="text-align: center; color: #666; font-size: 12px; padding-top: 20px; border-top: 1px solid #ddd;">
        <p>Thank you for your order!</p>
        <p style="margin: 5px 0;">Generated on ${new Date().toLocaleString()}</p>
      </div>
    </div>
  `;
}
