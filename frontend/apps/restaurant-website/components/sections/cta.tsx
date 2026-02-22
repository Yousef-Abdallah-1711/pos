'use client'

import React from 'react'

export interface CtaProps {
  title: string
  subtitle?: string
  buttonText?: string
  buttonLink?: string
  secondaryColor: string
  accentColor: string
}

export function Cta({ title, subtitle, buttonText, buttonLink = '#', secondaryColor, accentColor }: CtaProps) {
  return (
    <div className="cta-section p-12 md:p-20 text-white text-center">
      <h2 className="text-3xl md:text-5xl font-bold mb-4">{title}</h2>
      {subtitle && <p className="text-lg md:text-xl mb-6 opacity-90">{subtitle}</p>}
      {buttonText && (
        <a href={buttonLink}>
          <button className="cta-button px-8 py-3 rounded font-semibold hover:opacity-90 transition-opacity">{buttonText}</button>
        </a>
      )}

      <style jsx>{`
        .cta-section {
          background-color: ${secondaryColor};
        }

        .cta-button {
          background-color: ${accentColor};
          color: #000000;
        }
      `}</style>
    </div>
  )
}
