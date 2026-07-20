import { Store } from '@tanstack/store';

export interface OnboardingState {
  currentStep: number;
  isImportModalOpen: boolean;
}

export const onboardingStore = new Store<OnboardingState>({
  currentStep: 1,
  isImportModalOpen: false,
});

export const onboardingActions = {
  setStep: (step: number) => {
    onboardingStore.setState((state) => ({
      ...state,
      currentStep: Math.min(Math.max(step, 1), 3),
    }));
  },
  nextStep: () => {
    onboardingStore.setState((state) => ({
      ...state,
      currentStep: Math.min(state.currentStep + 1, 3),
    }));
  },
  prevStep: () => {
    onboardingStore.setState((state) => ({
      ...state,
      currentStep: Math.max(state.currentStep - 1, 1),
    }));
  },
  setImportModalOpen: (isOpen: boolean) => {
    onboardingStore.setState((state) => ({
      ...state,
      isImportModalOpen: isOpen,
    }));
  },
};
