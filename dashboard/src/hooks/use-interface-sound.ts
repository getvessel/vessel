import { useCallback, useEffect, useRef } from 'react';

let audioContext: AudioContext | null = null;

type AudioWindow = Window & typeof globalThis & { webkitAudioContext?: typeof AudioContext };

function getAudioContext() {
  if (typeof window === 'undefined') {
    return null;
  }

  const audioWindow = window as AudioWindow;
  const AudioContextConstructor = audioWindow.AudioContext || audioWindow.webkitAudioContext;

  if (!AudioContextConstructor) {
    return null;
  }

  audioContext ??= new AudioContextConstructor();

  return audioContext;
}

export function useInterfaceSound() {
  const lastPlayedAt = useRef(0);
  const canPlaySound = useRef(true);

  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)');

    const updatePreference = () => {
      canPlaySound.current = !mediaQuery.matches;
    };

    updatePreference();
    mediaQuery.addEventListener('change', updatePreference);

    return () => mediaQuery.removeEventListener('change', updatePreference);
  }, []);

  const playNavigationSound = useCallback(() => {
    if (!canPlaySound.current) {
      return;
    }

    const now = performance.now();

    if (now - lastPlayedAt.current < 80) {
      return;
    }

    lastPlayedAt.current = now;

    try {
      const context = getAudioContext();

      if (!context) {
        return;
      }

      if (context.state === 'suspended') {
        void context.resume();
      }

      const startedAt = context.currentTime;
      const oscillator = context.createOscillator();
      const gain = context.createGain();

      oscillator.type = 'sine';
      oscillator.frequency.setValueAtTime(880, startedAt);
      oscillator.frequency.exponentialRampToValueAtTime(620, startedAt + 0.06);

      gain.gain.setValueAtTime(0.0001, startedAt);
      gain.gain.exponentialRampToValueAtTime(0.035, startedAt + 0.012);
      gain.gain.exponentialRampToValueAtTime(0.0001, startedAt + 0.06);

      oscillator.connect(gain);
      gain.connect(context.destination);
      oscillator.start(startedAt);
      oscillator.stop(startedAt + 0.07);
    } catch {
      canPlaySound.current = false;
    }
  }, []);

  return { playNavigationSound };
}
