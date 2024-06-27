import { useCallback, useContext, useEffect, useState } from 'react';
import { useQuery } from 'react-query';
import { QueryClientContext, SessionContext } from '@/providers';
import { RemainingTopic, UnansweredQuestion } from '@/generated/api';

/**
 * Returns the initial topic to fetch the first question from: If the last answered topic has more
 * questions available, then it returns the last answered topic. If not, then it falls back to the
 * first available topic (if any). This helps to "remember" the last preferred topic in case the
 * user opens the app from some other device or refreshes the page.
 * TODO: this should probably be handled on the server side.
 */
function getInitialTopic(lastAnsweredTopic: string | undefined, remainingTopics: RemainingTopic[]) {
  const lastAnsweredTopicHasMoreQuestions = remainingTopics.some(
    ({ topic }) => topic === lastAnsweredTopic,
  );
  return lastAnsweredTopicHasMoreQuestions ? lastAnsweredTopic : remainingTopics?.[0]?.topic;
}

export function useCurrentQuestion() {
  const { session } = useContext(SessionContext);
  const { answersApi, questionsApi } = useContext(QueryClientContext);
  const [isLoadingQuestion, setIsLoadingQuestion] = useState(true);
  const [question, setQuestion] = useState<UnansweredQuestion | undefined>();

  const { data: lastAnsweredLogItem, isLoading: isLoadingLastAnsweredLogItem } = useQuery({
    queryKey: 'last-answered-log-item',
    queryFn: () =>
      answersApi
        .getAnswersLog({ userId: session!.user.id, page: 0, pageSize: 1 })
        .then((log) => log[0]),
    refetchOnWindowFocus: false,
  });

  const {
    data: remainingTopics,
    isLoading: isLoadingRemainingTopics,
    refetch: handleRefetchRemainingTopics,
  } = useQuery({
    queryKey: 'remaining-topics',
    queryFn: () => questionsApi.getRemainingTopics(),
    refetchOnWindowFocus: false,
  });

  const handleGetNextQuestion = useCallback(
    async (topic: string | undefined) => {
      if (!topic) {
        setQuestion(undefined);
        return;
      }

      setIsLoadingQuestion(true);
      const question = await questionsApi.getNextQuestion({ topic });
      setQuestion(question ?? undefined);
      setIsLoadingQuestion(false);
    },
    [questionsApi],
  );

  useEffect(() => {
    const lastAnsweredTopic = lastAnsweredLogItem?.question.topic;
    const fallbackTopic = getInitialTopic(lastAnsweredTopic, remainingTopics ?? []);
    if (!question && fallbackTopic) {
      handleGetNextQuestion(fallbackTopic);
    }
  }, [question, handleGetNextQuestion, remainingTopics, lastAnsweredLogItem?.question.topic]);

  return {
    question,
    remainingTopics,
    isLoading: isLoadingQuestion || isLoadingRemainingTopics || isLoadingLastAnsweredLogItem,
    handleGetNextQuestion,
    handleRefetchRemainingTopics,
  };
}
