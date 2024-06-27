import { useContext, useEffect, useState } from 'react';
import { useQuery } from 'react-query';
import { QueryClientContext, SessionContext } from '@/providers';
import { RemainingTopic } from '@/generated/api';

function getTopicToQuery(
  selectedTopic: string | undefined,
  lastAnsweredTopic: string | undefined,
  remainingTopics: RemainingTopic[],
) {
  const hasRemainingQuestionsFor = (forTopic: string) => {
    return remainingTopics.some(({ topic }) => topic === forTopic);
  };

  if (selectedTopic && hasRemainingQuestionsFor(selectedTopic)) {
    return selectedTopic;
  }

  if (lastAnsweredTopic && hasRemainingQuestionsFor(lastAnsweredTopic)) {
    return lastAnsweredTopic;
  }

  return remainingTopics?.[0]?.topic ?? lastAnsweredTopic;
}

export function useCurrentQuestion() {
  const { session } = useContext(SessionContext);
  const { answersApi, questionsApi } = useContext(QueryClientContext);
  const [selectedTopic, setSelectedTopic] = useState<string>();

  const {
    data: availableTopics,
    isLoading: isLoadingAvailableTopics,
    refetch: handleRefetchAvailableTopics,
  } = useQuery({
    queryKey: 'available-topics',
    queryFn: async () => {
      const [lastAnsweredTopic, remainingTopics] = await Promise.all([
        answersApi
          .getAnswersLog({ userId: session!.user.id, page: 0, pageSize: 1 })
          .then((log) => log[0]?.question.topic),
        questionsApi.getRemainingTopics(),
      ]);
      return { lastAnsweredTopic, remainingTopics };
    },
    onSuccess: ({ lastAnsweredTopic, remainingTopics }) => {
      if (!selectedTopic) {
        setSelectedTopic(getTopicToQuery(selectedTopic, lastAnsweredTopic, remainingTopics ?? []));
      }
    },
  });

  const lastAnsweredTopic = availableTopics?.lastAnsweredTopic;
  const remainingTopics = availableTopics?.remainingTopics;
  const {
    data: question,
    isLoading: isLoadingQuestion,
    refetch: handleRefetchCurrentQuestion,
  } = useQuery({
    queryKey: ['current-question', selectedTopic],
    queryFn: () => questionsApi.getNextQuestion({ topic: selectedTopic! }),
    refetchOnWindowFocus: false,
    enabled: !!selectedTopic,
  });

  useEffect(() => {
    const topicToQuery = getTopicToQuery(selectedTopic, lastAnsweredTopic, remainingTopics ?? []);
    setSelectedTopic(topicToQuery);
  }, [lastAnsweredTopic, remainingTopics, selectedTopic]);

  return {
    question: question ?? undefined,
    remainingTopics,
    isLoading: isLoadingQuestion || isLoadingAvailableTopics,
    handleChangeTopic: setSelectedTopic,
    handleRefetchCurrentQuestion,
    handleRefetchRemainingTopics: handleRefetchAvailableTopics,
  };
}
