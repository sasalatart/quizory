import { useContext } from 'react';
import { useMutation } from 'react-query';
import { QueryClientContext } from '@/providers';
import { Feedback } from '../feedback-card';

interface Props {
  onSubmit: (feedback: Feedback) => unknown;
}

export function useSubmitAnswer({ onSubmit }: Props) {
  const { answersApi } = useContext(QueryClientContext);

  const { mutateAsync: handleSubmitAnswer } = useMutation(
    async ({ choiceId }: { choiceId: string }) => {
      const { correctChoiceId, moreInfo } = await answersApi.submitAnswer({
        submitAnswerRequest: { choiceId },
      });
      await onSubmit({
        correctChoiceId,
        selectedChoiceId: choiceId,
        moreInfo,
      });
    },
  );

  return {
    handleSubmitAnswer,
  };
}
