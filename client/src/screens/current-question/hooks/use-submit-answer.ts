import { useContext } from 'react';
import { useMutation } from 'react-query';
import { QueryClientContext } from '@/providers';
import { Feedback } from '../feedback-card';

export function useSubmitAnswer({ onSubmit }: { onSubmit: (feedback: Feedback) => unknown }) {
  const { answersApi } = useContext(QueryClientContext);

  const { mutateAsync: handleSubmitAnswer } = useMutation(
    async ({ choiceId }: { choiceId: string }) => {
      const { correctChoiceId, moreInfo } = await answersApi.submitAnswer({
        submitAnswerRequest: { choiceId },
      });
      await onSubmit({ correctChoiceId, selectedChoiceId: choiceId, moreInfo });
    },
  );

  return {
    handleSubmitAnswer,
  };
}
