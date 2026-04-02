import React, { useEffect, useRef } from 'react';

export default function ConfettiLayer({ bursts }) {
  const canvasRef = useRef(null);
  const particlesRef = useRef([]);

  useEffect(() => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    let animationId = 0;

    const resize = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
    };

    const animate = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      particlesRef.current = particlesRef.current.filter((particle) => particle.life > 0);

      particlesRef.current.forEach((particle) => {
        particle.life -= 0.016;
        particle.x += particle.vx;
        particle.y += particle.vy;
        particle.vy += 0.05;
        particle.rotation += particle.spin;

        ctx.save();
        ctx.translate(particle.x, particle.y);
        ctx.rotate(particle.rotation);
        ctx.fillStyle = particle.color;
        ctx.globalAlpha = Math.max(particle.life, 0);
        ctx.fillRect(-particle.size / 2, -particle.size / 2, particle.size, particle.size * 0.6);
        ctx.restore();
      });

      animationId = requestAnimationFrame(animate);
    };

    resize();
    animate();
    window.addEventListener('resize', resize);
    return () => {
      cancelAnimationFrame(animationId);
      window.removeEventListener('resize', resize);
    };
  }, []);

  useEffect(() => {
    if (!bursts.length) return;
    const palette = ['#FF3CAC', '#FFD93D', '#00F5A0', '#00D2FF', '#ffffff'];
    const width = window.innerWidth / 2;
    const height = window.innerHeight / 3;
    const fresh = bursts.slice(-1);
    fresh.forEach(() => {
      particlesRef.current.push(
        ...Array.from({ length: 90 }, (_, index) => ({
          x: width,
          y: height,
          vx: (Math.random() - 0.5) * 12,
          vy: Math.random() * -8 - 4,
          size: Math.random() * 8 + 4,
          color: palette[index % palette.length],
          rotation: Math.random() * Math.PI,
          spin: (Math.random() - 0.5) * 0.4,
          life: 1
        }))
      );
    });
  }, [bursts]);

  return (
    <canvas
      ref={canvasRef}
      aria-hidden="true"
      style={{ position: 'fixed', inset: 0, width: '100%', height: '100%', pointerEvents: 'none', zIndex: 999 }}
    />
  );
}
