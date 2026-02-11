import { NextRequest, NextResponse } from 'next/server';
import { debugLog, debugError } from '@/lib/debug';

/**
 * API endpoint for receiving Web Vitals and performance metrics
 * POST /api/metrics
 *
 * Expected payload:
 * {
 *   name: 'LCP' | 'FID' | 'CLS' | 'FCP' | 'TTFB',
 *   value: number (milliseconds),
 *   rating: 'good' | 'needs-improvement' | 'poor',
 *   delta: number,
 *   url: string,
 *   timestamp: ISO string,
 *   ...other fields
 * }
 */

// In-memory storage for metrics (development only)
// TODO: Replace with database in production
const metricsBuffer: any[] = [];
const MAX_BUFFER_SIZE = 1000;

export async function POST(request: NextRequest) {
  try {
    const metric = await request.json();

    // Validate required fields
    if (!metric.name || metric.value === undefined) {
      return NextResponse.json(
        { error: 'Missing required fields: name, value' },
        { status: 400 }
      );
    }

    // Add to buffer
    metricsBuffer.push({
      ...metric,
      receivedAt: new Date().toISOString(),
      environment: process.env.NODE_ENV,
    });

    // Keep buffer size manageable
    if (metricsBuffer.length > MAX_BUFFER_SIZE) {
      metricsBuffer.shift();
    }

    // Log in development
    if (process.env.NODE_ENV === 'development') {
      debugLog('[📊 Metric]', {
        name: metric.name,
        value: `${metric.value?.toFixed?.(2) ?? metric.value}ms`,
        rating: metric.rating,
        url: metric.url?.split('?')[0], // Remove query params
      });
    }

    // TODO: In production, send to external service:
    // - Google Analytics: gtag('event', 'web_vital', { ... })
    // - Datadog: Send to datadog collector
    // - CloudWatch: Send to AWS CloudWatch
    // - Custom backend: Save to database

    return NextResponse.json({ success: true });
  } catch (error) {
    debugError('[❌ Metrics Error]', error);
    return NextResponse.json(
      { error: 'Failed to process metric' },
      { status: 500 }
    );
  }
}

/**
 * GET /api/metrics
 * Returns collected metrics (development/debugging only)
 */
export async function GET(request: NextRequest) {
  // Only allow in development
  if (process.env.NODE_ENV !== 'development') {
    return NextResponse.json(
      { error: 'Not available in production' },
      { status: 403 }
    );
  }

  const format = request.nextUrl.searchParams.get('format');

  if (format === 'summary') {
    // Return aggregated metrics
    const summary = metricsBuffer.reduce((acc: any, metric: any) => {
      if (!acc[metric.name]) {
        acc[metric.name] = {
          count: 0,
          values: [],
          ratings: {},
        };
      }

      acc[metric.name].count += 1;
      acc[metric.name].values.push(metric.value);
      acc[metric.name].ratings[metric.rating] =
        (acc[metric.name].ratings[metric.rating] || 0) + 1;

      return acc;
    }, {});

    // Calculate statistics
    const stats = Object.entries(summary).reduce((acc: any, [name, data]: any) => {
      const values = data.values;
      const sorted = [...values].sort((a: number, b: number) => a - b);

      acc[name] = {
        count: data.count,
        avg: (values.reduce((a: number, b: number) => a + b, 0) / values.length).toFixed(2),
        min: Math.min(...values).toFixed(2),
        max: Math.max(...values).toFixed(2),
        median: sorted[Math.floor(sorted.length / 2)].toFixed(2),
        p95: sorted[Math.floor(sorted.length * 0.95)].toFixed(2),
        p99: sorted[Math.floor(sorted.length * 0.99)].toFixed(2),
        ratings: data.ratings,
      };

      return acc;
    }, {});

    return NextResponse.json(stats);
  }

  // Return raw metrics
  return NextResponse.json({
    count: metricsBuffer.length,
    metrics: metricsBuffer.slice(-100), // Last 100
  });
}

/**
 * DELETE /api/metrics
 * Clear metrics buffer (development only)
 */
export async function DELETE(request: NextRequest) {
  if (process.env.NODE_ENV !== 'development') {
    return NextResponse.json(
      { error: 'Not available in production' },
      { status: 403 }
    );
  }

  const cleared = metricsBuffer.length;
  metricsBuffer.length = 0;

  return NextResponse.json({
    message: `Cleared ${cleared} metrics`,
  });
}
