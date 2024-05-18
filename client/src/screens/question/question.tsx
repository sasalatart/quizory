import { useContext, useState } from 'react';
import { useQuery } from 'react-query';
import { QueryClientContext } from '@/providers';
import { QuestionFormCard } from './form-card';
import { Feedback, QuestionFeedbackCard } from './feedback-card';

export function Question(): JSX.Element {
  const [feedback, setFeedback] = useState<Feedback | undefined>();

  const apiClient = useContext(QueryClientContext);
  const {
    data: question,
    isLoading,
    isRefetching,
    refetch: refetchNextQuestion,
  } = useQuery('current-question', {
    queryFn: () => apiClient.getNextQuestion(),
    refetchOnWindowFocus: false,
  });

  if (isLoading || !question) {
    return <span className="loading loading-spinner loading-xs"></span>;
  }

  const shouldShowFeedback = question.choices.some(({ id }) => id === feedback?.selectedChoiceId);
  if (feedback && shouldShowFeedback) {
    return (
      <QuestionFeedbackCard
        question={question}
        feedback={feedback}
        isLoadingNext={isRefetching}
        onNext={refetchNextQuestion}
      />
    );
  }

  return <QuestionFormCard question={question} onSubmit={(feedback) => setFeedback(feedback)} />;
}
