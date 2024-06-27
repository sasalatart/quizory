import { useContext } from 'react';
import { useMutation } from 'react-query';
import { QueryClientContext } from '@/providers';
import { Feedback } from '../feedback-card';
import { UnansweredQuestion } from '@/generated/api';

interface Props {
  question: UnansweredQuestion | undefined;
  onSubmit: (feedback: Feedback) => unknown;
}

export function useSubmitAnswer({ question, onSubmit }: Props) {
  const { answersApi } = useContext(QueryClientContext);

  const { mutateAsync: handleSubmitAnswer } = useMutation(
    async ({ choiceId }: { choiceId: string }) => {
      const { correctChoiceId, moreInfo } = await answersApi.submitAnswer({
        submitAnswerRequest: { choiceId },
      });
      await onSubmit({
        question: question!,
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
