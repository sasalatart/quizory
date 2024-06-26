import { useContext, useState } from 'react';
import { useQuery } from 'react-query';
import { QueryClientContext } from '@/providers';
import { CenteredSpinner } from '@/layout';
import { QuestionFormCard } from './form-card';
import { Feedback, QuestionFeedbackCard } from './feedback-card';
import { NoQuestionsLeftCard } from './no-questions-left-card';

export function Question(): JSX.Element {
  const [feedback, setFeedback] = useState<Feedback | undefined>();

  const { questionsApi } = useContext(QueryClientContext);

  const {
    data: question,
    isLoading: isLoadingCurrentQuestion,
    isRefetching: isRefetchingCurrentQuestion,
    refetch: refetchCurrentQuestion,
  } = useQuery('current-question', {
    queryFn: async () => {
      const remainingTopics = await questionsApi.getRemainingTopics();
      if (remainingTopics.length === 0) {
        return null;
      }

      // TODO: allow users to choose the actual topic
      return questionsApi.getNextQuestion({ topic: remainingTopics[0].topic });
    },
    refetchOnWindowFocus: false,
  });

  if (isLoadingCurrentQuestion) {
    return <CenteredSpinner />;
  }

  if (!question) {
    return <NoQuestionsLeftCard />;
  }


  const shouldShowFeedback = question?.choices.some(({ id }) => id === feedback?.selectedChoiceId);
  if (question && feedback && shouldShowFeedback) {
    return (
      <QuestionFeedbackCard
        question={question}
        feedback={feedback}
        isLoadingNext={isRefetchingCurrentQuestion}
        onNext={refetchCurrentQuestion}
      />
    );
  }

  return <QuestionFormCard question={question} onSubmit={(feedback) => setFeedback(feedback)} />;
}
