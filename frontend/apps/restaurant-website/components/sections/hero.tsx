'use client'

import React from 'react'

export interface HeroProps {
  title: string
  subtitle: string
  buttonText?: string
  primaryColor: string
  secondaryColor: string
  accentColor: string
}

export function Hero({ title, subtitle, buttonText, primaryColor, secondaryColor, accentColor }: HeroProps) {
  return (
    <div className="hero-section p-12 md:p-20 text-white text-center">
      <h1 className="text-4xl md:text-5xl font-bold mb-4">{title}</h1>
      <p className="text-lg md:text-xl mb-6 opacity-90">{subtitle}</p>
      {buttonText && <button className="hero-button px-8 py-3 rounded font-semibold hover:opacity-90 transition-opacity">{buttonText}</button>}

      <style jsx>{`
        .hero-section {
          background: linear-gradient(135deg, ${primaryColor} 0%, ${secondaryColor} 100%);
        }

        .hero-button {
          background-color: ${accentColor};
          color: #000000;
        }
      `}</style>
    </div>
  )
}
